package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func BsonToJsonStr(bson interface{}, depth int) {
	Bj(bson, depth, false, true)
}

func Bj(bson interface{}, depth int, isKey bool, isLast bool) {

	//fmt.Printf("vvvvv %T", bson)
	str := JsonToken(bson, depth, isKey, isLast)
	if len(str) > 0 {
		fmt.Fprintf(os.Stdout, str)
	} else if m, ok := bson.(map[string]interface{}); ok {
		startBrace := TabStr{"{", depth, true}
		fmt.Fprintf(os.Stdout, startBrace.String())
		mapLen := len(m)
		count := 0
		for key := range m {
			count++
			last := false
			if count == mapLen {
				last = true
			}
			str := JsonToken(key, depth+1, true, true)
			fmt.Fprintf(os.Stdout, str)
			next := m[key]
			Bj(next, depth+1, false, last)
		}
		endBrace := TabStr{"}", depth, isLast}
		fmt.Fprintf(os.Stdout, endBrace.String())
	} else if s, ok := bson.([]interface{}); ok {
		startBracket := TabStr{"[", depth, true}
		fmt.Fprintf(os.Stdout, startBracket.String())
		mapLen := len(s)
		count := 0
		for key := range s {
			count++
			last := false
			if count == mapLen {
				last = true
			}
			next := s[key]
			Bj(next, depth+1, false, last)
		}
		endBracket := TabStr{"]", depth, isLast}
		fmt.Fprintf(os.Stdout, endBracket.String())
	}
}

func JsonToken(i interface{}, depth int, isKey bool, isLast bool) string {
	v, ok := i.(int32)
	str := ""
	needQuote := false
	if ok {
		str = strconv.FormatInt(int64(v), 10)
	} else if f, ok := i.(float64); ok {
		str = strconv.FormatFloat(f, 'g', -1, 64)
	} else if u6, ok := i.(uint64); ok {
		str = strconv.FormatUint(u6, 10)
	} else if i6, ok := i.(int64); ok {
		str = strconv.FormatInt(int64(i6), 10)
	} else if i8, ok := i.(int8); ok {
		str = string(i8)
	} else if s, ok := i.(string); ok {
		str = s
		needQuote = true
	} else if in, ok := i.(int); ok {
		str = strconv.Itoa(in)
	} else if b, ok := i.(bool); ok {
		if b {
			str = "true"
		} else {
			str = "false"
		}
	} else if t, ok := i.(time.Time); ok {
		str = t.Format(time.RFC3339)
		needQuote = true
	}

	if len(str) > 0 {
		js := JsonStr{&str, depth, isKey, isLast, needQuote}
		return js.String()
	} else {
		fmt.Sprintf("unknow type %T", i)
	}

	return str
}

type JsonStr struct {
	value     *string
	depth     int
	isKey     bool
	isLast    bool
	needQuote bool
}

type TabStr struct {
	value  string
	tabs   int
	isLast bool
}

func (ts *TabStr) String() string {
	temp := make([]byte, ts.tabs)
	for i := 0; i < len(temp); i++ {
		temp[i] = '\t'
	}

	str := string(temp) + ts.value
	if ts.isLast == false {
		str += ","
	}
	str += "\n"

	return str
}

func (js *JsonStr) String() string {
	tabs := ""
	if js.isKey {
		for i := 0; i < js.depth; i++ {
			tabs += "\t"
		}
		tabs += "\"" + *js.value + "\":"
	} else if js.needQuote {
		tabs += "\"" + *js.value + "\""
	} else {
		tabs = *js.value
	}
	if js.isKey == false {
		if js.isLast {
			tabs += "\n"
		} else {
			tabs += ",\n"
		}
	}
	return tabs
}
