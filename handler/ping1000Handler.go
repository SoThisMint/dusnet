package handler

import (
	"dusnet/logger"
	"dusnet/packet"
)

type Ping1000Handler struct {
	baseHandler
}

// HandleMsg ping handler只做业务应答，连接续租由基础处理器完成
func (p Ping1000Handler) HandleMsg(pkt packet.IPacket) error {
	logger.Debug("handle ping1000 msg with pkt:%+v", pkt)
	ackPkt := packet.Packet{}
	ackPkt.ID = pkt.GetID()
	ackPkt.Data = []byte("pong")
	return p.write(&ackPkt)
}
