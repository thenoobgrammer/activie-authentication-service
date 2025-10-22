package database

import (
	"database/sql"
	"log"
	"sync"

	_ "github.com/lib/pq"
)

var (
	client *sql.DB
	once   sync.Once
)

func InitializeDB(dsn string) {
	log.Println("DSN", dsn)
	once.Do(func() {
		var err error
		client, err = sql.Open("postgres", dsn)
		if err != nil {
			log.Fatalf("Error opening database: %v", err)
		}

		if err = client.Ping(); err != nil {
			log.Fatalf("Error pinging database: %v", err)
		} else {
			log.Println("Service connected to Database.")
		}
	})
}

func GetClient() *sql.DB {
	if client == nil {
		log.Fatal("Database has not been initialized. Call database.Initialize first.")
	}
	return client
}

func Close() {
	client.Close()
}
