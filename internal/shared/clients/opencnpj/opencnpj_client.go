package opencnpj

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/danyele/podp/internal/shared/types"
)

type OpenCNPJClient struct {
	baseURL string
	client  *http.Client
}

func NovoOpenCNPJClient(baseURL string) *OpenCNPJClient {
	if baseURL == "" {
		baseURL = "https://api.opencnpj.org/%s"
	}
	return &OpenCNPJClient{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

func (o *OpenCNPJClient) Buscar(ctx context.Context, cnpj string) (*types.OpenCNPJResponse, error) {
	url := fmt.Sprintf(o.baseURL, cnpj)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := o.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("opencnpj status %d", resp.StatusCode)
	}
	var out types.OpenCNPJResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}
