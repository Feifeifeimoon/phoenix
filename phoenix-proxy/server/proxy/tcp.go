package proxy

import (
	log "github.com/sirupsen/logrus"
	"net"
	"phoenix-proxy/common/forward"
	"phoenix-proxy/common/msg"
)

type TCPProxy struct {
	*BaseProxy

	listener *net.TCPListener
	//外网的连接
	acceptCh chan net.Conn
	//内网
	localCh chan net.Conn
}

func (pxy *TCPProxy) Run() {
	defer close(pxy.acceptCh)
	defer close(pxy.localCh)

	addr, err := net.ResolveTCPAddr("tcp", ":"+pxy.remotePort)
	if err != nil {
		log.Errorf("ResolveTCPAddr Failed:", err.Error())
		return
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Errorf("ListenTCP Failed:", err.Error())
		return
	}
	pxy.listener = l
	go pxy.acceptRoutine()
	log.Debugf("[%s]:Listening on %s", pxy.proxyName, addr.String())
	for {
		select {
		case <-pxy.ctx.Done():
			log.Debugf("[%s]:ctx Done", pxy.proxyName)
			_ = pxy.listener.Close()
			return
		case src, ok := <-pxy.acceptCh:
			if !ok {
				log.Error("acceptCh Failed")
				return
			}
			log.Debugf("[%s]:Recv New Conn", pxy.proxyName)
			pxy.clientCh <- msg.NewProxyConn{ProxyName: pxy.proxyName}
			dst := <-pxy.localCh
			log.Debugf("[%s]:Start Forward", pxy.proxyName)
			go forward.Forward(src, dst)
		}
	}
}

func (pxy *TCPProxy) acceptRoutine() {
	for {
		conn, err := pxy.listener.Accept()
		if err != nil {
			log.Errorf("[%s]:Accept Failed:%s", pxy.proxyName, err.Error())
			return
		}
		pxy.acceptCh <- conn
	}
}

func (pxy *TCPProxy) NewProxyConn() chan<- net.Conn {
	return pxy.localCh
}
