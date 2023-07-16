package handler

import (
	zcodec "dusnet/codec"
	"dusnet/connect"
	"dusnet/logger"
	"dusnet/packet"
	"errors"
	"fmt"
)

// IBaseHandler 处理器基类接口
type IBaseHandler interface {
	SetConn(connect.IConnection)
	SetCodec(zcodec.Icodec)
	BindConn(connect.IConnection)
	write(packet.IPacket) error
}

// IRouteHandler 路由型处理器基类接口
type IRouteHandler interface {
	IBaseHandler
	HandleMsg0() error
}

// IBuilder 路由处理器构建接口
type IBuilder interface { // 默认路由handler构造器接口
	Codec(zcodec.Icodec) IBuilder
	Conn(connect.IConnection) IBuilder
	Build() IRouteHandler
}

type builder struct {
	handler IRouteHandler
}

type baseHandler struct {
	codec0 zcodec.Icodec       //编解码器
	conn   connect.IConnection // 路由handler和连接一对一绑定
}

func (h *baseHandler) write(pkt packet.IPacket) error {
	if h.codec0 == nil {
		// 使用默认编解码器
		h.codec0 = zcodec.Default()
	}
	buf, err := h.codec0.Encode(pkt)
	if err != nil {
		return err
	}
	return h.conn.Write(buf)
}

// 路由处理器，较baseHandler多实现了路由的函数HandleMsg0
type routerHandler struct {
	baseHandler
}

func (hr *routerHandler) HandleMsg0() error {
	if hr.conn == nil || !hr.conn.Alive() {
		logger.Error("connection not alive with baseHandler[%+v]", hr)
		return errors.New(fmt.Sprintf("connection[%+v] not alive with baseHandler[%+v]", hr.conn, hr))
	}
	pkt, err := hr.codec0.Decode(hr.conn)
	if err != nil {
		logger.Error("msg decode error,error:%+v", err)
		return err
	}
	logger.Debug("Receive msg[Head{id:%d,type:%d,length:%d}-Body{%s}] from address[%s:%d]",
		pkt.GetID(), pkt.GetType(), pkt.GetBodyLen(), string(pkt.GetData()), hr.conn.GetRemoteHost(), hr.conn.GetRemotePort())
	// todo renewal the connection
	if h, ok := childHandlerMap[pkt.GetID()]; ok {
		h.SetConn(hr.conn)
		return h.HandleMsg(pkt)
	}
	return errors.New(fmt.Sprintf("No childHandler to handle this msg[type:%d,id:%d]", pkt.GetType(), pkt.GetID()))
}

func (h *baseHandler) BindConn(conn connect.IConnection) {
	h.conn = conn
}

func (h *baseHandler) SetCodec(codec zcodec.Icodec) {
	h.codec0 = codec
}

func (b builder) Codec(c zcodec.Icodec) IBuilder {
	b.handler.SetCodec(c)
	return b
}

func (b builder) Build() IRouteHandler {
	return b.handler
}

func RouteBuilder() IBuilder {
	return builder{handler: &routerHandler{}}
}

func (b builder) Conn(conn connect.IConnection) IBuilder {
	b.handler.SetConn(conn)
	return b
}

func (h *baseHandler) SetConn(conn connect.IConnection) {
	h.conn = conn
}
