//  Copyright (c) 2018 Couchbase, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package stempel

import (
	"bufio"
	"os"
	"strings"
	"testing"

	"github.com/rogpeppe/go-charset/charset"
	_ "github.com/rogpeppe/go-charset/data"
)

func TestEmpty(t *testing.T) {
	trie, err := Open("pl/stemmer_20000.tbl")
	if err != nil {
		t.Fatal(err)
	}

	buff := []rune("")
	diff := trie.GetLastOnPath(buff)
	if len(diff) > 0 {
		t.Fatalf("expected empty diff, got %v", diff)
	}
	buff = Diff(buff, diff)
	if len(buff) > 0 {
		t.Fatalf("expected empty buff, got %v", buff)
	}
}

// TestStemp only tests that we can successfully stem everything in the
// dictionary without crashing.  It does not attempt to assert correct output.
func TestStem(t *testing.T) {
	trie, err := Open("pl/stemmer_20000.tbl")
	if err != nil {
		t.Fatal(err)
	}

	wordfile, err := os.Open("pl/pl_PL.dic")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		cerr := wordfile.Close()
		if cerr != nil {
			t.Fatal(cerr)
		}
	}()

	cr, err := charset.NewReader("iso-8859-2", wordfile)
	if err != nil {
		t.Fatal(err)
	}

	scanner := bufio.NewScanner(cr)
	for scanner.Scan() {
		before := scanner.Text()
		hasSlash := strings.Index(before, "/")
		if hasSlash > 0 {
			before = before[0:hasSlash]
		}
		buff := []rune(before)
		diff := trie.GetLastOnPath(buff)
		_ = Diff(buff, diff)
	}

	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
}
