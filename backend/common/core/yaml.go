package core

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// ReadFileYAML はYAMLを読み込んでoutputに入れる
func ReadFileYAML(file string, output interface{}) error {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(bytes, output)
}
