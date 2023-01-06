package serviceregdisc

import (
	"context"
	"encoding/json"
	"log"
	"sync"
)

// NewServerDiscover 创建对应服务的发现对象
// 每个服务独立的发现对象，依赖哪些服务时，按需创建即可
func NewServerDiscover(ctx context.Context, path string, discover RegisterDiscovery) (*Server, error) {
	eventChan, err := discover.Discovery(ctx, path)
	if err != nil {
		return nil, err
	}
	ser := &Server{
		path:          path,
		discoverEvent: eventChan,
		servers:       make([]*ServerInfo, 0),
		index:         0,
		total:         0,
	}
	ser.run()
	return ser, nil
}

// Server 对应一个服务对象
type Server struct {
	sync.Mutex
	index         int
	total         int
	path          string
	servers       []*ServerInfo // 服务节点
	discoverEvent <-chan *DiscoverEvent
}

// GetServer 获取一个服务链接
// 从服务的所有节点里面获取一个节点
func (s *Server) GetServer() string {
	if s.total == 0 {
		return ""
	}
	// 循环返回节点
	if s.index >= s.total {
		s.index = 0
	}
	addr := s.servers[s.index].GetAddress()
	s.index++
	return addr
}

// run 接收服务变化事件，更新server服务信息
func (s *Server) run() {
	go func() {
		for ser := range s.discoverEvent {
			// 拿到空数据，清空服务列表
			s.updateServers(ser.Server)
		}
	}()
}

func (s *Server) updateServers(data [][]byte) {
	servers := make([]*ServerInfo, 0)
	for _, val := range data {
		serverInfo := &ServerInfo{}
		err := json.Unmarshal(val, serverInfo)
		if err != nil {
			log.Println("unmarshal server info error:", err)
		}
		servers = append(servers, serverInfo)
	}
	s.Lock()
	s.servers = servers
	s.total = len(s.servers)
	s.Unlock()
	log.Printf("path: %s, update servers\n", s.path)
	for _, sv := range s.servers {
		log.Printf("path: %s, %#v", s.path, sv)
	}
}
