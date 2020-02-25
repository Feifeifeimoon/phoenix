package proxy

import (
	log "github.com/sirupsen/logrus"
	"net"
	"phoenix-proxy/common/forward"
)

type TCPProxy struct {
	BaseProxy
	//对应内网的连接
	localConn net.Conn
}

func (pxy *TCPProxy) Run() {
	if err := pxy.loginProxy(); err != nil {
		log.Errorf("[%s]:Login To Server Failed: %s", pxy.proxyName, err.Error())
		return
	}

	tcpAddr, _ := net.ResolveTCPAddr("tcp", pxy.localIP+":"+pxy.localPort)
	c, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Errorf("[%s]:Conn To %s Failed", pxy.proxyName, tcpAddr.String())
		return
	}
	pxy.localConn = c
	forward.Forward(pxy.remoteConn, pxy.localConn)
	log.Infof("[%s]:Exit", pxy.proxyName)
}
