package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	APIKey   string
	Endpoint string
	client   *http.Client
}

type Model struct {
	Name          string            `json:"name"`
	ModelProvider string            `json:"model_provider"`
	ModelName     string            `json:"model_name"`
	APIBase       string            `json:"api_base,omitempty"`
	APIKey        string            `json:"api_key,omitempty"`
	Metadata      map[string]string `json:"metadata,omitempty"`
}

type Key struct {
	KeyAlias  string   `json:"key_alias"`
	TeamID    string   `json:"team_id"`
	Models    []string `json:"models,omitempty"`
	MaxBudget float64  `json:"max_budget,omitempty"`
	ExpiresAt string   `json:"expires_at,omitempty"`
	Key       string   `json:"key,omitempty"`
}

func NewClient(apiKey, endpoint string) *Client {
	return &Client{
		APIKey:   apiKey,
		Endpoint: endpoint,
		client:   &http.Client{},
	}
}

func (c *Client) doRequest(method, path string, body interface{}) (*http.Response, error) {
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, fmt.Errorf("failed to encode request body: %w", err)
		}
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", c.Endpoint, path), &buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "terraform-provider-litellm")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read error response: %w", err)
		}
		return nil, parseError(resp.StatusCode, body)
	}

	return resp, nil
}

// Model operations
func (c *Client) CreateModel(model *Model) error {
	if err := validateModel(model); err != nil {
		return err
	}

	resp, err := c.doRequest("POST", "/api/models", model)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(model)
}

func (c *Client) GetModel(name string) (*Model, error) {
	if name == "" {
		return nil, fmt.Errorf("model name cannot be empty")
	}

	resp, err := c.doRequest("GET", fmt.Sprintf("/api/models/%s", name), nil)
	if err != nil {
		if apiErr, ok := err.(*APIError); ok && apiErr.StatusCode == http.StatusNotFound {
			return nil, nil
		}
		return nil, err
	}
	defer resp.Body.Close()

	var model Model
	if err := json.NewDecoder(resp.Body).Decode(&model); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &model, nil
}

func (c *Client) UpdateModel(model *Model) error {
	if err := validateModel(model); err != nil {
		return err
	}

	resp, err := c.doRequest("PUT", fmt.Sprintf("/api/models/%s", model.Name), model)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(model)
}

func (c *Client) DeleteModel(name string) error {
	if name == "" {
		return fmt.Errorf("model name cannot be empty")
	}

	resp, err := c.doRequest("DELETE", fmt.Sprintf("/api/models/%s", name), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// Key operations
func (c *Client) CreateKey(key *Key) error {
	if err := validateKey(key); err != nil {
		return err
	}

	resp, err := c.doRequest("POST", "/api/keys", key)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(key)
}

func (c *Client) GetKey(keyAlias string) (*Key, error) {
	if keyAlias == "" {
		return nil, fmt.Errorf("key alias cannot be empty")
	}

	resp, err := c.doRequest("GET", fmt.Sprintf("/api/keys/%s", keyAlias), nil)
	if err != nil {
		if apiErr, ok := err.(*APIError); ok && apiErr.StatusCode == http.StatusNotFound {
			return nil, nil
		}
		return nil, err
	}
	defer resp.Body.Close()

	var key Key
	if err := json.NewDecoder(resp.Body).Decode(&key); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &key, nil
}

func (c *Client) UpdateKey(key *Key) error {
	if err := validateKey(key); err != nil {
		return err
	}

	resp, err := c.doRequest("PUT", fmt.Sprintf("/api/keys/%s", key.KeyAlias), key)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(key)
}

func (c *Client) DeleteKey(keyAlias string) error {
	if keyAlias == "" {
		return fmt.Errorf("key alias cannot be empty")
	}

	resp, err := c.doRequest("DELETE", fmt.Sprintf("/api/keys/%s", keyAlias), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func validateModel(model *Model) error {
	if model == nil {
		return fmt.Errorf("model cannot be nil")
	}
	if model.Name == "" {
		return fmt.Errorf("model name cannot be empty")
	}
	if model.ModelProvider == "" {
		return fmt.Errorf("model provider cannot be empty")
	}
	if model.ModelName == "" {
		return fmt.Errorf("model name cannot be empty")
	}
	return nil
}

func validateKey(key *Key) error {
	if key == nil {
		return fmt.Errorf("key cannot be nil")
	}
	if key.KeyAlias == "" {
		return fmt.Errorf("key alias cannot be empty")
	}
	if key.TeamID == "" {
		return fmt.Errorf("team ID cannot be empty")
	}
	return nil
}
