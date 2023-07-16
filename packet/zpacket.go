package packet

type IPacket interface {
	GetHeadLen() uint32
	GetBodyLen() uint32

	GetID() uint32
	GetData() []byte
	SetID(uint32)
	SetData([]byte)
	GetType() uint16
	SetType(uint16)
}

// Packet TLV结构
// tag + length + value
//
//	or
//
// tag + length + value + crc
type Packet struct {
	PacketHead
	PacketBody
}

func (p *Packet) GetType() uint16 {
	return p.Type
}

func (p *Packet) SetType(u uint16) {
	p.Type = u
}

func (p *Packet) GetID() uint32 {
	return p.ID
}

func (p *Packet) GetData() []byte {
	return p.Data
}

func (p *Packet) SetID(id uint32) {
	p.ID = id
}

func (p *Packet) SetData(bytes []byte) {
	p.Data = bytes
}

func (p *Packet) GetHeadLen() uint32 {
	return 9
}

func (p *Packet) GetBodyLen() uint32 {
	if p.Length == 0 {
		p.Length = uint32(len(p.PacketBody.Data))
	}
	return p.Length
}

type PacketHead struct {
	ID     uint32 // 包id
	Type   uint16 // 包类型
	Length uint32 // 包体长度
}

type PacketBody struct {
	Data []byte // 包体
}
