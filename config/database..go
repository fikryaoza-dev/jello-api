package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	kivik "github.com/go-kivik/kivik/v4"
	_ "github.com/go-kivik/kivik/v4/couchdb"
)

type Database struct {
	Client *kivik.Client
	DB     *kivik.DB
}

func ConnectDB() *Database {
	couchURL := getEnv("COUCHDB_URL", "http://admin:password@localhost:5984/")
	dbName := getEnv("COUCHDB_NAME", "resto-app")

	var client *kivik.Client
	var err error

	maxRetries := 5

	for i := 0; i < maxRetries; i++ {
		client, err = kivik.New("couch", couchURL)
		if err == nil {
			break
		}

		log.Printf("Failed to connect to CouchDB (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("Failed to connect to CouchDB after %d attempts: %v", maxRetries, err)
	}

	ctx := context.Background()
	version, err := client.Version(ctx)
	if err != nil {
		log.Fatalf("Failed to get CouchDB version: %v", err)
	}
	log.Printf("✓ Connected to CouchDB version: %s", version)

	// Create database if not exists
	exists, err := client.DBExists(ctx, dbName)
	if err != nil {
		log.Fatalf("Failed to check database existence: %v", err)
	}

	if !exists {
		if err := client.CreateDB(ctx, dbName); err != nil {
			log.Fatalf("Failed to create database: %v", err)
		}
		log.Printf("✓ Database '%s' created", dbName)
	} else {
		log.Printf("✓ Using existing database '%s'", dbName)
	}

	db := client.DB(dbName)

	return &Database{
		Client: client,
		DB:     db,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetPort() string {
	port := getEnv("APP_PORT", "3000")
	return fmt.Sprintf(":%s", port)
}

func GetEnv() string {
	return getEnv("APP_ENV", "development")
}
