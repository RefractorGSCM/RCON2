package packet

import (
	"bytes"
	"encoding/binary"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"
)

func Test_idInArr(t *testing.T) {
	id := atomic.Int32{}
	id.Store(1000)

	arr := []int32{999, 1001, 1002, 1003}
	assert.False(t, idInArr(arr, id))

	arr = []int32{999, 1000, 1002, 1003}
	assert.True(t, idInArr(arr, id))
}

func Test_getNextID(t *testing.T) {
	nextClientPacketID.Store(1)

	restrictedIDs := []int32{1000, 1001, 1002, 1003, 1005, 1006, 1007, 1008, 1009}

	// Test correct incrementing
	assert.Equal(t, int32(2), nextID(restrictedIDs))
	assert.Equal(t, int32(3), nextID(restrictedIDs))
	assert.Equal(t, int32(4), nextID(restrictedIDs))
	assert.Equal(t, int32(5), nextID(restrictedIDs))

	// Test max int value reset
	nextClientPacketID.Store(math.MaxInt32-1)
	assert.Equal(t, int32(1), nextID(restrictedIDs))

	// Test skipping restricted values
	nextClientPacketID.Store(999)
	assert.Equal(t, int32(1004), nextID(restrictedIDs))
	assert.Equal(t, int32(1010), nextID(restrictedIDs))
	assert.Equal(t, int32(1011), nextID(restrictedIDs))
}

func TestNew(t *testing.T) {
	mode := binary.LittleEndian
	pType := TypeCommand
	body := []byte("hello world!")

	nextClientPacketID.Store(0)
	packet := New(mode, pType, body, nil)
	assert.Equal(t, int32(1), packet.ID)
	assert.Equal(t, pType, packet.Type)
	assert.Equal(t, mode, packet.Mode)
	assert.Equal(t, body, packet.Body)
}

func TestPacketBuildAndDecode(t *testing.T) {
	const expectedBody = "test command string"
	const expectedSize = int32(len(expectedBody) + headerBytes + padBytes)
	const expectedID = int32(10)
	const expectedType = int32(TypeCommand)

	nextClientPacketID.Store(expectedID - 1)
	packet := New(binary.LittleEndian, TypeCommand, []byte(expectedBody), []int32{})

	// 1. Test packet build
	out, err := packet.Encode()
	assert.Nil(t, err)

	buffer := bytes.NewBuffer(out)
	
	var size int32
	var id int32
	var pType int32

	// 1.5 Check packet fields
	err = binary.Read(buffer, packet.Mode, &size)
	assert.Nil(t, err)
	assert.Equal(t, expectedSize, size)

	err = binary.Read(buffer, packet.Mode, &id)
	assert.Nil(t, err)
	assert.Equal(t, expectedID, id)

	err = binary.Read(buffer, packet.Mode, &pType)
	assert.Nil(t, err)
	assert.Equal(t, expectedType, pType)

	// 2. Test packet decode
	out, err = packet.Encode()
	assert.Nil(t, err)
	buffer = bytes.NewBuffer(out)

	packet, err = Decode(packet.Mode, buffer)
	assert.Nil(t, err)

	assert.Equal(t, expectedSize, packet.Size)
	assert.Equal(t, expectedID, packet.ID)
	assert.Equal(t, expectedType, int32(packet.Type))
	assert.Equal(t, expectedBody, string(packet.Body))
}