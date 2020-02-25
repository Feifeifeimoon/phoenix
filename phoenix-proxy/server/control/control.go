package control

import (
	"context"
	log "github.com/sirupsen/logrus"
	"net"
	"phoenix-proxy/common/msg"
	"phoenix-proxy/server/proxy"
	"time"
)

type control struct {
	//控制control下的proxy退出
	ctx    context.Context
	cancel context.CancelFunc

	controlID uint64
	//对应的client
	conn net.Conn

	readCh chan msg.Message

	writeCh chan msg.Message

	closeCh chan struct{}

	pxyMap map[string]proxy.Proxy
}

func newControl(controlID uint64, conn net.Conn) *control {
	ctx, cancel := context.WithCancel(context.Background())
	return &control{
		ctx:       ctx,
		cancel:    cancel,
		controlID: controlID,
		conn:      conn,
		readCh:    make(chan msg.Message),
		writeCh:   make(chan msg.Message),
		closeCh:   make(chan struct{}),
		pxyMap:    make(map[string]proxy.Proxy),
	}
}

func (c *control) run() {
	defer c.conn.Close()
	defer close(c.writeCh)
	defer close(c.closeCh)
	defer close(c.readCh)

	go c.readRoutine()
	go c.writeRoutine()

	_ = c.conn.SetDeadline(time.Now().Add(time.Second * 10))
	for {
		select {
		case <-c.closeCh:
			log.Debugf("[%x]:Work closeCh Active", c.controlID)
			c.cancel()
			return
		case t, ok := <-c.readCh:
			if !ok {
				return
			}
			switch m := t.(type) {
			case *msg.RegisterProxy:
				c.handleRegister(m)
			case *msg.HeartBeat:
				log.Debugf("[%x]:Recv HeartBeat From Client", c.controlID)
				_ = c.conn.SetDeadline(time.Now().Add(time.Second * 10))
			default:
				log.Error("Unknown Request Type")
				return
			}
		}
	}
}
func (c *control) readRoutine() {
	defer c.close()
	for {
		m, err := msg.Read(c.conn)
		if err != nil {
			log.Debugf("[%x]:Recv From Client Failed:%s", c.controlID, err.Error())
			return
		}
		c.readCh <- m
	}

}

func (c *control) writeRoutine() {
	for {
		select {
		case <-c.closeCh:
			log.Debugf("[%x]:writeRoutine closeCh Active", c.controlID)
			return
		case m, ok := <-c.writeCh:
			if !ok {
				log.Warn("")
				return
			} else {
				if err := msg.Write(c.conn, m); err != nil {
					log.Error("Write Message Failed:", err.Error())
					return
				}
			}
		}
	}

}

func (c *control) handleRegister(m *msg.RegisterProxy) {
	log.Infof("[%x]New Register Request: ProxyName[%s], ProxyType[%d], RemotePort[%s]",
		c.controlID, m.ProxyName, m.ProxyType, m.RemotePort)
	pxy := proxy.NewProxy(c.ctx, m.ProxyName, m.ProxyType, m.RemotePort, c.writeCh)
	c.pxyMap[m.ProxyName] = pxy
	go pxy.Run()
	r := msg.RegisterProxyResp{Err: ""}
	c.writeCh <- r
	return
}

func (c *control) close() {
	// 向close chan发生两次，一次被write协程，一次被work协程
	c.closeCh <- struct{}{}
	c.closeCh <- struct{}{}
}

func (c *control) NewProxyConn(proxyName string, conn net.Conn) {
	pxy, ok := c.pxyMap[proxyName]
	if !ok {
		log.Error("No Such Proxy ", proxyName)
		return
	}
	pxy.NewProxyConn() <- conn
}
