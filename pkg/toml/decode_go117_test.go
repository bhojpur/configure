//go:build go1.17
// +build go1.17

package toml

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"fmt"
	"testing"
	"testing/fstest"
)

func TestDecodeFS(t *testing.T) {
	fsys := fstest.MapFS{
		"test.toml": &fstest.MapFile{
			Data: []byte("a = 42"),
		},
	}

	var i struct{ A int }
	meta, err := DecodeFS(fsys, "test.toml", &i)
	if err != nil {
		t.Fatal(err)
	}
	have := fmt.Sprintf("%v %v %v", i, meta.Keys(), meta.Type("a"))
	want := "{42} [a] Integer"
	if have != want {
		t.Errorf("\nhave: %s\nwant: %s", have, want)
	}
}
