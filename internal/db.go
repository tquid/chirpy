package database

import (
	"sync"

	m "codeberg.org/tquid/chirpy/models"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]m.Chirp `json:"chirps"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {

}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (m.Chirp, error)

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]m.Chirp, error)

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error)

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error
