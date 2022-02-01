package rcon

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/refractorgscm/rcon2/packet"
)

type Client struct {
	config *Config

	conn *net.TCPConn

	log Logger

	cmdMutex sync.Mutex
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

	RestrictedPacketIDs []int32
}


var DefaultConfig = &Config{
	ConnTimeout: time.Second*5,
	EndianMode: binary.LittleEndian,
	ReadDeadline: time.Second*2,
	WriteDeadline: time.Second*2,
	RestrictedPacketIDs: []int32{},
}

func NewClient(host string, port uint16, password string) *Client {
	config := DefaultConfig
	config.Host = host
	config.Port = port
	config.Password = password

	return NewClientFromConfig(DefaultConfig)
}

func NewClientFromConfig(config *Config) *Client {
	c := &Client{
		config: config,
	}

	if c.log == nil {
		c.log = &DefaultLogger{}
	}

	if c.config.EndianMode == nil {
		c.config.EndianMode = binary.LittleEndian
	}

	if c.config.ConnTimeout == 0 {
		c.config.ConnTimeout = DefaultConfig.ConnTimeout
	}

	if c.config.ReadDeadline == 0 {
		c.config.ConnTimeout = DefaultConfig.ReadDeadline
	}

	if c.config.WriteDeadline == 0 {
		c.config.ConnTimeout = DefaultConfig.WriteDeadline
	}

	return c
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

	if err := c.authenticate(); err != nil {
		c.log.Debug("Authentication failed", err)
		return err
	}

	return nil
}

func (c *Client) authenticate() error {
	p := c.newPacket(packet.TypeAuth, c.config.Password)

	if err := c.sendPacket(p); err != nil {
		return errors.Wrap(err, "could not send packet")
	}

	res, err := c.readPacketTimeout()
	if err != nil {
		return errors.Wrap(err, "could not get auth response")
	}

	if res.Type != packet.TypeAuthRes {
		return errors.New("packet was not of the type auth response")
	}

	if res.ID == packet.AuthFailedID {
		return errors.Wrap(ErrAuthentication, "authentication failed")
	}

	c.log.Debug("Authenticated")

	return nil
}

func (c *Client) RunCommand(cmd string) (string, error) {
	c.cmdMutex.Lock()
	defer c.cmdMutex.Unlock()

	p := c.newPacket(packet.TypeCommand, cmd)

	if err := c.sendPacket(p); err != nil {
		return "", err
	}

	res, err := c.readPacket()
	if err != nil {
		return "", err
	}

	return string(res.Body), nil
}

func (c *Client) newPacket(pType packet.PacketType, body string) *packet.Packet {
	return packet.New(c.config.EndianMode, pType, []byte(body), c.config.RestrictedPacketIDs)
}