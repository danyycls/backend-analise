package repositorio

import (
	"context"

	"github.com/google/uuid"

	"github.com/danyele/laceu/internal/shared/types"
)

func (r *Repositorio) EleicoesBuscarPorID(ctx context.Context, id uuid.UUID) (*types.Eleicao, error) {
	row := r.db.QueryRow(ctx, scanEleicaoQuery+` WHERE id = $1 AND deleted_at IS NULL`, id)
	return scanEleicaoRow(row)
}

func (r *Repositorio) EleicoesBuscarPorIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]*types.Eleicao, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	rows, err := r.db.Query(ctx, scanEleicaoQuery+` WHERE id = ANY($1) AND deleted_at IS NULL`, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[uuid.UUID]*types.Eleicao)
	for rows.Next() {
		e, err := scanEleicaoRow(rows)
		if err != nil {
			return nil, err
		}
		result[e.ID] = e
	}
	return result, rows.Err()
}

const scanEleicaoQuery = `
	SELECT id, created_at, updated_at, deleted_at,
	       ano, codigo_tse, codigo_tipo_eleicao, nome_tipo_eleicao,
	       descricao, data_eleicao
	FROM eleicao
`

type eleicaoScanner interface {
	Scan(dest ...any) error
}

func scanEleicaoRow(row eleicaoScanner) (*types.Eleicao, error) {
	var e types.Eleicao
	err := row.Scan(
		&e.ID, &e.CreatedAt, &e.UpdatedAt, &e.DeletedAt,
		&e.Ano, &e.CodigoTSE, &e.CodigoTipoEleicao, &e.NomeTipoEleicao,
		&e.Descricao, &e.DataEleicao,
	)
	if err != nil {
		return nil, err
	}
	return &e, nil
}
