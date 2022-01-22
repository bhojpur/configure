package tomltest

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
	"math"
	"reflect"
)

// CompareTOML compares the given arguments.
//
// The returned value is a copy of Test with Failure set to a (human-readable)
// description of the first element that is unequal. If both arguments are equal
// Test is returned unchanged.
//
// Reflect.DeepEqual could work here, but it won't tell us how the two
// structures are different.
func (r Test) CompareTOML(want, have interface{}) Test {
	if isTomlValue(want) {
		if !isTomlValue(have) {
			return r.fail("Type for key '%s' differs:\n"+
				"  Expected:     %[2]v (%[2]T)\n"+
				"  Your encoder: %[3]v (%[3]T)",
				r.Key, want, have)
		}

		if !deepEqual(want, have) {
			return r.fail("Values for key '%s' differ:\n"+
				"  Expected:     %[2]v (%[2]T)\n"+
				"  Your encoder: %[3]v (%[3]T)",
				r.Key, want, have)
		}
		return r
	}

	switch w := want.(type) {
	case map[string]interface{}:
		return r.cmpTOMLMap(w, have)
	case []interface{}:
		return r.cmpTOMLArrays(w, have)
	default:
		return r.fail("Unrecognized TOML structure: %T", want)
	}
}

func (r Test) cmpTOMLMap(want map[string]interface{}, have interface{}) Test {
	haveMap, ok := have.(map[string]interface{})
	if !ok {
		return r.mismatch("table", want, haveMap)
	}

	// Check that the keys of each map are equivalent.
	for k := range want {
		if _, ok := haveMap[k]; !ok {
			bunk := r.kjoin(k)
			return bunk.fail("Could not find key '%s' in encoder output", bunk.Key)
		}
	}
	for k := range haveMap {
		if _, ok := want[k]; !ok {
			bunk := r.kjoin(k)
			return bunk.fail("Could not find key '%s' in expected output", bunk.Key)
		}
	}

	// Okay, now make sure that each value is equivalent.
	for k := range want {
		if sub := r.kjoin(k).CompareTOML(want[k], haveMap[k]); sub.Failed() {
			return sub
		}
	}
	return r
}

func (r Test) cmpTOMLArrays(want []interface{}, have interface{}) Test {
	// Slice can be decoded to []interface{} for an array of primitives, or
	// []map[string]interface{} for an array of tables.
	//
	// TODO: it would be nicer if it could always decode to []interface{}?
	haveSlice, ok := have.([]interface{})
	if !ok {
		tblArray, ok := have.([]map[string]interface{})
		if !ok {
			return r.mismatch("array", want, have)
		}

		haveSlice = make([]interface{}, len(tblArray))
		for i := range tblArray {
			haveSlice[i] = tblArray[i]
		}
	}

	if len(want) != len(haveSlice) {
		return r.fail("Array lengths differ for key '%s'"+
			"  Expected:     %[2]v (len=%[4]d)\n"+
			"  Your encoder: %[3]v (len=%[5]d)",
			r.Key, want, haveSlice, len(want), len(haveSlice))
	}
	for i := 0; i < len(want); i++ {
		if sub := r.CompareTOML(want[i], haveSlice[i]); sub.Failed() {
			return sub
		}
	}
	return r
}

// reflect.DeepEqual() that deals with NaN != NaN
func deepEqual(want, have interface{}) bool {
	var wantF, haveF float64
	switch f := want.(type) {
	case float32:
		wantF = float64(f)
	case float64:
		wantF = f
	}
	switch f := have.(type) {
	case float32:
		haveF = float64(f)
	case float64:
		haveF = f
	}
	if math.IsNaN(wantF) && math.IsNaN(haveF) {
		return true
	}

	return reflect.DeepEqual(want, have)
}

func isTomlValue(v interface{}) bool {
	switch v.(type) {
	case map[string]interface{}, []interface{}:
		return false
	}
	return true
}
