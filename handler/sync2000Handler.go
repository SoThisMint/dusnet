package handler

import (
	"dusnet/logger"
	"dusnet/packet"
)

type sync2000Handler struct {
	baseHandler
}

func (p sync2000Handler) HandleMsg(packet packet.IPacket) error {
	logger.Debug("handle sync2000 msg with pkt:%+v", packet)
	return nil
}
