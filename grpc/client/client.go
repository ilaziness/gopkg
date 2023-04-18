package main

import (
	"context"
	"io"
	"log"
	"strconv"
	"time"

	pb "github.com/ilaziness/gopkg/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	addr := "127.0.0.1:7000"
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewHelloClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 调用rpc服务方法
	r, err := client.Ping(ctx, &pb.Req{Msg: "is client"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMsg())

	// 服务端流
	req := &pb.Req{Msg: "client call sss"}
	stream, err := client.Sss(ctx, req)
	if err != nil {
		log.Fatalln(err)
	}
	for {
		recv, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln("receive error", err)
		}
		log.Println("sss receve:", recv.Msg)
	}

	//客户端流
	cstream, err := client.Css(ctx)
	if err != nil {
		log.Fatalln("css error:", err)
	}
	for i := 1; i <= 5; i++ {
		if err := cstream.Send(&pb.Req{Msg: strconv.Itoa(i)}); err != nil {
			log.Fatalln("css send error:", err)
		}
	}
	reply, err := cstream.CloseAndRecv()
	if err != nil {
		log.Println("css close error:", err)
	}
	log.Println("css reply:", reply.Msg)

	//双向流
	bstream, err := client.Bs(ctx)
	waitc := make(chan struct{})
	go func() {
		for {
			in, err := bstream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalln("bs failed to receive:", err)
			}
			log.Println("bs received:", in.Msg)
		}
	}()
	for i := 1; i <= 5; i++ {
		if err := bstream.Send(&pb.Req{Msg: strconv.Itoa(i)}); err != nil {
			log.Fatalln("bs send error:", err)
		}
	}
	bstream.CloseSend()
	<-waitc
}
