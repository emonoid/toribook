package db

import ( 
	"database/sql" 
)

// Store provides all functions to execute SQL queries and transactions.
// It wraps the generated Queries type and provides a database connection.
type Store struct {
	*Queries         // embedded Queries pointer as like inheritance
	db       *sql.DB // same
}

// NewStore creates a new Store
func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db), // initialize Queries with the provided db connection
		db:      db,      // store the db connection
	}
}