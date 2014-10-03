package system

import (
	"encoding/json"
	"io/ioutil"
)

type ConfigurationDatabase struct {
	Hosts    string `json:"hosts"`
	Database string `json:"database"`
}

type Configuration struct {
	Secret       string `json:"secret"`
	PublicPath   string `json:"public_path"`
	TemplatePath string `json:"template_path"`
	Database     ConfigurationDatabase
}

func (configuration *Configuration) Load(filename string) (err error) {
	data, err := ioutil.ReadFile(filename)

	if err != nil {
		return
	}

	err = configuration.Parse(data)

	return
}

func (configuration *Configuration) Parse(data []byte) (err error) {
	err = json.Unmarshal(data, &configuration)

	return
}
