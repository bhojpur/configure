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

// Command tomlv validates TOML documents and prints each key's type.

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"text/tabwriter"

	"github.com/bhojpur/configure/pkg/toml"
)

var (
	flagTypes = false
)

func init() {
	log.SetFlags(0)
	flag.BoolVar(&flagTypes, "types", flagTypes, "Show the types for every key.")
	flag.Usage = usage
	flag.Parse()
}

func usage() {
	log.Printf("Usage: %s toml-file [ toml-file ... ]\n", path.Base(os.Args[0]))
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	if flag.NArg() < 1 {
		flag.Usage()
	}
	for _, f := range flag.Args() {
		var tmp interface{}
		md, err := toml.DecodeFile(f, &tmp)
		if err != nil {
			log.Fatalf("Error in '%s': %s", f, err)
		}
		if flagTypes {
			printTypes(md)
		}
	}
}

func printTypes(md toml.MetaData) {
	tabw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for _, key := range md.Keys() {
		fmt.Fprintf(tabw, "%s%s\t%s\n",
			strings.Repeat("    ", len(key)-1), key, md.Type(key...))
	}
	tabw.Flush()
}
