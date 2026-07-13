package pncp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/danyele/podp/internal/shared/logger"
)

const (
	cbLimiarFalhas  = 5
	cbTimeout       = 30 * time.Second
	maxAttempts     = 3
	backoffBase     = 1 * time.Second
	backoffBaseFast = 500 * time.Millisecond
)

type circuitBreaker struct {
	mu           sync.Mutex
	falhasConsec int
	abertoAte    time.Time
}

type PNCPClient struct {
	baseURL          string
	baseContratacoes string
	client           *http.Client
	cb               circuitBreaker
}

func NovoPNCPClient(baseURL string) *PNCPClient {
	return &PNCPClient{
		baseURL:          baseURL + "/contratos",
		baseContratacoes: baseURL + "/contratacoes/publicacao",
		client:           &http.Client{Timeout: 300 * time.Second},
	}
}

func (p *PNCPClient) BuscarContratos(ctx context.Context, cnpj, dataInicial, dataFinal string, pagina, tamanho int) (*ContratoResponse, error) {
	if p.cb.isAberto() {
		return nil, fmt.Errorf("circuit breaker aberto para PNCP")
	}

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
	for i := 1; i <= maxAttempts; i++ {
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
			p.cb.falha()
		} else {
			if resp.StatusCode == http.StatusNoContent {
				resp.Body.Close()
				p.cb.sucesso()
				log.Info("PNCP: 204 sem conteudo")
				return &ContratoResponse{}, nil
			}
			if resp.StatusCode != http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				resp.Body.Close()
				lastErr = fmt.Errorf("pncp status %d: %s", resp.StatusCode, string(body))
				log.Error("PNCP: status nao 200", "attempt", i, "error", lastErr)
				if !isRetryable(resp.StatusCode, nil) {
					p.cb.falha()
					return nil, lastErr
				}
				p.cb.falha()
			} else {
				p.cb.sucesso()
				body, rerr := io.ReadAll(resp.Body)
				resp.Body.Close()
				if rerr != nil {
					lastErr = rerr
					log.Error("PNCP: erro ler body", "attempt", i, "error", rerr)
					if !isRetryable(0, rerr) {
						return nil, lastErr
					}
					continue
				}

				var pubResp ContratoResponse
				derr := json.Unmarshal(body, &pubResp)
				if derr == nil {
					return &pubResp, nil
				}
				lastErr = derr
				log.Error("PNCP: erro decode", "attempt", i, "error", derr)
				if !isRetryable(0, lastErr) {
					return nil, lastErr
				}
			}
		}

		if err := waitWithBackoff(ctx, i, backoffBaseFast); err != nil {
			return nil, err
		}
	}
	return nil, lastErr
}

func (p *PNCPClient) BuscarContratosPorMunicipio(ctx context.Context, codigoMunicipio string, dataInicial, dataFinal, codigoModalidade string, pagina, tamanho int) (*ContratoResponse, error) {
	if codigoModalidade == "" {
		codigoModalidade = "8"
	}
	params := map[string]string{
		"codigoModalidadeContratacao": codigoModalidade,
		"codigoMunicipioIbge":         codigoMunicipio,
		"dataInicial":                 dataInicial,
		"dataFinal":                   dataFinal,
	}
	return p.parsebuscarContratosResponse(ctx, params, pagina, tamanho)
}

func (p *PNCPClient) BuscarContratosPorUF(ctx context.Context, uf, dataInicial, dataFinal, codigoModalidade string, pagina, tamanho int) (*ContratoResponse, error) {
	if codigoModalidade == "" {
		codigoModalidade = "8"
	}
	params := map[string]string{
		"codigoModalidadeContratacao": codigoModalidade,
		"uf":                          uf,
		"dataInicial":                 dataInicial,
		"dataFinal":                   dataFinal,
	}
	return p.parsebuscarContratosResponse(ctx, params, pagina, tamanho)
}

func (p *PNCPClient) parsebuscarContratosResponse(ctx context.Context, params map[string]string, pagina, tamanho int) (*ContratoResponse, error) {
	if p.cb.isAberto() {
		return nil, fmt.Errorf("circuit breaker aberto para PNCP")
	}

	log := logger.New("Clients: Client: parsebuscarContratosResponse")
	u, _ := url.Parse(p.baseContratacoes)
	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	q.Set("pagina", fmt.Sprintf("%d", pagina))
	q.Set("tamanhoPagina", fmt.Sprintf("%d", tamanho))
	u.RawQuery = q.Encode()

	var lastErr error
	for i := 1; i <= maxAttempts; i++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("User-Agent", "Liceu/1.0")
		req.Header.Set("Accept", "application/json")

		log.Info("PNCP contratos: solicitando", "pagina", pagina, "url", u.String(), "attempt", i)
		resp, err := p.client.Do(req)

		shouldRetry := false
		switch {
		case err != nil:
			lastErr = err
			log.Error("PNCP publicacao: erro na requisicao", "attempt", i, "error", err)
			p.cb.falha()
			shouldRetry = true
		case resp.StatusCode == http.StatusNoContent:
			resp.Body.Close()
			p.cb.sucesso()
			log.Info("PNCP publicacao: 204 sem conteudo para municipio")
			return &ContratoResponse{}, nil
		case resp.StatusCode != http.StatusOK:
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			lastErr = fmt.Errorf("pncp publicacao status %d: %s", resp.StatusCode, string(body))
			log.Error("PNCP publicacao: status nao 200", "attempt", i, "error", lastErr)
			if isRetryable(resp.StatusCode, nil) {
				shouldRetry = true
				p.cb.falha()
			} else {
				p.cb.falha()
				return nil, lastErr
			}
		default:
			p.cb.sucesso()
			body, rerr := io.ReadAll(resp.Body)
			resp.Body.Close()
			if rerr != nil {
				lastErr = rerr
				log.Error("PNCP publicacao: erro ler body", "attempt", i, "error", rerr)
				shouldRetry = true
				continue
			}
			var pubResp ContratoResponse
			derr := json.Unmarshal(body, &pubResp)
			if derr == nil {
				return &pubResp, nil
			}
			lastErr = derr
			log.Error("PNCP publicacao: erro decode", "attempt", i, "error", derr)
			shouldRetry = true
		}

		if !shouldRetry {
			return nil, lastErr
		}

		if err := waitWithBackoff(ctx, i, backoffBase); err != nil {
			return nil, err
		}
	}
	return nil, lastErr
}

func isRetryable(statusCode int, err error) bool {
	if err != nil {
		return true
	}
	return statusCode == http.StatusServiceUnavailable ||
		statusCode == http.StatusTooManyRequests ||
		statusCode == http.StatusGatewayTimeout
}

func waitWithBackoff(ctx context.Context, attempt int, base time.Duration) error {
	backoff := float64(base) * math.Pow(2, float64(attempt-1))
	jitter := (rand.Float64()*2 - 1) * backoff * 0.3
	total := time.Duration(backoff + jitter)
	if total < 0 {
		total = 0
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(total):
		return nil
	}
}

func (cb *circuitBreaker) isAberto() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return time.Now().Before(cb.abertoAte)
}

func (cb *circuitBreaker) sucesso() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.falhasConsec = 0
}

func (cb *circuitBreaker) falha() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.falhasConsec++
	if cb.falhasConsec >= cbLimiarFalhas {
		cb.abertoAte = time.Now().Add(cbTimeout)
	}
}
