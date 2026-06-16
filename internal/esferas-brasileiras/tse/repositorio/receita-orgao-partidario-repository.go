package repositorio

import (
	"context"

	"github.com/google/uuid"

	"github.com/danyele/laceu/internal/shared/types"
)

func (r *Repositorio) ReceitasOrgaoBuscarPorDoadorID(ctx context.Context, doadorID uuid.UUID) ([]*types.ReceitaOrgaoPartidario, error) {
	rows, err := r.db.Query(ctx, scanReceitaOrgaoQuery+`
		WHERE doador_id = $1 AND deleted_at IS NULL
	`, doadorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var receitas []*types.ReceitaOrgaoPartidario
	for rows.Next() {
		ro, err := scanReceitaOrgaoRow(rows)
		if err != nil {
			return nil, err
		}
		receitas = append(receitas, ro)
	}
	return receitas, rows.Err()
}

func (r *Repositorio) ReceitasOrgaoBuscarPorDoadorIDComPrestacao(ctx context.Context, doadorID string) ([]*types.ReceitaOrgaoPartidario, error) {
	rows, err := r.db.Query(ctx, scanReceitaOrgaoQuery+`
		WHERE doador_id = $1 AND deleted_at IS NULL
	`, doadorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var receitas []*types.ReceitaOrgaoPartidario
	for rows.Next() {
		ro, err := scanReceitaOrgaoRow(rows)
		if err != nil {
			return nil, err
		}
		receitas = append(receitas, ro)
	}
	return receitas, rows.Err()
}

func (r *Repositorio) ReceitaOrgaoBuscarPorSQ(ctx context.Context, sq int64) (*types.ReceitaOrgaoPartidario, error) {
	row := r.db.QueryRow(ctx, scanReceitaOrgaoQuery+` WHERE sq_receita = $1 AND deleted_at IS NULL`, sq)
	return scanReceitaOrgaoRow(row)
}

const scanReceitaOrgaoQuery = `
	SELECT id, created_at, updated_at, deleted_at,
	       prestacao_contas_id, partido_id, doador_id,
	       sq_receita,
	       fonte_receita_codigo, fonte_receita_descricao,
	       origem_receita_codigo, origem_receita_descricao,
	       natureza_receita_codigo, natureza_receita_descricao,
	       especie_receita_codigo, especie_receita_descricao,
	       numero_recibo_doacao, numero_documento_doacao,
	       data_receita, descricao, valor
	FROM receita_orgao_partidario
`

type receitaOrgaoScanner interface {
	Scan(dest ...any) error
}

func scanReceitaOrgaoRow(row receitaOrgaoScanner) (*types.ReceitaOrgaoPartidario, error) {
	var r types.ReceitaOrgaoPartidario
	err := row.Scan(
		&r.ID, &r.CreatedAt, &r.UpdatedAt, &r.DeletedAt,
		&r.PrestacaoContasID, &r.PartidoID, &r.DoadorID,
		&r.SQReceita,
		&r.FonteReceitaCodigo, &r.FonteReceitaDescricao,
		&r.OrigemReceitaCodigo, &r.OrigemReceitaDescricao,
		&r.NaturezaReceitaCodigo, &r.NaturezaReceitaDescricao,
		&r.EspecieReceitaCodigo, &r.EspecieReceitaDescricao,
		&r.NumeroReciboDoacao, &r.NumeroDocumentoDoacao,
		&r.DataReceita, &r.Descricao, &r.Valor,
	)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
