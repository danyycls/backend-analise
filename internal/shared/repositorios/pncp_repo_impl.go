package repositorios

import (
	"context"
	"fmt"
	"strings"

	"github.com/danyele/podp/internal/shared/database"
	"github.com/jackc/pgx/v5"
)

const colunasContrato = `
	numero_controle_pncp, cnpj_orgao, ug_uf_sigla, ug_codigo_ibge,
	data_publicacao_pncp, data_assinatura, data_inicio_vigencia, data_termino_vigencia,
	valor_global, valor_inicial, valor_total_estimado, valor_total_homologado,
	ni_fornecedor, codigo_amparo_legal,
	numero_contrato, codigo_contrato, codigo_tipo_contrato, tipo_contrato_nome,
	codigo_ug, nome_ug, ug_municipio_nome, ug_uf_nome,
	modalidade_nome, codigo_orgao, nome_orgao, nome_orgao_sub,
	objeto_contrato, numero_licitacao, origem_licitacao, produto, subtipo_contrato,
	ano_contrato, nome_razao_social_fornecedor, dados_completos
`

const numColunasContrato = 34
const chunkSize = 100

type pncpRepositoryImpl struct {
	db database.DB
}

func NovoPNCPRepository(db database.DB) PNCPRepository {
	return &pncpRepositoryImpl{db: db}
}

func (r *pncpRepositoryImpl) SalvarContratos(ctx context.Context, contratos []ContratoPersistido) error {
	if len(contratos) == 0 {
		return nil
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	for i := 0; i < len(contratos); i += chunkSize {
		fim := i + chunkSize
		if fim > len(contratos) {
			fim = len(contratos)
		}
		chunk := contratos[i:fim]

		if err := r.salvarContratosBatch(ctx, tx, chunk); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *pncpRepositoryImpl) salvarContratosBatch(ctx context.Context, tx pgx.Tx, contratos []ContratoPersistido) error {
	n := len(contratos)
	placeholders := make([]string, n)
	args := make([]any, 0, n*numColunasContrato)

	for j, c := range contratos {
		base := j * numColunasContrato
		placeholders[j] = fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			base+1, base+2, base+3, base+4, base+5, base+6, base+7, base+8,
			base+9, base+10, base+11, base+12, base+13, base+14,
			base+15, base+16, base+17, base+18, base+19, base+20,
			base+21, base+22, base+23, base+24, base+25, base+26,
			base+27, base+28, base+29, base+30, base+31, base+32, base+33, base+34,
		)
		args = append(args,
			c.NumeroControlePNCP, c.CNPJOrgao, nullStr(c.UGUFSigla), nullStr(c.UGCodigoIbge),
			c.DataPublicacaoPncp, c.DataAssinatura, c.DataInicioVigencia, c.DataTerminoVigencia,
			c.ValorGlobal, c.ValorInicial, c.ValorTotalEstimado, c.ValorTotalHomologado,
			c.NIFornecedor, c.CodigoAmparoLegal,
			c.NumeroContrato, c.CodigoContrato, c.CodigoTipoContrato, c.TipoContratoNome,
			c.CodigoUG, c.NomeUG, c.UGMunicipioNome, c.UGUFNome,
			c.ModalidadeNome, c.CodigoOrgao, c.NomeOrgao, c.NomeOrgaoSub,
			c.ObjetoContrato, c.NumeroLicitacao, c.OrigemLicitacao, c.Produto, c.SubtipoContrato,
			c.AnoContrato, c.NomeRazaoSocialFornecedor, c.DadosCompletos,
		)
	}

	query := fmt.Sprintf(`
		INSERT INTO licitacao_contrato (%s)
		VALUES %s
		ON CONFLICT (numero_controle_pncp) DO UPDATE SET updated_at = NOW()
	`, colunasContrato, strings.Join(placeholders, ", "))

	_, err := tx.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("salvar contratos batch [%d]: %w", n, err)
	}
	return nil
}

