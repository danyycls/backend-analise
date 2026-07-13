package deputados

import (
	"context"
	"fmt"
)

func (c *DeputadosClient) ListarReferencias(ctx context.Context, tipo string) ([]Referencia, error) {
	var out []Referencia
	if err := c.doGet(ctx, "/referencias/"+tipo, nil, &out); err != nil {
		return nil, fmt.Errorf("listar referencias %s: %w", tipo, err)
	}
	return out, nil
}
