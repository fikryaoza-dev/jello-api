package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"jello-api/internal/domain"
	"jello-api/internal/model"
	"jello-api/internal/repository/mapper"
	"jello-api/pkg/couchdb"
	"net/http"
)

type ITableRepository interface {
	Create(table *domain.Table) error
	GetAll() ([]domain.Table, error)
}

type couchTableRepo struct {
	client *couchdb.Client
}

func NewCouchTableRepo(client *couchdb.Client) ITableRepository {
	return &couchTableRepo{client: client}
}

func (r *couchTableRepo) Create(table *domain.Table) error {

	doc := mapper.ToModel(*table)

	body, _ := json.Marshal(doc)

	url := fmt.Sprintf("%s/%s", r.client.BaseURL, r.client.DBName)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.SetBasicAuth(r.client.Username, r.client.Password)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}

func (r *couchTableRepo) GetAll() ([]domain.Table, error) {
	query := map[string]interface{}{
		"selector": map[string]interface{}{
			"type": "table",
		},
	}

	body, _ := json.Marshal(query)

	url := fmt.Sprintf("%s/%s/_find", r.client.BaseURL, r.client.DBName)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.SetBasicAuth(r.client.Username, r.client.Password)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)

	var result struct {
		Docs []model.Table `json:"docs"`
	}

	json.Unmarshal(data, &result)

	tables := []domain.Table{}

	for _, doc := range result.Docs {
		tables = append(tables, mapper.ToDomain(doc))
	}

	return tables, nil
}