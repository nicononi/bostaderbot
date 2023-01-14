package database

import (
	"database/sql"
	"log"
	"os"
	"sync"
	"time"
)

var (
	db   *sql.DB
	once sync.Once
)

// getDB lazily instantiates a database connection pool. Users of Cloud Run or
// Cloud Functions may wish to skip this lazy instantiation and connect as soon
// as the function is loaded. This is primarily to help testing.
func GetDB() *sql.DB {
	once.Do(func() {
		db = mustConnect()
	})
	return db
}

// mustConnect creates a connection to the database based on environment
// variables. Setting one of INSTANCE_HOST, INSTANCE_UNIX_SOCKET, or
// INSTANCE_CONNECTION_NAME will establish a connection using a TCP socket, a
// Unix socket, or a connector respectively.
func mustConnect() *sql.DB {
	var (
		db  *sql.DB
		err error
	)

	// Use a TCP socket when INSTANCE_HOST (e.g., 127.0.0.1) is defined
	if os.Getenv("INSTANCE_HOST") != "" {
		db, err = connectTCPSocket()
		if err != nil {
			log.Fatalf("connectTCPSocket: unable to connect: %s", err)
		}
	}

	// Use the connector when INSTANCE_CONNECTION_NAME (proj:region:instance) is defined.
	if os.Getenv("INSTANCE_CONNECTION_NAME") != "" {
		db, err = connectWithConnector()
		if err != nil {
			log.Fatalf("connectConnector: unable to connect: %s", err)
		}
	}

	if db == nil {
		log.Fatal("Missing database connection type. Please define one of INSTANCE_HOST, INSTANCE_UNIX_SOCKET, or INSTANCE_CONNECTION_NAME")
	}

	if err := migrateDB(db); err != nil {
		log.Fatalf("unable to create table: %s", err)
	}

	return db
}

// configureConnectionPool sets database connection pool properties.
// For more information, see https://golang.org/pkg/database/sql
func configureConnectionPool(db *sql.DB) {
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(7)
	db.SetConnMaxLifetime(1800 * time.Second)
}

// migrateDB creates the visited and registered tables if they do not already exist.
func migrateDB(db *sql.DB) error {
	createVisited := `CREATE TABLE IF NOT EXISTS visited (
		chat_id NUMERIC NOT NULL,
		rooms NUMERIC NOT NULL,
		apartment_id NUMERIC NOT NULL,
		PRIMARY KEY (chat_id, rooms, apartment_id)
	);`
	_, err := db.Exec(createVisited)

	if err != nil {
		return err
	}

	createRegistered := `CREATE TABLE IF NOT EXISTS registered (
		chat_id NUMERIC UNIQUE NOT NULL,
		PRIMARY KEY (chat_id)
	);`
	_, err2 := db.Exec(createRegistered)
	return err2
}
