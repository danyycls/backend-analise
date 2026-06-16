package deputados

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

func (c *DeputadosClient) ListarInfoDeputadosAtivos(ctx context.Context, params map[string]string) ([]Deputado, error) {
	query := make(url.Values)
	for k, v := range params {
		query.Set(k, v)
	}
	var out []Deputado
	if err := c.doGet(ctx, "/deputados", query, &out); err != nil {
		return nil, fmt.Errorf("listar deputados: %w", err)
	}
	return out, nil
}

func (c *DeputadosClient) BuscarDeputado(ctx context.Context, id int) (*DeputadoDetalhe, error) {
	var out DeputadoDetalhe
	if err := c.doGet(ctx, "/deputados/"+strconv.Itoa(id), nil, &out); err != nil {
		return nil, fmt.Errorf("buscar deputado: %w", err)
	}
	return &out, nil
}

func (c *DeputadosClient) ListarDespesasPorDeputado(ctx context.Context, idDeputado int, params map[string]string) ([]DeputadoDespesa, error) {
	query := make(url.Values)
	for k, v := range params {
		query.Set(k, v)
	}
	var out []DeputadoDespesa
	path := fmt.Sprintf("/deputados/%d/despesas", idDeputado)
	if err := c.doGet(ctx, path, query, &out); err != nil {
		return nil, fmt.Errorf("listar despesas: %w", err)
	}
	return out, nil
}

func (c *DeputadosClient) ListarTodasDespesasPorDeputado(ctx context.Context, idDeputado int, params map[string]string) ([]DeputadoDespesa, error) {
	query := make(url.Values)
	for k, v := range params {
		query.Set(k, v)
	}

	var todas []DeputadoDespesa
	path := fmt.Sprintf("/deputados/%d/despesas", idDeputado)

	for {
		resultado, err := c.doGetRaw(ctx, path, query)
		if err != nil {
			return nil, fmt.Errorf("listar todas despesas: %w", err)
		}

		var pagina []DeputadoDespesa
		if err := json.Unmarshal(resultado.Dados, &pagina); err != nil {
			return nil, fmt.Errorf("listar todas despesas: erro decode: %w", err)
		}

		todas = append(todas, pagina...)

		nextURL := ""
		for _, link := range resultado.Links {
			if link.Rel == "next" {
				nextURL = link.Href
				break
			}
		}

		if nextURL == "" {
			break
		}

		parsed, err := url.Parse(nextURL)
		if err != nil {
			break
		}
		query = parsed.Query()
	}

	return todas, nil //nolint:nilerr
}

func (c *DeputadosClient) ListarFrentesDeputado(ctx context.Context, idDeputado int) ([]Frente, error) {
	var out []Frente
	path := fmt.Sprintf("/deputados/%d/frentes", idDeputado)
	if err := c.doGet(ctx, path, nil, &out); err != nil {
		return nil, fmt.Errorf("listar frentes deputado: %w", err)
	}
	return out, nil
}

func (c *DeputadosClient) ListarHistorico(ctx context.Context, idDeputado int) ([]DeputadoHistorico, error) {
	var out []DeputadoHistorico
	path := fmt.Sprintf("/deputados/%d/historico", idDeputado)
	if err := c.doGet(ctx, path, nil, &out); err != nil {
		return nil, fmt.Errorf("listar historico: %w", err)
	}
	return out, nil
}

func (c *DeputadosClient) ListarMandatosExternos(ctx context.Context, idDeputado int) ([]DeputadoMandatoExterno, error) {
	var out []DeputadoMandatoExterno
	path := fmt.Sprintf("/deputados/%d/mandatosExternos", idDeputado)
	if err := c.doGet(ctx, path, nil, &out); err != nil {
		return nil, fmt.Errorf("listar mandatos externos: %w", err)
	}
	return out, nil
}

func (c *DeputadosClient) ListarOrgaos(ctx context.Context, idDeputado int, params map[string]string) ([]DeputadoOrgao, error) {
	query := make(url.Values)
	for k, v := range params {
		query.Set(k, v)
	}
	var out []DeputadoOrgao
	path := fmt.Sprintf("/deputados/%d/orgaos", idDeputado)
	if err := c.doGet(ctx, path, query, &out); err != nil {
		return nil, fmt.Errorf("listar orgaos: %w", err)
	}
	return out, nil
}
