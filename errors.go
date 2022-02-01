package rcon

import "errors"

var ErrNotConnected = errors.New("not connected")
var ErrAuthentication = errors.New("authentication failed")
var ErrQueueTimeout = errors.New("queue timeout")
var ErrReadTimeout = errors.New("read timeout")