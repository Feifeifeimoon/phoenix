package proxy

import (
	"context"
	log "github.com/sirupsen/logrus"
	"net"
	"phoenix-proxy/common/msg"
)

type Proxy interface {
	Run()
}

type BaseProxy struct {
	ctx context.Context

	controlID uint64

	proxyName string

	proxyType byte
	//服务器的地址
	remoteAddr *net.TCPAddr
	remoteConn net.Conn
	//内网服务的地址
	localIP   string
	localPort string
}

func NewProxy(ctx context.Context, controlID uint64, cfg *ProxyConf, remoteAddr *net.TCPAddr) Proxy {
	base := BaseProxy{
		ctx:        ctx,
		controlID:  controlID,
		proxyName:  cfg.ProxyName,
		proxyType:  cfg.ProxyType,
		remoteAddr: remoteAddr,
		remoteConn: nil,
		localIP:    cfg.LocalIP,
		localPort:  cfg.LocalPort,
	}
	switch cfg.ProxyType {
	case msg.ProxyTypeTCP:
		return &TCPProxy{
			BaseProxy: base,
			localConn: nil,
		}
	}
	return nil
}

func (pxy *BaseProxy) loginProxy() (err error) {
	r, err := net.DialTCP("tcp", nil, pxy.remoteAddr)
	if err != nil {
		log.Debugf("[%s]:Conn To Server Failed: %s", pxy.proxyName, err.Error())
		return
	}
	pxy.remoteConn = r

	l := msg.Login{
		ControlID: pxy.controlID,
		ProxyName: pxy.proxyName,
	}
	if err = msg.Write(r, l); err != nil {
		log.Debugf("[%s]:Send Login Request Failed: %s", pxy.proxyName, err.Error())
		return
	}

	var resp msg.LoginResp
	if err = msg.ReadInto(r, &resp); err != nil {
		log.Debugf("[%s]:Read Login Response Failed: %s", pxy.proxyName, err.Error())
		return
	}

	return
}
