package senado

import (
	"context"
	"fmt"
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
