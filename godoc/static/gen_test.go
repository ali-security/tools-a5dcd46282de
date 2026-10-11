// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package static

import (
	"bytes"
	"io/ioutil"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"unicode"
)

func TestStaticIsUpToDate(t *testing.T) {
	if runtime.GOOS == "android" {
		t.Skip("files not available on android")
	}
	oldBuf, err := ioutil.ReadFile("static.go")
	if err != nil {
		t.Errorf("error while reading static.go: %v\n", err)
	}

	newBuf, err := Generate()
	if err != nil {
		t.Errorf("error while generating static.go: %v\n", err)
	}

	// gen.go stamps time.Now().Year() into the first (copyright) line, so
	// compare everything after it to keep the test from failing every new year.
	if bytes.Compare(dropFirstLine(oldBuf), dropFirstLine(newBuf)) != 0 {
		t.Error(`static.go is stale.  Run:
  $ go generate golang.org/x/tools/godoc/static
  $ git diff
to see the differences.`)

	}
}

// dropFirstLine returns b without its first line.
func dropFirstLine(b []byte) []byte {
	if parts := bytes.SplitN(b, []byte("\n"), 2); len(parts) == 2 {
		return parts[1]
	}
	return b
}

// TestAppendQuote ensures that AppendQuote produces a valid literal.
func TestAppendQuote(t *testing.T) {
	var in, out bytes.Buffer
	for r := rune(0); r < unicode.MaxRune; r++ {
		in.WriteRune(r)
	}
	appendQuote(&out, in.Bytes())
	in2, err := strconv.Unquote(out.String())
	if err != nil {
		t.Fatalf("AppendQuote produced invalid string literal: %v", err)
	}
	if got, want := in2, in.String(); got != want {
		t.Fatal("AppendQuote modified string") // no point printing got/want: huge
	}
}

// TestJQueryNoSelfClosingTagExpansion ensures the bundled jQuery does not
// rewrite self-closing tags (e.g. "<style/>") into open/close pairs before
// inserting HTML into the DOM. That expansion lets sanitized markup such as
// "<style><style/><img src=x onerror=alert(1)>" turn into executable script
// when passed to .html()/.append() (CVE-2020-11022, CVE-2020-11023).
func TestJQueryNoSelfClosingTagExpansion(t *testing.T) {
	js, ok := Files["jquery.js"]
	if !ok {
		t.Fatal("jquery.js missing from static files")
	}
	if strings.Contains(js, `"<$1></$2>"`) {
		t.Error(`jquery.js expands self-closing tags with "<$1></$2>"; HTML passed to DOM manipulation methods is not safe`)
	}
}
