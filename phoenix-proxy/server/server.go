package server

import (
	log "github.com/sirupsen/logrus"
	"net"
	"phoenix-proxy/common/msg"
	"phoenix-proxy/server/control"
	"phoenix-proxy/server/dao"
	"phoenix-proxy/server/rpc"
	"time"
)

type Server struct {
	TCPAddr *net.TCPAddr

	listener net.Listener

	controlManage *control.ControlManage
}

func NewServer(TCPAddr *net.TCPAddr, web bool) (*Server, error) {
	l, err := net.ListenTCP("tcp", TCPAddr)
	if err != nil {
		return nil, err
	}

	cm := control.NewControlManage()
	if web {
		cm.SetMaxClientNum(5)
		rpc.Init(cm)
		dao.Init()
	}
	return &Server{
		TCPAddr:       TCPAddr,
		listener:      l,
		controlManage: cm,
	}, nil
}

func (s *Server) Run() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Fatalln("Server Accept Error:", err.Error())
		}
		log.Debug("New Client From:", conn.RemoteAddr().String())
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(c net.Conn) {
	//10s内若没有收到客户端登陆信息则关闭客户端
	_ = c.SetReadDeadline(time.Now().Add(time.Second * 10))
	defer c.SetReadDeadline(time.Time{})
	var l msg.Login
	if err := msg.ReadInto(c, &l); err != nil {
		log.Warn("Read Login Request Failed:", err.Error())
		return
	}
	log.Infof("Client Login: hostname[%s] os[%s] arch[%s] controlID[%x] proxyName[%s]",
		l.HostName, l.Os, l.Arch, l.ControlID, l.ProxyName)

	var r msg.LoginResp
	if l.ControlID != 0 && len(l.ProxyName) != 0 {
		s.controlManage.NewProxyConn(l.ControlID, l.ProxyName, c)
	} else {
		id, err := s.controlManage.AddControl(c)
		if err != nil {
			r.Err = err.Error()
		} else {
			r.ControlID = id
			dao.AddClient(id, c.RemoteAddr().String(), l.HostName, l.Os, l.Arch)
		}

	}
	if err := msg.Write(c, r); err != nil {
		log.Error("Send LoginResp Failed:", err.Error())
	}
	return
}
