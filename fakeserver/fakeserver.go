package fakeserver

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"net"

	"github.com/refractorgscm/rcon2/packet"
)

var AuthShouldFail = false

type Server struct {
	port string
	mode binary.ByteOrder
	ctx context.Context
	terminate bool
}

func New(port uint16, mode binary.ByteOrder) *Server {
	portStr := fmt.Sprintf(":%d", port)

	return &Server{
		port: portStr,
		mode: mode,
	}
}

func (s *Server) Listen() {
	l, err := net.Listen("tcp4", s.port)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	fmt.Printf("Fake server listening on port %s\n", s.port)

	for {
		if s.terminate {
			break
		}

		c, err := l.Accept()
		if err != nil {
			panic(err)
		}

		go s.handleConnection(c)
	}
}

func (s *Server) Stop() {
	s.terminate = true
}

func (s *Server) handleConnection(conn net.Conn) {
	for {
		reader := bufio.NewReader(conn)

		p, err := packet.Decode(s.mode, reader)
		if err != nil {
			panic(err)
		}

		fmt.Println(p)

		var res *packet.Packet

		switch p.Type {
		case packet.TypeAuth:
			res = s.handleAuthReq(p)
		case packet.TypeCommand:
			switch string(p.Body) {
			case "help":
				res = packet.New(p.Mode, packet.TypeCommandRes, []byte("firstplayer\notherplayer\nlastplayer"), nil)
			default:
				res = packet.New(p.Mode, packet.TypeCommandRes, []byte("Unknown command"), nil)
			}
		}

		out, _ := res.Encode()

		printBytes(out)

		conn.Write(out)
	}
}

func (s *Server) handleAuthReq(p *packet.Packet) *packet.Packet {
	res := packet.New(s.mode, packet.TypeAuthRes, []byte{}, []int32{})
	res.ID = p.ID

	if AuthShouldFail {
		res.ID = packet.AuthFailedID
		return res
	}

	return res
}

func printBytes(arr []byte) {
	fmt.Printf("Bytes (%d): ", len(arr))
	for _, b := range arr {
		fmt.Printf("%x ", b)
	}
	fmt.Print("\n")
}