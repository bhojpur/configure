package markup_test

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
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	cfgsvr "github.com/bhojpur/configure/pkg/markup"
)

func TestUnmatchedKeyInJsonConfigFile(t *testing.T) {
	type configStruct struct {
		Name string
	}
	type configFile struct {
		Name string
		Test string
	}
	config := configFile{Name: "test", Test: "ATest"}

	file, err := ioutil.TempFile("/tmp", "bhojpur")
	if err != nil {
		t.Fatal("Could not create temp file")
	}
	defer os.Remove(file.Name())
	defer file.Close()

	filename := file.Name()

	if err := json.NewEncoder(file).Encode(config); err == nil {

		var result configStruct

		// Do not return error when there are unmatched keys but ErrorOnUnmatchedKeys is false
		if err := cfgsvr.New(&cfgsvr.Config{}).Load(&result, filename); err != nil {
			t.Errorf("Should NOT get error when loading configuration with extra keys. Error: %v", err)
		}

		// Return an error when there are unmatched keys and ErrorOnUnmatchedKeys is true
		if err := cfgsvr.New(&cfgsvr.Config{ErrorOnUnmatchedKeys: true}).Load(&result, filename); err == nil || !strings.Contains(err.Error(), "json: unknown field") {

			t.Errorf("Should get unknown field error when loading configuration with extra keys. Instead got error: %v", err)
		}

	} else {
		t.Errorf("failed to marshal config")
	}

	// Add .json to the file name and test
	err = os.Rename(file.Name(), file.Name()+".json")
	if err != nil {
		t.Errorf("Could not add suffix to file")
	}
	filename = file.Name() + ".json"
	defer os.Remove(filename)

	var result configStruct

	// Do not return error when there are unmatched keys but ErrorOnUnmatchedKeys is false
	if err := cfgsvr.New(&cfgsvr.Config{}).Load(&result, filename); err != nil {
		t.Errorf("Should NOT get error when loading configuration with extra keys. Error: %v", err)
	}

	// Return an error when there are unmatched keys and ErrorOnUnmatchedKeys is true
	if err := cfgsvr.New(&cfgsvr.Config{ErrorOnUnmatchedKeys: true}).Load(&result, filename); err == nil || !strings.Contains(err.Error(), "json: unknown field") {

		t.Errorf("Should get unknown field error when loading configuration with extra keys. Instead got error: %v", err)
	}

}
