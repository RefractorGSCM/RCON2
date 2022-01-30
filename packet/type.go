package packet

type PacketType int32

const (
	TypeAuth = PacketType(3)
	TypeAuthRes = PacketType(2)
	TypeCommand = PacketType(2)
	TypeCommandRes = PacketType(0)
	
	AuthFailedID = -1
)