func (r *pncpRepositoryImpl) SalvarFornecedores(ctx context.Context, fornecedores []FornecedorPersistido) error {
	if len(fornecedores) == 0 {
		return nil
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	for _, f := range fornecedores {
		_, err := tx.Exec(ctx, `
			INSERT INTO licitacao_fornecedor (cnpj, razao_social, dados_completos)
			VALUES ($1, $2, $3)
			ON CONFLICT (cnpj) DO UPDATE SET
				razao_social = EXCLUDED.razao_social,
				dados_completos = EXCLUDED.dados_completos
		`, f.CNPJ, f.RazaoSocial, f.DadosCompletos)
		if err != nil {
			return fmt.Errorf("salvar fornecedor %s: %w", f.CNPJ, err)
		}
	}

	return tx.Commit(ctx)
}

func (r *pncpRepositoryImpl) SalvarSocio(ctx context.Context, socio SocioPersistido) (string, error) {
	row := r.db.QueryRow(ctx, `
		INSERT INTO socio (cnpj_cpf_socio, nome_socio)
		VALUES ($1, $2)
		ON CONFLICT (cnpj_cpf_socio) DO UPDATE SET nome_socio = EXCLUDED.nome_socio
		RETURNING id
	`, socio.CNPJCPFSocio, socio.NomeSocio)

	var id string
	err := row.Scan(&id)
	if err != nil {
		return "", fmt.Errorf("salvar socio: %w", err)
	}
	return id, nil
}

func (r *pncpRepositoryImpl) SalvarFornecedorSocios(ctx context.Context, vinculos []FornecedorSocioPersistido) error {
	if len(vinculos) == 0 {
		return nil
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	fornecedores := make(map[string]bool)
	for _, v := range vinculos {
		fornecedores[v.CNPJFornecedor] = true
	}

	cnpjs := make([]string, 0, len(fornecedores))
	for cnpj := range fornecedores {
		cnpjs = append(cnpjs, cnpj)
	}

	_, err = tx.Exec(ctx, `DELETE FROM fornecedor_socio WHERE cnpj_fornecedor = ANY($1)`, cnpjs)
	if err != nil {
		return fmt.Errorf("limpar socios dos fornecedores: %w", err)
	}

	for _, v := range vinculos {
		_, err := tx.Exec(ctx, `
			INSERT INTO fornecedor_socio (
				cnpj_fornecedor, socio_id,
				data_entrada_sociedade, identificador_socio, nome_socio,
				qualificacao_socio, nome_representante, qualificacao_representante,
				representante_legal, faixa_etaria, pais_codigo, pais_descricao
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		`,
			v.CNPJFornecedor, v.SocioID,
			v.DataEntradaSociedade, v.IdentificadorSocio, v.NomeSocio,
			v.QualificacaoSocio, v.NomeRepresentante, v.QualificacaoRepresentante,
			v.RepresentanteLegal, v.FaixaEtaria, v.PaisCodigo, v.PaisDescricao,
		)
		if err != nil {
			return fmt.Errorf("vincular socio ao fornecedor: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (r *pncpRepositoryImpl) SalvarAmparosLegais(ctx context.Context, amparos []AmparoLegalPersistido) error {
	if len(amparos) == 0 {
		return nil
	}
	for _, a := range amparos {
		_, err := r.db.Exec(ctx, `
			INSERT INTO amparo_legal (codigo, nome, descricao)
			VALUES ($1, $2, $3)
			ON CONFLICT (codigo) DO NOTHING
		`, a.Codigo, a.Nome, a.Descricao)
		if err != nil {
			return fmt.Errorf("salvar amparo legal %d: %w", a.Codigo, err)
		}
	}
	return nil
}

func (r *pncpRepositoryImpl) BuscarContratosPorFiltro(ctx context.Context, tipo, valor string, ano, mes int) ([]ContratoPersistido, error) {
	where, args := montarFiltroMes(tipo, valor, ano, mes)
	return r.buscarContratos(ctx, where, args...)
}

func (r *pncpRepositoryImpl) BuscarContratosPorFiltroEPeriodo(ctx context.Context, tipo, valor, dataInicial, dataFinal string) ([]ContratoPersistido, error) {
	where, args := montarFiltroPeriodo(tipo, valor, dataInicial, dataFinal)
	return r.buscarContratos(ctx, where, args...)
}

func (r *pncpRepositoryImpl) BuscarContratoPorNumeroControle(ctx context.Context, numeroControle string) (*ContratoPersistido, error) {
	rows, err := r.buscarContratos(ctx, "WHERE numero_controle_pncp = $1", numeroControle)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil
	}
	return &rows[0], nil
}

func (r *pncpRepositoryImpl) buscarContratos(ctx context.Context, where string, args ...any) ([]ContratoPersistido, error) {
	query := fmt.Sprintf(`
		SELECT id, %s
		FROM licitacao_contrato
		%s
	`, colunasContrato, where)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("buscar contratos: %w", err)
	}
	defer rows.Close()

	var resultados []ContratoPersistido
	for rows.Next() {
		var c ContratoPersistido
		err := rows.Scan(
			&c.ID, &c.NumeroControlePNCP, &c.CNPJOrgao, &c.UGUFSigla, &c.UGCodigoIbge,
			&c.DataPublicacaoPncp, &c.DataAssinatura, &c.DataInicioVigencia, &c.DataTerminoVigencia,
			&c.ValorGlobal, &c.ValorInicial, &c.ValorTotalEstimado, &c.ValorTotalHomologado,
			&c.NIFornecedor, &c.CodigoAmparoLegal,
			&c.NumeroContrato, &c.CodigoContrato, &c.CodigoTipoContrato, &c.TipoContratoNome,
			&c.CodigoUG, &c.NomeUG, &c.UGMunicipioNome, &c.UGUFNome,
			&c.ModalidadeNome, &c.CodigoOrgao, &c.NomeOrgao, &c.NomeOrgaoSub,
			&c.ObjetoContrato, &c.NumeroLicitacao, &c.OrigemLicitacao, &c.Produto, &c.SubtipoContrato,
			&c.AnoContrato, &c.NomeRazaoSocialFornecedor, &c.DadosCompletos,
		)
		if err != nil {
			return nil, fmt.Errorf("scan contrato: %w", err)
		}
		resultados = append(resultados, c)
	}

	return resultados, nil
}

func (r *pncpRepositoryImpl) BuscarFornecedor(ctx context.Context, cnpj string) (*FornecedorPersistido, error) {
	row := r.db.QueryRow(ctx, `
		SELECT cnpj, razao_social, dados_completos
		FROM licitacao_fornecedor
		WHERE cnpj = $1
	`, cnpj)

	var f FornecedorPersistido
	err := row.Scan(&f.CNPJ, &f.RazaoSocial, &f.DadosCompletos)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("buscar fornecedor: %w", err)
	}
	return &f, nil
}

func (r *pncpRepositoryImpl) BuscarSociosPorFornecedor(ctx context.Context, cnpj string) ([]FornecedorSocioPersistido, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			fs.cnpj_fornecedor, fs.socio_id,
			fs.data_entrada_sociedade, fs.identificador_socio, fs.nome_socio,
			fs.qualificacao_socio, fs.nome_representante, fs.qualificacao_representante,
			fs.representante_legal, fs.faixa_etaria, fs.pais_codigo, fs.pais_descricao,
			s.cnpj_cpf_socio, s.nome_socio
		FROM fornecedor_socio fs
		JOIN socio s ON s.id = fs.socio_id
		WHERE fs.cnpj_fornecedor = $1
	`, cnpj)
	if err != nil {
		return nil, fmt.Errorf("buscar socios do fornecedor: %w", err)
	}
	defer rows.Close()

	return scanFornecedorSocios(rows)
}

func (r *pncpRepositoryImpl) BuscarFornecedoresPorSocio(ctx context.Context, cnpjCpfSocio string) ([]FornecedorSocioPersistido, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			fs.cnpj_fornecedor, fs.socio_id,
			fs.data_entrada_sociedade, fs.identificador_socio, fs.nome_socio,
			fs.qualificacao_socio, fs.nome_representante, fs.qualificacao_representante,
			fs.representante_legal, fs.faixa_etaria, fs.pais_codigo, fs.pais_descricao,
			s.cnpj_cpf_socio, s.nome_socio
		FROM fornecedor_socio fs
		JOIN socio s ON s.id = fs.socio_id
		WHERE s.cnpj_cpf_socio = $1
	`, cnpjCpfSocio)
	if err != nil {
		return nil, fmt.Errorf("buscar fornecedores do socio: %w", err)
	}
	defer rows.Close()

	return scanFornecedorSocios(rows)
}

func (r *pncpRepositoryImpl) BuscaJaRealizada(ctx context.Context, tipo, valor string, ano, mes int) (bool, error) {
	var total int
	err := r.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM licitacao_busca_controle
		WHERE tipo_busca = $1 AND valor_busca = $2 AND ano = $3 AND mes = $4
	`, tipo, valor, ano, mes).Scan(&total)
	if err != nil {
		return false, fmt.Errorf("verificar busca ja realizada: %w", err)
	}
	return total > 0, nil
}

func (r *pncpRepositoryImpl) RegistrarBusca(ctx context.Context, controle BuscaControlePersistido) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO licitacao_busca_controle (tipo_busca, valor_busca, ano, mes, data_inicial, data_final, total_contratos_encontrados)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (tipo_busca, valor_busca, ano, mes) DO NOTHING
	`, controle.TipoBusca, controle.ValorBusca, controle.Ano, controle.Mes,
		controle.DataInicial, controle.DataFinal, controle.TotalContratosEncontrados)
	return err
}

