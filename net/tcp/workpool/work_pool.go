package main

import (
	"log"
	"net"
	"sync"
)

type pool struct {
	worknum   int
	taskQueue chan net.Conn

	mu     sync.Mutex
	closed bool
	done   chan struct{}
}

func newPool(n int) *pool {
	return &pool{
		worknum:   n,
		taskQueue: make(chan net.Conn, 100),
		done:      make(chan struct{}),
	}
}

func (p *pool) start() {
	for i := 0; i < p.worknum; i++ {
		go p.startWorker()
	}
}

func (p *pool) startWorker() {
	for {
		select {
		case <-p.done:
			log.Println("worker exit")
			return
		case conn := <-p.taskQueue:
			if conn != nil {
				hancleConn(conn)
			}
		}
	}
}

func (p *pool) close() {
	p.mu.Lock()
	p.closed = true
	close(p.done)
	close(p.taskQueue)
	p.mu.Unlock()
}

func (p *pool) pushTask(conn net.Conn) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return
	}
	p.taskQueue <- conn
}
