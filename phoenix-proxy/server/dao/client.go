package dao

import (
	log "github.com/sirupsen/logrus"
	"time"
)

type Client struct {
	ClientID  uint64 `gorm:"PRIMARY_KEY;NOT NULL;"`
	IPAddr    string
	HostName  string
	Os        string
	Arch      string
	CreatedAt time.Time
	DeletedAt *time.Time
}

func AddClient(id uint64, ip, host, os, arch string) {
	if !enable {
		return
	}
	c := Client{
		ClientID: id,
		IPAddr:   ip,
		HostName: host,
		Os:       os,
		Arch:     arch,
	}
	if err := db.Create(&c).Error; err != nil {
		log.Error("Insert Client Failed:", err.Error())
	}
	return
}

func DelClient(id uint64) {
	if !enable {
		return
	}
	var c Client
	if db.Where(&Client{ClientID: id}).First(&c).RecordNotFound() {
		log.Warn("Delete Client Failed, No Such ID:", id)
		return
	}
	if err := db.Delete(&c).Error; err != nil {
		log.Error("Delete Client Failed:", err.Error())
		return
	}
}
