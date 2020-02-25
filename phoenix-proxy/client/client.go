package client

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	"phoenix-proxy/client/proxy"
	"phoenix-proxy/common/msg"
	"runtime"
	"time"
)

type Client struct {
	//服务器地址
	remoteAddr *net.TCPAddr
	//Proxy配置
	cfg []proxy.ProxyConf
	//与服务器连接的句柄
	conn net.Conn
	//对应服务器的control ID
	controlID uint64

	readCh chan msg.Message

	writeCh chan msg.Message

	closeCh chan struct{}
}

func NewClient(remoteAddr *net.TCPAddr, cfg []proxy.ProxyConf) *Client {
	return &Client{
		remoteAddr: remoteAddr,
		cfg:        cfg,
		conn:       nil,
		controlID:  0,
		readCh:     make(chan msg.Message),
		writeCh:    make(chan msg.Message),
		closeCh:    make(chan struct{}),
	}
}

func (c *Client) Run() {
	for {
		if err := c.login(); err != nil {
			log.Error("Login to Server Failed:", err.Error())
			time.Sleep(time.Second * 5)
			continue
		}
		log.Info("Success Login to Server")
		if err := c.register(); err != nil {
			log.Error("Register Proxy Failed:", err.Error())
			return
		}
		log.Info("Success Register Proxy Server")
		go c.readRoutine()
		go c.writeRoutine()
		c.work()
		log.Warn("Work Function Return, ReConnect to Server")
	}
}

func (c *Client) login() (err error) {
	//conn to server
	conn, err := net.DialTCP("tcp", nil, c.remoteAddr)
	if err != nil {
		log.Debug("Conn To Remote Failed:", err.Error())
		return err
	}
	c.conn = conn

	hostname, _ := os.Hostname()
	l := msg.Login{
		HostName:  hostname,
		Os:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		ControlID: 0,
		ProxyName: "",
	}

	if err = msg.Write(c.conn, l); err != nil {
		log.Error("Send login Request failed:", err.Error())
		return
	}

	var resp msg.LoginResp
	if err = msg.ReadInto(c.conn, &resp); err != nil {
		return
	}
	if len(resp.Err) != 0 {
		err = fmt.Errorf("%s", resp.Err)
		return
	}

	c.controlID = resp.ControlID
	return
}

func (c *Client) register() (err error) {
	var resp msg.RegisterProxyResp
	for _, val := range c.cfg {
		m := msg.RegisterProxy{
			ProxyName:  val.ProxyName,
			ProxyType:  val.ProxyType,
			RemotePort: val.RemotePort,
		}
		if err = msg.Write(c.conn, m); err != nil {
			log.Error("Send Register Request Failed:", err.Error())
			return
		}
		log.Info("Send Register Request Done")
		if err = msg.ReadInto(c.conn, &resp); err != nil {
			log.Error("Register Response Failed:", err.Error())
			return
		}
	}
	return
}

func (c *Client) readRoutine() {
	defer c.close()
	for {
		m, err := msg.Read(c.conn)
		if err != nil {
			log.Debug("Recv From Server Failed:", err.Error())
			return
		}
		c.readCh <- m
	}

}

func (c *Client) writeRoutine() {
	for {
		select {
		case <-c.closeCh:
			log.Debug("writeRoutine closeCh Active")
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

func (c *Client) work() {
	ticker := time.NewTicker(time.Second * 9)
	for {
		select {
		case <-c.closeCh:
			log.Warn("work closeCh Active")
			return
		case <-ticker.C:
			c.writeCh <- msg.HeartBeat{}
		case t, ok := <-c.readCh:
			if !ok {
				return
			}
			switch m := t.(type) {
			case *msg.NewProxyConn:
				c.handleNewProxyConn(m)
			}
		}
	}
}

func (c *Client) handleNewProxyConn(m *msg.NewProxyConn) {
	var cfg proxy.ProxyConf
	for _, val := range c.cfg {
		if val.ProxyName == m.ProxyName {
			cfg = val
		}
	}
	pxy := proxy.NewProxy(context.TODO(), c.controlID, &cfg, c.remoteAddr)
	go pxy.Run()
}

func (c *Client) close() {
	c.closeCh <- struct{}{}
	c.closeCh <- struct{}{}
}
