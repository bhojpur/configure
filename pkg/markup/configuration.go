package markup

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
	"os"
	"reflect"
	"regexp"
	"time"
)

type Configure struct {
	*Config
	configModTimes map[string]time.Time
}

type Config struct {
	Environment        string
	ENVPrefix          string
	Debug              bool
	Verbose            bool
	Silent             bool
	AutoReload         bool
	AutoReloadInterval time.Duration
	AutoReloadCallback func(config interface{})

	// In case of json files, this field will be used only when compiled with
	// go 1.10 or later.
	// This field will be ignored when compiled with go versions lower than 1.10.
	ErrorOnUnmatchedKeys bool
}

// New initialize a Configure
func New(config *Config) *Configure {
	if config == nil {
		config = &Config{}
	}

	if os.Getenv("CONFIGURE_DEBUG_MODE") != "" {
		config.Debug = true
	}

	if os.Getenv("CONFIGURE_VERBOSE_MODE") != "" {
		config.Verbose = true
	}

	if os.Getenv("CONFIGURE_SILENT_MODE") != "" {
		config.Silent = true
	}

	if config.AutoReload && config.AutoReloadInterval == 0 {
		config.AutoReloadInterval = time.Second
	}

	return &Configure{Config: config}
}

var testRegexp = regexp.MustCompile("_test|(\\.test$)")

// GetEnvironment get environment
func (configure *Configure) GetEnvironment() string {
	if configure.Environment == "" {
		if env := os.Getenv("CONFIGURE_ENV"); env != "" {
			return env
		}

		if testRegexp.MatchString(os.Args[0]) {
			return "test"
		}

		return "development"
	}
	return configure.Environment
}

// GetErrorOnUnmatchedKeys returns a boolean indicating if an error should be
// thrown if there are keys in the config file that do not correspond to the
// config struct
func (configure *Configure) GetErrorOnUnmatchedKeys() bool {
	return configure.ErrorOnUnmatchedKeys
}

// Load will unmarshal configurations to struct from files that you provide
func (configure *Configure) Load(config interface{}, files ...string) (err error) {
	defaultValue := reflect.Indirect(reflect.ValueOf(config))
	if !defaultValue.CanAddr() {
		return fmt.Errorf("Config %v should be addressable", config)
	}
	err, _ = configure.load(config, false, files...)

	if configure.Config.AutoReload {
		go func() {
			timer := time.NewTimer(configure.Config.AutoReloadInterval)
			for range timer.C {
				reflectPtr := reflect.New(reflect.ValueOf(config).Elem().Type())
				reflectPtr.Elem().Set(defaultValue)

				var changed bool
				if err, changed = configure.load(reflectPtr.Interface(), true, files...); err == nil && changed {
					reflect.ValueOf(config).Elem().Set(reflectPtr.Elem())
					if configure.Config.AutoReloadCallback != nil {
						configure.Config.AutoReloadCallback(config)
					}
				} else if err != nil {
					fmt.Printf("Failed to reload configuration from %v, got error %v\n", files, err)
				}
				timer.Reset(configure.Config.AutoReloadInterval)
			}
		}()
	}
	return
}

// ENV return environment
func ENV() string {
	return New(nil).GetEnvironment()
}

// Load will unmarshal configurations to struct from files that you provide
func Load(config interface{}, files ...string) error {
	return New(nil).Load(config, files...)
}
