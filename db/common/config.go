package common

type Config struct {
	Hostname string `yaml:"hostname,omitempty"`
	Socket   string `yaml:"socket,omitempty"`
	Port     int    `yaml:"port,omitempty"`
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
	Database string `yaml:"database,omitempty"`
	Compress bool   `yaml:"compress,omitempty"`
}
