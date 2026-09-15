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

package files

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/apache/skywalking-cli/pkg/aiagent/view"
	"github.com/apache/skywalking-cli/pkg/contextkey"
)

type storedFile struct {
	name    string
	content []byte
}

// stored are the Session Data files of session "s": a small transcript; a provider file whose one
// record line is past any scanner's default buffer; a file without a final newline; and an empty
// one.
var stored = map[int64]storedFile{
	1: {"s/streams/main/transcript-20260101T000000.000000000Z-000001.sd", []byte("{\"h\":1}\n{\"ord\":1}\n{\"t\":\"end\"}\n")},
	2: {"s/provider_body/provider_body-20260101T000000.000000000Z-000002.sd", []byte(
		"{\"h\":1}\n{\"ord\":1,\"parts\":[{\"k\":\"data\",\"data\":\"" + strings.Repeat("x", 3<<20) + "\"}]}\n{\"t\":\"end\"}\n")},
	3: {"s/unknown-000003.sd", []byte("{\"h\":1}\n{\"t\":\"end\"}")},
	4: {"s/unknown-000004.sd", nil},
}

func naming(f storedFile, seq int64) string {
	sum := sha256.Sum256(f.content)
	return fmt.Sprintf(`{"file":%q,"seq":%d,"lines":%d,"bytes":%d,"digest":%q}`+"\n",
		f.name, seq, bytes.Count(f.content, []byte("\n")), len(f.content), hex.EncodeToString(sum[:]))
}

// frame writes one file the way the OAP does: its naming line, its bytes, and a newline after a
// non-empty file that does not end with one.
func frame(f storedFile, seq int64) string {
	out := naming(f, seq) + string(f.content)
	if len(f.content) > 0 && f.content[len(f.content)-1] != '\n' {
		out += "\n"
	}
	return out
}

// server answers the files route the way the OAP does. It records how many seqs each request
// carried.
func server(t *testing.T, batches *[]int, mu *sync.Mutex) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("service") != "agent" || q.Get("instance") != "sender" || q.Get("session") != "s" ||
			r.Header.Get("Accept") != MediaType {
			w.WriteHeader(http.StatusTeapot)
			return
		}
		if r.URL.Path != "/ai-agent/conversations/c/v1/files" {
			w.Header().Set("Content-Type", "application/problem+json")
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"type":"about:blank","title":"Not Found","status":404,"detail":"no round"}`))
			return
		}
		mu.Lock()
		*batches = append(*batches, len(q["seq"]))
		mu.Unlock()
		w.Header().Set("Content-Type", MediaType+"; charset=utf-8")
		for _, text := range q["seq"] {
			n, _ := strconv.ParseInt(text, 10, 64)
			if f, ok := stored[n]; ok {
				_, _ = w.Write([]byte(frame(f, n)))
			}
		}
	}))
}

func testContext(serverURL string) context.Context {
	return context.WithValue(context.Background(), contextkey.BaseURL{}, serverURL+"/graphql")
}

func TestEveryChosenFileComesBackByteForByte(t *testing.T) {
	var batches []int
	var mu sync.Mutex
	srv := server(t, &batches, &mu)
	defer srv.Close()

	got := map[string][]byte{}
	take := func(f File, content []byte) error {
		got[f.ID] = content
		return nil
	}
	if err := Read(testContext(srv.URL), "c", "agent", "sender", "s", []int64{1, 2, 3, 4, 9}, take); err != nil {
		t.Fatal(err)
	}
	if len(got) != len(stored) {
		t.Fatalf("%d files, want %d", len(got), len(stored))
	}
	for _, f := range stored {
		if !bytes.Equal(got[f.name], f.content) {
			t.Errorf("%s: %d bytes, want %d", f.name, len(got[f.name]), len(f.content))
		}
	}
}

func TestSeqsAreAskedForInBatchesUnderTheRouteLimit(t *testing.T) {
	var batches []int
	var mu sync.Mutex
	srv := server(t, &batches, &mu)
	defer srv.Close()

	many := make([]int64, 0, 2*MaxSeqs+1)
	for i := int64(0); i < 2*MaxSeqs+1; i++ {
		many = append(many, 1000+i)
	}
	if err := Read(testContext(srv.URL), "c", "agent", "sender", "s", many, func(File, []byte) error { return nil }); err != nil {
		t.Fatal(err)
	}
	if fmt.Sprint(batches) != fmt.Sprint([]int{MaxSeqs, MaxSeqs, 1}) {
		t.Fatalf("batches %v", batches)
	}
}

func TestParseRefusesWhatDoesNotMatchItsNamingLine(t *testing.T) {
	f := storedFile{"a", []byte("{\"h\":1}\n{\"t\":\"end\"}\n")}
	unended := storedFile{"b", []byte("{\"h\":1}")}
	cases := map[string]string{
		"a wrong digest":             strings.Replace(naming(f, 1), `"digest":"`, `"digest":"00`, 1) + string(f.content),
		"a short file":               naming(f, 1) + "{\"h\":1}\n",
		"no naming line":             "not json\n",
		"a cut naming line":          `{"file":"a"`,
		"no newline after unended":   naming(unended, 2) + string(unended.content),
		"a negative size":            `{"file":"a","bytes":-1,"lines":0,"digest":""}` + "\n",
		"a size no response carries": `{"file":"a","seq":1,"lines":0,"bytes":9223372036854775807,"digest":""}` + "\n",
		"another byte after unended": naming(unended, 2) + string(unended.content) + "x",
	}
	for what, body := range cases {
		if err := Parse(strings.NewReader(body), func(File, []byte) error { return nil }); err == nil {
			t.Errorf("%s: no error", what)
		}
	}
	// an empty file is followed by nothing, and the next naming line comes straight after its own
	empty := storedFile{"e", nil}
	whole := frame(f, 1) + frame(empty, 3) + frame(unended, 2)
	var got [][]byte
	if err := Parse(strings.NewReader(whole), func(_ File, content []byte) error {
		got = append(got, content)
		return nil
	}); err != nil || len(got) != 3 || len(got[1]) != 0 || !bytes.Equal(got[2], unended.content) {
		t.Errorf("a whole stream: %v, %d files", err, len(got))
	}
}

func TestTheSenderTheSessionAndASeqAreRequired(t *testing.T) {
	noop := func(File, []byte) error { return nil }
	for what, err := range map[string]error{
		"no instance": Read(context.Background(), "c", "agent", "", "s", []int64{1}, noop),
		"no session":  Read(context.Background(), "c", "agent", "sender", "", []int64{1}, noop),
		"no seq":      Read(context.Background(), "c", "agent", "sender", "s", nil, noop),
	} {
		if err == nil {
			t.Errorf("%s: no error", what)
		}
	}
}

func TestAProblemDocumentIsTheError(t *testing.T) {
	var batches []int
	var mu sync.Mutex
	srv := server(t, &batches, &mu)
	defer srv.Close()

	var problem *view.Problem
	err := Read(testContext(srv.URL), "missing", "agent", "sender", "s", []int64{1}, func(File, []byte) error { return nil })
	if !errors.As(err, &problem) || problem.Status != 404 {
		t.Fatalf("problem: %v", err)
	}
}
