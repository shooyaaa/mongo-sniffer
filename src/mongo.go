package main

type MsgHeader struct {
	Len        int32
	Id         int32
	ResponseTo int32
	OpCode     int32
}

type Document struct {
}
type OpUpdate struct {
	header         MsgHeader
	zero           int32
	collectionName string
	flags          int32
	selector       Document
	update         Document
}

func OpType(opcode int32) string {
	opType := "ERROR_TYPE"
	switch opcode {
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
