package couchdb

import (
	"context"
	"fmt"

	"github.com/go-kivik/kivik/v4"
)

// Find executes a Mango query
func (c *Client) Find(ctx context.Context, query interface{}) (*kivik.ResultSet, error) {
	rows := c.DB.Find(ctx, query)
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}
	return rows, nil
}

// // Query executes a query using a design document view
// func (c *Client) Query(ctx context.Context, ddoc, view string, options map[string]interface{}) (*kivik.Rows, error) {
// 	rows := c.DB.Query(ctx, ddoc, view, options)
// 	if err := rows.Err(); err != nil {
// 		return nil, fmt.Errorf("failed to query: %w", err)
// 	}
// 	return rows, nil
// }

// // QueryWithOptions executes a query with kivik.Options
// func (c *Client) QueryWithOptions(ctx context.Context, ddoc, view string, options kivik.Options) (*kivik.Rows, error) {
// 	rows := c.DB.Query(ctx, ddoc, view, options)
// 	if err := rows.Err(); err != nil {
// 		return nil, fmt.Errorf("failed to query: %w", err)
// 	}
// 	return rows, nil
// }

// // FindOne executes a Mango query and returns a single document
func (c *Client) FindOne(ctx context.Context, query interface{}, result interface{}) error {
	rows := c.DB.Find(ctx, query)
	if err := rows.Err(); err != nil {
		return fmt.Errorf("failed to find document: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return fmt.Errorf("document not found")
	}

	if err := rows.ScanDoc(result); err != nil {
		return fmt.Errorf("failed to scan document: %w", err)
	}

	return nil
}

// // Explain returns the query plan for a Mango query
// func (c *Client) Explain(ctx context.Context, query interface{}) (*kivik.QueryPlan, error) {
// 	plan, err := c.DB.Explain(ctx, query)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to explain query: %w", err)
// 	}
// 	return plan, nil
// }
