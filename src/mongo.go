package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

type MongoOp struct {
	Payload *bytes.Reader
	Ip      IpPair
	Port    PortPair
	Header  MsgHeader
}

//check if payload is valid protocal
func (o *MongoOp) validate() bool {
	return true
}

func (o *MongoOp) decode() {
	//fmt.Printf("all payload %v\n", o.Payload)
	o.ReadHeader(&o.Header)
	switch o.Header.OpCode {
	case 2004:
		decodeOpQuery(o)
	case 2013:
		decodeOpMsg(o)
	default:
		fmt.Printf("unhandled type %v", o.Header.OpCode)
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
		fmt.Printf("error read byte from MongoOp: %s ", err)
		panic("111")
	}
	return b
}

func (m *MongoOp) ReadBytes(n uint32) []byte {
	bts := []byte{}

	for i := uint32(0); i < n; i++ {
		temp := m.ReadByte()
		bts = append(bts, temp)
		//fmt.Printf("%v, %q --- %v\n", i, temp, n)
	}
	return bts
}

func (m *MongoOp) ReadCstr() string {
	count := 0
	b := byte(1)
	bts := make([]byte, 1)

	for b != 0 {
		err := binary.Read(m.Payload, binary.LittleEndian, &b)
		if err != nil {
			panic("error read byte from payload")
		}
		bts = append(bts, b)
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
	kind := m.ReadByte()

	fmt.Printf("op_msg[%v](%q:%q)----(%q:%v)", m.Header.Len, m.Ip.SrcIp, m.Port.SrcPort, m.Ip.DstIp, m.Port.DstPort)
	if kind == 0 {
		sectionLen := m.Payload.Len()
		fmt.Printf("expect bson len %v\n", sectionLen)
		rawDoc := m.ReadBytes(uint32(sectionLen))
		s := rawDoc[:4]
		rawBson := DecodeBson(rawDoc)
		fmt.Printf("raw doc %v[%v]\n", len(rawDoc), s)
		BsonToJsonStr(rawBson, 0)
	}
}

func decodeOpQuery(m *MongoOp) {
	q := OpQuery{}
	q.Flags = m.ReadInt32()
	q.CollectionName = m.ReadCstr()
	q.NumberToSkip = m.ReadInt32()
	q.NumberToReturn = m.ReadInt32()
	fmt.Printf("op_query[%v]--->%v\n", m.Header.Len, q.CollectionName)
	sectionLen := m.Payload.Len()
	if sectionLen != 0 {
		rawDoc := m.ReadBytes(uint32(sectionLen))
		rawBson := DecodeBson(rawDoc)
		BsonToJsonStr(rawBson, 0)
	}
}

func DecodeBson(b []byte) map[string]interface{} {
	//fmt.Printf("doc  %v", b)
	var rawBson map[string]interface{}
	err := bson.Unmarshal(b, &rawBson)
	if err != nil {
		fmt.Printf("unmarshal error %q\n", err)
	}
	return rawBson
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
