package client

import "net"

type Client struct {
	config *Config

	conn *net.TCPConn
}

type Config struct {
	Host string
	Port uint16
	Password string
}

func NewClient(host string, port uint16, password string) *Client {
	return &Client{
		config: &Config{
			Host: host,
			Port: port,
			Password: password,
		},
	}
}

func NewClientFromConfig(config *Config) *Client {
	return &Client{
		config: config,
	}
}