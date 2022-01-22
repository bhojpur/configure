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
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/bhojpur/configure/pkg/toml"
	"gopkg.in/yaml.v2"
)

type Anonymous struct {
	Description string
}

type testConfig struct {
	APPName string `default:"bhojpur" json:",omitempty"`
	Hosts   []string

	DB struct {
		Name     string
		User     string `default:"root"`
		Password string `required:"true" env:"DBPassword"`
		Port     uint   `default:"3306" json:",omitempty"`
		SSL      bool   `default:"true" json:",omitempty"`
	}

	Contacts []struct {
		Name  string
		Email string `required:"true"`
	}

	Anonymous `anonymous:"true"`

	private string
}

func generateDefaultConfig() testConfig {
	return testConfig{
		APPName: "bhojpur",
		Hosts:   []string{"http://example.org", "http://bhojpur.net"},
		DB: struct {
			Name     string
			User     string `default:"root"`
			Password string `required:"true" env:"DBPassword"`
			Port     uint   `default:"3306" json:",omitempty"`
			SSL      bool   `default:"true" json:",omitempty"`
		}{
			Name:     "bhojpur",
			User:     "bhojpur",
			Password: "bhojpur",
			Port:     3306,
			SSL:      true,
		},
		Contacts: []struct {
			Name  string
			Email string `required:"true"`
		}{
			{
				Name:  "Shashi Bhushan Rai",
				Email: "shashi.rai@bhojpur.net",
			},
		},
		Anonymous: Anonymous{
			Description: "This is an anonymous embedded struct whose environment variables should NOT include 'ANONYMOUS'",
		},
	}
}

func TestLoadNormaltestConfig(t *testing.T) {
	config := generateDefaultConfig()
	if bytes, err := json.Marshal(config); err == nil {
		if file, err := ioutil.TempFile("/tmp", "bhojpur"); err == nil {
			defer file.Close()
			defer os.Remove(file.Name())
			file.Write(bytes)

			var result testConfig
			Load(&result, file.Name())
			if !reflect.DeepEqual(result, config) {
				t.Errorf("result should equal to original configuration")
			}
		}
	} else {
		t.Errorf("failed to marshal config")
	}
}

func TestLoadtestConfigFromTomlWithExtension(t *testing.T) {
	var (
		config = generateDefaultConfig()
		buffer bytes.Buffer
	)

	if err := toml.NewEncoder(&buffer).Encode(config); err == nil {
		if file, err := ioutil.TempFile("/tmp", "bhojpur.toml"); err == nil {
			defer file.Close()
			defer os.Remove(file.Name())
			file.Write(buffer.Bytes())

			var result testConfig
			Load(&result, file.Name())
			if !reflect.DeepEqual(result, config) {
				t.Errorf("result should equal to original configuration")
			}
		}
	} else {
		t.Errorf("failed to marshal config")
	}
}

func TestLoadtestConfigFromTomlWithoutExtension(t *testing.T) {
	var (
		config = generateDefaultConfig()
		buffer bytes.Buffer
	)

	if err := toml.NewEncoder(&buffer).Encode(config); err == nil {
		if file, err := ioutil.TempFile("/tmp", "bhojpur"); err == nil {
			defer file.Close()
			defer os.Remove(file.Name())
			file.Write(buffer.Bytes())

			var result testConfig
			Load(&result, file.Name())
			if !reflect.DeepEqual(result, config) {
				t.Errorf("result should equal to original configuration")
			}
		}
	} else {
		t.Errorf("failed to marshal config")
	}
}

func TestDefaultValue(t *testing.T) {
	config := generateDefaultConfig()
	config.APPName = ""
	config.DB.Port = 0
	config.DB.SSL = false

	if bytes, err := json.Marshal(config); err == nil {
		if file, err := ioutil.TempFile("/tmp", "bhojpur"); err == nil {
			defer file.Close()
			defer os.Remove(file.Name())
			file.Write(bytes)

			var result testConfig
			Load(&result, file.Name())

			if !reflect.DeepEqual(result, generateDefaultConfig()) {
				t.Errorf("result should be set default value correctly")
			}
		}
	} else {
		t.Errorf("failed to marshal config")
	}
}

