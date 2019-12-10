package config

type Config struct {
	Addr 			string		`yaml:"addr"`
	Dialect 		string		`yaml:"dialect"`
	DSN				string		`yaml:"dsn"`
	MaxIdleConn		int			`yaml:"max_idle_conn"`
}

var config *Config

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