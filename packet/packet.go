package packet

import (
	"encoding/binary"
	"math"

	"go.uber.org/atomic"
)

type Packet struct {
	mode binary.ByteOrder
	pType PacketType
	body []byte
	id int32
}

var nextClientPacketID = atomic.Int32{}

func init() {
	nextClientPacketID.Store(1)
}

func idInArr(arr []int32, id atomic.Int32) bool {
	for _, v := range arr {
		if v == id.Load() {
			return true
		}
	}

	return false
}

func getNextID(restrictedIDs []int32) int32 {
	if nextClientPacketID.Load()+1 == math.MaxInt32 {
		nextClientPacketID.Store(1)
	} else {
		nextClientPacketID.Inc()
	}

	// Check if the current nextClientPacketID is a restricted ID and increment it until it no longer is
	for idInArr(restrictedIDs, nextClientPacketID) {
		if nextClientPacketID.Load()+1 == math.MaxInt32 {
			nextClientPacketID.Store(1)
		} else {
			nextClientPacketID.Inc()
		}
	}

	return nextClientPacketID.Load()
}