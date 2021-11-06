package internal

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Step struct {
	Name string `yaml:"name"`
	Exec string `yaml:"exec"`
}

type Migration struct {
	Up   []Step `yaml:"up"`
	Down []Step `yaml:"down"`
}

func NewMigration(file string) (*Migration, error) {
	migration, err := readFile(file)
	if err != nil {
		return nil, fmt.Errorf("err: %v", err)
	}

	return migration, nil
}

func readFile(file string) (*Migration, error) {
	m := Migration{}
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("%v #%v ", file, err)
	}

	err = yaml.Unmarshal(yamlFile, &m)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal: %v", err)
	}

	return &m, nil
}
