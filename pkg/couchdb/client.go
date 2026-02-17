package couchdb

import (
	"context"
	"fmt"
	"jello-api/config"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/go-kivik/kivik/v4"
	_ "github.com/go-kivik/kivik/v4/couchdb"
)

// Client wraps both Kivik and direct HTTP access
type Client struct {
	// Kivik client and database
	Client *kivik.Client
	DB     *kivik.DB

	// Connection details for direct HTTP access
	BaseURL  string
	Username string
	Password string
	DBName   string
}

// Config holds the configuration for CouchDB client
type Config struct {
	URL    string
	DBName string
}

// NewClient creates a new CouchDB client from environment variables
func NewClient(ctx context.Context) (*Client, error) {
	cfg := Config{
		URL:    config.GetEnv("COUCHDB_URL", "http://admin:password@localhost:5984"),
		DBName: config.GetEnv("COUCHDB_NAME", "resto-app"),
	}
	return NewClientFromConfig(ctx, cfg)
}

// NewClientFromConfig creates a new CouchDB client from configuration with retry logic
func NewClientFromConfig(ctx context.Context, cfg Config) (*Client, error) {
	// Parse URL to extract credentials
	parsedURL, err := url.Parse(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid CouchDB URL: %w", err)
	}

	username := ""
	password := ""
	if parsedURL.User != nil {
		username = parsedURL.User.Username()
		password, _ = parsedURL.User.Password()
	}

	baseURL := fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)
	baseURL = strings.TrimSuffix(baseURL, "/")

	// Connect to CouchDB with retry logic
	var kivikClient *kivik.Client
	maxRetries := 5

	for i := 0; i < maxRetries; i++ {
		kivikClient, err = kivik.New("couch", cfg.URL)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to CouchDB (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect after %d attempts: %w", maxRetries, err)
	}

	// Verify connection
	version, err := kivikClient.Version(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get CouchDB version: %w", err)
	}
	log.Printf("✓ Connected to CouchDB version: %s", version)

	// Check if database exists
	exists, err := kivikClient.DBExists(ctx, cfg.DBName)
	if err != nil {
		return nil, fmt.Errorf("failed to check database existence: %w", err)
	}

	// Create database if it doesn't exist
	if !exists {
		if err := kivikClient.CreateDB(ctx, cfg.DBName); err != nil {
			return nil, fmt.Errorf("failed to create database: %w", err)
		}
		log.Printf("✓ Database '%s' created", cfg.DBName)
	} else {
		log.Printf("✓ Using existing database '%s'", cfg.DBName)
	}

	db := kivikClient.DB(cfg.DBName)

	return &Client{
		Client:   kivikClient,
		DB:       db,
		BaseURL:  baseURL,
		Username: username,
		Password: password,
		DBName:   cfg.DBName,
	}, nil
}

// Close closes the client connection
func (c *Client) Close() error {
	return c.Client.Close()
}