func TestMissingRequiredValue(t *testing.T) {
	config := generateDefaultConfig()
	config.DB.Password = ""

	if bytes, err := json.Marshal(config); err == nil {
		if file, err := ioutil.TempFile("/tmp", "bhojpur"); err == nil {
			defer file.Close()
			defer os.Remove(file.Name())
			file.Write(bytes)

			var result testConfig
			if err := Load(&result, file.Name()); err == nil {
				t.Errorf("Should got error when load configuration missing db password")
			}
		}
	} else {
		t.Errorf("failed to marshal config")
	}
}

func TestUnmatchedKeyInTomltestConfigFile(t *testing.T) {
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

	if err := toml.NewEncoder(file).Encode(config); err == nil {

		var result configStruct

		// Do not return error when there are unmatched keys but ErrorOnUnmatchedKeys is false
		if err := New(&Config{}).Load(&result, filename); err != nil {
			t.Errorf("Should NOT get error when loading configuration with extra keys")
		}

		// Return an error when there are unmatched keys and ErrorOnUnmatchedKeys is true
		err := New(&Config{ErrorOnUnmatchedKeys: true}).Load(&result, filename)
		if err == nil {
			t.Errorf("Should get error when loading configuration with extra keys")
		}

		// The error should be of type UnmatchedTomlKeysError
		tomlErr, ok := err.(*UnmatchedTomlKeysError)
		if !ok {
			t.Errorf("Should get UnmatchedTomlKeysError error when loading configuration with extra keys")
		}

		// The error.Keys() function should return the "Test" key
		keys := GetStringTomlKeys(tomlErr.Keys)
		if len(keys) != 1 || keys[0] != "Test" {
			t.Errorf("The UnmatchedTomlKeysError should contain the Test key")
		}

	} else {
		t.Errorf("failed to marshal config")
	}

	// Add .toml to the file name and test again
	err = os.Rename(filename, filename+".toml")
	if err != nil {
		t.Errorf("Could not add suffix to file")
	}
	filename = filename + ".toml"
	defer os.Remove(filename)

	var result configStruct

	// Do not return error when there are unmatched keys but ErrorOnUnmatchedKeys is false
	if err := New(&Config{}).Load(&result, filename); err != nil {
		t.Errorf("Should NOT get error when loading configuration with extra keys. Error: %v", err)
	}

	// Return an error when there are unmatched keys and ErrorOnUnmatchedKeys is true
	err = New(&Config{ErrorOnUnmatchedKeys: true}).Load(&result, filename)
	if err == nil {
		t.Errorf("Should get error when loading configuration with extra keys")
	}

	// The error should be of type UnmatchedTomlKeysError
	tomlErr, ok := err.(*UnmatchedTomlKeysError)
	if !ok {
		t.Errorf("Should get UnmatchedTomlKeysError error when loading configuration with extra keys")
	}

	// The error.Keys() function should return the "Test" key
	keys := GetStringTomlKeys(tomlErr.Keys)
	if len(keys) != 1 || keys[0] != "Test" {
		t.Errorf("The UnmatchedTomlKeysError should contain the Test key")
	}

}

func TestUnmatchedKeyInYamltestConfigFile(t *testing.T) {
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

	if data, err := yaml.Marshal(config); err == nil {
		file.WriteString(string(data))

		var result configStruct

		// Do not return error when there are unmatched keys but ErrorOnUnmatchedKeys is false
		if err := New(&Config{}).Load(&result, filename); err != nil {
			t.Errorf("Should NOT get error when loading configuration with extra keys. Error: %v", err)
		}

		// Return an error when there are unmatched keys and ErrorOnUnmatchedKeys is true
		if err := New(&Config{ErrorOnUnmatchedKeys: true}).Load(&result, filename); err == nil {
			t.Errorf("Should get error when loading configuration with extra keys")

			// The error should be of type *yaml.TypeError
		} else if _, ok := err.(*yaml.TypeError); !ok {
			// || !strings.Contains(err.Error(), "not found in struct") {
			t.Errorf("Error should be of type yaml.TypeError. Instead error is %v", err)
		}

	} else {
		t.Errorf("failed to marshal config")
	}

	// Add .yaml to the file name and test again
	err = os.Rename(filename, filename+".yaml")
	if err != nil {
		t.Errorf("Could not add suffix to file")
	}
	filename = filename + ".yaml"
	defer os.Remove(filename)

	var result configStruct

	// Do not return error when there are unmatched keys but ErrorOnUnmatchedKeys is false
	if err := New(&Config{}).Load(&result, filename); err != nil {
		t.Errorf("Should NOT get error when loading configuration with extra keys. Error: %v", err)
	}

	// Return an error when there are unmatched keys and ErrorOnUnmatchedKeys is true
	if err := New(&Config{ErrorOnUnmatchedKeys: true}).Load(&result, filename); err == nil {
		t.Errorf("Should get error when loading configuration with extra keys")

		// The error should be of type *yaml.TypeError
	} else if _, ok := err.(*yaml.TypeError); !ok {
		// || !strings.Contains(err.Error(), "not found in struct") {
		t.Errorf("Error should be of type yaml.TypeError. Instead error is %v", err)
	}
}

