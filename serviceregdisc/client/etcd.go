package client

import (
	"context"
	"log"
	"time"

	"github.com/ilaziness/gopkg/serviceregdisc"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type serverInfo map[string][]byte

type EtcdClient struct {
	client     *clientv3.Client
	user, pass string
	ctx        context.Context
}

var _ serviceregdisc.Client = &EtcdClient{}

func NewEtcdClient(ctx context.Context, endpoints []string, user, pass string) (*EtcdClient, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	return &EtcdClient{
		client: cli,
		user:   user,
		pass:   pass,
		ctx:    ctx,
	}, nil
}

// Register 服务注册
// path 注册路径， data节点数据
// example:
// path: /cc/service/endpoint/user/192.168.2.1
// path作为数据存储key，带租约存储数据，定时续期，续期的相当于心跳的作用，服务失效了，租约过期，数据会被删除，节点信息也就不存在了
func (e *EtcdClient) Register(ctx context.Context, path string, data []byte) error {
	lease, err := e.createLease()
	if err != nil {
		return err
	}
	kv := clientv3.NewKV(e.client)
	_, err = kv.Put(ctx, path, string(data), clientv3.WithLease(lease.ID))
	return err
}

// createLease 创建租约
func (e *EtcdClient) createLease() (*clientv3.LeaseGrantResponse, error) {
	lease := clientv3.NewLease(e.client)
	// 创建10s租约
	leaseResp, err := lease.Grant(e.ctx, 10)
	if err != nil {
		return nil, err
	}

	// 续期租约
	karChan, err := lease.KeepAlive(e.ctx, leaseResp.ID)
	go func() {
		for {
			select {
			case <-e.ctx.Done():
				log.Println("ctx cancel")
				lease.Close()
				e.client.Close()
				return
			case resp := <-karChan:
				if resp == nil {
					log.Println("租约已关闭")
					return
				}
				log.Println("租约续租成功")
			}
		}
	}()

	if err != nil {
		return nil, err
	}
	return leaseResp, nil
}

// Discovery 服务发现
// 第一次先用前缀key主动获取一次服务信息
// 之后用watch监视前缀key，监视数据变化更新服务信息
func (e *EtcdClient) Discovery(ctx context.Context, path string, event chan *serviceregdisc.DiscoverEvent) {
	//第一次取出服务的所有节点信息，存储在si
	si := e.getServerInfo(path)
	watcher := clientv3.NewWatcher(e.client)
	watchChan := watcher.Watch(ctx, path, clientv3.WithPrefix())

	event <- convServerInfo(si)

	for {
		select {
		case <-ctx.Done():
			close(event)
			log.Println("discover exit")
			return
		case ev := <-watchChan:
			if ev.Err() != nil {
				log.Println("watchChan error:", ev.Err())
				continue
			}
			log.Printf("watch service %s change\n", path)
			// 监视到变化在更新si，之后通知发现服务
			event <- e.updateServerInfo(si, ev)
		}
	}
}

// getServerInfo 获取服务所有节点
func (e *EtcdClient) getServerInfo(path string) serverInfo {
	info := make(serverInfo)
	resp, err := e.client.Get(e.ctx, path, clientv3.WithPrefix())
	if err != nil {
		log.Println("etcd Discovery get key error:", err)
		return info
	}
	if resp == nil || resp.Kvs == nil {
		log.Println("etcd Discovery get key error: resp is nil")
		return info
	}
	for i := range resp.Kvs {
		if v := resp.Kvs[i].Value; v != nil {
			info[string(resp.Kvs[i].Key)] = v
		}
	}

	return info
}

// updateServerInfo 更新服务节点信息
func (e *EtcdClient) updateServerInfo(si serverInfo, wresp clientv3.WatchResponse) *serviceregdisc.DiscoverEvent {
	for _, ev := range wresp.Events {
		switch ev.Type {
		case mvccpb.PUT:
			si[string(ev.Kv.Key)] = ev.Kv.Value
		case mvccpb.DELETE:
			delete(si, string(ev.Kv.Key))
		}
	}
	return convServerInfo(si)
}

// convServerInfo 转换服务器节点信息为所需的event数据
func convServerInfo(si serverInfo) *serviceregdisc.DiscoverEvent {
	eventData := &serviceregdisc.DiscoverEvent{Server: make([][]byte, 0)}
	for _, v := range si {
		eventData.Server = append(eventData.Server, v)
	}
	return eventData
}
