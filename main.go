package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"
)

type IpsAlt struct {
	ID                  int            `json:"id"`
	PackageUUID         string         `json:"packageUUID"`
	TimeStamp           string         `json:"timeStamp"`
	PatientName         string         `json:"patientName"`
	PatientGiven        string         `json:"patientGiven"`
	PatientDob          string         `json:"patientDob"`
	PatientGender       sql.NullString `json:"patientGender"`
	PatientNation       string         `json:"patientNation"`
	PatientPractitioner string         `json:"patientPractitioner"`
	PatientOrganization sql.NullString `json:"patientOrganization"`
	CreatedAt           string         `json:"createdAt"`
	UpdatedAt           string         `json:"updatedAt"`
}

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found or could not be loaded")
	}
}

func getDBConnection() (*sql.DB, error) {
	// Read database connection details from environment variables
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT") // Read port from environment variable, or use a default value
	if port == "" {
		port = "3306" // Default port if not specified
	}
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbName)

	fmt.Println("Connecting to DB with DSN:", dsn)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func getIpsAltHandler(w http.ResponseWriter, r *http.Request) {
	db, err := getDBConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM ipsAlt")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []IpsAlt
	for rows.Next() {
		var ips IpsAlt
		err := rows.Scan(&ips.ID, &ips.PackageUUID, &ips.TimeStamp, &ips.PatientName,
			&ips.PatientGiven, &ips.PatientDob, &ips.PatientGender, &ips.PatientNation,
			&ips.PatientPractitioner, &ips.PatientOrganization, &ips.CreatedAt, &ips.UpdatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		results = append(results, ips)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func main() {
	// Serve static files from the 'static' directory
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	// API endpoint for fetching IPS Alt records
	http.HandleFunc("/ipsAlt", getIpsAltHandler)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
