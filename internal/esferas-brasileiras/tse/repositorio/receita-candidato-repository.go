package repositorio

import (
	"context"

	"github.com/google/uuid"

	"github.com/danyele/podp/internal/shared/types"
)

func (r *Repositorio) ReceitasCandidatoBuscarPorDoadorID(ctx context.Context, doadorID uuid.UUID) ([]*types.ReceitaCandidato, error) {
	rows, err := r.db.Query(ctx, scanReceitaCandidatoQuery+`
		WHERE doador_id = $1 AND deleted_at IS NULL
	`, doadorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var receitas []*types.ReceitaCandidato
	for rows.Next() {
		rc, err := scanReceitaCandidatoRow(rows)
		if err != nil {
			return nil, err
		}
		receitas = append(receitas, rc)
	}
	return receitas, rows.Err()
}

func (r *Repositorio) ReceitasCandidatoBuscarPorDoadorIDComPrestacao(ctx context.Context, doadorID string) ([]*types.ReceitaCandidato, error) {
	rows, err := r.db.Query(ctx, scanReceitaCandidatoQuery+`
		WHERE doador_id = $1 AND deleted_at IS NULL
	`, doadorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var receitas []*types.ReceitaCandidato
	for rows.Next() {
		rc, err := scanReceitaCandidatoRow(rows)
		if err != nil {
			return nil, err
		}
		receitas = append(receitas, rc)
	}
	return receitas, rows.Err()
}

func (r *Repositorio) ReceitaCandidatoBuscarPorSQ(ctx context.Context, sq int64) (*types.ReceitaCandidato, error) {
	row := r.db.QueryRow(ctx, scanReceitaCandidatoQuery+` WHERE sq_receita = $1 AND deleted_at IS NULL`, sq)
	return scanReceitaCandidatoRow(row)
}

const scanReceitaCandidatoQuery = `
	SELECT id, created_at, updated_at, deleted_at,
	       prestacao_contas_id, candidato_id, doador_id,
	       sq_receita,
	       fonte_receita_codigo, fonte_receita_descricao,
	       origem_receita_codigo, origem_receita_descricao,
	       natureza_receita_codigo, natureza_receita_descricao,
	       especie_receita_codigo, especie_receita_descricao,
	       numero_recibo_doacao, numero_documento_doacao,
	       data_receita, descricao, valor,
	       natureza_recurso_estimavel, genero, cor_raca
	FROM receita_candidato
`

type receitaCandidatoScanner interface {
	Scan(dest ...any) error
}

func scanReceitaCandidatoRow(row receitaCandidatoScanner) (*types.ReceitaCandidato, error) {
	var r types.ReceitaCandidato
	err := row.Scan(
		&r.ID, &r.CreatedAt, &r.UpdatedAt, &r.DeletedAt,
		&r.PrestacaoContasID, &r.CandidatoID, &r.DoadorID,
		&r.SQReceita,
		&r.FonteReceitaCodigo, &r.FonteReceitaDescricao,
		&r.OrigemReceitaCodigo, &r.OrigemReceitaDescricao,
		&r.NaturezaReceitaCodigo, &r.NaturezaReceitaDescricao,
		&r.EspecieReceitaCodigo, &r.EspecieReceitaDescricao,
		&r.NumeroReciboDoacao, &r.NumeroDocumentoDoacao,
		&r.DataReceita, &r.Descricao, &r.Valor,
		&r.NaturezaRecursoEstimavel, &r.Genero, &r.CorRaca,
	)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
