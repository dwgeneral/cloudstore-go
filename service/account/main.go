package main

import (
	"log"
	"time"

	"github.com/micro/go-micro"

	"cloudstore-go/common"
	"cloudstore-go/service/account/handler"
	proto "cloudstore-go/service/account/proto"
	dbproxy "cloudstore-go/service/dbproxy/client"
)

func main() {
	service := micro.NewService(
		micro.Name("go.micro.service.user"),
		micro.RegisterTTL(time.Second*10),
		micro.RegisterInterval(time.Second*5),
		micro.Flags(common.CustomFlags...),
	)

	// 初始化service, 解析命令行参数等
	service.Init()

	// 初始化dbproxy client
	dbproxy.Init(service)

	proto.RegisterUserServiceHandler(service.Server(), new(handler.User))
	if err := service.Run(); err != nil {
		log.Println(err)
	}
}
