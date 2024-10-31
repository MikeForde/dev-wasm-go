package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	// Get current working directory to debug path issues
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Current working directory:", dir)

	// Log the contents of the 'static' folder using an absolute path
	staticPath := filepath.Join(dir, "static")
	log.Println("Listing contents of the 'static' folder at:", staticPath)

	err = filepath.Walk(staticPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		log.Println(path)
		return nil
	})

	if err != nil {
		log.Printf("Error listing files in the 'static' folder: %v\n", err)
	}

	// Start the server using the absolute path to 'static'
	fs := http.FileServer(http.Dir(staticPath))
	http.Handle("/", fs)

	log.Println("Server listening on :8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
