package senado

import (
	"context"
	"fmt"
)

func (c *SenadoClient) ListarTodasComissoes(ctx context.Context) ([]ComissaoResumo, error) {
	var resultado ColegiadoListaWrapper
	if err := c.doGetJSON(ctx, "/comissao/lista/colegiados", nil, &resultado); err != nil {
		return nil, fmt.Errorf("listar todas comissoes: %w", err)
	}
	return resultado.ListaColegiados.Colegiados.Colegiado, nil
}

func (c *SenadoClient) BuscarComissao(ctx context.Context, codigo string) (*ComissaoDetalhe, error) {
	var resultado ComissaoDetalheWrapper
	path := fmt.Sprintf("/comissao/%s", codigo)
	if err := c.doGetJSON(ctx, path, nil, &resultado); err != nil {
		return nil, fmt.Errorf("buscar comissao: %w", err)
	}

	raw := resultado.ComissoesCongressoNacional.Colegiados.Colegiado
	if len(raw) == 0 {
		return nil, fmt.Errorf("comissao nao encontrada: %s", codigo)
	}
	col := raw[0]

	detalhe := &ComissaoDetalhe{
		Codigo: col.CodigoColegiado,
		Sigla:  col.SiglaColegiado,
		Nome:   col.NomeColegiado,
	}

	if col.MembrosBlocoSF != nil {
		for _, bloco := range col.MembrosBlocoSF.PartidoBloco {
			if bloco.MembrosSF == nil {
				continue
			}
			for _, m := range bloco.MembrosSF.Membro {
				detalhe.Membros = append(detalhe.Membros, MembroComissao{
					CodigoParlamentar:       m.CodigoParlamentar,
					NomeParlamentar:         m.NomeParlamentar,
					SiglaPartidoParlamentar: m.Partido,
					UfParlamentar:           m.SiglaUf,
					DescricaoParticipacao:   m.TipoVaga,
				})
			}
		}
	}

	return detalhe, nil
}
