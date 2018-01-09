package main

import (
	"fmt"
	"os"
	"reflect"
)

func BsonToJsonStr(bson interface{}, depth int) {
	rf := reflect.ValueOf(bson)
	Bj(rf, depth)
}

func Bj(v reflect.Value, depth int) {

	fmt.Printf("vvv %T\n", v.Kind())
	switch v.Kind() {
	case reflect.String:
		str := string(v.String())
		js := JsonStr{&str, depth + 1, false}
		fmt.Fprintf(os.Stdout, js.String())
	case reflect.Uint:
		ji := JsonInt{v.Uint(), depth + 1, false}
		fmt.Fprintf(os.Stdout, ji.String())
	case reflect.Map:
		mapKeys := v.MapKeys()
		fmt.Fprintf(os.Stdout, "{")
		for _, vv := range mapKeys {
			str := string(vv.String())
			js := JsonStr{&str, depth + 1, true}
			fmt.Fprintf(os.Stdout, js.String())
			token := v.MapIndex(vv)
			fmt.Printf("token %T", token)
			Bj(token, depth+1)
		}
		fmt.Fprintf(os.Stdout, "}")
	case reflect.Interface:
		//fmt.Printf("interface %t, %v", v)
	}
}

type JsonInt struct {
	value uint64
	depth int
	isKey bool
}

func (ji *JsonInt) String() string {
	tabs := ""
	for i := 0; i < ji.depth; i++ {
		tabs += "\t"
	}
	tabs += string(ji.value)
	if ji.isKey == false {
		tabs += "\n"
	}
	return tabs
}

type JsonStr struct {
	value *string
	depth int
	isKey bool
}

func (js *JsonStr) String() string {
	tabs := ""
	for i := 0; i < js.depth; i++ {
		tabs += "\t"
	}
	tabs += "\"" + *js.value + "\":"
	if js.isKey == false {
		tabs += "\n"
	}
	return tabs
}
