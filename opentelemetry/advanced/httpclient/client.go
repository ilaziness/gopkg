package httpclient

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var host = map[string]string{
	"as": "http://127.0.0.1:8088",
	"bs": "http://127.0.0.1:8089",
}

var client http.Client

func init() {
	client = http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
}

func Get(path string) {

}
