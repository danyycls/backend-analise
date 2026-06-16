package pncp

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

type PNCPClient struct {
	baseURL        string
	basePublicacao string
	client         *http.Client
}

func NovoPNCPClient(baseURL string) *PNCPClient {
	if baseURL == "" {
		baseURL = "https://pncp.gov.br/pncp-consulta/v1"
	}
	return &PNCPClient{
		baseURL:        baseURL + "/contratos",
		basePublicacao: baseURL + "/contratacoes/publicacao",
		client:         &http.Client{Timeout: 300 * time.Second},
	}
}

func (p *PNCPClient) BuscarContratos(ctx context.Context, cnpj, dataInicial, dataFinal string, pagina, tamanho int) ([]Contrato, error) {
	log := logger.New("Clients: Client: BuscarContratos")
	u, _ := url.Parse(p.baseURL)
	q := u.Query()
	q.Set("cnpjOrgao", cnpj)
	q.Set("dataInicial", dataInicial)
	q.Set("dataFinal", dataFinal)
	q.Set("pagina", fmt.Sprintf("%d", pagina))
	q.Set("tamanhoPagina", fmt.Sprintf("%d", tamanho))
	u.RawQuery = q.Encode()

	var lastErr error
	attempts := 3
	for i := 1; i <= attempts; i++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("User-Agent", "Liceu/1.0")
		req.Header.Set("Accept", "application/json")

		log.Info("PNCP: solicitando", "pagina", pagina, "tamanho", tamanho, "url", u.String(), "attempt", i)
		resp, err := p.client.Do(req)
		if err != nil {
			lastErr = err
			log.Error("PNCP: erro na requisicao", "attempt", i, "error", err)
		} else {
			if resp.StatusCode != 200 {
				body, _ := io.ReadAll(resp.Body)
				resp.Body.Close()
				lastErr = fmt.Errorf("pncp status %d: %s", resp.StatusCode, string(body))
				log.Error("PNCP: status nao 200", "attempt", i, "error", lastErr)
			} else {
				body, rerr := io.ReadAll(resp.Body)
				resp.Body.Close()
				if rerr != nil {
					lastErr = rerr
					log.Error("PNCP: erro ler body", "attempt", i, "error", rerr)
				} else {
					var contratos []Contrato
					if derr := json.Unmarshal(body, &contratos); derr == nil {
						return contratos, nil
					}

					var obj map[string]json.RawMessage
					if derr2 := json.Unmarshal(body, &obj); derr2 == nil {
						for _, raw := range obj {
							var arr []Contrato
							if derr3 := json.Unmarshal(raw, &arr); derr3 == nil {
								return arr, nil
							}
						}
						lastErr = fmt.Errorf("pncp: nao encontrou array de contratos no objeto retornado")
						log.Error("PNCP: decode objeto", "attempt", i, "error", lastErr)
					} else {
						lastErr = derr2
						log.Error("PNCP: erro decode", "attempt", i, "error", derr2)
					}
				}
			}
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(time.Duration(i) * 500 * time.Millisecond):
		}
	}
	return nil, lastErr
}

func (p *PNCPClient) buscarPublicacaoResponse(ctx context.Context, params map[string]string, pagina, tamanho int) (*PublicacaoResponse, error) {
	log := logger.New("Clients: Client: buscarPublicacaoResponse")
	u, _ := url.Parse(p.basePublicacao)
	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	q.Set("pagina", fmt.Sprintf("%d", pagina))
	q.Set("tamanhoPagina", fmt.Sprintf("%d", tamanho))
	u.RawQuery = q.Encode()
	var lastErr error
	attempts := 1
	for i := 1; i <= attempts; i++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("User-Agent", "Liceu/1.0")
		req.Header.Set("Accept", "application/json")

		log.Info("PNCP publicacao: solicitando", "pagina", pagina, "url", u.String(), "attempt", i)
		resp, err := p.client.Do(req)
		if err != nil {
			lastErr = err
			log.Error("PNCP publicacao: erro na requisicao", "attempt", i, "error", err)
		} else if resp.StatusCode == 204 {
			resp.Body.Close()
			log.Info("PNCP publicacao: 204 sem conteudo para municipio")
			return &PublicacaoResponse{}, nil
		} else if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			lastErr = fmt.Errorf("pncp publicacao status %d: %s", resp.StatusCode, string(body))
			log.Error("PNCP publicacao: status nao 200", "attempt", i, "error", lastErr)
		} else {
			body, rerr := io.ReadAll(resp.Body)
			resp.Body.Close()
			if rerr != nil {
				lastErr = rerr
				log.Error("PNCP publicacao: erro ler body", "attempt", i, "error", rerr)
			} else {
				var resp PublicacaoResponse
				if derr := json.Unmarshal(body, &resp); derr == nil {
					return &resp, nil
				} else {
					lastErr = derr
					log.Error("PNCP publicacao: erro decode", "attempt", i, "error", derr)
				}
			}
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(time.Duration(i) * 500 * time.Millisecond):
		}
	}
	return nil, lastErr
}

func (p *PNCPClient) BuscarContratacoesPorMunicipio(ctx context.Context, codigoMunicipio string, dataInicial, dataFinal string, pagina, tamanho int) (*PublicacaoResponse, error) {
	params := map[string]string{
		"codigoModalidadeContratacao": "8",
		"codigoMunicipioIbge":         codigoMunicipio,
		"dataInicial":                 dataInicial,
		"dataFinal":                   dataFinal,
	}
	return p.buscarPublicacaoResponse(ctx, params, pagina, tamanho)
}

func (p *PNCPClient) BuscarContratacoesPorUF(ctx context.Context, uf, dataInicial, dataFinal string, pagina, tamanho int) (*PublicacaoResponse, error) {
	params := map[string]string{
		"codigoModalidadeContratacao": "8",
		"uf":                          uf,
		"dataInicial":                 dataInicial,
		"dataFinal":                   dataFinal,
	}
	return p.buscarPublicacaoResponse(ctx, params, pagina, tamanho)
}
