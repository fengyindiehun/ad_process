package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/go-ini/ini"
	"time"
)

type redisProxy struct {
	host, port                                string
	connectTimeout, readTimeout, writeTimeout int
	conn                                      redis.Conn
	isConnected                               bool
}

func (proxy *redisProxy) doConnect() bool {
	if proxy.isConnected == true {
		return true
	}

	var err error
	proxy.conn, err = redis.DialTimeout(
		"tcp",
		proxy.host+":"+proxy.port,
		time.Duration(proxy.connectTimeout)*time.Millisecond,
		time.Duration(proxy.readTimeout)*time.Millisecond,
		time.Duration(proxy.writeTimeout)*time.Millisecond)
	if err != nil {
		proxy.isConnected = false
		return false
	} else {
		proxy.isConnected = true
		return true
	}
}

func (proxy *redisProxy) set(key string, value string) error {
	if proxy.doConnect() == false {
		return fmt.Errorf("redis connect error")
	}
	_, err := proxy.conn.Do("SET", key, value)
	if err != nil {
		return err
	}
	return nil
}

func (proxy *redisProxy) mSet(args []interface{}) error {
	if proxy.doConnect() == false {
		return fmt.Errorf("redis connect error")
	}
	log.Debug(args)
	_, err := proxy.conn.Do("MSET", args...)
	if err != nil {
		return err
	}
	return nil
}

func (proxy *redisProxy) get(key string) (value string, err error) {
	if proxy.doConnect() == false {
		err = fmt.Errorf("redis connect error")
		return
	}
	value, err = redis.String(proxy.conn.Do("GET", key))
	if err != nil {
		return
	}
	return
}

func newRedisProxy(section *ini.Section) (*redisProxy, error) {
	host := section.Key(REDIS_SECTION_HOST).String()
	if host == "" {
		return nil, fmt.Errorf("redis host parse error")
	}

	port := section.Key(REDIS_SECTION_PORT).String()
	if port == "" {
		return nil, fmt.Errorf("redis port parse error")
	}

	connectTimeout, err := section.Key(REDIS_SECTION_CONNECT_TIMEOUT).Int()
	if err != nil {
		return nil, fmt.Errorf("redis connect_timeout parse error")
	}

	readTimeout, err := section.Key(REDIS_SECTION_READ_TIMEOUT).Int()
	if err != nil {
		return nil, fmt.Errorf("redis read_timeout parse error")
	}

	writeTimeout, err := section.Key(REDIS_SECTION_WRITE_TIMEOUT).Int()
	if err != nil {
		return nil, fmt.Errorf("redis write_timeout parse error")
	}

	conn, err := redis.DialTimeout(
		"tcp",
		host+":"+port,
		time.Duration(connectTimeout)*time.Millisecond,
		time.Duration(readTimeout)*time.Millisecond,
		time.Duration(writeTimeout)*time.Millisecond)
	if err != nil {
		return nil, err
	}

	return &redisProxy{
		host:           host,
		port:           port,
		conn:           conn,
		connectTimeout: connectTimeout,
		readTimeout:    readTimeout,
		writeTimeout:   writeTimeout,
		isConnected:    true,
	}, nil
}
