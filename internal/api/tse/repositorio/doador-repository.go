package repositorio

import (
	"context"

	"github.com/danyele/podp/internal/shared/types"
)

func (r *Repositorio) DoadoresBuscarPorDocumento(ctx context.Context, documentos []string) ([]*types.Doador, error) {
	rows, err := r.db.Query(ctx, scanDoadorQuery+` WHERE cpf_cnpj = ANY($1) AND deleted_at IS NULL`, documentos)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var doadores []*types.Doador
	for rows.Next() {
		var d types.Doador
		if err := rows.Scan(
			&d.ID, &d.CreatedAt, &d.UpdatedAt, &d.DeletedAt,
			&d.CPFCNPJ, &d.Nome, &d.NomeRFB, &d.CNAECodigo, &d.CNAEDescricao,
			&d.EsferaPartidariaCodigo, &d.EsferaPartidariaDescricao,
			&d.UFSigla, &d.MunicipioNome,
			&d.SQCandidatoRelacionado, &d.NumeroCandidatoRelacionado,
			&d.CargoCodigoRelacionado, &d.CargoDescricaoRelacionada,
			&d.PartidoNumeroRelacionado, &d.PartidoSiglaRelacionado, &d.PartidoNomeRelacionado,
		); err != nil {
			return nil, err
		}
		doadores = append(doadores, &d)
	}
	return doadores, rows.Err()
}

func (r *Repositorio) DoadoresBuscarPorDocumentoParcial(ctx context.Context, pattern string, nome string) ([]*types.Doador, error) {
	rows, err := r.db.Query(ctx, scanDoadorQuery+`
		WHERE cpf_cnpj LIKE $1 AND nome ILIKE $2 AND deleted_at IS NULL
	`, pattern, "%"+nome+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var doadores []*types.Doador
	for rows.Next() {
		var d types.Doador
		if err := rows.Scan(
			&d.ID, &d.CreatedAt, &d.UpdatedAt, &d.DeletedAt,
			&d.CPFCNPJ, &d.Nome, &d.NomeRFB, &d.CNAECodigo, &d.CNAEDescricao,
			&d.EsferaPartidariaCodigo, &d.EsferaPartidariaDescricao,
			&d.UFSigla, &d.MunicipioNome,
			&d.SQCandidatoRelacionado, &d.NumeroCandidatoRelacionado,
			&d.CargoCodigoRelacionado, &d.CargoDescricaoRelacionada,
			&d.PartidoNumeroRelacionado, &d.PartidoSiglaRelacionado, &d.PartidoNomeRelacionado,
		); err != nil {
			return nil, err
		}
		doadores = append(doadores, &d)
	}
	return doadores, rows.Err()
}

func (r *Repositorio) DoadorBuscarPorCNPJExato(ctx context.Context, cnpj string) (*types.Doador, error) {
	row := r.db.QueryRow(ctx, scanDoadorQuery+` WHERE cpf_cnpj = $1 AND deleted_at IS NULL`, cnpj)
	var d types.Doador
	err := row.Scan(
		&d.ID, &d.CreatedAt, &d.UpdatedAt, &d.DeletedAt,
		&d.CPFCNPJ, &d.Nome, &d.NomeRFB, &d.CNAECodigo, &d.CNAEDescricao,
		&d.EsferaPartidariaCodigo, &d.EsferaPartidariaDescricao,
		&d.UFSigla, &d.MunicipioNome,
		&d.SQCandidatoRelacionado, &d.NumeroCandidatoRelacionado,
		&d.CargoCodigoRelacionado, &d.CargoDescricaoRelacionada,
		&d.PartidoNumeroRelacionado, &d.PartidoSiglaRelacionado, &d.PartidoNomeRelacionado,
	)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *Repositorio) DoadorBuscarPorID(ctx context.Context, id interface{}) (*types.Doador, error) {
	row := r.db.QueryRow(ctx, scanDoadorQuery+` WHERE id = $1 AND deleted_at IS NULL`, id)
	var d types.Doador
	err := row.Scan(
		&d.ID, &d.CreatedAt, &d.UpdatedAt, &d.DeletedAt,
		&d.CPFCNPJ, &d.Nome, &d.NomeRFB, &d.CNAECodigo, &d.CNAEDescricao,
		&d.EsferaPartidariaCodigo, &d.EsferaPartidariaDescricao,
		&d.UFSigla, &d.MunicipioNome,
		&d.SQCandidatoRelacionado, &d.NumeroCandidatoRelacionado,
		&d.CargoCodigoRelacionado, &d.CargoDescricaoRelacionada,
		&d.PartidoNumeroRelacionado, &d.PartidoSiglaRelacionado, &d.PartidoNomeRelacionado,
	)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

const scanDoadorQuery = `
	SELECT id, created_at, updated_at, deleted_at,
	       cpf_cnpj, nome, nome_rfb, cnae_codigo, cnae_descricao,
	       esfera_partidaria_codigo, esfera_partidaria_descricao,
	       sg_uf, municipio_nome,
	       sq_candidato_relacionado, numero_candidato_relacionado,
	       cargo_codigo_relacionado, cargo_descricao_relacionada,
	       partido_numero_relacionado, partido_sigla_relacionado, partido_nome_relacionado
	FROM doador
`
