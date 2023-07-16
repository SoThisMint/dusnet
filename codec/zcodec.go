package zcodec

import (
	"bytes"
	"dusnet/connect"
	"dusnet/logger"
	"dusnet/packet"
	"encoding/binary"
	"errors"
	"strconv"
)

const (
	TYPE_PING     = iota + 1234 // PING报文
	TYPE_SYNC                   // 同步报文
	TYPE_BUSINESS               // 业务报文 业务报文永远在队尾
)

type Icodec interface {
	Encode(packet.IPacket) ([]byte, error)
	Decode(connect.IConnection) (packet.IPacket, error)
}

func Default() Icodec {
	return &codec{}
}

type codec struct {
}

func (c *codec) Encode(p packet.IPacket) ([]byte, error) {
	buffer := bytes.NewBuffer(make([]byte, 0))
	id := p.GetID()
	err := binary.Write(buffer, binary.BigEndian, &id)
	if err != nil {
		logger.Error("binary.Write head/id to bytes error,error:%+v", err)
		return buffer.Bytes(), err
	}
	type0 := p.GetType()
	err = binary.Write(buffer, binary.BigEndian, &type0)
	if err != nil {
		logger.Error("binary.Write head/type to bytes error,error:%+v", err)
		return buffer.Bytes(), err
	}
	bodyLen := p.GetBodyLen()
	err = binary.Write(buffer, binary.BigEndian, &bodyLen)
	if err != nil {
		logger.Error("binary.Write head/length to bytes error,error:%+v", err)
		return buffer.Bytes(), err
	}

	err = binary.Write(buffer, binary.BigEndian, p.GetData())
	if err != nil {
		logger.Error("binary.Write body/data to bytes error,error:%+v", err)
		return buffer.Bytes(), err
	}
	return buffer.Bytes(), nil
}

func (c *codec) Decode(conn connect.IConnection) (packet.IPacket, error) {
	pkt := &packet.Packet{}
	// decode head/id
	idBuf := make([]byte, 4)
	err := conn.Read(idBuf)
	if err != nil {
		logger.Error("conn.Read head/id error,error:%+v", err)
		return pkt, err
	}
	err = binary.Read(bytes.NewBuffer(idBuf), binary.BigEndian, &pkt.ID)
	if err != nil {
		logger.Error("binary.Read head/id to pkt error,error:%+v", err)
		return pkt, err
	}
	if pkt.ID == 0 {
		return pkt, errors.New("id zero value")
	}

	typeBuf := make([]byte, 2)
	err = conn.Read(typeBuf)
	if err != nil {
		logger.Error("conn.Read head/type error,error:%+v", err)
		return pkt, err
	}
	err = binary.Read(bytes.NewBuffer(typeBuf), binary.BigEndian, &pkt.Type)
	if err != nil {
		logger.Error("binary.Read head/type to pkt error,error:%+v", err)
		return pkt, err
	}
	if pkt.Type < TYPE_PING || pkt.Type > TYPE_BUSINESS {
		return pkt, errors.New("pkt type " + strconv.Itoa(int(pkt.Type)) + " not defined")
	}

	// decode head/length
	lengthBuf := make([]byte, 4)
	err = conn.Read(lengthBuf)
	if err != nil {
		logger.Error("conn.Read head/length error,error:%+v", err)
		return pkt, err
	}
	err = binary.Read(bytes.NewBuffer(lengthBuf), binary.BigEndian, &pkt.Length)
	if err != nil {
		logger.Error("binary.Write head/length to pkt error,error:%+v", err)
		return pkt, err
	}

	// decode head body/data
	dataBuf := make([]byte, pkt.Length)
	err = conn.Read(dataBuf)
	if err != nil {
		logger.Error("conn.Read body/data error,error:%+v", err)
		return pkt, err
	}
	pkt.Data = dataBuf
	return pkt, nil
}