func TestLoadtestConfigurationByEnvironment(t *testing.T) {
	config := generateDefaultConfig()
	config2 := struct {
		APPName string
	}{
		APPName: "config2",
	}

	if file, err := ioutil.TempFile("/tmp", "bhojpur"); err == nil {
		defer file.Close()
		defer os.Remove(file.Name())
		configBytes, _ := yaml.Marshal(config)
		config2Bytes, _ := yaml.Marshal(config2)
		ioutil.WriteFile(file.Name()+".yaml", configBytes, 0644)
		defer os.Remove(file.Name() + ".yaml")
		ioutil.WriteFile(file.Name()+".production.yaml", config2Bytes, 0644)
		defer os.Remove(file.Name() + ".production.yaml")

		var result testConfig
		os.Setenv("CONFIGURE_ENV", "production")
		defer os.Setenv("CONFIGURE_ENV", "")
		if err := Load(&result, file.Name()+".yaml"); err != nil {
			t.Errorf("No error should happen when load configurations, but got %v", err)
		}

		var defaultConfig = generateDefaultConfig()
		defaultConfig.APPName = "config2"
		if !reflect.DeepEqual(result, defaultConfig) {
			t.Errorf("result should be load configurations by environment correctly")
		}
	}
}

func TestLoadtestConfigurationByEnvironmentSetBytestConfig(t *testing.T) {
	config := generateDefaultConfig()
	config2 := struct {
		APPName string
	}{
		APPName: "production_config2",
	}

	if file, err := ioutil.TempFile("/tmp", "bhojpur"); err == nil {
		defer file.Close()
		defer os.Remove(file.Name())
		configBytes, _ := yaml.Marshal(config)
		config2Bytes, _ := yaml.Marshal(config2)
		ioutil.WriteFile(file.Name()+".yaml", configBytes, 0644)
		defer os.Remove(file.Name() + ".yaml")
		ioutil.WriteFile(file.Name()+".production.yaml", config2Bytes, 0644)
		defer os.Remove(file.Name() + ".production.yaml")

		var result testConfig
		var bhojpur = New(&Config{Environment: "production"})
		if bhojpur.Load(&result, file.Name()+".yaml"); err != nil {
			t.Errorf("No error should happen when load configurations, but got %v", err)
		}

		var defaultConfig = generateDefaultConfig()
		defaultConfig.APPName = "production_config2"
		if !reflect.DeepEqual(result, defaultConfig) {
			t.Errorf("result should be load configurations by environment correctly")
		}

		if bhojpur.GetEnvironment() != "production" {
			t.Errorf("bhojpur's environment should be production")
		}
	}
}

func TestOverwritetestConfigurationWithEnvironmentWithDefaultPrefix(t *testing.T) {
	config := generateDefaultConfig()

	if bytes, err := json.Marshal(config); err == nil {
		if file, err := ioutil.TempFile("/tmp", "bhojpur"); err == nil {
			defer file.Close()
			defer os.Remove(file.Name())
			file.Write(bytes)
			var result testConfig
			os.Setenv("CONFIGURE_APPNAME", "config2")
			os.Setenv("CONFIGURE_HOSTS", "- http://example.org\n- http://bhojpur.net")
			os.Setenv("CONFIGURE_DB_NAME", "db_name")
			defer os.Setenv("CONFIGURE_APPNAME", "")
			defer os.Setenv("CONFIGURE_HOSTS", "")
			defer os.Setenv("CONFIGURE_DB_NAME", "")
			Load(&result, file.Name())

			var defaultConfig = generateDefaultConfig()
			defaultConfig.APPName = "config2"
			defaultConfig.Hosts = []string{"http://example.org", "http://bhojpur.net"}
			defaultConfig.DB.Name = "db_name"
			if !reflect.DeepEqual(result, defaultConfig) {
				t.Errorf("result should equal to original configuration")
			}
		}
	}
}

