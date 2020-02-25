package service

import (
	"app/rpc"
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"strings"
)

func ServerStatus() (status bool, err error) {
	cmd := exec.Command(bash, "-c", proxyStatus)
	ret, err := cmd.Output()
	if err != nil {
		log.Errorf("Exec CMD %s Failed:%s", proxyStatus, err.Error())
		return
	}
	if strings.Replace(string(ret), "\n", "", -1) == "RUNNING" {
		status = true
		return
	}
	return
}

func ManageServer(command string) (err error) {
	cmd := exec.Command(bash, "-c", fmt.Sprintf("%s %s %s", superCtl, command, server))
	if err = cmd.Run(); err != nil {
		log.Errorf("Manage Server %s Failed:%s", command, err.Error())
		return
	}
	return
}

func ServerMaxClient() int64 {
	r, err := rpc.RPCManageClient.MaxClientNum(context.TODO(), &empty.Empty{})
	if err != nil {
		log.Warn("RPC Failed", err.Error())
		return 0
	}
	return r.Value
}

func ServerCurClient() int64 {
	r, err := rpc.RPCManageClient.CurClientNum(context.TODO(), &empty.Empty{})
	if err != nil {
		log.Warn("RPC Failed", err.Error())
		return 0
	}
	return r.Value
}
