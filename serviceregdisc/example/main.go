package main

import (
	"context"
	"log"

	"github.com/ilaziness/gopkg/serviceregdisc"
	"github.com/ilaziness/gopkg/serviceregdisc/client"
)

func main() {
	zkHost := []string{"127.0.0.1:2181"}
	user := ""
	pass := ""
	clt, err := client.NewZKClient(zkHost, user, pass)
	if err != nil {
		log.Fatalln(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	regdisc := serviceregdisc.NewRegDisc("crm", clt)

	// discover
	servers := make(map[string]*serviceregdisc.Server)
	service := []string{"user", "product"}
	for _, id := range service {
		ser, err := serviceregdisc.NewServerDiscover(ctx, regdisc.GetServicePath(id), *regdisc)
		if err != nil {
			log.Fatalln(err)
		}
		servers[id] = ser
	}

	// register
	serverInfo := serviceregdisc.ServerInfo{
		IP:     "127.0.0.1",
		Port:   8080,
		Schema: "https",
		UUID:   "23424",
	}
	err = regdisc.Register(ctx, "user", serverInfo)
	if err != nil {
		log.Println("register service error ", err)
	}

	// 多注册几个服务
	go newReg("abc")
	go newReg("efg")
	go newReg("user")

	// 获取用户服务地址
	log.Println(servers["user"].GetServer())

	cancel()
}

func newReg(id string) {
	zkHost := []string{"127.0.0.1:2181"}
	user := ""
	pass := ""
	clt, err := client.NewZKClient(zkHost, user, pass)
	if err != nil {
		log.Fatalln(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	regdisc := serviceregdisc.NewRegDisc("crm", clt)

	// register
	serverInfo := serviceregdisc.ServerInfo{
		IP:     "127.0.0.1",
		Port:   8080,
		Schema: "https",
		UUID:   "23424",
	}
	err = regdisc.Register(ctx, id, serverInfo)
	if err != nil {
		log.Println("register service error ", err)
	}
}
