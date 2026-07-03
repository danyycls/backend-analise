package repositorios

import (
	"context"
	"errors"
	"fmt"
	"time"

	"sync"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/danyele/podp/internal/shared/types"
)

type ImportacaoResultado struct {
	RegistrosInseridos int64
	TempoCOPY          time.Duration
	TempoMerge         time.Duration
}

type Repositorio struct {
	pool *pgxpool.Pool
	// cache for placeholders to avoid repeated DB roundtrips
	placeholderMu             sync.Mutex
	prestacaoPlaceholderCache map[string]uuid.UUID
}

func Novo(pool *pgxpool.Pool) *Repositorio {
	return &Repositorio{pool: pool, prestacaoPlaceholderCache: make(map[string]uuid.UUID)}
}

func (r *Repositorio) Pool() *pgxpool.Pool {
	return r.pool
}

func (r *Repositorio) InserirEleicoesComRetorno(
	ctx context.Context, tx pgx.Tx, valores []*types.Eleicao, lote int,
) (map[uuid.UUID]uuid.UUID, error) {
	columns := []string{"id", "codigo_tse", "ano", "codigo_tipo_eleicao", "nome_tipo_eleicao", "descricao", "data_eleicao"}
	conflict := "(codigo_tse)"
	setClause := "ano = EXCLUDED.ano, codigo_tipo_eleicao = EXCLUDED.codigo_tipo_eleicao, nome_tipo_eleicao = EXCLUDED.nome_tipo_eleicao, descricao = EXCLUDED.descricao, data_eleicao = EXCLUDED.data_eleicao, updated_at = NOW()"
	returning := []string{"id", "codigo_tse"}

	return copyInsertReturning(ctx, tx, valores, lote, "eleicao", columns, conflict, setClause, returning,
		func(v *types.Eleicao) []any {
			return []any{v.ID, v.CodigoTSE, v.Ano, v.CodigoTipoEleicao, v.NomeTipoEleicao, v.Descricao, v.DataEleicao}
		},
		func(v *types.Eleicao) string { return fmt.Sprintf("%d", v.CodigoTSE) },
	)
}

func (r *Repositorio) InserirUnidadesEleitoraisComRetorno(
	ctx context.Context, tx pgx.Tx, valores []*types.UnidadeEleitoral, lote int,
) (map[uuid.UUID]uuid.UUID, error) {
	columns := []string{"id", "sg_uf", "codigo_tse", "nome"}
	conflict := "(sg_uf, codigo_tse)"
	setClause := "nome = EXCLUDED.nome, updated_at = NOW()"
	returning := []string{"id", "sg_uf", "codigo_tse"}

	return copyInsertReturning(ctx, tx, valores, lote, "unidade_eleitoral", columns, conflict, setClause, returning,
		func(v *types.UnidadeEleitoral) []any {
			return []any{v.ID, v.UFSigla, v.CodigoTSE, v.Nome}
		},
		func(v *types.UnidadeEleitoral) string { return v.UFSigla + "|" + v.CodigoTSE },
	)
}

func (r *Repositorio) InserirPartidosComRetorno(
	ctx context.Context, tx pgx.Tx, valores []*types.Partido, lote int,
) (map[uuid.UUID]uuid.UUID, error) {
	columns := []string{"id", "numero", "sigla", "nome", "federacao_codigo_tse", "federacao_sigla", "federacao_nome", "coligacao_codigo_tse", "coligacao_nome", "coligacao_composicao"}
	conflict := "(numero)"
	setClause := "sigla = EXCLUDED.sigla, nome = EXCLUDED.nome, federacao_codigo_tse = EXCLUDED.federacao_codigo_tse, federacao_sigla = EXCLUDED.federacao_sigla, federacao_nome = EXCLUDED.federacao_nome, coligacao_codigo_tse = EXCLUDED.coligacao_codigo_tse, coligacao_nome = EXCLUDED.coligacao_nome, coligacao_composicao = EXCLUDED.coligacao_composicao, updated_at = NOW()"
	returning := []string{"id", "numero"}

	return copyInsertReturning(ctx, tx, valores, lote, "partido", columns, conflict, setClause, returning,
		func(v *types.Partido) []any {
			return []any{v.ID, v.Numero, v.Sigla, v.Nome, v.FederacaoCodigoTSE, v.FederacaoSigla, v.FederacaoNome, v.ColigacaoCodigoTSE, v.ColigacaoNome, v.ColigacaoComposicao}
		},
		func(v *types.Partido) string { return fmt.Sprintf("%d", v.Numero) },
	)
}

