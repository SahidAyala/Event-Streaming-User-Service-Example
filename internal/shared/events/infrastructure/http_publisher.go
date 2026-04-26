package infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/SahidAyala/Event-Streaming-User-Service-Example/internal/shared/events"
)

type HTTPEventPublisher struct {
	baseURL string
	client  *http.Client
	apiKey  string
}

func NewHTTPEventPublisher(baseURL, apiKey string) *HTTPEventPublisher {
	return &HTTPEventPublisher{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 5 * time.Second},
		apiKey:  apiKey,
	}
}

func (p *HTTPEventPublisher) Publish(ctx context.Context, event events.Event) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/events", bytes.NewBuffer(body))

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", p.apiKey)
	resp, err := p.client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("Failed to publish event: %d", resp.StatusCode)
	}

	return nil
}

var _ events.Publisher = (*HTTPEventPublisher)(nil)