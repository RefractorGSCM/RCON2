package rcon

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"

	"github.com/pkg/errors"
)

type Client struct {
	config *Config

	conn *net.TCPConn

	log Logger
}

type Config struct {
	Host string
	Port uint16
	Password string

	// ConnTimeout is the timeout used when establishing a new RCON connection
	//
	// Default: 5 seconds
	ConnTimeout time.Duration

	// EndianMode represents the byte order being used by whatever game you're using this module with.
	// Valve games typically use little endian, but others may use big endian. ymmv.
	EndianMode binary.ByteOrder

	ReadDeadline time.Duration
	WriteDeadline time.Duration
}


var DefaultConfig = &Config{
	ConnTimeout: time.Second*5,
	EndianMode: binary.LittleEndian,
	ReadDeadline: time.Second*2,
	WriteDeadline: time.Second*2,
}

func NewClient(host string, port uint16, password string) *Client {
	config := DefaultConfig
	config.Host = host
	config.Port = port
	config.Password = password

	return &Client{
		config: config,
	}
}

func NewClientFromConfig(config *Config) *Client {
	return &Client{
		config: config,
	}
}

func (c *Client) SetLogger(logger Logger) {
	c.log = logger
}

func (c *Client) Connect() error {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", c.config.Host, c.config.Port), c.config.ConnTimeout)
	if err != nil {
		return errors.Wrap(err, "tcp dial error")
	}
	c.log.Debug("Connection established")

	var ok bool
	c.conn, ok = conn.(*net.TCPConn)
	if !ok {
		return errors.Wrap(err, "tcp dial error")
	}

	return nil
}