func (r *Repositorio) InserirCandidatosComRetorno(
	ctx context.Context, tx pgx.Tx, valores []*types.Candidato, lote int,
) (map[uuid.UUID]uuid.UUID, error) {
	columns := []string{"id", "sq_candidato", "eleicao_id", "sg_uf", "partido_id", "cargo_codigo", "cargo_nome", "genero_descricao", "cor_raca_descricao", "estado_civil_nome", "grau_instrucao_nome", "ocupacao_codigo", "ocupacao_nome", "numero_candidato", "cpf", "cpf_vice", "nome_completo", "nome_urna", "nome_social", "data_nascimento", "situacao_totalizacao_descricao"}
	conflict := "(sq_candidato)"
	setClause := "eleicao_id = EXCLUDED.eleicao_id, sg_uf = EXCLUDED.sg_uf, partido_id = EXCLUDED.partido_id, cargo_codigo = EXCLUDED.cargo_codigo, cargo_nome = EXCLUDED.cargo_nome, genero_descricao = EXCLUDED.genero_descricao, cor_raca_descricao = EXCLUDED.cor_raca_descricao, estado_civil_nome = EXCLUDED.estado_civil_nome, grau_instrucao_nome = EXCLUDED.grau_instrucao_nome, ocupacao_codigo = EXCLUDED.ocupacao_codigo, ocupacao_nome = EXCLUDED.ocupacao_nome, numero_candidato = EXCLUDED.numero_candidato, cpf = EXCLUDED.cpf, cpf_vice = EXCLUDED.cpf_vice, nome_completo = EXCLUDED.nome_completo, nome_urna = EXCLUDED.nome_urna, nome_social = EXCLUDED.nome_social, data_nascimento = EXCLUDED.data_nascimento, situacao_totalizacao_descricao = EXCLUDED.situacao_totalizacao_descricao, updated_at = NOW()"
	returning := []string{"id", "sq_candidato"}

	return copyInsertReturning(ctx, tx, valores, lote, "candidato", columns, conflict, setClause, returning,
		func(v *types.Candidato) []any {
			return []any{v.ID, v.SQCandidato, v.EleicaoID, v.UFSigla, v.PartidoID, v.CargoCodigo, v.CargoNome, v.GeneroDescricao, v.CorRacaDescricao, v.EstadoCivilNome, v.GrauInstrucaoNome, v.OcupacaoCodigo, v.OcupacaoNome, v.NumeroCandidato, v.CPF, v.CPFVice, v.NomeCompleto, v.NomeUrna, v.NomeSocial, v.DataNascimento, v.SituacaoTotalizacaoDescricao}
		},
		func(v *types.Candidato) string { return fmt.Sprintf("%d", v.SQCandidato) },
	)
}

