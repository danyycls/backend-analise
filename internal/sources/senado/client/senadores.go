package senado

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

func (c *SenadoClient) ListarSenadores(ctx context.Context) ([]ParlamentarResumo, error) {
	var resultado ListaParlamentarEmExercicio
	if err := c.doGetJSON(ctx, "/senador/lista/atual", nil, &resultado); err != nil {
		return nil, fmt.Errorf("listar senadores: %w", err)
	}
	return resultado.ListaParlamentarEmExercicio.Parlamentares.Parlamentar, nil
}

func (c *SenadoClient) BuscarSenador(ctx context.Context, codigo string) (*ParlamentarDetalhe, error) {
	var resultado DetalheParlamentar
	path := fmt.Sprintf("/senador/%s", codigo)
	if err := c.doGetJSON(ctx, path, nil, &resultado); err != nil {
		return nil, fmt.Errorf("buscar senador: %w", err)
	}
	return &resultado.DetalheParlamentar.Parlamentar, nil
}

func (c *SenadoClient) ListarCargos(ctx context.Context, codigo string) ([]Cargo, error) {
	var resultado CargoParlamentar
	path := fmt.Sprintf("/senador/%s/cargos", codigo)
	if err := c.doGetJSON(ctx, path, nil, &resultado); err != nil {
		return nil, fmt.Errorf("listar cargos: %w", err)
	}
	return resultado.CargoParlamentar.Parlamentar.Cargos.Cargo, nil
}

func (c *SenadoClient) ListarComissoes(ctx context.Context, codigo string) ([]ComissaoMembro, error) {
	var resultado MembroComissaoParlamentar
	path := fmt.Sprintf("/senador/%s/comissoes", codigo)
	if err := c.doGetJSON(ctx, path, nil, &resultado); err != nil {
		return nil, fmt.Errorf("listar comissoes: %w", err)
	}
	return resultado.MembroComissaoParlamentar.Parlamentar.MembroComissoes.Comissao, nil
}

func (c *SenadoClient) ListarMandatos(ctx context.Context, codigo string) ([]MandatoDetalhe, error) {
	var resultado MandatoParlamentar
	path := fmt.Sprintf("/senador/%s/mandatos", codigo)
	if err := c.doGetJSON(ctx, path, nil, &resultado); err != nil {
		return nil, fmt.Errorf("listar mandatos: %w", err)
	}
	return resultado.MandatoParlamentar.Parlamentar.Mandatos.Mandato, nil
}

func (c *SenadoClient) BaixarDocumentoEmenda(ctx context.Context, idDocumento int) ([]byte, string, error) {
	docURL := fmt.Sprintf("https://legis.senado.leg.br/sdleg-getter/documento?dm=%d", idDocumento)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, docURL, nil)
	if err != nil {
		return nil, "", fmt.Errorf("baixar documento: erro requisicao: %w", err)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("baixar documento: erro execucao: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, "", fmt.Errorf("baixar documento: status %d: %s", resp.StatusCode, string(body))
	}
	contentType := resp.Header.Get("Content-Type")
	dados, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("baixar documento: erro leitura: %w", err)
	}
	return dados, contentType, nil
}
