package msg

type Message interface{}

type Login struct {
	//主机名
	HostName string
	//操作系统类型
	Os string
	//系统架构
	Arch string
	//ControlID 第一次连接时为零
	ControlID uint64
	//proxyName
	ProxyName string
}

type LoginResp struct {
	//分配的control ID
	ControlID uint64
	//错误信息
	Err string
}

type RegisterProxy struct {
	//Proxy名称
	ProxyName string
	//Proxy 类型 TCP UDP...
	ProxyType byte
	//对应的端口
	RemotePort string
}

type RegisterProxyResp struct {
	Err string
}

//server -> client
type NewProxyConn struct {
	ProxyName string
}

type HeartBeat struct{}
