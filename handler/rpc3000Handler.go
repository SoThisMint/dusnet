package handler

import (
	"dusnet/logger"
	"dusnet/packet"
)

type rpc3000Handler struct {
	baseHandler
}

func (p rpc3000Handler) HandleMsg(packet packet.IPacket) error {
	logger.Debug("handle rpc3000 msg with pkt:%+v", packet)
	return nil
}
