package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type Pet struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func InitDatabase() {
	var err error
	db, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Create pets table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS pets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL
		)
	`)
	if err != nil {
		log.Fatal("Failed to create pets table:", err)
	}

	// Insert sample data
	_, err = db.Exec(`INSERT INTO pets (name) VALUES ('Fluffy'), ('Buddy'), ('Max')`)
	if err != nil {
		log.Fatal("Failed to insert sample data:", err)
	}

	log.Println("Database initialized successfully")
}

func GetAllPets() []Pet {
	rows, err := db.Query("SELECT id, name FROM pets")
	if err != nil {
		log.Println("Error querying pets:", err)
		return []Pet{}
	}
	defer rows.Close()

	var pets []Pet
	for rows.Next() {
		var pet Pet
		if err := rows.Scan(&pet.ID, &pet.Name); err != nil {
			log.Println("Error scanning pet:", err)
			continue
		}
		pets = append(pets, pet)
	}

	return pets
}

func GetPetByID(id string) *Pet {
	query := fmt.Sprintf("SELECT id, name FROM pets WHERE id = %s", id)
	row := db.QueryRow(query)

	var pet Pet
	if err := row.Scan(&pet.ID, &pet.Name); err != nil {
		log.Println("Error getting pet by ID:", err)
		return nil
	}

	return &pet
}

func CreatePetByName(name string) int {
	query := fmt.Sprintf("INSERT INTO pets (name) VALUES ('%s')", name)
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
