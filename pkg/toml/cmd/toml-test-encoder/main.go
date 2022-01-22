package main

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

// Command toml-test-encoder satisfies the toml-test interface for testing TOML
// encoders. Namely, it accepts JSON on stdin and outputs TOML on stdout.

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"path"

	"github.com/bhojpur/configure/pkg/toml"
	"github.com/bhojpur/configure/pkg/toml/internal/tag"
)

func init() {
	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()
}

func usage() {
	log.Printf("Usage: %s < json-file\n", path.Base(os.Args[0]))
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	if flag.NArg() != 0 {
		flag.Usage()
	}

	var tmp interface{}
	if err := json.NewDecoder(os.Stdin).Decode(&tmp); err != nil {
		log.Fatalf("Error decoding JSON: %s", err)
	}

	rm, err := tag.Remove(tmp)
	if err != nil {
		log.Fatalf("Error decoding JSON: %s", err)
	}

	if err := toml.NewEncoder(os.Stdout).Encode(rm); err != nil {
		log.Fatalf("Error encoding TOML: %s", err)
	}
}
