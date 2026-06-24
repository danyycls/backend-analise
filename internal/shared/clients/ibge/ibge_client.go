package ibge

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/danyele/podp/internal/shared/logger"
	"github.com/danyele/podp/internal/shared/types"
)

type IBGEClient struct {
	baseURL      string
	agregadosURL string
	client       *http.Client
}

func NovoIBGEClient(baseURL, agregadosURL string) *IBGEClient {
	return &IBGEClient{
		baseURL:      baseURL,
		agregadosURL: agregadosURL,
		client:       &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *IBGEClient) ListarMunicipios(ctx context.Context, uf string) ([]types.MunicipioIBGE, error) {
	log := logger.New("Clients: Client: ListarMunicipios")
	url := fmt.Sprintf("%s/estados/%s/municipios", c.baseURL, uf)
	log.Info("IBGE: solicitando municipios", "uf", uf)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("ibge: erro criar requisicao: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ibge: erro requisicao: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("ibge: status %d", resp.StatusCode)
	}

	var municipios []types.MunicipioIBGE
	if err := json.NewDecoder(resp.Body).Decode(&municipios); err != nil {
		return nil, fmt.Errorf("ibge: erro decode: %w", err)
	}

	return municipios, nil
}

func (c *IBGEClient) ListarEstados(ctx context.Context) ([]types.EstadoIBGE, error) {
	log := logger.New("Clients: Client: ListarEstados")
	url := fmt.Sprintf("%s/estados", c.baseURL)
	log.Info("IBGE: solicitando estados")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("ibge: erro criar requisicao: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ibge: erro requisicao: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("ibge: status %d", resp.StatusCode)
	}

	var estados []types.EstadoIBGE
	if err := json.NewDecoder(resp.Body).Decode(&estados); err != nil {
		return nil, fmt.Errorf("ibge: erro decode estados: %w", err)
	}

	return estados, nil
}

func (c *IBGEClient) ListarMunicipiosCompleto(ctx context.Context) ([]types.MunicipioDetalhadoIBGE, error) {
	log := logger.New("Clients: Client: ListarMunicipiosCompleto")
	url := fmt.Sprintf("%s/municipios", c.baseURL)
	log.Info("IBGE: solicitando todos os municipios com detalhes")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("ibge: erro criar requisicao: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ibge: erro requisicao: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ibge: status %d: %s", resp.StatusCode, string(body))
	}

	var municipios []types.MunicipioDetalhadoIBGE
	if err := json.NewDecoder(resp.Body).Decode(&municipios); err != nil {
		return nil, fmt.Errorf("ibge: erro decode municipios completo: %w", err)
	}

	return municipios, nil
}

type populacaoV3 struct {
	Resultados []struct {
		Series []struct {
			Localidade struct {
				ID string `json:"id"`
			} `json:"localidade"`
			Serie map[string]string `json:"serie"`
		} `json:"series"`
	} `json:"resultados"`
}

type populacaoV2 struct {
	Resultados []struct {
		Series [][]interface{} `json:"serie"`
	} `json:"resultados"`
}

func (c *IBGEClient) BuscarPopulacao(ctx context.Context, municipioIDs []int) (map[int]int64, error) {
	if len(municipioIDs) == 0 {
		return map[int]int64{}, nil
	}

	resultado := make(map[int]int64)
	batchSize := 50

	for i := 0; i < len(municipioIDs); i += batchSize {
		end := i + batchSize
		if end > len(municipioIDs) {
			end = len(municipioIDs)
		}
		batch := municipioIDs[i:end]

		ids := make([]string, len(batch))
		for j, id := range batch {
			ids[j] = strconv.Itoa(id)
		}

		localidades := "N6[" + strings.Join(ids, ",") + "]"

		// Tenta tabela nova (6579, var 9324) — formato v3: series[].localidade + series[].serie{ano: valor}
		url := fmt.Sprintf("%s/6579/periods/-6/variaveis/9324?localidades=%s", c.agregadosURL, localidades)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			continue
		}

		resp, err := c.client.Do(req)
		if err != nil || resp.StatusCode != 200 {
			if resp != nil {
				resp.Body.Close()
			}
			// Fallback: tenta tabela antiga (4700, var 93)
			url = fmt.Sprintf("%s/4700/periods/-6/variaveis/93?localidades=%s", c.agregadosURL, localidades)
			req, _ = http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
			resp, err = c.client.Do(req)
			if err != nil || resp.StatusCode != 200 {
				if resp != nil {
					resp.Body.Close()
				}
				continue
			}

			var v2 []populacaoV2
			if err := json.NewDecoder(resp.Body).Decode(&v2); err != nil {
				resp.Body.Close()
				continue
			}
			resp.Body.Close()

			for idx, item := range v2 {
				if len(item.Resultados) > 0 && len(item.Resultados[0].Series) > 0 {
					serie := item.Resultados[0].Series[0]
					if len(serie) > 1 {
						s := fmt.Sprintf("%v", serie[1])
						s = strings.ReplaceAll(s, ".", "")
						s = strings.TrimSpace(s)
						if pop, err := strconv.ParseInt(s, 10, 64); err == nil {
							resultado[batch[idx]] = pop
						}
					}
				}
			}
			continue
		}

		var v3 []populacaoV3
		if err := json.NewDecoder(resp.Body).Decode(&v3); err != nil {
			resp.Body.Close()
			continue
		}
		resp.Body.Close()

		for _, item := range v3 {
			if len(item.Resultados) == 0 || len(item.Resultados[0].Series) == 0 {
				continue
			}
			for _, s := range item.Resultados[0].Series {
				locID, err := strconv.Atoi(s.Localidade.ID)
				if err != nil {
					continue
				}
				// pega o valor do ano mais recente
				for _, v := range s.Serie {
					v = strings.TrimSpace(v)
					if pop, err := strconv.ParseInt(v, 10, 64); err == nil {
						resultado[locID] = pop
					}
					break
				}
			}
		}
	}

	return resultado, nil
}
