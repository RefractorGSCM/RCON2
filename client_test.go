package client

import (
	"testing"

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