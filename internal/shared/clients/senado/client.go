package senado

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type SenadoClient struct {
	baseURL string
	client  *http.Client
}

func makeQueryParams(params map[string]string) url.Values {
	if len(params) == 0 {
		return nil
	}
	q := make(url.Values)
	for k, v := range params {
		if v != "" {
			q.Set(k, v)
		}
	}
	return q
}

func NovoSenadoClient(baseURL string) *SenadoClient {
	return &SenadoClient{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 60 * time.Second},
	}
}

func (c *SenadoClient) doGetJSON(ctx context.Context, path string, query url.Values, dest any) error {
	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return fmt.Errorf("senado: erro url: %w", err)
	}
	if query != nil {
		u.RawQuery = query.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return fmt.Errorf("senado: erro requisicao: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("senado: erro execucao: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("senado: status %d: %s", resp.StatusCode, string(body))
	}
	if err := json.NewDecoder(resp.Body).Decode(dest); err != nil {
		return fmt.Errorf("senado: erro decode: %w", err)
	}
	return nil
}
