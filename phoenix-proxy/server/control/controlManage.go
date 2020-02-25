package control

import (
	log "github.com/sirupsen/logrus"
	"github.com/sony/sonyflake"
	"net"
	"phoenix-proxy/server/dao"
)

type ControlManage struct {
	//Generate Control id
	flake *sonyflake.Sonyflake
	//control map
	controlMap map[uint64]*control
	//最大client数量
	maxClientNum int64
}

func NewControlManage() *ControlManage {
	return &ControlManage{
		flake:        sonyflake.NewSonyflake(sonyflake.Settings{}),
		controlMap:   make(map[uint64]*control),
		maxClientNum: 9999, //命令行启动的话默认999
	}
}

//获取最大客户端数量
func (c *ControlManage) GetMaxClientNum() int64 {
	return c.maxClientNum
}

//设置最大client数量
func (c *ControlManage) SetMaxClientNum(num int64) {
	c.maxClientNum = num
}

//获取当前已经连接的client数量
func (c *ControlManage) GetCurClientNum() int64 {
	return int64(len(c.controlMap))
}

//  添加client连接
func (c *ControlManage) AddControl(conn net.Conn) (id uint64, err error) {
	id, err = c.flake.NextID()
	if err != nil {
		return
	}
	c.controlMap[id] = newControl(id, conn)
	go c.runControl(id)
	return
}

//启动control，并等待退出，用来更新map
func (c *ControlManage) runControl(id uint64) {
	control := c.controlMap[id]
	control.run() //返回即代表client关闭，清理map
	dao.DelClient(id)
	delete(c.controlMap, id)
}

//  client发起的对应proxy的连接
func (c *ControlManage) NewProxyConn(id uint64, proxyName string, conn net.Conn) {
	control, ok := c.controlMap[id]
	if !ok {
		log.Errorf("No Such ID Control ", id)
		return
	}
	control.NewProxyConn(proxyName, conn)
}
