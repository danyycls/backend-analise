package portaltransparencia

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/danyele/laceu/internal/shared/logger"
)

type PortalTransparenciaClient struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

func NovoPortalTransparenciaClient(apiKey, baseURL string) *PortalTransparenciaClient {
	if baseURL == "" {
		baseURL = "https://api.portaldatransparencia.gov.br"
	}
	return &PortalTransparenciaClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *PortalTransparenciaClient) doGet(ctx context.Context, path string, queryParams map[string]string, dest interface{}) error {
	log := logger.New("Clients: Client: doGet")
	fullURL, err := url.Parse(c.baseURL + path)
	if err != nil {
		return fmt.Errorf("portaltransparencia: erro ao parsear URL: %w", err)
	}

	q := fullURL.Query()
	for k, v := range queryParams {
		if v != "" {
			q.Set(k, v)
		}
	}
	fullURL.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL.String(), nil)
	if err != nil {
		return fmt.Errorf("portaltransparencia: erro ao criar request: %w", err)
	}

	req.Header.Set("chave-api-dados", c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("portaltransparencia: erro ao fazer request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Error("API portaltransparencia retornou status", "status", resp.StatusCode, "url", fullURL.Redacted(), "body", string(body))
		return fmt.Errorf("portaltransparencia: status %d - %s", resp.StatusCode, string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(dest); err != nil {
		return fmt.Errorf("portaltransparencia: erro ao decodificar response: %w", err)
	}

	return nil
}
