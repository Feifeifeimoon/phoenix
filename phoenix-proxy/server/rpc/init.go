package rpc

import (
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"phoenix-proxy/server/control"
	"phoenix-proxy/server/rpc/manage"
)

func Init(cm *control.ControlManage) {
	s := grpc.NewServer()
	manage.RegisterManageServiceServer(s, manage.NewManageService(cm))
	l, err := net.Listen("tcp", ":44568")
	if err != nil {
		log.Info("RPC Start Failed:", err.Error())
	}
	go func() {
		if err := s.Serve(l); err != nil {
			log.Errorf("RPC Server Failed: %v", err)
		}
	}()
}