func TestOverwritetestConfigurationWithEnvironment(t *testing.T) {
	config := generateDefaultConfig()

	if bytes, err := json.Marshal(config); err == nil {
		if file, err := ioutil.TempFile("/tmp", "bhojpur"); err == nil {
			defer file.Close()
			defer os.Remove(file.Name())
			file.Write(bytes)
			var result testConfig
			os.Setenv("CONFIGURE_ENV_PREFIX", "app")
			os.Setenv("APP_APPNAME", "config2")
			os.Setenv("APP_DB_NAME", "db_name")
			defer os.Setenv("CONFIGURE_ENV_PREFIX", "")
			defer os.Setenv("APP_APPNAME", "")
			defer os.Setenv("APP_DB_NAME", "")
			Load(&result, file.Name())

			var defaultConfig = generateDefaultConfig()
			defaultConfig.APPName = "config2"
			defaultConfig.DB.Name = "db_name"
			if !reflect.DeepEqual(result, defaultConfig) {
				t.Errorf("result should equal to original configuration")
			}
		}
	}
}

func TestOverwritetestConfigurationWithEnvironmentThatSetBytestConfig(t *testing.T) {
	config := generateDefaultConfig()

	if bytes, err := json.Marshal(config); err == nil {
		if file, err := ioutil.TempFile("/tmp", "bhojpur"); err == nil {
			defer file.Close()
			defer os.Remove(file.Name())
			file.Write(bytes)
			os.Setenv("APP1_APPName", "config2")
			os.Setenv("APP1_DB_Name", "db_name")
			defer os.Setenv("APP1_APPName", "")
			defer os.Setenv("APP1_DB_Name", "")

			var result testConfig
			var bhojpur = New(&Config{ENVPrefix: "APP1"})
			bhojpur.Load(&result, file.Name())

			var defaultConfig = generateDefaultConfig()
			defaultConfig.APPName = "config2"
			defaultConfig.DB.Name = "db_name"
			if !reflect.DeepEqual(result, defaultConfig) {
				t.Errorf("result should equal to original configuration")
			}
		}
	}
}

func TestResetPrefixToBlank(t *testing.T) {
	config := generateDefaultConfig()

	if bytes, err := json.Marshal(config); err == nil {
		if file, err := ioutil.TempFile("/tmp", "bhojpur"); err == nil {
			defer file.Close()
			defer os.Remove(file.Name())
			file.Write(bytes)
			var result testConfig
			os.Setenv("CONFIGURE_ENV_PREFIX", "-")
			os.Setenv("APPNAME", "config2")
			os.Setenv("DB_NAME", "db_name")
			defer os.Setenv("CONFIGURE_ENV_PREFIX", "")
			defer os.Setenv("APPNAME", "")
			defer os.Setenv("DB_NAME", "")
			Load(&result, file.Name())

			var defaultConfig = generateDefaultConfig()
			defaultConfig.APPName = "config2"
			defaultConfig.DB.Name = "db_name"
			if !reflect.DeepEqual(result, defaultConfig) {
				t.Errorf("result should equal to original configuration")
			}
		}
	}
}

func TestResetPrefixToBlank2(t *testing.T) {
	config := generateDefaultConfig()

	if bytes, err := json.Marshal(config); err == nil {
		if file, err := ioutil.TempFile("/tmp", "bhojpur"); err == nil {
			defer file.Close()
			defer os.Remove(file.Name())
			file.Write(bytes)
			var result testConfig
			os.Setenv("CONFIGURE_ENV_PREFIX", "-")
			os.Setenv("APPName", "config2")
			os.Setenv("DB_Name", "db_name")
			defer os.Setenv("CONFIGURE_ENV_PREFIX", "")
			defer os.Setenv("APPName", "")
			defer os.Setenv("DB_Name", "")
			Load(&result, file.Name())

			var defaultConfig = generateDefaultConfig()
			defaultConfig.APPName = "config2"
			defaultConfig.DB.Name = "db_name"
			if !reflect.DeepEqual(result, defaultConfig) {
				t.Errorf("result should equal to original configuration")
			}
		}
	}
}

