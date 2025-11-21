package siem

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kxrty/loggerv2/internal/models"
)

// HTTPForwarder отправляет события через HTTP/HTTPS (для Elastic, Splunk HEC, etc.)
type HTTPForwarder struct {
	url        string
	token      string
	client     *http.Client
	headers    map[string]string
}

// NewHTTPForwarder создает новый HTTP форвардер
func NewHTTPForwarder(url, token string, headers map[string]string) *HTTPForwarder {
	if headers == nil {
		headers = make(map[string]string)
	}

	return &HTTPForwarder{
		url:     url,
		token:   token,
		headers: headers,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Forward отправляет событие через HTTP
func (f *HTTPForwarder) Forward(event *models.GOSTEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("ошибка сериализации события: %w", err)
	}

	req, err := http.NewRequest("POST", f.url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("ошибка создания запроса: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	
	if f.token != "" {
		req.Header.Set("Authorization", "Bearer "+f.token)
	}

	for key, value := range f.headers {
		req.Header.Set(key, value)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка отправки в SIEM: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("SIEM вернул ошибку: %d", resp.StatusCode)
	}

	return nil
}

// ForwardBatch отправляет несколько событий одним запросом
func (f *HTTPForwarder) ForwardBatch(events []*models.GOSTEvent) error {
	payload, err := json.Marshal(events)
	if err != nil {
		return fmt.Errorf("ошибка сериализации событий: %w", err)
	}

	req, err := http.NewRequest("POST", f.url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("ошибка создания запроса: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	
	if f.token != "" {
		req.Header.Set("Authorization", "Bearer "+f.token)
	}

	for key, value := range f.headers {
		req.Header.Set(key, value)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка отправки в SIEM: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("SIEM вернул ошибку: %d", resp.StatusCode)
	}

	return nil
}
