package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptrace"
)

// httptrace 用来跟踪http请求中的事件

func main() {
	req, _ := http.NewRequest("GET", "http://baidu.com", nil)
	trace := &httptrace.ClientTrace{
		GotConn: func(connInfo httptrace.GotConnInfo) {
			fmt.Printf("Got Conn: %+v\n", connInfo)
		},
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			fmt.Printf("DNS Info: %+v\n", dnsInfo)
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	_, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	//_, err := http.DefaultTransport.RoundTrip(req)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
