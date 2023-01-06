package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ilaziness/gopkg/serviceregdisc"
	"github.com/ilaziness/gopkg/serviceregdisc/client"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	zkHost := []string{"127.0.0.1:2181"}
	etcdHost := []string{"127.0.0.1:2379"}
	user := "test"
	pass := "test"
	ctx, cancel := context.WithCancel(context.Background())

	// zk client
	clt, err := client.NewZKClient(zkHost, user, pass)
	if err != nil {
		log.Fatalln(err)
	}
	// etcd client
	clietcd, err := client.NewEtcdClient(ctx, etcdHost, user, pass)
	if err != nil {
		log.Fatalln(err)
	}
	regdisc := serviceregdisc.NewRegDisc("crm", clt)
	regdisc = serviceregdisc.NewRegDisc("crm", clietcd)

	// discover
	servers := make(map[string]*serviceregdisc.Server)
	// 要发现的服务
	service := []string{"servicea", "serviceb", "servicec"}
	for _, id := range service {
		ser, err := serviceregdisc.NewServerDiscover(ctx, regdisc.GetServicePath(id), *regdisc)
		if err != nil {
			log.Fatalln(err)
		}
		servers[id] = ser
	}

	// 获取用户服务地址
	// log.Println(servers["servicea"].GetServer())
	// log.Println(servers["serviceb"].GetServer())
	// log.Println(servers["servicec"].GetServer())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	cancel()
	time.Sleep(time.Second)
	log.Printf("Server quiting....\n")
}
