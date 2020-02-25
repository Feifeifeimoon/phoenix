package manage

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/wrappers"
	"phoenix-proxy/server/control"
)

type ManageServiceImpl struct {
	cm *control.ControlManage
}

func NewManageService(cm *control.ControlManage) *ManageServiceImpl {
	return &ManageServiceImpl{cm: cm}
}

func (s *ManageServiceImpl) MaxClientNum(context.Context, *empty.Empty) (*wrappers.Int64Value, error) {
	return &wrappers.Int64Value{Value: s.cm.GetMaxClientNum()}, nil
}

func (s *ManageServiceImpl) SetMaxClientNum(_ context.Context, args *wrappers.Int64Value) (*empty.Empty, error) {
	s.cm.SetMaxClientNum(args.Value)
	return nil, nil
}

func (s *ManageServiceImpl) CurClientNum(context.Context, *empty.Empty) (*wrappers.Int64Value, error) {
	return &wrappers.Int64Value{Value: s.cm.GetCurClientNum()}, nil
}
