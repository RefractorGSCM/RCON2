package packet

import (
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
	assert.Equal(t, int32(2), getNextID(restrictedIDs))
	assert.Equal(t, int32(3), getNextID(restrictedIDs))
	assert.Equal(t, int32(4), getNextID(restrictedIDs))
	assert.Equal(t, int32(5), getNextID(restrictedIDs))

	// Test max int value reset
	nextClientPacketID.Store(math.MaxInt32-1)
	assert.Equal(t, int32(1), getNextID(restrictedIDs))

	// Test skipping restricted values
	nextClientPacketID.Store(999)
	assert.Equal(t, int32(1004), getNextID(restrictedIDs))
	assert.Equal(t, int32(1010), getNextID(restrictedIDs))
	assert.Equal(t, int32(1011), getNextID(restrictedIDs))
}