package main

import (
	"fmt"
	"io/ioutil"
	"launchpad.net/goyaml"
)

type Config map[string]interface{}

func LoadConfig(file string) (Config, error) {
	yml, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var c Config
	err = goyaml.Unmarshal(yml, &c)
	if err != nil {
		return nil, err
	}

	err = c.Validate()
	if err != nil {
		return nil, err
	}

	return c, err
}

func (c Config) Validate() error {
	if _, ok := c["mysql"]; !ok {
		fmt.Errorf("mysql section of config not present")
	}

	if _, ok := c["destination"]; !ok {
		fmt.Errorf("destination section of config not present")
	}

	return nil
}
