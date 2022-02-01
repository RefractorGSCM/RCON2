package rcon

import (
	"bufio"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/refractorgscm/rcon2/packet"
)

func (c *Client) sendPacket(p *packet.Packet) error {
	out, err := p.Encode()
	if err != nil {
		return errors.Wrap(err, "could not encode packet")
	}

	if err := c.writeToConn(out); err != nil {
		return errors.Wrap(err, "could not send packet")
	}

	return nil
}

func (c *Client) readPacket() (*packet.Packet, error) {
	if c.conn == nil {
		return nil, ErrNotConnected
	}

	if err := c.conn.SetDeadline(time.Time{}); err != nil {
		if strings.HasSuffix(err.Error(), "use of closed network connection") {
			return nil, ErrNotConnected
		}

		return nil, errors.Wrap(err, "could not set connection deadline")
	}

	reader := bufio.NewReader(c.conn)

	res, err := packet.Decode(c.config.EndianMode, reader)
	if err != nil {
		if strings.HasSuffix(err.Error(), "use of closed network connection") {
			return nil, ErrNotConnected
		}

		return nil, errors.Wrap(err, "could not read packet")
	}

	c.log.Debug("Read packet ID: ", res.ID, ", Body: ", string(res.Body))

	return res, nil
}

func (c *Client) readPacketTimeout() (*packet.Packet, error) {
	if c.conn == nil {
		return nil, ErrNotConnected
	}

	if err := c.conn.SetDeadline(time.Now().Add(c.config.ReadDeadline)); err != nil {
		if strings.HasSuffix(err.Error(), "use of closed network connection") {
			return nil, ErrNotConnected
		}

		return nil, errors.Wrap(err, "could not set connection deadline")
	}

	reader := bufio.NewReader(c.conn)

	res, err := packet.Decode(c.config.EndianMode, reader)
	if err != nil {
		if strings.HasSuffix(err.Error(), "use of closed network connection") {
			return nil, ErrNotConnected
		}

		return nil, errors.Wrap(err, "could not read packet")
	}

	return res, nil
}

func (c *Client) writeToConn(data []byte) error {
	_, err := c.conn.Write(data)
	return err
}