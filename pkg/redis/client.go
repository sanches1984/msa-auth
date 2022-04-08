package redis

import (
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"net"
	"sync"
	"time"
)

var ErrRecordNotFound = errors.New("record not found")

type Config struct {
	Host              string
	Password          string
	Db                int
	ConnectionTimeout time.Duration
	OperationTimeout  time.Duration
}

type Client struct {
	sync.Mutex
	config Config
	conn   *redis.Conn
}

func NewClient(config Config) (*Client, error) {
	client := &Client{
		config: config,
	}
	err := client.doConnect()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (c *Client) Set(key string, value []byte) error {
	if value == nil {
		value = []byte{}
	}
	_, err := c.do("SET", key, value)
	return err
}

func (c *Client) Delete(key string) error {
	_, err := c.do("DEL", key)
	return err
}

func (c *Client) Get(key string) ([]byte, error) {
	data, err := c.do("GET", key)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, ErrRecordNotFound
	}
	return data.([]byte), nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return (*c.conn).Close()
	}
	return nil
}

func (c *Client) doConnect() error {
	redisConn, err := redis.Dial(
		"tcp",
		c.config.Host,
		redis.DialPassword(c.config.Password),
		redis.DialDatabase(c.config.Db),
		redis.DialReadTimeout(c.config.OperationTimeout),
		redis.DialWriteTimeout(c.config.OperationTimeout),
		redis.DialConnectTimeout(c.config.ConnectionTimeout),
	)
	if err != nil {
		return fmt.Errorf("can't connect to redis: %v", err)
	}
	c.conn = &redisConn
	return nil
}

func (c *Client) doReconnect() error {
	if c.conn != nil {
		_ = (*c.conn).Close()
	}
	c.conn = nil
	return c.doConnect()
}

func (c *Client) do(commandName string, args ...interface{}) (interface{}, error) {
	if c.conn == nil {
		return nil, fmt.Errorf("redis disconnected")
	}
	firstAttemp := true
	c.Lock()
TryAgain:
	reply, err := (*c.conn).Do(commandName, args...)
	if firstAttemp && err != nil {
		if _, ok := err.(*net.OpError); ok {
			err = c.doReconnect()
			if err == nil {
				firstAttemp = false
				goto TryAgain
			}
		}
	}
	c.Unlock()
	return reply, err
}
