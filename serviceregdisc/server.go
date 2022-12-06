package serviceregdisc

import (
	"context"
	"fmt"
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
			fmt.Println("discover received event", s.path)
			s.Lock()
			s.servers = ser.Server
			s.total = len(s.servers)
			s.Unlock()
		}
	}()
}