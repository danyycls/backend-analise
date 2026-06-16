package repositorio

import (
	"context"

	"github.com/google/uuid"

	"github.com/danyele/laceu/internal/shared/types"
)

func (r *Repositorio) CandidatosEleitosPorUF(ctx context.Context, ufSigla string, cargos []string) ([]*types.Candidato, error) {
	rows, err := r.db.Query(ctx, scanCandidatoQuery+`
		WHERE sg_uf = $1
		  AND cargo_nome = ANY($2)
		  AND situacao_totalizacao_descricao = ANY($3)
		  AND deleted_at IS NULL
		ORDER BY cargo_nome, nome_urna
	`, ufSigla, cargos, []string{"SUPLENTE", "ELEITO POR MÉDIA", "ELEITO POR QP", "ELEITO"})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var candidatos []*types.Candidato
	for rows.Next() {
		c, err := scanCandidato(rows)
		if err != nil {
			return nil, err
		}
		candidatos = append(candidatos, c)
	}
	return candidatos, rows.Err()
}

func (r *Repositorio) PartidoBuscarPorID(ctx context.Context, id uuid.UUID) (*types.Partido, error) {
	return r.PartidosBuscarPorID(ctx, id)
}

func (r *Repositorio) EleicaoBuscarPorID(ctx context.Context, id uuid.UUID) (*types.Eleicao, error) {
	return r.EleicoesBuscarPorID(ctx, id)
}
