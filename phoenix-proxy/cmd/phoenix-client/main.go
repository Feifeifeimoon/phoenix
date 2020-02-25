package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"net"
	"phoenix-proxy/client"
	"phoenix-proxy/client/proxy"
)

var (
	configName string
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:      false,
		DisableColors:    false,
		DisableTimestamp: false,
		FullTimestamp:    true,
		TimestampFormat:  "2006-01-02 15:04:05",
	})
	log.SetLevel(log.DebugLevel)
	flag.StringVar(&configName, "c", "./conf.ini", "config file")
	flag.Parse()
}

func main() {
	cfg, err := ini.Load(configName)
	if err != nil {
		log.Println("Load config file Failed:", err.Error())
		return
	}
	sec, err := cfg.GetSection("remote")
	if err != nil {
		log.Println("Config Error: Remote configuration not found")
		return
	}
	IPString := sec.Key("ip").String() + ":" + sec.Key("port").String()
	proxyConf := generatePxyConf(cfg)
	TCPAddr, err := net.ResolveTCPAddr("tcp", IPString)
	if err != nil {
		log.Fatalln("ResolveIPAddr Failed:", err.Error())
	}
	c := client.NewClient(TCPAddr, proxyConf)
	c.Run()
}

func generatePxyConf(cfg *ini.File) []proxy.ProxyConf {
	var proxyConfigs []proxy.ProxyConf
	var temp proxy.ProxyConf
	for _, val := range cfg.SectionStrings() {
		if val != ini.DefaultSection && val != "remote" {
			sec := cfg.Section(val)
			temp.ProxyName = val
			temp.ProxyType = proxy.Str2Type[sec.Key("type").String()]
			if err := sec.MapTo(&temp); err != nil {
				fmt.Println("MapTo Error:", err.Error())
			}
			proxyConfigs = append(proxyConfigs, temp)
			log.Info(temp)
		}
	}
	return proxyConfigs
}
