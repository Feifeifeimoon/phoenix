FROM ubuntu:19.04

# 切换源，安装软件
RUN sed -i 's/archive.ubuntu.com/mirrors.aliyun.com/g' /etc/apt/sources.list && \
    sed -i 's/security.ubuntu.com/mirrors.aliyun.com/g' /etc/apt/sources.list && \
    apt update && apt install -y redis nginx supervisor

# sed -i 's/PermitRootLogin without-password/PermitRootLogin yes/' /etc/ssh/sshd_config && \
# echo "root:996icu" | chpasswd

# Golang SDK
#ADD https://studygolang.com/dl/golang/go1.13.4.linux-amd64.tar.gz /tmp
#RUN tar -zxvf /tmp/go1.13.4.linux-amd64.tar.gz -C /usr/local/
#
#ENV GOROOT /usr/local/go
#ENV PATH $PATH:$GOROOT/bin
#ENV GOPROXY https://mirrors.aliyun.com/goproxy/

COPY build /phoenix
#COPY phoenix-proxy /phoenix/phoenix-proxy
#COPY phoenix-web /phoenix/phoenix-web

RUN mkdir -p /phoenix/var/log/supervisor/ /phoenix/var/log/nginx/

ENTRYPOINT [ "sh", "/phoenix/bin/init.sh"]
