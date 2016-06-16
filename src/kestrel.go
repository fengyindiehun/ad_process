package main

import (
	"encoding/json"
	"fmt"
	//"reflect"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-ini/ini"
)

type kestrelProxy struct {
	host, port, queue, timeout string
	c                          chan string
	conn                       *memcache.Client
}

//cause gomemcache.Get will judge whether it is connected to kestrel
//so we do not need to do connect anymore
func (proxy *kestrelProxy) receiveMessage() {
	for {
		it, err := proxy.conn.Get(proxy.queue + "/t=" + proxy.timeout)
		if err != nil {
			log.Error(err)
			continue
		}
		proxy.c <- string(it.Value)
	}
}

func (proxy *kestrelProxy) messages() chan string {
	return proxy.c
}

type kestrelMessage struct {
	Mid, Uid string
}

func (proxy *kestrelProxy) decodeMessage(message string) (feedId, custId string, err error) {
	var m kestrelMessage
	err = json.Unmarshal([]byte(message), &m)
	if err != nil {
		return
	}
	feedId = m.Mid
	custId = m.Uid
	return
}

func newKestrelProxy(section *ini.Section) (*kestrelProxy, error) {
	host := section.Key(KESTREL_SECTION_HOST).String()
	if host == "" {
		return nil, fmt.Errorf("kestrel host parse error")
	}

	port := section.Key(KESTREL_SECTION_PORT).String()
	if port == "" {
		return nil, fmt.Errorf("kestrel port parse error")
	}

	queue := section.Key(KESTREL_SECTION_QUEUE).String()
	if port == "" {
		return nil, fmt.Errorf("kestrel port parse error")
	}

	timeout := section.Key(KESTREL_SECTION_TIMEOUT).String()
	if timeout == "" {
		return nil, fmt.Errorf("kestrel timeout parse error")
	}

	c := make(chan string)

	conn := memcache.New(host + ":" + port)

	return &kestrelProxy{
		host:    host,
		port:    port,
		queue:   queue,
		timeout: timeout,
		c:       c,
		conn:    conn,
	}, nil
}
