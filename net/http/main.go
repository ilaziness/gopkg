package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	testHttpRespBodyClose()
}

// http body应该关闭，否则连接不会放回空闲池，不读取body不能复用连接
func testHttpRespBodyClose() {
	tr := &http.Transport{
		MaxIdleConns:    1,
		MaxConnsPerHost: 1,
	}
	client := &http.Client{
		Transport: tr,
	}

	for i := 0; i < 5; i++ {
		resp, err := client.Get("http://www.baidu.com")
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("resp: %v", i)
		//resp.Body.Close()
		io.ReadAll(resp.Body)
	}
}
