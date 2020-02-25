package proxy

import (
	"context"
	"net"
	"phoenix-proxy/common/msg"
)

type Proxy interface {
	Run()

	NewProxyConn() chan<- net.Conn
}

type BaseProxy struct {
	ctx context.Context

	proxyName string

	remotePort string
	//通知client有连接到达
	clientCh chan msg.Message
}

func NewProxy(ctx context.Context, proxyName string, proxyType byte, remotePort string, ch chan msg.Message) Proxy {
	base := BaseProxy{
		ctx:        ctx,
		proxyName:  proxyName,
		remotePort: remotePort,
		clientCh:   ch,
	}
	switch proxyType {
	case msg.ProxyTypeTCP:
		return &TCPProxy{
			BaseProxy: &base,
			listener:  nil,
			acceptCh:  make(chan net.Conn),
			localCh:   make(chan net.Conn),
		}
	}
	return nil
}
