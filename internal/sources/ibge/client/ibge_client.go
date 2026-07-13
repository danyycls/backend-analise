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
	ID         string `json:"id"`
	Variavel   string `json:"variavel"`
	Unidade    string `json:"unidade"`
	Resultados []struct {
		Series []struct {
			Localidade struct {
				ID string `json:"id"`
			} `json:"localidade"`
			Serie map[string]string `json:"serie"`
		} `json:"series"`
	} `json:"resultados"`
}

func (c *IBGEClient) BuscarPopulacao(ctx context.Context, municipioIDs []int) (map[int]int64, error) {
	log := logger.New("Clients: Client: BuscarPopulacao")

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

		if c.extrairPopulacaoEstimada(ctx, log, localidades, resultado) {
			continue
		}

		c.extrairPopulacaoCenso2022(ctx, log, localidades, resultado)
	}

	log.Info(
		"BuscarPopulacao concluido",
		"municipios_com_populacao", len(resultado),
		"total_solicitado", len(municipioIDs),
	)

	return resultado, nil
}

func (c *IBGEClient) extrairPopulacaoEstimada(
	ctx context.Context,
	log *logger.Logger,
	localidades string,
	resultado map[int]int64,
) bool {

	url := fmt.Sprintf(
		"%s/6579/periods/-6/variaveis/9324?localidades=%s",
		c.agregadosURL,
		localidades,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Warn("IBGE estimativa: erro criar requisicao", "erro", err)
		return false
	}

	resp, err := c.client.Do(req)
	if err != nil {
		log.Warn("IBGE estimativa: erro na requisicao", "erro", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Warn("IBGE estimativa: status inesperado", "status", resp.StatusCode)
		return false
	}

	var dados []populacaoV3
	if err := json.NewDecoder(resp.Body).Decode(&dados); err != nil {
		log.Warn("IBGE estimativa: erro decode", "erro", err)
		return false
	}

	for _, item := range dados {
		for _, r := range item.Resultados {
			for _, s := range r.Series {
				locID, err := strconv.Atoi(s.Localidade.ID)
				if err != nil {
					continue
				}

				var (
					maiorAno int
					pop      int64
				)

				for anoStr, valor := range s.Serie {
					ano, err := strconv.Atoi(anoStr)
					if err != nil {
						continue
					}

					valor = strings.ReplaceAll(valor, ".", "")
					valor = strings.TrimSpace(valor)

					p, err := strconv.ParseInt(valor, 10, 64)
					if err != nil {
						continue
					}

					if ano > maiorAno {
						maiorAno = ano
						pop = p
					}
				}

				if pop > 0 {
					resultado[locID] = pop
				}
			}
		}
	}

	return len(resultado) > 0
}

func (c *IBGEClient) extrairPopulacaoCenso2022(
	ctx context.Context,
	log *logger.Logger,
	localidades string,
	resultado map[int]int64,
) {

	url := fmt.Sprintf(
		"%s/8395/periods/2022/variaveis/12494?localidades=%s",
		c.agregadosURL,
		localidades,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Warn("IBGE censo 2022: erro criar requisicao", "erro", err)
		return
	}

	resp, err := c.client.Do(req)
	if err != nil {
		log.Warn("IBGE censo 2022: erro na requisicao", "erro", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Warn("IBGE censo 2022: status inesperado", "status", resp.StatusCode)
		return
	}

	var dados []populacaoV3
	if err := json.NewDecoder(resp.Body).Decode(&dados); err != nil {
		log.Warn("IBGE censo 2022: erro decode", "erro", err)
		return
	}

	for _, item := range dados {
		for _, r := range item.Resultados {
			for _, s := range r.Series {
				locID, err := strconv.Atoi(s.Localidade.ID)
				if err != nil {
					continue
				}

				for _, valor := range s.Serie {
					valor = strings.ReplaceAll(valor, ".", "")
					valor = strings.TrimSpace(valor)

					pop, err := strconv.ParseInt(valor, 10, 64)
					if err != nil {
						continue
					}

					resultado[locID] = pop
					break
				}
			}
		}
	}
}
