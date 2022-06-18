package configuration

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
)

func GetConfig() *Configuration {
	var configuration Configuration
	unmarshal("./configuration.yml", &configuration)
	return &configuration
}

func unmarshal(file string, out interface{}) {
	ymlFile, err := ioutil.ReadFile(file)

	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	ymlFile = []byte(os.ExpandEnv(string(ymlFile)))
	err = yaml.Unmarshal(ymlFile, out)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
}
