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

    Medications   []Medication   `json:"medications"`
    Allergies     []Allergy      `json:"allergies"`
    Conditions    []Condition    `json:"conditions"`
    Observations  []Observation  `json:"observations"`
    Immunizations []Immunization `json:"immunizations"`
}

type Medication struct {
    ID         int    `json:"id"`
    IPSModelID int    `json:"ipsModelId"`
    Name       string `json:"name"`
    Date       string `json:"date"`
    Dosage     string `json:"dosage"`
}

type Allergy struct {
    ID          int    `json:"id"`
    IPSModelID  int    `json:"ipsModelId"`
    Name        string `json:"name"`
    Criticality string `json:"criticality"`
    Date        string `json:"date"`
}

type Condition struct {
    ID         int    `json:"id"`
    IPSModelID int    `json:"ipsModelId"`
    Name       string `json:"name"`
    Date       string `json:"date"`
}

type Observation struct {
    ID         int    `json:"id"`
    IPSModelID int    `json:"ipsModelId"`
    Name       string `json:"name"`
    Date       string `json:"date"`
    Value      string `json:"value"`
}

type Immunization struct {
    ID         int    `json:"id"`
    IPSModelID int    `json:"ipsModelId"`
    Name       string `json:"name"`
    System     string `json:"system"`
    Date       string `json:"date"`
}


func init() {
	// Attempt to load the .env file, but continue without it if not found
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Println("Error loading .env file, proceeding with system environment variables:", err)
		} else {
			log.Println(".env file loaded successfully")
		}
	} else {
		log.Println("No .env file found, using system environment variables")
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

    // Fetch the main IPS Alt records
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

        // Fetch related data for each IPS Alt record
        ips.Medications, err = fetchMedications(db, fmt.Sprintf("%d", ips.ID))
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        ips.Allergies, err = fetchAllergies(db, fmt.Sprintf("%d", ips.ID))
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        ips.Conditions, err = fetchConditions(db, fmt.Sprintf("%d", ips.ID))
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        ips.Observations, err = fetchObservations(db, fmt.Sprintf("%d", ips.ID))
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        ips.Immunizations, err = fetchImmunizations(db, fmt.Sprintf("%d", ips.ID))
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        results = append(results, ips)
    }

    // Return the complete set of data including nested arrays
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(results)
}


func fetchMedications(db *sql.DB, ipsModelID string) ([]Medication, error) {
    query := `SELECT id, IPSModelId, name, date, dosage FROM Medications WHERE IPSModelId = ?`
    rows, err := db.Query(query, ipsModelID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var medications []Medication
    for rows.Next() {
        var med Medication
        err := rows.Scan(&med.ID, &med.IPSModelID, &med.Name, &med.Date, &med.Dosage)
        if err != nil {
            return nil, err
        }
        medications = append(medications, med)
    }
    return medications, nil
}

// Similarly create fetchAllergies, fetchConditions, fetchObservations, fetchImmunizations
func fetchAllergies(db *sql.DB, ipsModelID string) ([]Allergy, error) {
    query := `SELECT id, IPSModelId, name, criticality, date FROM Allergies WHERE IPSModelId = ?`
    rows, err := db.Query(query, ipsModelID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var allergies []Allergy
    for rows.Next() {
        var allergy Allergy
        err := rows.Scan(&allergy.ID, &allergy.IPSModelID, &allergy.Name, &allergy.Criticality, &allergy.Date)
        if err != nil {
            return nil, err
        }
        allergies = append(allergies, allergy)
    }
    return allergies, nil
}

func fetchConditions(db *sql.DB, ipsModelID string) ([]Condition, error) {
    query := `SELECT id, IPSModelId, name, date FROM Conditions WHERE IPSModelId = ?`
    rows, err := db.Query(query, ipsModelID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var conditions []Condition
    for rows.Next() {
        var condition Condition
        err := rows.Scan(&condition.ID, &condition.IPSModelID, &condition.Name, &condition.Date)
        if err != nil {
            return nil, err
        }
        conditions = append(conditions, condition)
    }
    return conditions, nil
}

func fetchObservations(db *sql.DB, ipsModelID string) ([]Observation, error) {
    query := `SELECT id, IPSModelId, name, date, value FROM Observations WHERE IPSModelId = ?`
    rows, err := db.Query(query, ipsModelID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var observations []Observation
    for rows.Next() {
        var observation Observation
        err := rows.Scan(&observation.ID, &observation.IPSModelID, &observation.Name, &observation.Date, &observation.Value)
        if err != nil {
            return nil, err
        }
        observations = append(observations, observation)
    }
    return observations, nil
}

func fetchImmunizations(db *sql.DB, ipsModelID string) ([]Immunization, error) {
    query := "SELECT id, IPSModelId, name, `system`, date FROM Immunizations WHERE IPSModelId = ?"
    rows, err := db.Query(query, ipsModelID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var immunizations []Immunization
    for rows.Next() {
        var immunization Immunization
        err := rows.Scan(&immunization.ID, &immunization.IPSModelID, &immunization.Name, &immunization.System, &immunization.Date)
        if err != nil {
            return nil, err
        }
        immunizations = append(immunizations, immunization)
    }
    return immunizations, nil
}



func getIpsAltByIDHandler(w http.ResponseWriter, r *http.Request) {
    // Get the ID from query parameters
    id := r.URL.Query().Get("id")
    if id == "" {
        http.Error(w, "Missing 'id' query parameter", http.StatusBadRequest)
        return
    }

    db, err := getDBConnection()
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    var ips IpsAlt
    // Fetch the main IpsAlt record
    query := `SELECT * FROM ipsAlt WHERE id = ?`
    err = db.QueryRow(query, id).Scan(
        &ips.ID, &ips.PackageUUID, &ips.TimeStamp, &ips.PatientName,
        &ips.PatientGiven, &ips.PatientDob, &ips.PatientGender, &ips.PatientNation,
        &ips.PatientPractitioner, &ips.PatientOrganization, &ips.CreatedAt, &ips.UpdatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "Record not found", http.StatusNotFound)
        } else {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    // Fetch related data
    ips.Medications, err = fetchMedications(db, id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    ips.Allergies, err = fetchAllergies(db, id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    ips.Conditions, err = fetchConditions(db, id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    ips.Observations, err = fetchObservations(db, id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    ips.Immunizations, err = fetchImmunizations(db, id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Return the assembled record
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(ips)
}


func main() {
	// Serve static files from the 'static' directory
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	// API endpoint for fetching IPS Alt records
	http.HandleFunc("/ipsAlt", getIpsAltHandler)
	http.HandleFunc("/ipsAltByID", getIpsAltByIDHandler)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
