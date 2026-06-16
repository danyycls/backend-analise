package repositorio

import (
	"context"

	"github.com/google/uuid"

	"github.com/danyele/laceu/internal/shared/types"
)

func (r *Repositorio) FornecedoresBuscarPorDocumento(ctx context.Context, documentos []string) ([]*types.Fornecedor, error) {
	rows, err := r.db.Query(ctx, scanFornecedorQuery+` WHERE cpf_cnpj = ANY($1) AND deleted_at IS NULL`, documentos)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fornecedores []*types.Fornecedor
	for rows.Next() {
		var f types.Fornecedor
		if err := rows.Scan(
			&f.ID, &f.CreatedAt, &f.UpdatedAt, &f.DeletedAt,
			&f.CPFCNPJ, &f.Nome, &f.NomeRFB, &f.TipoFornecedorCodigo, &f.TipoFornecedorDescricao,
			&f.CNAECodigo, &f.CNAEDescricao,
			&f.EsferaPartidariaCodigo, &f.EsferaPartidariaDescricao,
			&f.UFSigla, &f.MunicipioNome,
			&f.SQCandidatoRelacionado, &f.NumeroCandidatoRelacionado,
			&f.CargoCodigoRelacionado, &f.CargoDescricaoRelacionada,
			&f.PartidoNumeroRelacionado, &f.PartidoSiglaRelacionado, &f.PartidoNomeRelacionado,
		); err != nil {
			return nil, err
		}
		fornecedores = append(fornecedores, &f)
	}
	return fornecedores, rows.Err()
}

func (r *Repositorio) FornecedoresBuscarPorDocumentoParcialENome(ctx context.Context, pattern string, nome string) ([]*types.Fornecedor, error) {
	rows, err := r.db.Query(ctx, scanFornecedorQuery+`
		WHERE cpf_cnpj LIKE $1 AND nome ILIKE $2 AND deleted_at IS NULL
	`, pattern, "%"+nome+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fornecedores []*types.Fornecedor
	for rows.Next() {
		var f types.Fornecedor
		if err := rows.Scan(
			&f.ID, &f.CreatedAt, &f.UpdatedAt, &f.DeletedAt,
			&f.CPFCNPJ, &f.Nome, &f.NomeRFB, &f.TipoFornecedorCodigo, &f.TipoFornecedorDescricao,
			&f.CNAECodigo, &f.CNAEDescricao,
			&f.EsferaPartidariaCodigo, &f.EsferaPartidariaDescricao,
			&f.UFSigla, &f.MunicipioNome,
			&f.SQCandidatoRelacionado, &f.NumeroCandidatoRelacionado,
			&f.CargoCodigoRelacionado, &f.CargoDescricaoRelacionada,
			&f.PartidoNumeroRelacionado, &f.PartidoSiglaRelacionado, &f.PartidoNomeRelacionado,
		); err != nil {
			return nil, err
		}
		fornecedores = append(fornecedores, &f)
	}
	return fornecedores, rows.Err()
}

func (r *Repositorio) FornecedorBuscarPorCNPJExato(ctx context.Context, cnpj string) (*types.Fornecedor, error) {
	row := r.db.QueryRow(ctx, scanFornecedorQuery+` WHERE cpf_cnpj = $1 AND deleted_at IS NULL`, cnpj)
	var f types.Fornecedor
	err := row.Scan(
		&f.ID, &f.CreatedAt, &f.UpdatedAt, &f.DeletedAt,
		&f.CPFCNPJ, &f.Nome, &f.NomeRFB, &f.TipoFornecedorCodigo, &f.TipoFornecedorDescricao,
		&f.CNAECodigo, &f.CNAEDescricao,
		&f.EsferaPartidariaCodigo, &f.EsferaPartidariaDescricao,
		&f.UFSigla, &f.MunicipioNome,
		&f.SQCandidatoRelacionado, &f.NumeroCandidatoRelacionado,
		&f.CargoCodigoRelacionado, &f.CargoDescricaoRelacionada,
		&f.PartidoNumeroRelacionado, &f.PartidoSiglaRelacionado, &f.PartidoNomeRelacionado,
	)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *Repositorio) FornecedorBuscarPorID(ctx context.Context, id uuid.UUID) (*types.Fornecedor, error) {
	row := r.db.QueryRow(ctx, scanFornecedorQuery+` WHERE id = $1 AND deleted_at IS NULL`, id)
	var f types.Fornecedor
	err := row.Scan(
		&f.ID, &f.CreatedAt, &f.UpdatedAt, &f.DeletedAt,
		&f.CPFCNPJ, &f.Nome, &f.NomeRFB, &f.TipoFornecedorCodigo, &f.TipoFornecedorDescricao,
		&f.CNAECodigo, &f.CNAEDescricao,
		&f.EsferaPartidariaCodigo, &f.EsferaPartidariaDescricao,
		&f.UFSigla, &f.MunicipioNome,
		&f.SQCandidatoRelacionado, &f.NumeroCandidatoRelacionado,
		&f.CargoCodigoRelacionado, &f.CargoDescricaoRelacionada,
		&f.PartidoNumeroRelacionado, &f.PartidoSiglaRelacionado, &f.PartidoNomeRelacionado,
	)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

const scanFornecedorQuery = `
	SELECT id, created_at, updated_at, deleted_at,
	       cpf_cnpj, nome, nome_rfb, tipo_fornecedor_codigo, tipo_fornecedor_descricao,
	       cnae_codigo, cnae_descricao,
	       esfera_partidaria_codigo, esfera_partidaria_descricao,
	       sg_uf, municipio_nome,
	       sq_candidato_relacionado, numero_candidato_relacionado,
	       cargo_codigo_relacionado, cargo_descricao_relacionada,
	       partido_numero_relacionado, partido_sigla_relacionado, partido_nome_relacionado
	FROM fornecedor
`
