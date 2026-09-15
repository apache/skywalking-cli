// Licensed to Apache Software Foundation (ASF) under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Apache Software Foundation (ASF) licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

// Package files reads chosen Session Data files of an AI agent conversation's session from the
// OAP's GET /ai-agent/conversations/{conversation}/v1/files route, on the query host beside the view
// route. A file is chosen by its session and its landed seq; there is no read of every file. The
// body is application/vnd.skywalking.asz.files+ndjson: for each stored file a line naming it, then
// exactly as many bytes as it says, the file. Files are read as they arrive, never the whole
// response at once.
package files

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/apache/skywalking-cli/pkg/aiagent/view"
	"github.com/apache/skywalking-cli/pkg/contextkey"
	"github.com/apache/skywalking-cli/pkg/transport"
)

const (
	// MediaType is the route's one body format.
	MediaType = "application/vnd.skywalking.asz.files+ndjson"
	// MaxSeqs is the most seqs one request carries, the route's own limit: the Sessionizer cuts a
	// file at 2 MiB, so one response holds about 64 MiB at most.
	MaxSeqs = 32

	defaultBaseURL = "http://127.0.0.1:12800/graphql"
)

// File is the line naming one stored file.
type File struct {
	ID     string `json:"file"`
	Seq    int64  `json:"seq"`
	Lines  int    `json:"lines"`
	Bytes  int    `json:"bytes"`
	Digest string `json:"digest"`
}

// Path is the files route of one conversation.
func Path(conversation string) string {
	return "/ai-agent/conversations/" + url.PathEscape(conversation) + "/v1/files"
}

// Read asks for the chosen files of session, at most MaxSeqs a request, and hands each stored file
// to fn with its bytes, checked against the digest its naming line gives. A seq no stored file
// answers is left out. serviceName, instanceName, session and at least one seq are required.
func Read(ctx context.Context, conversation, serviceName, instanceName, session string, seqs []int64,
	fn func(File, []byte) error) error {
	switch {
	case serviceName == "" || instanceName == "":
		return errors.New("the service and the instance are both required")
	case session == "":
		return errors.New("the session is required")
	case len(seqs) == 0:
		return errors.New("at least one seq is required")
	}
	// one client for every batch, so a long read reuses its connections instead of opening one per batch
	client := transport.HTTPClient(ctx)
	defer client.CloseIdleConnections()
	for start := 0; start < len(seqs); start += MaxSeqs {
		end := start + MaxSeqs
		if end > len(seqs) {
			end = len(seqs)
		}
		query := url.Values{"service": {serviceName}, "instance": {instanceName}, "session": {session}}
		for _, n := range seqs[start:end] {
			query.Add("seq", strconv.FormatInt(n, 10))
		}
		if err := readBatch(ctx, client, conversation, query, fn); err != nil {
			return err
		}
	}
	return nil
}

func readBatch(ctx context.Context, client *http.Client, conversation string, query url.Values, fn func(File, []byte) error) error {
	full := view.CoreURL(transport.GetValue(ctx, contextkey.BaseURL{}, defaultBaseURL)) + Path(conversation) + "?" + query.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, full, http.NoBody)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", MediaType)
	if authorization := transport.AuthHeader(ctx); authorization != "" {
		req.Header.Set("Authorization", authorization)
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return view.ReadError(resp, full)
	}
	return Parse(resp.Body, fn)
}

// Parse reads a files body: a naming line, then exactly as many bytes as it says, the file, then
// the one newline that follows a file not ending with its own. A file can be as large as the
// largest stored, so it is read by its size, never line by line with a bounded scanner.
func Parse(body io.Reader, fn func(File, []byte) error) error {
	r := bufio.NewReaderSize(body, 256*1024)
	for {
		line, err := r.ReadBytes('\n')
		if errors.Is(err, io.EOF) && len(line) == 0 {
			return nil
		}
		if err != nil {
			return fmt.Errorf("the files response ends inside a naming line: %w", err)
		}
		var f File
		if err := json.Unmarshal(line, &f); err != nil {
			return fmt.Errorf("not a naming line: %w", err)
		}
		if f.Bytes < 0 {
			return fmt.Errorf("%s: a naming line with %d bytes", f.ID, f.Bytes)
		}
		// the buffer grows with the bytes that actually arrive, so a size no response carries never allocates
		var buf bytes.Buffer
		if n, err := io.Copy(&buf, io.LimitReader(r, int64(f.Bytes))); err != nil || n != int64(f.Bytes) {
			return fmt.Errorf("%s ends after %d of its %d bytes: %v", f.ID, n, f.Bytes, err)
		}
		content := buf.Bytes()
		if f.Bytes > 0 && content[f.Bytes-1] != '\n' {
			if b, err := r.ReadByte(); err != nil || b != '\n' {
				return fmt.Errorf("%s: no newline after a file that does not end with one", f.ID)
			}
		}
		sum := sha256.Sum256(content)
		if got := hex.EncodeToString(sum[:]); got != f.Digest {
			return fmt.Errorf("%s: its bytes hash to %s, the OAP names %s", f.ID, got, f.Digest)
		}
		if err := fn(f, content); err != nil {
			return err
		}
	}
}
