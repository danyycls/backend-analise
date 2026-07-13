package repositorio

import (
	"context"

	"github.com/google/uuid"

	"github.com/danyele/podp/internal/shared/types"
)

func (r *Repositorio) UnidadesEleitoraisBuscarPorID(ctx context.Context, id uuid.UUID) (*types.UnidadeEleitoral, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, created_at, updated_at, deleted_at, sg_uf, codigo_tse, nome
		FROM unidade_eleitoral
		WHERE id = $1 AND deleted_at IS NULL
	`, id)
	var ue types.UnidadeEleitoral
	err := row.Scan(&ue.ID, &ue.CreatedAt, &ue.UpdatedAt, &ue.DeletedAt, &ue.UFSigla, &ue.CodigoTSE, &ue.Nome)
	if err != nil {
		return nil, err
	}
	return &ue, nil
}
