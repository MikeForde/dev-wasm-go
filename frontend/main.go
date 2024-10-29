package main

import (
	"fmt"
)

func main() {
	// Optional: Initialization code
	fmt.Println("WASM TinyGo Initialized")
}

//export hello
func hello() {
	fmt.Println("Hello from TinyGo WASM!")
}