func (r *Repositorio) InserirFornecedoresComRetorno(
	ctx context.Context, tx pgx.Tx, valores []*types.Fornecedor, lote int,
) (map[uuid.UUID]uuid.UUID, error) {
	columns := []string{"id", "cpf_cnpj", "nome", "nome_rfb", "tipo_fornecedor_codigo", "tipo_fornecedor_descricao", "cnae_codigo", "cnae_descricao", "esfera_partidaria_codigo", "esfera_partidaria_descricao", "sg_uf", "municipio_nome", "sq_candidato_relacionado", "numero_candidato_relacionado", "cargo_codigo_relacionado", "cargo_descricao_relacionada", "partido_numero_relacionado", "partido_sigla_relacionado", "partido_nome_relacionado"}
	conflict := "(cpf_cnpj)"
	setClause := "nome = EXCLUDED.nome, nome_rfb = EXCLUDED.nome_rfb, tipo_fornecedor_codigo = EXCLUDED.tipo_fornecedor_codigo, tipo_fornecedor_descricao = EXCLUDED.tipo_fornecedor_descricao, cnae_codigo = EXCLUDED.cnae_codigo, cnae_descricao = EXCLUDED.cnae_descricao, esfera_partidaria_codigo = EXCLUDED.esfera_partidaria_codigo, esfera_partidaria_descricao = EXCLUDED.esfera_partidaria_descricao, sg_uf = EXCLUDED.sg_uf, municipio_nome = EXCLUDED.municipio_nome, sq_candidato_relacionado = EXCLUDED.sq_candidato_relacionado, numero_candidato_relacionado = EXCLUDED.numero_candidato_relacionado, cargo_codigo_relacionado = EXCLUDED.cargo_codigo_relacionado, cargo_descricao_relacionada = EXCLUDED.cargo_descricao_relacionada, partido_numero_relacionado = EXCLUDED.partido_numero_relacionado, partido_sigla_relacionado = EXCLUDED.partido_sigla_relacionado, partido_nome_relacionado = EXCLUDED.partido_nome_relacionado, updated_at = NOW()"
	returning := []string{"id", "cpf_cnpj"}

	return copyInsertReturning(ctx, tx, valores, lote, "fornecedor", columns, conflict, setClause, returning,
		func(v *types.Fornecedor) []any {
			return []any{v.ID, v.CPFCNPJ, v.Nome, v.NomeRFB, v.TipoFornecedorCodigo, v.TipoFornecedorDescricao, v.CNAECodigo, v.CNAEDescricao, v.EsferaPartidariaCodigo, v.EsferaPartidariaDescricao, v.UFSigla, v.MunicipioNome, v.SQCandidatoRelacionado, v.NumeroCandidatoRelacionado, v.CargoCodigoRelacionado, v.CargoDescricaoRelacionada, v.PartidoNumeroRelacionado, v.PartidoSiglaRelacionado, v.PartidoNomeRelacionado}
		},
		func(v *types.Fornecedor) string { return v.CPFCNPJ },
	)
}

func (r *Repositorio) InserirDoadoresComRetorno(
	ctx context.Context, tx pgx.Tx, valores []*types.Doador, lote int,
) (map[uuid.UUID]uuid.UUID, error) {
	columns := []string{"id", "cpf_cnpj", "nome", "nome_rfb", "cnae_codigo", "cnae_descricao", "esfera_partidaria_codigo", "esfera_partidaria_descricao", "sg_uf", "municipio_nome", "sq_candidato_relacionado", "numero_candidato_relacionado", "cargo_codigo_relacionado", "cargo_descricao_relacionada", "partido_numero_relacionado", "partido_sigla_relacionado", "partido_nome_relacionado"}
	conflict := "(cpf_cnpj)"
	setClause := "nome = EXCLUDED.nome, nome_rfb = EXCLUDED.nome_rfb, cnae_codigo = EXCLUDED.cnae_codigo, cnae_descricao = EXCLUDED.cnae_descricao, esfera_partidaria_codigo = EXCLUDED.esfera_partidaria_codigo, esfera_partidaria_descricao = EXCLUDED.esfera_partidaria_descricao, sg_uf = EXCLUDED.sg_uf, municipio_nome = EXCLUDED.municipio_nome, sq_candidato_relacionado = EXCLUDED.sq_candidato_relacionado, numero_candidato_relacionado = EXCLUDED.numero_candidato_relacionado, cargo_codigo_relacionado = EXCLUDED.cargo_codigo_relacionado, cargo_descricao_relacionada = EXCLUDED.cargo_descricao_relacionada, partido_numero_relacionado = EXCLUDED.partido_numero_relacionado, partido_sigla_relacionado = EXCLUDED.partido_sigla_relacionado, partido_nome_relacionado = EXCLUDED.partido_nome_relacionado, updated_at = NOW()"
	returning := []string{"id", "cpf_cnpj"}

	return copyInsertReturning(ctx, tx, valores, lote, "doador", columns, conflict, setClause, returning,
		func(v *types.Doador) []any {
			return []any{v.ID, v.CPFCNPJ, v.Nome, v.NomeRFB, v.CNAECodigo, v.CNAEDescricao, v.EsferaPartidariaCodigo, v.EsferaPartidariaDescricao, v.UFSigla, v.MunicipioNome, v.SQCandidatoRelacionado, v.NumeroCandidatoRelacionado, v.CargoCodigoRelacionado, v.CargoDescricaoRelacionada, v.PartidoNumeroRelacionado, v.PartidoSiglaRelacionado, v.PartidoNomeRelacionado}
		},
		func(v *types.Doador) string { return v.CPFCNPJ },
	)
}

