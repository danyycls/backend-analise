package senado

import (
	"context"
	"fmt"
)

func (c *SenadoClient) ListarOrcamento(ctx context.Context) ([]LoteEmendasOrcamento, error) {
	var resultado OrcamentoLista
	if err := c.doGetJSON(ctx, "/orcamento/lista", nil, &resultado); err != nil {
		return nil, fmt.Errorf("listar orcamento: %w", err)
	}
	return resultado.ListaLoteEmendas.LotesEmendasOrcamento.LoteEmendasOrcamento, nil
}
