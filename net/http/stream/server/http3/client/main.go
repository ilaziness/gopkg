package main

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

func main() {
	tr := &http3.Transport{
		// set a TLS client config, if desired
		TLSClientConfig: &tls.Config{
			NextProtos:         []string{http3.NextProtoH3}, // set the ALPN for HTTP/3
			InsecureSkipVerify: true,
		},
		QUICConfig: &quic.Config{}, // QUIC connection options
	}
	defer tr.Close()
	client := &http.Client{
		Transport: tr,
	}

	urls := []string{
		"https://127.0.0.1:443",
		"https://127.0.0.1:443/stream",
		"https://127.0.0.1:4431",
	}

	for _, url := range urls {
		rsp, err := client.Get(url)
		if err != nil {
			log.Printf("req error: %s", err)
			return
		}
		defer rsp.Body.Close()

		rb, err := io.ReadAll(rsp.Body)
		if err != nil {
			log.Println("read body err: ", err)
			return
		}
		log.Println("req:", url, "-- response:", string(rb))
	}
}
