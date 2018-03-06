package main

import (
	"fmt"
	"reflect"
)

func main() {
	var a int

	b := make(map[reflect.Type]string)
	b[reflect.TypeOf(a)] = "Hello"

	fmt.Println(b)

	str := b[]
	fmt.Println(str)
}
