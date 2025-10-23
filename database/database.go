package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

type Pet struct {
	ID    int    `json:"pet_id"`
	Name  string `json:"name"`
	Owner string `json:"owner"`
}

func InitDatabase() {
	var err error

	// Get PostgreSQL connection string from environment variable
	// Format: postgres://username:password@hostname:port/database?sslmode=disable
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/zen_demo_go?sslmode=disable"
		log.Println("DATABASE_URL not set, using default:", dbURL)
	}

	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Test connection
	if err = db.Ping(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Create pets table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS pets (
			pet_id SERIAL PRIMARY KEY,
			pet_name TEXT NOT NULL,
			owner TEXT NOT NULL
		)
	`)
	if err != nil {
		log.Fatal("Failed to create pets table:", err)
	}

	// Insert sample data (only if table is empty)
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM pets").Scan(&count)
	if err != nil {
		log.Fatal("Failed to check pet count:", err)
	}

	if count == 0 {
		_, err = db.Exec(`INSERT INTO pets (pet_name, owner) VALUES ('Fluffy', 'Aikido Security'), ('Buddy', 'Aikido Security'), ('Max', 'Aikido Security')`)
		if err != nil {
			log.Fatal("Failed to insert sample data:", err)
		}
	}

	log.Println("Database initialized successfully")
}

func GetAllPets() []Pet {
	rows, err := db.Query("SELECT pet_id, pet_name, owner FROM pets")
	if err != nil {
		log.Println("Error querying pets:", err)
		return []Pet{}
	}
	defer rows.Close()

	var pets []Pet
	for rows.Next() {
		var pet Pet
		if err := rows.Scan(&pet.ID, &pet.Name, &pet.Owner); err != nil {
			log.Println("Error scanning pet:", err)
			continue
		}
		pets = append(pets, pet)
	}

	return pets
}

func GetPetByID(id string) *Pet {
	query := fmt.Sprintf("SELECT pet_id, pet_name, owner FROM pets WHERE pet_id = '%s'", id)
	row := db.QueryRow(query)

	var pet Pet
	if err := row.Scan(&pet.ID, &pet.Name, &pet.Owner); err != nil {
		log.Println("Error getting pet by ID:", err)
		return nil
	}

	return &pet
}

func CreatePetByName(name string) int {
	query := fmt.Sprintf("INSERT INTO pets (pet_name, owner) VALUES ('%s', 'Aikido Security')", name)
	result, err := db.Exec(query)
	if err != nil {
		log.Println("Error creating pet:", err)
		return -1
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("Error getting rows affected:", err)
		return -1
	}

	return int(rows)
}

func ClearAll() {
	_, err := db.Exec("DELETE FROM pets")
	if err != nil {
		log.Println("Error clearing pets:", err)
	}
}