func (r *Repositorio) InserirPrestacoesComRetorno(
	ctx context.Context, tx pgx.Tx, valores []*types.PrestacaoContas, lote int,
) (map[uuid.UUID]uuid.UUID, error) {
	columns := []string{"id", "sq_prestador_contas", "eleicao_id", "candidato_id", "partido_id", "sg_uf", "unidade_eleitoral_id", "tipo_prestador", "tipo_prestacao", "data_prestacao", "turno", "cnpj_prestador_conta", "esfera_partidaria_codigo", "esfera_partidaria_descricao"}
	conflict := "(tipo_prestador, eleicao_id, sq_prestador_contas)"
	setClause := "candidato_id = EXCLUDED.candidato_id, partido_id = EXCLUDED.partido_id, sg_uf = EXCLUDED.sg_uf, unidade_eleitoral_id = EXCLUDED.unidade_eleitoral_id, tipo_prestacao = EXCLUDED.tipo_prestacao, data_prestacao = EXCLUDED.data_prestacao, turno = EXCLUDED.turno, cnpj_prestador_conta = EXCLUDED.cnpj_prestador_conta, esfera_partidaria_codigo = EXCLUDED.esfera_partidaria_codigo, esfera_partidaria_descricao = EXCLUDED.esfera_partidaria_descricao, updated_at = NOW()"
	returning := []string{"id", "tipo_prestador", "eleicao_id", "sq_prestador_contas"}

	return copyInsertReturning(ctx, tx, valores, lote, "prestacao_contas", columns, conflict, setClause, returning,
		func(v *types.PrestacaoContas) []any {
			return []any{v.ID, v.SQPrestadorContas, v.EleicaoID, v.CandidatoID, v.PartidoID, v.UFSigla, v.UnidadeEleitoralID, v.TipoPrestador, v.TipoPrestacao, v.DataPrestacao, v.Turno, v.CNPJPrestadorConta, v.EsferaPartidariaCodigo, v.EsferaPartidariaDescricao}
		},
		func(v *types.PrestacaoContas) string {
			return v.TipoPrestador + "|" + v.EleicaoID.String() + "|" + fmt.Sprintf("%d", v.SQPrestadorContas)
		},
	)
}

func (r *Repositorio) InserirDespesasCandidato(
	ctx context.Context, tx pgx.Tx, valores []*types.DespesaCandidato, lote int,
) (int64, error) {
	columns := []string{"id", "prestacao_contas_id", "candidato_id", "fornecedor_id", "sq_despesa", "tipo_registro", "tipo_documento", "numero_documento", "origem_despesa_codigo", "origem_despesa_descricao", "fonte_despesa_codigo", "fonte_despesa_descricao", "natureza_despesa_codigo", "natureza_despesa_descricao", "especie_recurso_codigo", "especie_recurso_descricao", "sq_parcelamento_despesa", "data_despesa", "descricao", "valor"}
	conflict := "(sq_despesa, tipo_registro)"

	return copyInsertEmLote(ctx, tx, valores, lote, "despesa_candidato", columns, conflict,
		func(v *types.DespesaCandidato) []any {
			return []any{v.ID, v.PrestacaoContasID, v.CandidatoID, v.FornecedorID, v.SQDespesa, v.TipoRegistro, v.TipoDocumento, v.NumeroDocumento, v.OrigemDespesaCodigo, v.OrigemDespesaDescricao, v.FonteDespesaCodigo, v.FonteDespesaDescricao, v.NaturezaDespesaCodigo, v.NaturezaDespesaDescricao, v.EspecieRecursoCodigo, v.EspecieRecursoDescricao, v.SQPlanoParcelamento, v.DataDespesa, v.Descricao, v.Valor}
		},
	)
}

