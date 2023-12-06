package agent

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type MetricAgent interface {
	FetchMetrics(query string) ([]byte, error)
}

type PromScaleWayAgent struct {
	client   *http.Client
	endpoint string
	token    string
}

func NewPromScaleWayAgent(client *http.Client, endpoint, token string) MetricAgent {
	return &PromScaleWayAgent{
		client:   client,
		endpoint: endpoint,
		token:    token,
	}
}

func (p *PromScaleWayAgent) FetchMetrics(query string) ([]byte, error) {
	fullURL := fmt.Sprintf("%s/prometheus/api/v1/query?query=%s", p.endpoint, url.QueryEscape(query))
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+p.token)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Log the response body for debugging
	log.Printf("Fetched Metrics size is: %d", len(body))

	return body, nil
}
