package repositorio

import (
	"context"

	"github.com/google/uuid"

	"github.com/danyele/podp/internal/shared/types"
)

func (r *Repositorio) PrestacaoBuscarPorID(ctx context.Context, id uuid.UUID) (*types.PrestacaoContas, error) {
	return r.prestacaoBuscar(ctx, `WHERE id = $1 AND deleted_at IS NULL`, id)
}

func (r *Repositorio) PrestacaoBuscarPorSQ(ctx context.Context, sq int64) (*types.PrestacaoContas, error) {
	return r.prestacaoBuscar(ctx, `WHERE sq_prestador_contas = $1 AND deleted_at IS NULL LIMIT 1`, sq)
}

func (r *Repositorio) prestacaoBuscar(ctx context.Context, where string, args ...any) (*types.PrestacaoContas, error) {
	row := r.db.QueryRow(ctx, scanPrestacaoQuery+` `+where, args...)
	var p types.PrestacaoContas
	err := row.Scan(
		&p.ID, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
		&p.SQPrestadorContas, &p.EleicaoID, &p.CandidatoID, &p.PartidoID,
		&p.UFSigla, &p.UnidadeEleitoralID, &p.TipoPrestador, &p.TipoPrestacao,
		&p.DataPrestacao, &p.Turno, &p.CNPJPrestadorConta,
		&p.EsferaPartidariaCodigo, &p.EsferaPartidariaDescricao,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

const scanPrestacaoQuery = `
	SELECT id, created_at, updated_at, deleted_at,
	       sq_prestador_contas, eleicao_id, candidato_id, partido_id,
	       sg_uf, unidade_eleitoral_id, tipo_prestador, tipo_prestacao,
	       data_prestacao, turno, cnpj_prestador_conta,
	       esfera_partidaria_codigo, esfera_partidaria_descricao
	FROM prestacao_contas
`