func (r *Repositorio) InserirDespesasOrgaoPartidario(
	ctx context.Context, tx pgx.Tx, valores []*types.DespesaOrgaoPartidario, lote int,
) (int64, error) {
	columns := []string{"id", "prestacao_contas_id", "partido_id", "fornecedor_id", "sq_despesa", "tipo_registro", "tipo_documento", "numero_documento", "origem_despesa_codigo", "origem_despesa_descricao", "fonte_despesa_codigo", "fonte_despesa_descricao", "natureza_despesa_codigo", "natureza_despesa_descricao", "especie_recurso_codigo", "especie_recurso_descricao", "sq_parcelamento_despesa", "data_despesa", "descricao", "valor"}
	conflict := "(sq_despesa, tipo_registro)"

	return copyInsertEmLote(ctx, tx, valores, lote, "despesa_orgao_partidario", columns, conflict,
		func(v *types.DespesaOrgaoPartidario) []any {
			return []any{v.ID, v.PrestacaoContasID, v.PartidoID, v.FornecedorID, v.SQDespesa, v.TipoRegistro, v.TipoDocumento, v.NumeroDocumento, v.OrigemDespesaCodigo, v.OrigemDespesaDescricao, v.FonteDespesaCodigo, v.FonteDespesaDescricao, v.NaturezaDespesaCodigo, v.NaturezaDespesaDescricao, v.EspecieRecursoCodigo, v.EspecieRecursoDescricao, v.SQPlanoParcelamento, v.DataDespesa, v.Descricao, v.Valor}
		},
	)
}

func (r *Repositorio) InserirReceitasCandidatoComRetorno(
	ctx context.Context, tx pgx.Tx, valores []*types.ReceitaCandidato, lote int,
) (map[uuid.UUID]uuid.UUID, error) {
	columns := []string{"id", "prestacao_contas_id", "candidato_id", "doador_id", "sq_receita", "fonte_receita_codigo", "fonte_receita_descricao", "origem_receita_codigo", "origem_receita_descricao", "natureza_receita_codigo", "natureza_receita_descricao", "especie_receita_codigo", "especie_receita_descricao", "numero_recibo_doacao", "numero_documento_doacao", "data_receita", "descricao", "valor", "natureza_recurso_estimavel", "genero", "cor_raca"}
	conflict := "(sq_receita)"
	setClause := "prestacao_contas_id = EXCLUDED.prestacao_contas_id, candidato_id = EXCLUDED.candidato_id, doador_id = EXCLUDED.doador_id, fonte_receita_codigo = EXCLUDED.fonte_receita_codigo, fonte_receita_descricao = EXCLUDED.fonte_receita_descricao, origem_receita_codigo = EXCLUDED.origem_receita_codigo, origem_receita_descricao = EXCLUDED.origem_receita_descricao, natureza_receita_codigo = EXCLUDED.natureza_receita_codigo, natureza_receita_descricao = EXCLUDED.natureza_receita_descricao, especie_receita_codigo = EXCLUDED.especie_receita_codigo, especie_receita_descricao = EXCLUDED.especie_receita_descricao, numero_recibo_doacao = EXCLUDED.numero_recibo_doacao, numero_documento_doacao = EXCLUDED.numero_documento_doacao, data_receita = EXCLUDED.data_receita, descricao = EXCLUDED.descricao, valor = EXCLUDED.valor, natureza_recurso_estimavel = EXCLUDED.natureza_recurso_estimavel, genero = EXCLUDED.genero, cor_raca = EXCLUDED.cor_raca, updated_at = NOW()"
	returning := []string{"id", "sq_receita"}

	return copyInsertReturning(ctx, tx, valores, lote, "receita_candidato", columns, conflict, setClause, returning,
		func(v *types.ReceitaCandidato) []any {
			return []any{v.ID, v.PrestacaoContasID, v.CandidatoID, v.DoadorID, v.SQReceita, v.FonteReceitaCodigo, v.FonteReceitaDescricao, v.OrigemReceitaCodigo, v.OrigemReceitaDescricao, v.NaturezaReceitaCodigo, v.NaturezaReceitaDescricao, v.EspecieReceitaCodigo, v.EspecieReceitaDescricao, v.NumeroReciboDoacao, v.NumeroDocumentoDoacao, v.DataReceita, v.Descricao, v.Valor, v.NaturezaRecursoEstimavel, v.Genero, v.CorRaca}
		},
		func(v *types.ReceitaCandidato) string { return fmt.Sprintf("%d", v.SQReceita) },
	)
}

