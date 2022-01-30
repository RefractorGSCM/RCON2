package packet

import (
	"bytes"
	"encoding/binary"
	"io"
	"math"

	"github.com/pkg/errors"
	"go.uber.org/atomic"
)

type PacketHeader struct {
	Size int32
	ID int32
	Type PacketType
}
type Packet struct {
	PacketHeader

	Mode binary.ByteOrder
	Body []byte
}

var nextClientPacketID = atomic.Int32{}

func init() {
	nextClientPacketID.Store(0)
}

func idInArr(arr []int32, id atomic.Int32) bool {
	for _, v := range arr {
		if v == id.Load() {
			return true
		}
	}

	return false
}

func nextID(restrictedIDs []int32) int32 {
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

const (
	headerBytes = 8
	padBytes = 2
)

func New(mode binary.ByteOrder, pType PacketType, body []byte, restrictedIDs []int32) *Packet {
	id := nextID(restrictedIDs)

	p := &Packet{
		PacketHeader: PacketHeader{
			ID: id,
			Type: pType,
			Size: int32(len(body) + headerBytes + padBytes),
		},
		Mode: mode,
		Body: body,
	}

	if len(body) == 0 {
		p.Body = []byte{}
	}

	return p
}

func (p *Packet) Encode() ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})

	order := p.Mode

	if err := binary.Write(buffer, order, p.Size); err != nil {
		return nil, errors.Wrap(err, "could not write packet size")
	}

	if err := binary.Write(buffer, order, p.ID); err != nil {
		return nil, errors.Wrap(err, "could not write packet id")
	}

	if err := binary.Write(buffer, order, p.Type); err != nil {
		return nil, errors.Wrap(err, "could not write packet type")
	}

	if err := binary.Write(buffer, order, p.Body); err != nil {
		return nil, errors.Wrap(err, "could not write packet body")
	}

	if err := binary.Write(buffer, order, [2]byte{}); err != nil {
		return nil, errors.Wrap(err, "could not write packet padding")
	}

	return buffer.Bytes(), nil
}

var malformedPacketErr = errors.New("malformed packet")
var badPacketTypeErr = errors.New("bad packet type")

func Decode(mode binary.ByteOrder, reader io.Reader) (*Packet, error) {
	header := PacketHeader{}

	if err := binary.Read(reader, mode, &header); err != nil {
		return nil, errors.Wrap(err, malformedPacketErr.Error())
	}

	payload := make([]byte, header.Size - headerBytes)
	_, err := io.ReadFull(reader, payload)
	if err != nil {
		return nil, errors.Wrap(err, malformedPacketErr.Error())
	}

	if header.Type != TypeAuthRes && header.Type != TypeCommandRes {
		return nil, badPacketTypeErr
	}

	return &Packet{
		PacketHeader: header,
		Body: payload[:len(payload)-2],
	}, nil
}