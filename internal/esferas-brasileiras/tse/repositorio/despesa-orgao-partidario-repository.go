package repositorio

import (
	"context"

	"github.com/google/uuid"

	"github.com/danyele/laceu/internal/shared/types"
)

func (r *Repositorio) DespesasPartidoBuscarPorPrestacaoID(ctx context.Context, prestacaoID uuid.UUID) ([]*types.DespesaOrgaoPartidario, error) {
	rows, err := r.db.Query(ctx, scanDespesaPartidoQuery+`
		WHERE prestacao_contas_id = $1 AND deleted_at IS NULL
	`, prestacaoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var despesas []*types.DespesaOrgaoPartidario
	for rows.Next() {
		d, err := scanDespesaPartidoRow(rows)
		if err != nil {
			return nil, err
		}
		despesas = append(despesas, d)
	}
	return despesas, rows.Err()
}

func (r *Repositorio) DespesasPartidoBuscarPorFornecedorID(ctx context.Context, fornecedorID uuid.UUID) ([]*types.DespesaOrgaoPartidario, error) {
	rows, err := r.db.Query(ctx, scanDespesaPartidoQuery+`
		WHERE fornecedor_id = $1 AND deleted_at IS NULL
	`, fornecedorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var despesas []*types.DespesaOrgaoPartidario
	for rows.Next() {
		d, err := scanDespesaPartidoRow(rows)
		if err != nil {
			return nil, err
		}
		despesas = append(despesas, d)
	}
	return despesas, rows.Err()
}

func (r *Repositorio) DespesasPartidoBuscarPorFornecedorIDComPrestacao(ctx context.Context, fornecedorID string) ([]*types.DespesaOrgaoPartidario, error) {
	rows, err := r.db.Query(ctx, scanDespesaPartidoQuery+`
		WHERE fornecedor_id = $1 AND deleted_at IS NULL
	`, fornecedorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var despesas []*types.DespesaOrgaoPartidario
	for rows.Next() {
		d, err := scanDespesaPartidoRow(rows)
		if err != nil {
			return nil, err
		}
		despesas = append(despesas, d)
	}
	return despesas, rows.Err()
}

func (r *Repositorio) DespesaOrgaoBuscarPorSQ(ctx context.Context, sq int64) (*types.DespesaOrgaoPartidario, error) {
	row := r.db.QueryRow(ctx, scanDespesaPartidoQuery+` WHERE sq_despesa = $1 AND deleted_at IS NULL`, sq)
	return scanDespesaPartidoRow(row)
}

const scanDespesaPartidoQuery = `
	SELECT id, created_at, updated_at, deleted_at,
	       prestacao_contas_id, partido_id, fornecedor_id,
	       sq_despesa, tipo_registro, tipo_documento, numero_documento,
	       origem_despesa_codigo, origem_despesa_descricao,
	       fonte_despesa_codigo, fonte_despesa_descricao,
	       natureza_despesa_codigo, natureza_despesa_descricao,
	       especie_recurso_codigo, especie_recurso_descricao,
	       sq_parcelamento_despesa, data_despesa, descricao, valor
	FROM despesa_orgao_partidario
`

type despesaPartidoScanner interface {
	Scan(dest ...any) error
}

func scanDespesaPartidoRow(row despesaPartidoScanner) (*types.DespesaOrgaoPartidario, error) {
	var d types.DespesaOrgaoPartidario
	err := row.Scan(
		&d.ID, &d.CreatedAt, &d.UpdatedAt, &d.DeletedAt,
		&d.PrestacaoContasID, &d.PartidoID, &d.FornecedorID,
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
