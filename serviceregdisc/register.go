package serviceregdisc

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// ServerInfo 服务信息
type ServerInfo struct {
	IP     string `json:"ip"`
	Port   uint   `json:"port"`
	Schema string `json:"schema"` //协议，比如https
	UUID   string `json:"uuid"`
}

func (s *ServerInfo) GetAddress() string {
	str := strings.Builder{}
	str.WriteString(s.Schema)
	str.WriteString("://")
	str.WriteString(s.IP)
	str.WriteString(":")
	str.WriteString(strconv.Itoa(int(s.Port)))

	return str.String()
}

// Client 数据存储客户端接口
type Client interface {
	Register(ctx context.Context, path string, data []byte) error
	Discovery(ctx context.Context, path string, event chan *DiscoverEvent)
}

const (
	RootPath = "service/endpoint"
)

// NewRegDisc 创建服务注册和发现对象
func NewRegDisc(prefix string, cli Client) *RegisterDiscovery {
	return &RegisterDiscovery{
		client: cli,
		prefix: prefix,
	}
}

// RegisterDiscovery 服务注册和发现对象
type RegisterDiscovery struct {
	client Client
	prefix string //通常是系统id
}

// GetServicePath 获取服务路径
func (rd *RegisterDiscovery) GetServicePath(id string) string {
	return fmt.Sprintf("%s/%s/%s", rd.prefix, RootPath, id)
}

// Register 服务注册
// id 服务ID，info 服务信息
func (rd *RegisterDiscovery) Register(ctx context.Context, id string, info ServerInfo) error {
	// 注册路径
	regPath := fmt.Sprintf("%s/%s/%s/%s", rd.prefix, RootPath, id, info.IP)
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}
	return rd.client.Register(ctx, regPath, data)
}

type DiscoverEvent struct {
	Server []*ServerInfo
}

// Discovery 服务发现
// DiscoverEvent 有变化是通过chan DiscoverEvent 通知
// path:要发现的服务，比如/cc/service/endpoint/user 用户服务
func (rd *RegisterDiscovery) Discovery(ctx context.Context, path string) (<-chan *DiscoverEvent, error) {
	event := make(chan *DiscoverEvent, 1)
	go rd.client.Discovery(ctx, path, event)

	return event, nil
}