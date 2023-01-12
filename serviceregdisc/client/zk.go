package client

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/go-zookeeper/zk"
	"github.com/ilaziness/gopkg/serviceregdisc"
)

var _ serviceregdisc.Client = &ZKClient{}

type ZKClient struct {
	conn       *zk.Conn
	host       []string
	user, pass string
	acl        []zk.ACL
}

// NewZKClient
// host []string{"192.168.1.2:2181", "192.168.1.3:2181"}
func NewZKClient(host []string, user, pass string) (*ZKClient, error) {
	c := &ZKClient{
		host: host,
		user: user,
		pass: pass,
		acl:  zk.DigestACL(zk.PermAll, user, pass), //给节点添加用户验证
	}

	conn, _, err := zk.Connect(host, time.Second*60)
	if err != nil {
		return nil, err
	}
	c.conn = conn

	if err := c.addAuth(); err != nil {
		c.conn.Close()
		return nil, err
	}

	return c, nil
}

// Register 注册节点
// example:
// path: /cc/service/endpoint/user/192.168.2.1
// 先创建节点/cc/service/endpoint/user，再创建临时节点/cc/service/endpoint/user/192.168.2.1
// 服务信息创建临时节点的作用是如果服务不可用了，会话超时被销毁，注册的节点也会被销毁，起到了监控服务存活的目的
// CreateProtectedEphemeralSequential 是创建受保护的临时顺序节点，作用是如果服务器崩溃了，重连到其他服务器可以继续保持前一个服务器的会话
func (zkc *ZKClient) Register(ctx context.Context, path string, data []byte) error {
	go func() {
		var (
			regPath  string
			exists   bool
			err      error
			watchEvt <-chan zk.Event
		)
		for {
			if exists {
				// 添加一个watch，监控变化
				exists, _, watchEvt, err = zkc.conn.ExistsW(regPath)
				if err != nil {
					switch err {
					case zk.ErrClosing, zk.ErrConnectionClosed:
						time.Sleep(time.Second * 5)
					default:
						zkc.conn.Delete(regPath, -1)
						exists = false
						time.Sleep(time.Second * 5)
						continue
					}
				}
			}

			if !exists {
				// 注册节点
				regPath, err = zkc.createRegNode(path, data)
				if err != nil {
					log.Println("register service error:", err)
					time.Sleep(time.Second * 1)
					continue
				}
				log.Printf("register path: %s, zk path: %s\n", path, regPath)
				exists = true
				continue
			}

			select {
			case <-ctx.Done():
				log.Println("service register exit")
				zkc.conn.Delete(regPath, -1)
				zkc.conn.Close()
				return
			case e := <-watchEvt:
				log.Printf("watch register node(%s) exist changed, event(%v)\n", path, e)
				continue
			}
		}
	}()

	return nil
}

// createRegNode 创建注册节点
func (zkc *ZKClient) createRegNode(path string, data []byte) (string, error) {
	pathPart := strings.Split(path, "/")
	parentPath := strings.Join(pathPart[:len(pathPart)-1], "/")
	exists, _, err := zkc.conn.Exists(parentPath)
	if err != nil && err == zk.ErrNoAuth {
		zkc.addAuth()
		exists, _, err = zkc.conn.Exists(parentPath)
	}
	if err != nil {
		return "", err
	}
	if !exists {
		// 创建父节点
		tmpPath := ""
		for _, p := range pathPart[1 : len(pathPart)-1] {
			tmpPath += "/" + p
			e, _, err := zkc.conn.Exists(tmpPath)
			if err != nil {
				return "", err
			}
			if e {
				continue
			}
			_, err = zkc.conn.Create(tmpPath, []byte{}, 0, zkc.acl)
			if err != nil && err != zk.ErrNodeExists {
				return "", err
			}
		}
	}
	return zkc.conn.CreateProtectedEphemeralSequential(path, data, zkc.acl)
}

// addAuth 添加用户验证
func (zkc *ZKClient) addAuth() error {
	auth := zkc.user + ":" + zkc.pass
	return zkc.conn.AddAuth("digest", []byte(auth))
}

// Discovery 服务发现
// 添加一个watch，当服务信息发生变化时重新获取
// 服务器信息通过chan传递
func (zkc *ZKClient) Discovery(ctx context.Context, path string, event chan *serviceregdisc.DiscoverEvent) {
	for {
		// 设置一个watch
		_, _, watchEvt, err := zkc.conn.ChildrenW(path)
		if err != nil && err == zk.ErrNoAuth {
			zkc.addAuth()
			_, _, watchEvt, err = zkc.conn.ChildrenW(path)
		}
		if err != nil {
			if err == zk.ErrNoNode {
				// 要监控的服务不存在，等待后重新获取
				log.Printf("path: %s is not exists, will watch after 5s\n", path)
			} else {
				log.Printf("path: %s, discover watch error: %s\n", path, err)
			}
			time.Sleep(time.Second * 5)
			continue
		}
		// 第一次主动获取一次服务信息
		// 后面当服务信息有变化时，将会发生watchEvt事件，触发再次获取
		event <- zkc.getServerInfoByPath(path)

		select {
		case <-ctx.Done():
			close(event)
			log.Println("discover exit")
			return
		case e := <-watchEvt:
			log.Printf("watch found the children of path(%s) change. event type:%s, event err:%v\n", path, e.Type.String(), e.Err)
		}
	}
}

// getServerInfoByPath 获取指定路径下的服务信息
// 比如/cc/service/endpoint/user 用户服务，会获取到所有user服务的所有服务节点信息
func (zkc *ZKClient) getServerInfoByPath(path string) *serviceregdisc.DiscoverEvent {
	eventData := &serviceregdisc.DiscoverEvent{Server: make([][]byte, 0)}
	nodes, _, err := zkc.conn.Children(path)
	if err != nil {
		log.Printf("discover get service node error: %s\n", err)
		return eventData
	}
	for _, node := range nodes {
		serverPath := path + "/" + node
		nodeData, _, err := zkc.conn.Get(serverPath)
		if err != nil {
			log.Printf("discover get server info node error: %s\n", err)
			continue
		}
		eventData.Server = append(eventData.Server, nodeData)
	}
	return eventData
}
