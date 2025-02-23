//go:build js && wasm
// working version of the code

package main

import (
	"syscall/js"
)

func main() {
	js.Global().Set("hello", js.FuncOf(hello))
	select {} // Keep the main function running
}

func hello(this js.Value, args []js.Value) interface{} {
	js.Global().Get("alert").Invoke("Hello from TinyGo WASM!")
	return nil
}