func (r *Repositorio) InserirReceitasOrgaoComRetorno(
	ctx context.Context, tx pgx.Tx, valores []*types.ReceitaOrgaoPartidario, lote int,
) (map[uuid.UUID]uuid.UUID, error) {
	columns := []string{"id", "prestacao_contas_id", "partido_id", "doador_id", "sq_receita", "fonte_receita_codigo", "fonte_receita_descricao", "origem_receita_codigo", "origem_receita_descricao", "natureza_receita_codigo", "natureza_receita_descricao", "especie_receita_codigo", "especie_receita_descricao", "numero_recibo_doacao", "numero_documento_doacao", "data_receita", "descricao", "valor"}
	conflict := "(sq_receita)"
	setClause := "prestacao_contas_id = EXCLUDED.prestacao_contas_id, partido_id = EXCLUDED.partido_id, doador_id = EXCLUDED.doador_id, fonte_receita_codigo = EXCLUDED.fonte_receita_codigo, fonte_receita_descricao = EXCLUDED.fonte_receita_descricao, origem_receita_codigo = EXCLUDED.origem_receita_codigo, origem_receita_descricao = EXCLUDED.origem_receita_descricao, natureza_receita_codigo = EXCLUDED.natureza_receita_codigo, natureza_receita_descricao = EXCLUDED.natureza_receita_descricao, especie_receita_codigo = EXCLUDED.especie_receita_codigo, especie_receita_descricao = EXCLUDED.especie_receita_descricao, numero_recibo_doacao = EXCLUDED.numero_recibo_doacao, numero_documento_doacao = EXCLUDED.numero_documento_doacao, data_receita = EXCLUDED.data_receita, descricao = EXCLUDED.descricao, valor = EXCLUDED.valor, updated_at = NOW()"
	returning := []string{"id", "sq_receita"}

	return copyInsertReturning(ctx, tx, valores, lote, "receita_orgao_partidario", columns, conflict, setClause, returning,
		func(v *types.ReceitaOrgaoPartidario) []any {
			return []any{v.ID, v.PrestacaoContasID, v.PartidoID, v.DoadorID, v.SQReceita, v.FonteReceitaCodigo, v.FonteReceitaDescricao, v.OrigemReceitaCodigo, v.OrigemReceitaDescricao, v.NaturezaReceitaCodigo, v.NaturezaReceitaDescricao, v.EspecieReceitaCodigo, v.EspecieReceitaDescricao, v.NumeroReciboDoacao, v.NumeroDocumentoDoacao, v.DataReceita, v.Descricao, v.Valor}
		},
		func(v *types.ReceitaOrgaoPartidario) string { return fmt.Sprintf("%d", v.SQReceita) },
	)
}

func (r *Repositorio) InserirEmLote(
	ctx context.Context, tx pgx.Tx, valores interface{}, lote int,
) (int64, error) {
	switch v := valores.(type) {
	case []*types.Convenio:
		return r.inserirConvenios(ctx, tx, v, lote)
	case []*types.ReceitaDoadorOriginarioCandidato:
		return r.inserirReceitasDoadorOriginarioCandidato(ctx, tx, v, lote)
	case []*types.ReceitaDoadorOriginarioOrgaoPartidario:
		return r.inserirReceitasDoadorOriginarioOrgaoPartidario(ctx, tx, v, lote)
	case []*types.BemCandidato:
		return r.inserirBensCandidato(ctx, tx, v, lote)
	default:
		return 0, fmt.Errorf("tipo nao suportado para InserirEmLote: %T", valores)
	}
}

