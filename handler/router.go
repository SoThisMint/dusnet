package handler

import (
	zcodec "dusnet/codec"
	"dusnet/packet"
)

// IHandler 业务处理器接口，所有业务子类handler均实现此接口
type IHandler interface {
	IBaseHandler
	HandleMsg(packet.IPacket) error
}

// RegisterChildHandler 注册路由handler公用函数
func RegisterChildHandler(routerId uint32, handler IHandler) {
	childHandlerMap[routerId] = handler
}

func AllChildHandlers() map[uint32]IHandler {
	return childHandlerMap
}

var childHandlerMap = map[uint32]IHandler{
	2000: &sync2000Handler{},
	3000: &rpc3000Handler{
		baseHandler{
			codec0: zcodec.Default(),
		},
	},
}
