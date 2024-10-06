package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type APIClient struct {
	BaseURL string
	Client  *http.Client
}

func NewAPIClient() *APIClient {
	return &APIClient{
		BaseURL: "http://localhost:8082",
		Client:  &http.Client{},
	}
}

func (api *APIClient) Get(key string) ([]byte, error) {
	url := fmt.Sprintf("%s/GET", api.BaseURL)

	data := map[string]string{
		"key": key,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := api.Client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (api *APIClient) Set(key, value string) ([]byte, error) {
	url := fmt.Sprintf("%s/SET", api.BaseURL)

	data := map[string]string{
		"key":   key,
		"value": value,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := api.Client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (api *APIClient) Expire(key string, seconds int) ([]byte, error) {
	url := fmt.Sprintf("%s/EXPIRE", api.BaseURL)

	data := map[string]interface{}{
		"key":     key,
		"seconds": seconds,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := api.Client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