func (r *Repositorio) inserirConvenios(
	ctx context.Context, tx pgx.Tx, valores []*types.Convenio, lote int,
) (int64, error) {
	columns := []string{
		"id", "numero_convenio", "uf", "codigo_siafi_municipio", "nome_municipio",
		"situacao_convenio", "numero_original", "numero_processo", "objeto_convenio",
		"codigo_orgao_superior", "nome_orgao_superior",
		"codigo_orgao_concedente", "nome_orgao_concedente",
		"codigo_ug_concedente", "nome_ug_concedente",
		"codigo_convenente", "tipo_convenente", "nome_convenente", "tipo_ente_convenente",
		"tipo_instrumento",
		"valor_convenio", "valor_liberado",
		"data_publicacao", "data_inicio_vigencia", "data_final_vigencia",
		"valor_contrapartida", "data_ultima_liberacao", "valor_ultima_liberacao",
	}
	conflict := "(numero_convenio)"

	return copyInsertEmLote(ctx, tx, valores, lote, "convenio", columns, conflict,
		func(v *types.Convenio) []any {
			return []any{
				v.ID, v.NumeroConvenio, strNil(v.UF), strNil(v.CodigoSIAFIMunicipio),
				strNil(v.NomeMunicipio), strNil(v.SituacaoConvenio),
				strNil(v.NumeroOriginal), strNil(v.NumeroProcesso),
				strNil(v.ObjetoConvenio),
				strNil(v.CodigoOrgaoSuperior), strNil(v.NomeOrgaoSuperior),
				strNil(v.CodigoOrgaoConcedente), strNil(v.NomeOrgaoConcedente),
				strNil(v.CodigoUGConcedente), strNil(v.NomeUGConcedente),
				strNil(v.CodigoConvenente), strNil(v.TipoConvenente),
				strNil(v.NomeConvenente), strNil(v.TipoEnteConvenente),
				strNil(v.TipoInstrumento),
				v.ValorConvenio, v.ValorLiberado,
				v.DataPublicacao, v.DataInicioVigencia, v.DataFinalVigencia,
				v.ValorContrapartida, v.DataUltimaLiberacao, v.ValorUltimaLiberacao,
			}
		},
	)
}

func (r *Repositorio) inserirReceitasDoadorOriginarioCandidato(
	ctx context.Context, tx pgx.Tx, valores []*types.ReceitaDoadorOriginarioCandidato, lote int,
) (int64, error) {
	columns := []string{"id", "prestacao_contas_id", "receita_candidato_id", "sq_receita", "documento_doador", "nome_doador", "nome_doador_rfb", "tipo_doador", "cnae_codigo", "cnae_descricao", "data_receita", "descricao", "valor"}
	conflict := "(sq_receita, documento_doador)"

	return copyInsertEmLote(ctx, tx, valores, lote, "receita_doador_originario_candidato", columns, conflict,
		func(v *types.ReceitaDoadorOriginarioCandidato) []any {
			return []any{v.ID, v.PrestacaoContasID, v.ReceitaCandidatoID, v.SQReceita, v.DocumentoDoador, v.NomeDoador, v.NomeDoadorRFB, v.TipoDoador, v.CNAECodigo, v.CNAEDescricao, v.DataReceita, v.Descricao, v.Valor}
		},
	)
}

func (r *Repositorio) inserirReceitasDoadorOriginarioOrgaoPartidario(
	ctx context.Context, tx pgx.Tx, valores []*types.ReceitaDoadorOriginarioOrgaoPartidario, lote int,
) (int64, error) {
	columns := []string{"id", "prestacao_contas_id", "receita_orgao_partidario_id", "sq_receita", "documento_doador", "nome_doador", "nome_doador_rfb", "tipo_doador", "cnae_codigo", "cnae_descricao", "data_receita", "descricao", "valor"}
	conflict := "(sq_receita, documento_doador)"

	return copyInsertEmLote(ctx, tx, valores, lote, "receita_doador_originario_orgao_partidario", columns, conflict,
		func(v *types.ReceitaDoadorOriginarioOrgaoPartidario) []any {
			return []any{v.ID, v.PrestacaoContasID, v.ReceitaOrgaoPartidarioID, v.SQReceita, v.DocumentoDoador, v.NomeDoador, v.NomeDoadorRFB, v.TipoDoador, v.CNAECodigo, v.CNAEDescricao, v.DataReceita, v.Descricao, v.Valor}
		},
	)
}

