package dao

import (
	"app/conf"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
	"time"
)

var DB *gorm.DB

//初始化连接数据库
func init() {
	log.Info("初始化数据库连接")
	connString := fmt.Sprintf("%s:%s@tcp(%s)/app?charset=utf8&parseTime=True&loc=Local",
		conf.Conf.DB.User, conf.Conf.DB.PassWd, conf.Conf.DB.IPAddr)

	db, err := gorm.Open("mysql", connString)
	if err != nil {
		log.Error("mysql连接失败：", err.Error())
		return
	}
	log.Info("mysql连接成功")

	db.LogMode(true)
	//设置连接池
	db.DB().SetMaxIdleConns(20)
	//打开
	db.DB().SetMaxOpenConns(100)
	//超时时间
	db.DB().SetConnMaxLifetime(time.Second * 30)

}
