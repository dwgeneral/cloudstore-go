#!/bin/bash

ROOT_DIR=$GOPATH/cloudstore-go
services="
dbproxy
upload
download
transfer
account
apigw
"

# 编译service可执行文件
build_service() {
    rm -f ${ROOT_DIR}/deploy/bin/$1
    go build -o ${ROOT_DIR}/deploy/bin/$1 ${ROOT_DIR}/service/$1/main.go
    echo -e "\033[32m编译完成: \033[0m ${ROOT_DIR}/deploy/bin/$1"
}

# 打包镜像
build_image() {
    sudo docker build -t cloudstore-go/$1 -f ./service/$1/Dockerfile .
    echo -e "\033[32m镜像打包完成: \033[0m cloudstore-go/$1\n"
}

# 切换到工程根目录
cd ${ROOT_DIR}

# 打包静态资源
mkdir ${ROOT_DIR}/assets -p && go-bindata-assetfs -pkg assets -o ${ROOT_DIR}/assets/asset.go static/...

# 执行编译service
mkdir -p ${ROOT_DIR}/deploy/bin && rm -f ${ROOT_DIR}/deploy/bin/*
for service in $services
do
    build_service $service
done

echo -e "\033[32m编译完毕, 开始构建docker镜像... \033[0m"

# 打包微服务镜像
cd ${ROOT_DIR}/deploy/
for service in $services
do
    build_image $service
done

echo -e "\033[32mdocker镜像构建完毕.\033[0m"

# 容器启动示例
# 启动account service
# docker run -it -e PARAMS="--registry=consul --registry_address=192.168.123.138:8500" cloudstore-go/account