func TestReadFromEnvironmentWithSpecifiedEnvName(t *testing.T) {
	config := generateDefaultConfig()

	if bytes, err := json.Marshal(config); err == nil {
		if file, err := ioutil.TempFile("/tmp", "bhojpur"); err == nil {
			defer file.Close()
			defer os.Remove(file.Name())
			file.Write(bytes)
			var result testConfig
			os.Setenv("DBPassword", "db_password")
			defer os.Setenv("DBPassword", "")
			Load(&result, file.Name())

			var defaultConfig = generateDefaultConfig()
			defaultConfig.DB.Password = "db_password"
			if !reflect.DeepEqual(result, defaultConfig) {
				t.Errorf("result should equal to original configuration")
			}
		}
	}
}

func TestAnonymousStruct(t *testing.T) {
	config := generateDefaultConfig()

	if bytes, err := json.Marshal(config); err == nil {
		if file, err := ioutil.TempFile("/tmp", "bhojpur"); err == nil {
			defer file.Close()
			defer os.Remove(file.Name())
			file.Write(bytes)
			var result testConfig
			os.Setenv("CONFIGURE_DESCRIPTION", "environment description")
			defer os.Setenv("CONFIGURE_DESCRIPTION", "")
			Load(&result, file.Name())

			var defaultConfig = generateDefaultConfig()
			defaultConfig.Anonymous.Description = "environment description"
			if !reflect.DeepEqual(result, defaultConfig) {
				t.Errorf("result should equal to original configuration")
			}
		}
	}
}

func TestENV(t *testing.T) {
	if ENV() != "test" {
		t.Errorf("Env should be test when running `go test`, instead env is %v", ENV())
	}

	os.Setenv("CONFIGURE_ENV", "production")
	defer os.Setenv("CONFIGURE_ENV", "")
	if ENV() != "production" {
		t.Errorf("Env should be production when set it with CONFIGURE_ENV")
	}
}

type slicetestConfig struct {
	Test1 int
	Test2 []struct {
		Test2Ele1 int
		Test2Ele2 int
	}
}

func TestSliceFromEnv(t *testing.T) {
	var tc = slicetestConfig{
		Test1: 1,
		Test2: []struct {
			Test2Ele1 int
			Test2Ele2 int
		}{
			{
				Test2Ele1: 1,
				Test2Ele2: 2,
			},
			{
				Test2Ele1: 3,
				Test2Ele2: 4,
			},
		},
	}

	var result slicetestConfig
	os.Setenv("CONFIGURE_TEST1", "1")
	os.Setenv("CONFIGURE_TEST2_0_TEST2ELE1", "1")
	os.Setenv("CONFIGURE_TEST2_0_TEST2ELE2", "2")

	os.Setenv("CONFIGURE_TEST2_1_TEST2ELE1", "3")
	os.Setenv("CONFIGURE_TEST2_1_TEST2ELE2", "4")
	err := Load(&result)
	if err != nil {
		t.Fatalf("load from env err:%v", err)
	}

	if !reflect.DeepEqual(result, tc) {
		t.Fatalf("unexpected result:%+v", result)
	}
}

func TestConfigFromEnv(t *testing.T) {
	type config struct {
		LineBreakString string `required:"true"`
		Count           int64
		Slient          bool
	}

	cfg := &config{}

	os.Setenv("CONFIGURE_ENV_PREFIX", "CONFIGURE")
	os.Setenv("CONFIGURE_LineBreakString", "Line one\nLine two\nLine three\nAnd more lines")
	os.Setenv("CONFIGURE_Slient", "1")
	os.Setenv("CONFIGURE_Count", "10")
	Load(cfg)

	if os.Getenv("CONFIGURE_LineBreakString") != cfg.LineBreakString {
		t.Error("Failed to load value has line break from env")
	}

	if !cfg.Slient {
		t.Error("Failed to load bool from env")
	}

	if cfg.Count != 10 {
		t.Error("Failed to load number from env")
	}
}

type Menu struct {
	Key      string `json:"key" yaml:"key"`
	Name     string `json:"name" yaml:"name"`
	Icon     string `json:"icon" yaml:"icon"`
	Children []Menu `json:"children" yaml:"children"`
}

type MenuList struct {
	Top []Menu `json:"top"  yaml:"top"`
}

func TestLoadNestedConfig(t *testing.T) {
	adminConfig := MenuList{}
	New(&Config{Verbose: true}).Load(&adminConfig, "admin.yml")
}