func (r *pncpRepositoryImpl) AtualizarBusca(ctx context.Context, controle BuscaControlePersistido) error {
	_, err := r.db.Exec(ctx, `
		UPDATE licitacao_busca_controle
		SET total_contratos_encontrados = $1, ultima_atualizacao = NOW()
		WHERE tipo_busca = $2 AND valor_busca = $3 AND ano = $4 AND mes = $5
	`, controle.TotalContratosEncontrados, controle.TipoBusca, controle.ValorBusca, controle.Ano, controle.Mes)
	return err
}

func montarFiltroMes(tipo, valor string, ano, mes int) (string, []any) {
	coluna := colunaPorTipo(tipo)
	where := fmt.Sprintf(
		"WHERE %s = $1 AND data_publicacao_pncp >= make_date($2, $3, 1) AND data_publicacao_pncp < make_date($2, $3, 1) + INTERVAL '1 month'",
		coluna,
	)
	return where, []any{valor, ano, mes}
}

func montarFiltroPeriodo(tipo, valor, dataInicial, dataFinal string) (string, []any) {
	coluna := colunaPorTipo(tipo)
	where := fmt.Sprintf("WHERE %s = $1 AND data_publicacao_pncp BETWEEN $2 AND $3", coluna)
	return where, []any{valor, dataInicial, dataFinal}
}

func colunaPorTipo(tipo string) string {
	switch tipo {
	case "uf":
		return "ug_uf_sigla"
	case "municipio":
		return "ug_codigo_ibge"
	case "orgao":
		return "cnpj_orgao"
	default:
		return "cnpj_orgao"
	}
}

func scanFornecedorSocios(rows pgx.Rows) ([]FornecedorSocioPersistido, error) {
	var resultados []FornecedorSocioPersistido
	for rows.Next() {
		var v FornecedorSocioPersistido
		err := rows.Scan(
			&v.CNPJFornecedor, &v.SocioID,
			&v.DataEntradaSociedade, &v.IdentificadorSocio, &v.NomeSocio,
			&v.QualificacaoSocio, &v.NomeRepresentante, &v.QualificacaoRepresentante,
			&v.RepresentanteLegal, &v.FaixaEtaria, &v.PaisCodigo, &v.PaisDescricao,
			&v.CNPJCPFSocio, &v.NomeSocioGlobal,
		)
		if err != nil {
			return nil, fmt.Errorf("scan fornecedor_socio: %w", err)
		}
		resultados = append(resultados, v)
	}
	return resultados, nil
}

func nullStr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

var _ PNCPRepository = (*pncpRepositoryImpl)(nil)
