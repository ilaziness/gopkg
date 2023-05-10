package httpclient

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const (
	ServiceNameAs = "as"
	ServiceNameBs = "bs"
)

var host = map[string]string{
	"as": "http://127.0.0.1:8088",
	"bs": "http://127.0.0.1:8089",
}

var client http.Client

func init() {
	client = http.Client{
		Transport: otelhttp.NewTransport(nil),
		Timeout:   time.Second * 2,
	}
}

// Get http get请求服务接口
func Get(ctx context.Context, serviceName string, path string) ([]byte, error) {
	if _, ok := host[serviceName]; !ok {
		return nil, errors.New("服务不存在")
	}
	url := host[serviceName] + path

	//..... 直接Get
	//otelhttp.DefaultClient.Timeout = time.Second * 2
	//resp, err := otelhttp.Get(ctx, url)

	//.... 自定义client对象
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
