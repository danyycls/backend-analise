package repositorio

import (
	"context"

	"github.com/google/uuid"

	"github.com/danyele/laceu/internal/shared/types"
)

func (r *Repositorio) DespesasCandidatoBuscarPorPrestacaoID(ctx context.Context, prestacaoID uuid.UUID) ([]*types.DespesaCandidato, error) {
	rows, err := r.db.Query(ctx, scanDespesaCandidatoQuery+`
		WHERE prestacao_contas_id = $1 AND deleted_at IS NULL
	`, prestacaoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var despesas []*types.DespesaCandidato
	for rows.Next() {
		d, err := scanDespesaCandidatoRow(rows)
		if err != nil {
			return nil, err
		}
		despesas = append(despesas, d)
	}
	return despesas, rows.Err()
}

func (r *Repositorio) DespesasCandidatoBuscarPorFornecedorID(ctx context.Context, fornecedorID uuid.UUID) ([]*types.DespesaCandidato, error) {
	rows, err := r.db.Query(ctx, scanDespesaCandidatoQuery+`
		WHERE fornecedor_id = $1 AND deleted_at IS NULL
	`, fornecedorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var despesas []*types.DespesaCandidato
	for rows.Next() {
		d, err := scanDespesaCandidatoRow(rows)
		if err != nil {
			return nil, err
		}
		despesas = append(despesas, d)
	}
	return despesas, rows.Err()
}

func (r *Repositorio) DespesasCandidatoBuscarPorFornecedorIDComPrestacao(ctx context.Context, fornecedorID string) ([]*types.DespesaCandidato, error) {
	rows, err := r.db.Query(ctx, scanDespesaCandidatoQuery+`
		WHERE fornecedor_id = $1 AND deleted_at IS NULL
	`, fornecedorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var despesas []*types.DespesaCandidato
	for rows.Next() {
		d, err := scanDespesaCandidatoRow(rows)
		if err != nil {
			return nil, err
		}
		despesas = append(despesas, d)
	}
	return despesas, rows.Err()
}

func (r *Repositorio) DespesaCandidatoBuscarPorSQ(ctx context.Context, sq int64) (*types.DespesaCandidato, error) {
	row := r.db.QueryRow(ctx, scanDespesaCandidatoQuery+` WHERE sq_despesa = $1 AND deleted_at IS NULL`, sq)
	return scanDespesaCandidatoRow(row)
}

const scanDespesaCandidatoQuery = `
	SELECT id, created_at, updated_at, deleted_at,
	       prestacao_contas_id, candidato_id, fornecedor_id,
	       sq_despesa, tipo_registro, tipo_documento, numero_documento,
	       origem_despesa_codigo, origem_despesa_descricao,
	       fonte_despesa_codigo, fonte_despesa_descricao,
	       natureza_despesa_codigo, natureza_despesa_descricao,
	       especie_recurso_codigo, especie_recurso_descricao,
	       sq_parcelamento_despesa, data_despesa, descricao, valor
	FROM despesa_candidato
`

type despesaCandidatoScanner interface {
	Scan(dest ...any) error
}

func scanDespesaCandidatoRow(row despesaCandidatoScanner) (*types.DespesaCandidato, error) {
	var d types.DespesaCandidato
	err := row.Scan(
		&d.ID, &d.CreatedAt, &d.UpdatedAt, &d.DeletedAt,
		&d.PrestacaoContasID, &d.CandidatoID, &d.FornecedorID,
		&d.SQDespesa, &d.TipoRegistro, &d.TipoDocumento, &d.NumeroDocumento,
		&d.OrigemDespesaCodigo, &d.OrigemDespesaDescricao,
		&d.FonteDespesaCodigo, &d.FonteDespesaDescricao,
		&d.NaturezaDespesaCodigo, &d.NaturezaDespesaDescricao,
		&d.EspecieRecursoCodigo, &d.EspecieRecursoDescricao,
		&d.SQPlanoParcelamento, &d.DataDespesa, &d.Descricao, &d.Valor,
	)
	if err != nil {
		return nil, err
	}
	return &d, nil
}
