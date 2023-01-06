package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ilaziness/gopkg/serviceregdisc"
	"github.com/ilaziness/gopkg/serviceregdisc/client"
)

func main() {
	newReg("servicec")
}

func newReg(id string) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	port := ""
	fmt.Println("input port:")
	fmt.Scan(&port)

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

	// register
	serverInfo := serviceregdisc.ServerInfo{
		IP:     "127.0.0.1",
		Port:   port,
		Schema: "https",
		UUID:   "23424",
	}
	err = regdisc.Register(ctx, id, serverInfo)
	if err != nil {
		log.Println("register service error ", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	cancel()
	time.Sleep(time.Second)
	log.Printf("Server Shutdown....\n")
}
