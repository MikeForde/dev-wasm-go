package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// Struct to hold data from the ipsAlt table
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

// Connect to MySQL and fetch data
func getIpsAltHandler(w http.ResponseWriter, r *http.Request) {
	// Replace with your MySQL connection details
	db, err := sql.Open("mysql", "root:password@tcp(host.docker.internal:3306)/test")
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
