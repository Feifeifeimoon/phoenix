package dao

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
	"time"
)

var enable bool
var db *gorm.DB

func Init() {
	connString := "phoenix:phoenix996icu@tcp(172.18.0.2)/app?charset=utf8&parseTime=True&loc=Local"
	t, err := gorm.Open("mysql", connString)
	if err != nil {
		log.Println("mysql连接失败：", err.Error())
		return
	}
	enable = true
	db = t

	db.LogMode(true)
	//设置连接池
	//空闲
	db.DB().SetMaxIdleConns(20)
	//打开
	db.DB().SetMaxOpenConns(100)
	//超时时间
	db.DB().SetConnMaxLifetime(time.Second * 30)

	db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8").AutoMigrate(&Client{})
}
