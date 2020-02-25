package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"phoenix-proxy/server"
)

var (
	configName string
	cpuProfile bool
	web        bool
)

func init() {
	//初始化日志输出格式以及读取配置文件
	log.SetFormatter(&log.TextFormatter{
		ForceColors:      false,
		DisableColors:    false,
		DisableTimestamp: false,
		FullTimestamp:    true,
		TimestampFormat:  "2006-01-02 15:04:05",
	})
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
	flag.StringVar(&configName, "c", "./conf.ini", "config file")
	flag.BoolVar(&web, "rpc", false, "enable rpc")
	flag.BoolVar(&cpuProfile, "cpu", false, "enable cpuProfile")
	flag.Parse()
}

func main() {
	//pprof 调试
	if cpuProfile {
		go func() {
			log.Println(http.ListenAndServe(":8080", nil))
		}()
	}

	cfg, err := ini.Load(configName)
	if err != nil {
		log.Fatalln("Load config file Failed:", err.Error())
		return
	}
	sec, err := cfg.GetSection("server")
	if err != nil {
		log.Fatalln("Config Error: Lack server")
		return
	}
	TCPAddr, err := net.ResolveTCPAddr("tcp", ":"+sec.Key("port").String())
	if err != nil {
		log.Fatalln("ResolveIPAddr Failed:", err.Error())
	}
	s, err := server.NewServer(TCPAddr, web)
	if err != nil {
		log.Fatalln("Server Start Failed:", err.Error())
		return
	}
	log.Info("Phoenix-Proxy Start Success, Listen on", TCPAddr.String())
	s.Run()
}
