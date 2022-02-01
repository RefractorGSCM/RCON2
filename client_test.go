package rcon

import (
	"encoding/binary"
	"testing"

	"github.com/refractorgscm/rcon2/fakeserver"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	client := NewClient("localhost", 1234, "suPerSecure123")

	assert.Equal(t, "localhost", client.config.Host)
	assert.Equal(t, uint16(1234), client.config.Port)
	assert.Equal(t, "suPerSecure123", client.config.Password)
}

func TestNewClientFromConfig(t *testing.T) {
	config := &Config{
		Host: "localhost",
		Port: 1234,
		Password: "suPerSecure123",
	}

	client := NewClientFromConfig(config)

	assert.Equal(t, "localhost", client.config.Host)
	assert.Equal(t, uint16(1234), client.config.Port)
	assert.Equal(t, "suPerSecure123", client.config.Password)
}

func TestClient_Connect(t *testing.T) {
	fakeServer := fakeserver.New(9898, binary.LittleEndian)
	go fakeServer.Listen()

	client := NewClient("localhost", 9898, "suPerSecure123")
	err := client.Connect()
	assert.Nil(t, err)
}

func TestClient_RunCommand(t *testing.T) {
	fakeServer := fakeserver.New(9899, binary.LittleEndian)
	go fakeServer.Listen()

	client := NewClient("localhost", 9899, "suPerSecure123")
	err := client.Connect()
	assert.Nil(t, err)

	res, err := client.RunCommand("help")
	assert.Nil(t, err)
	assert.Equal(t, "firstplayer\notherplayer\nlastplayer", res)
}