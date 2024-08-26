package rediscache

import (
	"sync"
	"time"
)

type Command string

const (
	SET Command = "SET"
	GET Command = "GET"
	DEL Command = "DEL"
)

type Request struct {
	Command Command
	Key     string
	Value   interface{}
	TTL     time.Duration
	Result  chan interface{}
}

type Cache struct {
	data       map[string]interface{}
	ttl        map[string]time.Time
	mu         sync.RWMutex
	requests   chan Request
	workerPool int
}

func NewRedisCache(workerPool int) *Cache {
	cache := &Cache{
		data:       make(map[string]interface{}),
		ttl:        make(map[string]time.Time),
		requests:   make(chan Request),
		workerPool: workerPool,
	}
	cache.startWorkers()
	go cache.cleanExpiredKeys()
	return cache
}

func (c *Cache) startWorkers() {
	for i := 0; i < c.workerPool; i++ {
		go func() {
			for req := range c.requests {
				c.handleRequest(req)
			}
		}()
	}
}

func (c *Cache) handleRequest(req Request) {
	switch req.Command {
	case SET:
		c.mu.Lock()
		c.data[req.Key] = req.Value
		if req.TTL > 0 {
			c.ttl[req.Key] = time.Now().Add(req.TTL)
		}
		c.mu.Unlock()
		req.Result <- true

	case GET:
		c.mu.RLock()
		if value, found := c.data[req.Key]; found {
			if expiry, ok := c.ttl[req.Key]; ok && time.Now().After(expiry) {
				c.mu.RUnlock()
				c.Del(req.Key)
				req.Result <- nil
			} else {
				c.mu.RUnlock()
				req.Result <- value
			}
		} else {
			c.mu.RUnlock()
			req.Result <- nil
		}

	case DEL:
		c.Del(req.Key)
		req.Result <- true
	}
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	req := Request{
		Command: SET,
		Key:     key,
		Value:   value,
		TTL:     ttl,
		Result:  make(chan interface{}),
	}
	c.requests <- req
	<-req.Result
}

func (c *Cache) Get(key string) interface{} {
	req := Request{
		Command: GET,
		Key:     key,
		Result:  make(chan interface{}),
	}
	c.requests <- req
	return <-req.Result
}

func (c *Cache) Del(key string) {
	req := Request{
		Command: DEL,
		Key:     key,
		Result:  make(chan interface{}),
	}
	c.requests <- req
	<-req.Result
}

func (c *Cache) cleanExpiredKeys() {
	for {
		time.Sleep(1 * time.Second)
		c.mu.Lock()
		for key, expiry := range c.ttl {
			if time.Now().After(expiry) {
				delete(c.data, key)
				delete(c.ttl, key)
			}
		}
		c.mu.Unlock()
	}
}
