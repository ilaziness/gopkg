package rpcall

import (
	"context"
	"log"

	"github.com/ilaziness/gopkg/opentelemetry/advanced/api"
)

// CallAsTest 调用as服务的test
func CallAsTest(ctx context.Context, c api.AsRpcClient) {
	resp, err := c.Test(ctx, &api.Req{A: "a", B: 34})
	if err != nil {
		log.Fatal("call as test error:", err)
	}
	log.Println("call as response:", resp.S)
}

// CallBsTest 调用bs服务的test
func CallBsTest(ctx context.Context, c api.BsRpcClient) {
	resp, err := c.Test(ctx, &api.Req{A: "b", B: 35})
	if err != nil {
		log.Fatal("call bs test error:", err)
	}
	log.Println("call bs response:", resp.S)
}
