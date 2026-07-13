package repositorio

import (
	"context"

	"github.com/google/uuid"

	"github.com/danyele/podp/internal/shared/types"
)

func (r *Repositorio) PartidosBuscarPorID(ctx context.Context, id uuid.UUID) (*types.Partido, error) {
	row := r.db.QueryRow(ctx, scanPartidoQuery+` WHERE id = $1 AND deleted_at IS NULL`, id)
	return scanPartidoRow(row)
}

func (r *Repositorio) PartidosBuscarPorIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]*types.Partido, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	rows, err := r.db.Query(ctx, scanPartidoQuery+` WHERE id = ANY($1) AND deleted_at IS NULL`, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[uuid.UUID]*types.Partido)
	for rows.Next() {
		p, err := scanPartidoRow(rows)
		if err != nil {
			return nil, err
		}
		result[p.ID] = p
	}
	return result, rows.Err()
}

func (r *Repositorio) PartidosListarDistintos(ctx context.Context) ([]*types.Partido, error) {
	rows, err := r.db.Query(ctx, `
		SELECT DISTINCT ON (numero, sigla, nome) id, created_at, updated_at, deleted_at,
		       numero, sigla, nome,
		       federacao_codigo_tse, federacao_sigla, federacao_nome,
		       coligacao_codigo_tse, coligacao_nome, coligacao_composicao
		FROM partido
		WHERE numero IS NOT NULL AND sigla != '' AND deleted_at IS NULL
		ORDER BY numero
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var partidos []*types.Partido
	for rows.Next() {
		p, err := scanPartidoRow(rows)
		if err != nil {
			return nil, err
		}
		partidos = append(partidos, p)
	}
	return partidos, rows.Err()
}

const scanPartidoQuery = `
	SELECT id, created_at, updated_at, deleted_at,
	       numero, sigla, nome,
	       federacao_codigo_tse, federacao_sigla, federacao_nome,
	       coligacao_codigo_tse, coligacao_nome, coligacao_composicao
	FROM partido
`

type partidoScanner interface {
	Scan(dest ...any) error
}

func scanPartidoRow(row partidoScanner) (*types.Partido, error) {
	var p types.Partido
	err := row.Scan(
		&p.ID, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
		&p.Numero, &p.Sigla, &p.Nome,
		&p.FederacaoCodigoTSE, &p.FederacaoSigla, &p.FederacaoNome,
		&p.ColigacaoCodigoTSE, &p.ColigacaoNome, &p.ColigacaoComposicao,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
