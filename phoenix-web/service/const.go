package service

const (
	server      = "phoenix-proxy"
	bash        = "/bin/bash"
	superCtl    = "supervisorctl"
	proxyStatus = superCtl + " status | grep phoenix-proxy | awk '{print $2}' "
)
