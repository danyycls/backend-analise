package senado

import (
	"context"
	"fmt"
)

func (c *SenadoClient) ListarAgendaDia(ctx context.Context, data string, params map[string]string) ([]Reuniao, error) {
	var resultado PlenarioAgendaDia
	path := fmt.Sprintf("/plenario/agenda/dia/%s", data)
	query := makeQueryParams(params)
	if err := c.doGetJSON(ctx, path, query, &resultado); err != nil {
		return nil, fmt.Errorf("listar agenda dia: %w", err)
	}
	return resultado.AgendaDia.Reunioes.Reuniao, nil
}

func (c *SenadoClient) ListarAgendaMes(ctx context.Context, data string, params map[string]string) ([]Reuniao, error) {
	var resultado PlenarioAgendaMes
	path := fmt.Sprintf("/plenario/agenda/mes/%s", data)
	query := makeQueryParams(params)
	if err := c.doGetJSON(ctx, path, query, &resultado); err != nil {
		return nil, fmt.Errorf("listar agenda mes: %w", err)
	}
	return resultado.AgendaMes.Reunioes.Reuniao, nil
}

func (c *SenadoClient) BuscarEncontro(ctx context.Context, codigo string, params map[string]string) (*PlenarioEncontro, error) {
	var resultado PlenarioEncontro
	path := fmt.Sprintf("/plenario/encontro/%s", codigo)
	query := makeQueryParams(params)
	if err := c.doGetJSON(ctx, path, query, &resultado); err != nil {
		return nil, fmt.Errorf("buscar encontro: %w", err)
	}
	return &resultado, nil
}
