#!/bin/sh

cd /root/phoenix/phoenix-proxy/cmd/phoenix-server
go build -o /root/phoenix/build/bin/phoenix-proxy main.go


cd /root/phoenix/phoenix-web/main
go build -o /root/phoenix/build/bin/phoenix-web main.go

docker exec phoenix supervisorctl restart phoenix-proxy
docker exec phoenix supervisorctl restart phoenix-web
