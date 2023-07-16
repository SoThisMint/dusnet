package server

import (
	zcodec "dusnet/codec"
	"dusnet/connect"
	"dusnet/handler"
	"dusnet/logger"
	"fmt"
	"net"
	"reflect"
)

type Option func()

// IServer server端抽象
type IServer interface {
	Start(opts ...Option) error // 启动
	Stop() error                // 停止
}

type mServer struct {
	name         string                 // server名称
	network      string                 // 网络
	host         string                 // host
	port         int                    // 端口
	routeHandler handler.IRouteHandler  // 基础处理器
	connMgr      connect.IConnectionMgr // 连接管理器
}

// Default 返回默认的server实现
func Default(name string, network string, host string, port int) IServer {
	return mServer{
		name:         name,
		network:      network,
		host:         host,
		port:         port,
		routeHandler: handler.RouteBuilder().Codec(zcodec.Default()).Build(), //使用默认内置路由handler
		connMgr:      connect.DefaultConnMgr(),
	}
}

// New 自定义server实现，可以启用自己的handler及编解码器、连接管理器
func New(name string, network string, host string, port int, routeHandler handler.IRouteHandler, connMgr connect.IConnectionMgr) IServer {
	return mServer{
		name:         name,
		network:      network,
		host:         host,
		port:         port,
		routeHandler: routeHandler, //使用自定义路由handler及编解码器
		connMgr:      connMgr,      //使用自定义连接管理器
	}
}

func (m mServer) Start(opts ...Option) error {
	printServerEnv(m)
	addr, err := net.ResolveTCPAddr(m.network, fmt.Sprintf("%s:%d", m.host, m.port))
	if err != nil {
		logger.Error("net.ResolveTCPAddr error,error:%+v", err)
		return err
	}
	listener, err := net.ListenTCP(m.network, addr)
	if err != nil {
		logger.Error("net.ListenTCP error,error:%+v", err)
		return err
	}
	logger.Info("server[%s] started on %s:%d", m.name, m.host, m.port)
	logger.Debug("")
	for _, opt := range opts {
		opt()
	}
	go func() {
		for {
			conn := connect.New(listener, m.connMgr)
			go func() {
				for {
					m.routeHandler.BindConn(conn)
					err := m.routeHandler.HandleMsg0()
					if err != nil {
						logger.Error("routeHandler.HandleMsg0() error,error:%+v", err)
						if conn != nil {
							err := m.connMgr.RemoveConnByID(conn.GetID())
							if err != nil {
								logger.Error("remove conn error,error:%+v", err)
							}
							logger.Warn("One connection[id=%d,laddr:%s:%d,raddr:%s:%d] released",
								conn.GetID(), conn.GetLocalHost(), conn.GetLocalPort(), conn.GetRemoteHost(), conn.GetRemotePort())
						}
						return
					}
				}
			}()
		}
	}()
	return err
}

// 打印服务端配置信息
func printServerEnv(m mServer) {
	logger.Debug("=========================================[Server-Info]=========================================")
	logger.Debug("")
	logger.Debug("[Name]:[%s]", m.name)
	logger.Debug("[RouterHandler]:[%+v]", reflect.TypeOf(m.routeHandler))
	logger.Debug("")
	logger.Debug("[ChildHandlers]")
	for id, h := range handler.AllChildHandlers() {
		logger.Debug("[%d]:[%+v]", id, reflect.TypeOf(h))
	}
	logger.Debug("")
	logger.Debug("=========================================[Server-Info]=========================================")

}

func (m mServer) Stop() error {
	// close all connections for now
	all := m.connMgr.All()
	for _, conn := range all {
		if conn != nil && conn.Alive() {
			err := m.connMgr.RemoveConnByID(conn.GetID())
			if err != nil {
				return err
			}
		}
	}
	return nil
}
