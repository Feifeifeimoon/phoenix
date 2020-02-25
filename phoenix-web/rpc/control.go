package rpc

import (
	"app/rpc/manage"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var RPCManageClient manage.ManageServiceClient

func init() {
	log.Info("初始化RPC连接")
	conn, err := grpc.Dial(":44568", grpc.WithInsecure())
	if err != nil {
		return
	}
	RPCManageClient = manage.NewManageServiceClient(conn)
}
