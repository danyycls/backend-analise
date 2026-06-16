package repositorio

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/danyele/laceu/internal/shared/logger"
	"github.com/danyele/laceu/internal/shared/types"
)

type CargoDistinto struct {
	CargoCodigo *int
	CargoNome   string
}

func (r *Repositorio) CargosDistintos(ctx context.Context) ([]CargoDistinto, error) {
	rows, err := r.db.Query(ctx, `
		SELECT DISTINCT cargo_codigo, cargo_nome
		FROM candidato
		WHERE cargo_codigo IS NOT NULL AND cargo_nome != ''
		  AND deleted_at IS NULL
		ORDER BY cargo_nome
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []CargoDistinto
	for rows.Next() {
		var c CargoDistinto
		if err := rows.Scan(&c.CargoCodigo, &c.CargoNome); err != nil {
			return nil, err
		}
		result = append(result, c)
	}
	return result, rows.Err()
}

func (r *Repositorio) CandidatoBuscarPorFiltros(ctx context.Context, cargoNome *string, partidoID *string, ufSigla *string, situacaoTotalizacao *string) ([]*types.Candidato, error) {
	log := logger.New("TSE: Repositorio: CandidatoBuscarPorFiltros")
	var conditions []string
	var args []any
	argIdx := 1

	if cargoNome != nil && *cargoNome != "" {
		conditions = append(conditions, fmt.Sprintf("cargo_nome = $%d", argIdx))
		args = append(args, *cargoNome)
		argIdx++
	}
	if partidoID != nil {
		pid, err := uuid.Parse(*partidoID)
		if err == nil {
			conditions = append(conditions, fmt.Sprintf("partido_id = $%d", argIdx))
			args = append(args, pid)
			argIdx++
		}
	}
	if ufSigla != nil && *ufSigla != "" {
		conditions = append(conditions, fmt.Sprintf("sg_uf = $%d", argIdx))
		args = append(args, *ufSigla)
		argIdx++
	}
	if situacaoTotalizacao != nil && *situacaoTotalizacao != "" {
		conditions = append(conditions, fmt.Sprintf("situacao_totalizacao_descricao = $%d", argIdx))
		args = append(args, *situacaoTotalizacao)
	}

	conditions = append(conditions, "deleted_at IS NULL")

	query := fmt.Sprintf(`
		SELECT id, created_at, updated_at, deleted_at,
		       sq_candidato, eleicao_id, sg_uf, partido_id, cargo_codigo, cargo_nome,
		       genero_descricao, cor_raca_descricao, estado_civil_nome, grau_instrucao_nome,
		       ocupacao_codigo, ocupacao_nome, numero_candidato, cpf, cpf_vice,
		       nome_completo, nome_urna, nome_social, data_nascimento, situacao_totalizacao_descricao
		FROM candidato
		WHERE %s
		ORDER BY sg_uf, cargo_nome, nome_completo
		LIMIT 500
	`, strings.Join(conditions, " AND "))

	log.Info("busca de candidatos por filtros",
		"cargo_nome", cargoNome, "partido_id", partidoID, "uf_sigla", ufSigla, "situacao", situacaoTotalizacao)

	rows, err := r.db.Query(ctx, query, args...)
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
	if err := rows.Err(); err != nil {
		return nil, err
	}

	log.Info("resultado da busca de candidatos",
		"encontrados", len(candidatos))
	return candidatos, nil
}

func (r *Repositorio) CandidatosBuscarPorCPF(ctx context.Context, cpf string) ([]*types.Candidato, error) {
	rows, err := r.db.Query(ctx, scanCandidatoQuery+` WHERE (cpf = $1 OR cpf_vice = $2) AND deleted_at IS NULL`, cpf, cpf)
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
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for _, c := range candidatos {
		bens, err := r.bensBuscarPorCandidatoID(ctx, c.ID)
		if err != nil {
			return nil, err
		}
		c.Bens = bens
	}

	return candidatos, nil
}

func (r *Repositorio) CandidatosBuscarPorCPFParcialENome(ctx context.Context, pattern string, nome string) ([]*types.Candidato, error) {
	rows, err := r.db.Query(ctx, scanCandidatoQuery+`
		WHERE (cpf LIKE $1 OR cpf_vice LIKE $2) AND nome_completo ILIKE $3 AND deleted_at IS NULL
	`, pattern, pattern, "%"+nome+"%")
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
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for _, c := range candidatos {
		bens, err := r.bensBuscarPorCandidatoID(ctx, c.ID)
		if err != nil {
			return nil, err
		}
		c.Bens = bens
	}

	return candidatos, nil
}

func (r *Repositorio) CandidatoBuscarPorSQCandidato(ctx context.Context, sq int64) (*types.Candidato, error) {
	row := r.db.QueryRow(ctx, scanCandidatoQuery+` WHERE sq_candidato = $1 AND deleted_at IS NULL`, sq)
	c, err := scanCandidatoRow(row)
	if err != nil {
		return nil, err
	}

	bens, err := r.bensBuscarPorCandidatoID(ctx, c.ID)
	if err != nil {
		return nil, err
	}
	c.Bens = bens

	return c, nil
}

func (r *Repositorio) CandidatoBuscarPorID(ctx context.Context, id uuid.UUID) (*types.Candidato, error) {
	row := r.db.QueryRow(ctx, scanCandidatoQuery+` WHERE id = $1 AND deleted_at IS NULL`, id)
	c, err := scanCandidatoRow(row)
	if err != nil {
		return nil, err
	}

	bens, err := r.bensBuscarPorCandidatoID(ctx, c.ID)
	if err != nil {
		return nil, err
	}
	c.Bens = bens

	return c, nil
}

const scanCandidatoQuery = `
	SELECT id, created_at, updated_at, deleted_at,
	       sq_candidato, eleicao_id, sg_uf, partido_id, cargo_codigo, cargo_nome,
	       genero_descricao, cor_raca_descricao, estado_civil_nome, grau_instrucao_nome,
	       ocupacao_codigo, ocupacao_nome, numero_candidato, cpf, cpf_vice,
	       nome_completo, nome_urna, nome_social, data_nascimento, situacao_totalizacao_descricao
	FROM candidato
`

func scanCandidatoRow(row pgx.Row) (*types.Candidato, error) {
	var c types.Candidato
	err := row.Scan(
		&c.ID, &c.CreatedAt, &c.UpdatedAt, &c.DeletedAt,
		&c.SQCandidato, &c.EleicaoID, &c.UFSigla, &c.PartidoID, &c.CargoCodigo, &c.CargoNome,
		&c.GeneroDescricao, &c.CorRacaDescricao, &c.EstadoCivilNome, &c.GrauInstrucaoNome,
		&c.OcupacaoCodigo, &c.OcupacaoNome, &c.NumeroCandidato, &c.CPF, &c.CPFVice,
		&c.NomeCompleto, &c.NomeUrna, &c.NomeSocial, &c.DataNascimento, &c.SituacaoTotalizacaoDescricao,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func scanCandidato(rows pgx.Rows) (*types.Candidato, error) {
	return scanCandidatoRow(rows)
}

func (r *Repositorio) bensBuscarPorCandidatoID(ctx context.Context, candidatoID uuid.UUID) ([]types.BemCandidato, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, created_at, updated_at, deleted_at,
		       candidato_id, tipo_bem_codigo, tipo_bem_nome, numero_ordem, descricao,
		       valor, data_ultima_atualizacao, hora_ultima_atualizacao
		FROM bem_candidato
		WHERE candidato_id = $1 AND deleted_at IS NULL
	`, candidatoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bens []types.BemCandidato
	for rows.Next() {
		var b types.BemCandidato
		if err := rows.Scan(
			&b.ID, &b.CreatedAt, &b.UpdatedAt, &b.DeletedAt,
			&b.CandidatoID, &b.TipoBemCodigo, &b.TipoBemNome, &b.NumeroOrdem, &b.Descricao,
			&b.Valor, &b.DataUltimaAtualizacao, &b.HoraUltimaAtualizacao,
		); err != nil {
			return nil, err
		}
		bens = append(bens, b)
	}
	return bens, rows.Err()
}
