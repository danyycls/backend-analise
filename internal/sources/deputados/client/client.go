package deputados

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type DeputadosClient struct {
	baseURL string
	client  *http.Client
}

func NovoDeputadosClient(baseURL string) *DeputadosClient {
	return &DeputadosClient{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *DeputadosClient) doGet(ctx context.Context, path string, query url.Values, dest any) error {
	resultado, err := c.doGetRaw(ctx, path, query)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(resultado.Dados, dest); err != nil {
		return fmt.Errorf("deputados: erro dados: %w", err)
	}
	return nil
}

func (c *DeputadosClient) doGetRaw(ctx context.Context, path string, query url.Values) (*Resultado, error) {
	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return nil, fmt.Errorf("deputados: erro url: %w", err)
	}
	if query != nil {
		u.RawQuery = query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("deputados: erro requisicao: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("deputados: erro execucao: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("deputados: status %d: %s", resp.StatusCode, string(body))
	}

	var resultado Resultado
	if err := json.NewDecoder(resp.Body).Decode(&resultado); err != nil {
		return nil, fmt.Errorf("deputados: erro decode: %w", err)
	}

	return &resultado, nil
}
