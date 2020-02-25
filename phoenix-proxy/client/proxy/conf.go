package proxy

import "phoenix-proxy/common/msg"

var Str2Type = map[string]byte{
	"tcp":  msg.ProxyTypeTCP,
	"TCP":  msg.ProxyTypeTCP,
	"udp":  msg.ProxyTypeUDP,
	"UDP":  msg.ProxyTypeUDP,
	"Http": msg.ProxyTypeHTTP,
	"HTTP": msg.ProxyTypeHTTP,
}

type ProxyConf struct {
	ProxyName string

	ProxyType byte

	LocalIP string `ini:"local_ip"`

	LocalPort string `ini:"local_port"`

	RemotePort string `ini:"remote_port"`
}
