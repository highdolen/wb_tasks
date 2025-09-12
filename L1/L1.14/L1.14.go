package main

import (
	"fmt"
	"reflect"
)

func detect(value any) string {

	if value == nil {
		return "nil"
	}

	switch value.(type) {
	case int:
		return "int"
	case string:
		return "string"
	case bool:
		return "bool"
	default:
		if reflect.ValueOf(value).Kind() == reflect.Chan {
			return "chan"
		} else {
			return "no name type"
		}
	}

}

func main() {
	fmt.Println(detect(52))
	fmt.Println(detect("word"))
	fmt.Println(detect(false))
	fmt.Println(detect(make(chan struct{})))
}
