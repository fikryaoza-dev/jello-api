package couchdb

import (
	"context"
	"fmt"
	"time"
)

// CreateDoc creates a new document using Kivik
func (c *Client) CreateDoc(ctx context.Context, doc interface{}) (string, string, error) {
	id, rev, err := c.DB.CreateDoc(ctx, doc)
	if err != nil {
		return "", "", fmt.Errorf("failed to create document: %w", err)
	}
	return id, rev, nil
}

// CreateDocWithID creates a new document with a specific ID
func (c *Client) CreateDocWithID(ctx context.Context, id string, doc interface{}) (string, error) {
	rev, err := c.DB.Put(ctx, id, doc)
	if err != nil {
		return "", fmt.Errorf("failed to create document: %w", err)
	}
	return rev, nil
}

// GetDoc retrieves a document by ID
func (c *Client) GetDoc(ctx context.Context, id string, result interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	row := c.DB.Get(ctx, id)
	if err := row.Err; err != nil {
		return fmt.Errorf("failed to get document: %w", err)
	}

	if err := row.ScanDoc(result); err != nil {
		return fmt.Errorf("failed to scan document: %w", err)
	}

	return nil
}

// UpdateDoc updates an existing document
func (c *Client) UpdateDoc(ctx context.Context, id string, doc interface{}) (string, error) {
	rev, err := c.DB.Put(ctx, id, doc)
	if err != nil {
		return "", fmt.Errorf("failed to update document: %w", err)
	}
	return rev, nil
}

// // UpdateDocWithRev updates an existing document with explicit revision
// func (c *Client) UpdateDocWithRev(ctx context.Context, id, rev string, doc interface{}) (string, error) {
// 	newRev, err := c.DB.Put(ctx, id, doc, kivik.Rev(rev))
// 	if err != nil {
// 		return "", fmt.Errorf("failed to update document: %w", err)
// 	}
// 	return newRev, nil
// }

// // DeleteDoc deletes a document
// func (c *Client) DeleteDoc(ctx context.Context, id, rev string) error {
// 	_, err := c.DB.Delete(ctx, id, rev)
// 	if err != nil {
// 		return fmt.Errorf("failed to delete document: %w", err)
// 	}
// 	return nil
// }

// // ListDocs lists all documents in the database
// func (c *Client) ListDocs(ctx context.Context) (*kivik.Rows, error) {
// 	rows := c.DB.AllDocs(ctx, kivik.Options{
// 		"include_docs": true,
// 	})
// 	if err := rows.Err(); err != nil {
// 		return nil, fmt.Errorf("failed to list documents: %w", err)
// 	}
// 	return rows, nil
// }

// // DocExists checks if a document exists
// func (c *Client) DocExists(ctx context.Context, id string) (bool, error) {
// 	row := c.DB.Get(ctx, id)
// 	if row.Err != nil {
// 		// Check if it's a not found error
// 		if kivik.StatusCode(row.Err) == 404 {
// 			return false, nil
// 		}
// 		return false, fmt.Errorf("failed to check document existence: %w", row.Err)
// 	}
// 	return true, nil
// }

// // GetDocRev retrieves only the revision of a document
// func (c *Client) GetDocRev(ctx context.Context, id string) (string, error) {
// 	row := c.DB.Get(ctx, id)
// 	if err := row.Err; err != nil {
// 		return "", fmt.Errorf("failed to get document: %w", err)
// 	}
// 	return row.Rev, nil
// }