func (r *Repositorio) inserirBensCandidato(
	ctx context.Context, tx pgx.Tx, valores []*types.BemCandidato, lote int,
) (int64, error) {
	columns := []string{"id", "candidato_id", "tipo_bem_codigo", "tipo_bem_nome", "numero_ordem", "descricao", "valor", "data_ultima_atualizacao", "hora_ultima_atualizacao"}
	conflict := "(candidato_id, numero_ordem)"

	return copyInsertEmLote(ctx, tx, valores, lote, "bem_candidato", columns, conflict,
		func(v *types.BemCandidato) []any {
			var horaUltimaAtualizacao any = v.HoraUltimaAtualizacao
			if v.HoraUltimaAtualizacao == "" {
				horaUltimaAtualizacao = nil
			}
			return []any{v.ID, v.CandidatoID, v.TipoBemCodigo, v.TipoBemNome, v.NumeroOrdem, v.Descricao, v.Valor, v.DataUltimaAtualizacao, horaUltimaAtualizacao}
		},
	)
}

func (r *Repositorio) ArquivoJaImportado(ctx context.Context, tx pgx.Tx, caminhoRelativo string) (bool, error) {
	var total int64
	err := tx.QueryRow(ctx, `SELECT COUNT(*) FROM arquivo_importado WHERE caminho_relativo = $1`, caminhoRelativo).Scan(&total)
	if err != nil {
		return false, fmt.Errorf("verificar arquivo importado: %w", err)
	}
	return total > 0, nil
}

func (r *Repositorio) ListarTodosArquivosImportados(ctx context.Context) (map[string]bool, error) {
	importados := make(map[string]bool)
	rows, err := r.pool.Query(ctx, `SELECT nome FROM arquivo_importado`)
	if err != nil {
		return nil, fmt.Errorf("listar arquivos importados: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var nome string
		if err := rows.Scan(&nome); err != nil {
			return nil, fmt.Errorf("scan nome arquivo: %w", err)
		}
		importados[nome] = true
	}
	return importados, rows.Err()
}

func (r *Repositorio) RegistrarArquivoImportado(ctx context.Context, tx pgx.Tx, caminhoRelativo, nome, tipo, uf string, totalRegistros int) error {
	_, err := tx.Exec(ctx,
		`INSERT INTO arquivo_importado (caminho_relativo, nome, tipo, uf, total_registros) VALUES ($1,$2,$3,$4,$5)`,
		caminhoRelativo, nome, tipo, uf, totalRegistros,
	)
	if err != nil {
		return fmt.Errorf("registrar arquivo importado: %w", err)
	}
	return nil
}

func (r *Repositorio) AdquirirConexao(ctx context.Context) (*pgxpool.Conn, error) {
	return r.pool.Acquire(ctx)
}

// BuscarIDCandidatoPorSQ retorna o UUID de um candidato ja persistido.
func (r *Repositorio) BuscarIDCandidatoPorSQ(ctx context.Context, sq int64) (uuid.UUID, error) {
	var id uuid.UUID
	err := r.pool.QueryRow(ctx, `SELECT id FROM candidato WHERE sq_candidato = $1`, sq).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, fmt.Errorf("candidato SQ %d nao encontrado no banco", sq)
		}
		return uuid.Nil, fmt.Errorf("buscar candidato SQ %d: %w", sq, err)
	}
	return id, nil
}

// CarregarCandidatosNoMapa preenche o mapa de candidatos para resolucao em arquivos dependentes.
func (r *Repositorio) CarregarCandidatosNoMapa(ctx context.Context, dest map[int64]*types.Candidato) (int, error) {
	if dest == nil {
		return 0, fmt.Errorf("mapa de candidatos nao informado")
	}
	rows, err := r.pool.Query(ctx, `SELECT id, sq_candidato FROM candidato`)
	if err != nil {
		return 0, fmt.Errorf("carregar candidatos: %w", err)
	}
	defer rows.Close()

	total := 0
	for rows.Next() {
		var id uuid.UUID
		var sq int64
		if err := rows.Scan(&id, &sq); err != nil {
			return total, fmt.Errorf("scan candidato: %w", err)
		}
		dest[sq] = &types.Candidato{ModeloBase: types.ModeloBase{ID: id}, SQCandidato: sq}
		total++
	}
	return total, rows.Err()
}
