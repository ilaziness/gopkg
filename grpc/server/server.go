package main

import (
	"context"
	"io"
	"log"
	"net"
	"strconv"

	pb "github.com/ilaziness/gopkg/grpc/proto"
	"google.golang.org/grpc"
)

// 实现 .proto定义的rpc方法
type helloServer struct {
	pb.UnimplementedHelloServer
}

func (h *helloServer) Ping(ctx context.Context, req *pb.Req) (*pb.Resp, error) {
	log.Println("sample rpc received:", req.Msg, req.GetMsg())
	return &pb.Resp{Msg: "hello, is server"}, nil
}

// 服务端流
func (h *helloServer) Sss(req *pb.Req, stream pb.Hello_SssServer) error {
	log.Println("sss received:", req.Msg)
	for i := 1; i <= 5; i++ {
		if err := stream.Send(&pb.Resp{Msg: strconv.Itoa(i)}); err != nil {
			return err
		}
	}
	return nil
}

// 客户端流
func (h *helloServer) Css(stream pb.Hello_CssServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.Resp{Msg: "client stream close"})
		}
		if err != nil {
			return err
		}
		log.Println("client stream received:", req.Msg)
	}
}

// 双向流
func (h *helloServer) Bs(stream pb.Hello_BsServer) error {
	wait := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				log.Println("bs receive over")
				close(wait)
				return
			}
			if err != nil {
				log.Println("bs receive error:", err)
				return
			}
			log.Println("bs receive:", in.Msg)
		}
	}()
	for i := 1; i <= 5; i++ {
		if err := stream.Send(&pb.Resp{Msg: strconv.Itoa(i)}); err != nil {
			return err
		}
	}
	<-wait
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":7000")
	if err != nil {
		log.Fatalln("fiailed to listen:", err)
	}
	s := grpc.NewServer()
	pb.RegisterHelloServer(s, &helloServer{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
