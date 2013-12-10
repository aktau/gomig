package main

import (
	"fmt"
	"io/ioutil"
	"launchpad.net/goyaml"
)

/* type Config map[string]interface{} */

type ConnConfig struct {
	Hostname string `yaml:"hostname,omitempty"`
	Port     int    `yaml:"port,omitempty"`
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
	Database string `yaml:"database,omitempty"`
	Compress bool   `yaml:"compress,omitempty"`
}

type DestinationConfig struct {
	File     string      `yaml:"file,omitempty"`
	Postgres *ConnConfig `yaml:"postgres,omitempty"`
}

type Config struct {
	Mysql        *ConnConfig                  `yaml:"mysql,omitempty"`
	Destination  *DestinationConfig           `yaml:"destination,omitempty"`
	Views        map[string]string            `yaml:"views,omitempty"`
	Tables       map[string]map[string]string `yaml:"tables,omitempty"`
	SuppressData bool                         `yaml:"supress_data"`
	SuppressDdl  bool                         `yaml:"supress_ddl"`
	Truncate     bool                         `yaml:"force_truncate"`
	Timezone     bool                         `yaml:"timezone"`
}

func LoadConfig(file string) (*Config, error) {
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

	return &c, err
}

func (c *Config) Validate() error {
	if c.Mysql == nil {
		return fmt.Errorf("config: mysql section of config not present")
	}

	if c.Destination == nil {
		return fmt.Errorf("config: destination section of config not present or complete, %v", c)
	}

	if c.Destination.File == "" && c.Destination.Postgres == nil {
		return fmt.Errorf("config: either file or postgres has to be specified in the destination field", c)
	}

	return nil
}
