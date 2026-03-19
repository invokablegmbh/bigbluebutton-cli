package api

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// Client is the BigBlueButton API client.
type Client struct {
	BaseURL    string
	Secret     string
	HTTPClient *http.Client
	Verbose    bool
}

// NewClient creates a new BigBlueButton API client.
func NewClient(baseURL, secret string, verbose bool) *Client {
	// Ensure baseURL ends with /api
	baseURL = strings.TrimRight(baseURL, "/")
	if !strings.HasSuffix(baseURL, "/api") {
		if !strings.HasSuffix(baseURL, "/bigbluebutton") {
			baseURL += "/bigbluebutton"
		}
		baseURL += "/api"
	}
	return &Client{
		BaseURL:    baseURL,
		Secret:     secret,
		HTTPClient: http.DefaultClient,
		Verbose:    verbose,
	}
}

// doGet makes a GET request to the BBB API endpoint with the given parameters.
func (c *Client) doGet(apiCall string, params url.Values) ([]byte, error) {
	qs := params.Encode()
	checksum := Checksum(apiCall, qs, c.Secret)
	if qs != "" {
		qs += "&"
	}
	qs += "checksum=" + checksum

	fullURL := fmt.Sprintf("%s/%s?%s", c.BaseURL, apiCall, qs)
	if c.Verbose {
		log.Printf("-> GET %s/%s?%s&checksum=***", c.BaseURL, apiCall, params.Encode())
	}

	resp, err := c.HTTPClient.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("HTTP request to %s failed: %w", apiCall, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response from %s: %w", apiCall, err)
	}

	if c.Verbose {
		log.Printf("<- %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
		log.Printf("<- %s", string(body))
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d from %s", resp.StatusCode, apiCall)
	}

	return body, nil
}

// doPost makes a POST request to the BBB API endpoint.
func (c *Client) doPost(apiCall string, params url.Values, contentType string, bodyContent io.Reader) ([]byte, error) {
	qs := params.Encode()
	checksum := Checksum(apiCall, qs, c.Secret)
	if qs != "" {
		qs += "&"
	}
	qs += "checksum=" + checksum

	fullURL := fmt.Sprintf("%s/%s?%s", c.BaseURL, apiCall, qs)
	if c.Verbose {
		log.Printf("-> POST %s/%s?%s&checksum=***", c.BaseURL, apiCall, params.Encode())
	}

	req, err := http.NewRequest("POST", fullURL, bodyContent)
	if err != nil {
		return nil, fmt.Errorf("creating request for %s: %w", apiCall, err)
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request to %s failed: %w", apiCall, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response from %s: %w", apiCall, err)
	}

	if c.Verbose {
		log.Printf("<- %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
		log.Printf("<- %s", string(body))
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d from %s", resp.StatusCode, apiCall)
	}

	return body, nil
}

// parseXMLResponse unmarshals XML response and checks for API errors.
func parseXMLResponse[T interface{ IsSuccess() bool }](data []byte, resp T) (T, error) {
	if err := xml.Unmarshal(data, resp); err != nil {
		return resp, fmt.Errorf("parsing XML response: %w", err)
	}
	return resp, nil
}

// parseJSONResponse unmarshals JSON response.
func parseJSONResponse[T any](data []byte, resp T) (T, error) {
	if err := json.Unmarshal(data, resp); err != nil {
		return resp, fmt.Errorf("parsing JSON response: %w", err)
	}
	return resp, nil
}
