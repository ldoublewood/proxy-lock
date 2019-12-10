package config

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Addr 			string		`yaml:"addr"`
	Dialect 		string		`yaml:"dialect"`
	DSN				string		`yaml:"dsn"`
	MaxIdleConn		int			`yaml:"max_idle_conn"`
}

var config *Config

// not used!
func Load(path string) error {
	result, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(result, &config)
}

func Init(dialect string, source string, addr string) {
	config = &Config{
		Addr: addr,
		Dialect:     dialect,
		DSN:         source,
		MaxIdleConn: 100,
	}
}
func Get() *Config {
	return config
}