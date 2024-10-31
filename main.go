package main

import (
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	// Ensure the correct MIME type is set for .wasm files
	mime.AddExtensionType(".wasm", "application/wasm")

	dir, err1 := os.Getwd()
	if err1 != nil {
		log.Fatal(err1)
	}
	log.Println("Current working directory:", dir)

	// Log the contents of the 'static' directory
	log.Println("Listing contents of the 'static' folder:")

	err2 := filepath.Walk("static", func(path string, info os.FileInfo, err2 error) error {
		if err2 != nil {
			return err2
		}
		log.Println(path)
		return nil
	})

	if err2 != nil {
		log.Printf("Error listing files in the 'static' folder: %v\n", err2)
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	log.Println(http.Dir("./static"))

	log.Println("Server listening on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
