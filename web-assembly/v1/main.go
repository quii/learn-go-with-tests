// +build js,wasm

package main

import (
	"fmt"
	"syscall/js"
)

var (
	global    = js.Global()
	null      = js.Null()
	undefined = js.Undefined()
)

func main() {
	fmt.Println("Hello World")
}
