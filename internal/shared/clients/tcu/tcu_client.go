package tcu

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/danyele/podp/internal/shared/logger"
)

type TCUClient struct {
	baseURL string
	client  *http.Client
}

func NovoTCUClient(baseURL string) *TCUClient {
	return &TCUClient{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 60 * time.Second},
	}
}

func (c *TCUClient) doPostJSON(ctx context.Context, path string, body any, dest any) error {
	log := logger.New("Clients: Client: doPostJSON")
	url := c.baseURL + path

	var lastErr error
	attempts := 3
	for i := 1; i <= attempts; i++ {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("tcu: erro marshal body: %w", err)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
		if err != nil {
			return fmt.Errorf("tcu: erro criando requisicao: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")

		log.Info("TCU: POST", "url", url, "attempt", i)
		resp, err := c.client.Do(req)
		if err != nil {
			lastErr = err
			log.Error("TCU: erro requisicao", "attempt", i, "error", err)
		} else {
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				rBody, _ := io.ReadAll(resp.Body)
				lastErr = fmt.Errorf("tcu: status %d: %s", resp.StatusCode, string(rBody))
				log.Error("TCU: status nao 200", "attempt", i, "error", lastErr)
			} else {
				rBody, rerr := io.ReadAll(resp.Body)
				if rerr != nil {
					lastErr = rerr
					log.Error("TCU: erro ler body", "attempt", i, "error", rerr)
				} else {
					derr := json.Unmarshal(rBody, dest)
					if derr == nil {
						return nil
					}
					lastErr = derr
					log.Error("TCU: erro decode", "attempt", i, "error", derr)
				}
			}
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Duration(i) * 500 * time.Millisecond):
		}
	}
	return lastErr
}

func (c *TCUClient) BuscarContasIrregulares(ctx context.Context, filter TCUQueryParams) ([]ContasIrregulares, error) {
	var result []ContasIrregulares
	err := c.doPostJSON(ctx, "/responsive-contas-irregulares", filter, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *TCUClient) BuscarFinsEleitorais(ctx context.Context, filter TCUQueryParams) ([]FinsEleitorais, error) {
	var result []FinsEleitorais
	err := c.doPostJSON(ctx, "/responsive-fins-eleitorais", filter, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *TCUClient) BuscarInabilitados(ctx context.Context, filter TCUQueryParams) ([]Sancoes, error) {
	var result []Sancoes
	err := c.doPostJSON(ctx, "/responsive-inabilitados", filter, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *TCUClient) BuscarInidoneos(ctx context.Context, filter TCUQueryParams) ([]Sancoes, error) {
	var result []Sancoes
	err := c.doPostJSON(ctx, "/responsive-inidoneos", filter, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
