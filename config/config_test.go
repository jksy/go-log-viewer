package config

import (
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/jksy/go-log-viewer/config"
)

type withTempConfigFileFunc func(filename string)

func TestLoadConfig(t *testing.T) {

	withTempConfigFile(t, func(filename string) {
		actual, err := config.LoadConfig(filename)
		if err != nil {
			t.Error(err.Error())
			return
		}

		inputs := []string{"file:///var/log/nginx/access.log", "file:///var/log/nginx/error.log"}
		expected := &config.Config{Inputs: inputs}

		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("got %v\nwant %v", expected, actual)
			return
		}
	})
}

func withTempConfigFile(t *testing.T, f withTempConfigFileFunc) {
	content := `{
  "inputs":
  [
    "file:///var/log/nginx/access.log",
    "file:///var/log/nginx/error.log"
  ]
}
  `
	file, err := ioutil.TempFile(os.TempDir(), "test")
	defer func() {
		os.Remove(file.Name())
	}()

	file.Write([]byte(content))
	fi, err := file.Stat()
	if err != nil {
		t.Error(err.Error())
		return
	}
	filename := path.Join(os.TempDir(), fi.Name())
	file.Close()
	f(filename)

}
