package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type MongoOp struct {
	Payload *bytes.Reader
	Ip      IpPair
	Port    PortPair
}

//check if payload is valid protocal
func (o *MongoOp) validate() bool {
	return true
}

func (o *MongoOp) decode() {
	header := MsgHeader{}
	o.ReadHeader(&header)
	switch header.OpCode {
	case 2004:
		decodeOpQuery(o)
	case 2013:
		decodeOpMsg(o)
	}
}

func (o *MongoOp) ReadHeader(h *MsgHeader) error {
	err := binary.Read(o.Payload, binary.LittleEndian, h)
	return err
}

type MsgHeader struct {
	Len        int32
	Id         int32
	ResponseTo int32
	OpCode     int32
}

type OpQuery struct {
	Flags          int32
	CollectionName string
	NumberToSkip   int32
	NumberToReturn int32
}

type OpMsg struct {
	Flags    uint32
	sections interface{}
}

type SectionBody struct {
	test int32
}

type SectionSeq struct {
	size int32
	seq  string
}

func (m *MongoOp) ReadByte() byte {
	b := byte(1)
	err := binary.Read(m.Payload, binary.LittleEndian, &b)
	if err != nil {
		panic("error read byte from MongoOp")
	}
	return b
}

func (m *MongoOp) ReadCstr() string {
	count := 0
	b := byte(1)
	bts := [1000]byte{}

	for b != 0 {
		err := binary.Read(m.Payload, binary.LittleEndian, &b)
		if err != nil {
			fmt.Println(err)
			panic("error read byte from payload")
		}
		bts[count] = b
		count++
	}

	return string(bts[:count])

}

func (m *MongoOp) ReadInt32() int32 {
	i := int32(0)
	binary.Read(m.Payload, binary.LittleEndian, &i)
	return i
}

func (m *MongoOp) ReadUint32() uint32 {
	return uint32(m.ReadInt32())
}

func decodeOpMsg(m *MongoOp) {
	msg := OpMsg{}
	msg.Flags = m.ReadUint32()
	kind = m.ReadByte()
	fmt.Printf("op_msg(%q:%q)----(%q:%q), flags:%q\n", m.Ip.SrcIp, m.Port.SrcPort, m.Ip.DstIp, m.Port.DstPort, msg.Flags)
}

func decodeOpQuery(m *MongoOp) {
	q := OpQuery{}
	q.Flags = m.ReadInt32()
	q.CollectionName = m.ReadCstr()
	q.NumberToSkip = m.ReadInt32()
	q.NumberToReturn = m.ReadInt32()
	fmt.Printf("op_query--->%s\n", q.CollectionName)
}

func (m MsgHeader) OpType() string {
	opType := "ERROR_TYPE"
	switch m.OpCode {
	case 1:
		opType = "OP_REPLY"
	case 2001:
		opType = "OP_UPDATE"
	case 2002:
		opType = "OP_INSERT"
	case 2003:
		opType = "RESERVED"
	case 2004:
		opType = "OP_QUERY"
	case 2005:
		opType = "OP_GET_MORE"
	case 2006:
		opType = "OP_DELETE"
	case 2007:
		opType = "OP_KILL_CURSORS"
	case 2010:
		opType = "OP_COMMAND"
	case 2011:
		opType = "OP_COMMANDREPLY"
	case 2013:
		opType = "OP_MSG"

	}
	return opType
}
