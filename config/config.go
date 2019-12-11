package config

import "gopkg.in/urfave/cli.v1"

type Config struct {
	Addr 			string
	Dialect 		string
	DSN				string
	MaxIdleConn		int
	TestDestUrl		string
}

var config *Config

func Init(c *cli.Context) {
	config = &Config{
		Addr: c.String("addr"),
		Dialect:     c.String("dialect"),
		DSN:         c.String("source"),
		TestDestUrl: c.String("testdesturl"),
		MaxIdleConn: 100,
	}
}
func Get() *Config {
	return config
}