package webclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	BaseURL string
	Client  *http.Client
}

func New(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		Client:  &http.Client{},
	}
}

func (c *Client) get(path string, out interface{}) error {
	resp, err := c.Client.Get(c.BaseURL + path)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("server error: %s", resp.Status)
	}

	return json.NewDecoder(resp.Body).Decode(out)
}

func (c *Client) delete(path string, out interface{}) error {
	req, _ := http.NewRequest("DELETE", c.BaseURL+path, nil)
	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(out)
}

func (c *Client) post(path string, payload interface{}, out interface{}) error {
	body, _ := json.Marshal(payload)
	resp, err := c.Client.Post(c.BaseURL+path, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(out)
}

func (c *Client) put(path string, payload interface{}, out interface{}) error {
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("PUT", c.BaseURL+path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(out)
}
