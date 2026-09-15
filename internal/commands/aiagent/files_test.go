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

package aiagent

import (
	"os"
	"path/filepath"
	"testing"
)

// An export lands inside its directory whatever a file's name says: a name climbing out is refused,
// and so is one that would reach outside through a symbolic link.
func TestExportStaysInsideItsDirectory(t *testing.T) {
	dir := t.TempDir()
	outside := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(dir, "linked")); err != nil {
		t.Skip("no symbolic links here:", err)
	}
	root, err := os.OpenRoot(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer root.Close()

	if _, err := export(root, "s/streams/main/transcript-20260101T000000.000000000Z-000001.sd", []byte("x\n")); err != nil {
		t.Fatal(err)
	}
	written := filepath.Join(dir, "s", "streams", "main", "transcript-20260101T000000.000000000Z-000001.sd")
	if b, err := os.ReadFile(written); err != nil || string(b) != "x\n" {
		t.Fatalf("written: %q, %v", b, err)
	}
	for _, name := range []string{"../escaped.sd", "linked/escaped.sd", "/abs/escaped.sd"} {
		if _, err := export(root, name, []byte("x\n")); err == nil {
			t.Errorf("%s: written", name)
		}
	}
	if entries, _ := os.ReadDir(outside); len(entries) != 0 {
		t.Errorf("the outside directory holds %d entries", len(entries))
	}
}
