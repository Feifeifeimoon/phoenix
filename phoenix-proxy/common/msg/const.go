package msg

import "reflect"

//Msg Type
const (
	TypeLogin = 1 << iota
	TypeLoginResp
	TypeRegisterProxy
	TypeRegisterProxyResp
	TypeNewProxyConn
	TypeHeartBeat
)

var byte2type = map[byte]reflect.Type{
	TypeLogin:             reflect.TypeOf(Login{}),
	TypeLoginResp:         reflect.TypeOf(LoginResp{}),
	TypeRegisterProxy:     reflect.TypeOf(RegisterProxy{}),
	TypeRegisterProxyResp: reflect.TypeOf(RegisterProxyResp{}),
	TypeNewProxyConn:      reflect.TypeOf(NewProxyConn{}),
	TypeHeartBeat:         reflect.TypeOf(HeartBeat{}),
}

var type2byte = map[reflect.Type]byte{
	reflect.TypeOf(Login{}):             TypeLogin,
	reflect.TypeOf(LoginResp{}):         TypeLoginResp,
	reflect.TypeOf(RegisterProxy{}):     TypeRegisterProxy,
	reflect.TypeOf(RegisterProxyResp{}): TypeRegisterProxyResp,
	reflect.TypeOf(NewProxyConn{}):      TypeNewProxyConn,
	reflect.TypeOf(HeartBeat{}):         TypeHeartBeat,
}

//Proxy Type
const (
	ProxyTypeTCP = 1 << iota
	ProxyTypeUDP
	ProxyTypeHTTP
)
