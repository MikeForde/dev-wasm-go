package main

import (
	"log"
	"mime"
	"net/http"
)

func main() {
	// Ensure the correct MIME type is set for .wasm files
	mime.AddExtensionType(".wasm", "application/wasm")

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	log.Println("Server listening on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
