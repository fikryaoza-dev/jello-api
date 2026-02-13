package couchdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	BaseURL  string
	Username string
	Password string
	DBName   string
}

func (c *Client) CreateDoc(doc interface{}) error {
	url := fmt.Sprintf("%s/%s", c.BaseURL, c.DBName)

	body, _ := json.Marshal(doc)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.SetBasicAuth(c.Username, c.Password)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	_, err := client.Do(req)
	return err
}