package conf

import (
	"flag"
	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
)

var (
	ConfPath string
	Conf     = &Config{}
)

type Config struct {
	DB     *DB     `toml:"database"`
	Server *Server `toml:"servers"`
}

type DB struct {
	IPAddr string
	User   string
	PassWd string
}

type Server struct {
	IPAddr string
}

func init() {
	flag.StringVar(&ConfPath, "c", "./conf.toml", "default config path")
	flag.Parse()
	if _, err := toml.DecodeFile(ConfPath, &Conf); err != nil {
		log.Error("Decode Config File Error", err.Error())
	}
}
