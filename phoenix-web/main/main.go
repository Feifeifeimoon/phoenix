package main

import (
	"app/conf"
	"app/dao"
	"app/router"
	log "github.com/sirupsen/logrus"
	"os"
)

var _ = dao.DB

//初始化log
func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	//gin.SetMode(gin.ReleaseMode)
	//gin.DefaultWriter = log.Out
}

func main() {
	log.Info("==========start==========")
	r := router.NewRouter()
	if err := r.Run(conf.Conf.Server.IPAddr); err != nil {
		log.Info("start error")
	}
}
