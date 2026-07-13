package parse

import (
	"context"
	"fmt"
	"time"

	"github.com/danyele/podp/internal/shared/logger"
	"github.com/danyele/podp/internal/shared/types"
	"github.com/google/uuid"

	repositorios "github.com/danyele/podp/internal/sources/tse/importacao/repositorios"
	tipos "github.com/danyele/podp/internal/sources/tse/importacao/types"
	"github.com/jackc/pgx/v5"
)

func lockDimensao(ctx context.Context, tx pgx.Tx, nome string) error {
	sql := `SELECT pg_advisory_xact_lock(hashtext($1))`
	_, err := tx.Exec(ctx, sql, nome)
	if err != nil {
		return fmt.Errorf("lock %s: %w", nome, err)
	}
	return nil
}

func PersistirDadosImportacaoPgCopy(
	ctx context.Context,
	tx pgx.Tx,
	repo *repositorios.Repositorio,
	dados *tipos.DadosImportacao,
	lote int,
	resultado *repositorios.ImportacaoResultado,
) error {
	log := logger.New("LeitorCSV: Utils: PersistirDadosImportacaoPgCopy")
	inicio := time.Now()

	if resultado.Etapas == nil {
		resultado.Etapas = make(map[string]time.Duration)
	}

	// Nivel 0: Convenios (independence)
	if len(dados.Convenios) > 0 {
		if resultado.SetEntidade != nil {
			resultado.SetEntidade("n0_convenios")
		}
		qtdOperacoes := 0
		copyAntes := resultado.TempoCOPY
		mergeAntes := resultado.TempoMerge
		inicioNivel := time.Now()

		log.Debug("persistindo", "nivel", "0", "entidade", "Convenios", "registros", len(dados.Convenios))
		if _, err := repo.InserirEmLote(ctx, tx, dados.Convenios, lote, resultado); err != nil {
			return fmt.Errorf("convenio: %w", err)
		}
		qtdOperacoes++

		durCopy := resultado.TempoCOPY - copyAntes
		durMerge := resultado.TempoMerge - mergeAntes
		durTotal := time.Since(inicioNivel)
		durParse := durTotal - durCopy - durMerge
		if durParse < 0 {
			durParse = 0
		}
		resultado.RegistrarNivel("n0_convenios", repositorios.NivelTiming{
			Copia:     durCopy,
			Mesclar:   durMerge,
			Parse:     durParse,
			Total:     durTotal,
			Registros: int64(len(dados.Convenios)),
			Operacoes: qtdOperacoes,
		})
		if resultado.SetEntidade != nil {
			resultado.SetEntidade("")
		}
	}

	if err := lockDimensao(ctx, tx, "importacao"); err != nil {
		return err
	}

	// Nivel 1: Eleicoes, UnidadesEleitorais, Partidos (serial — pgx.Tx nao suporta uso concorrente)
	if resultado.SetEntidade != nil {
		resultado.SetEntidade("n1_dimensoes")
	}
	qtdOperacoes1 := 0
	copyAntes1 := resultado.TempoCOPY
	mergeAntes1 := resultado.TempoMerge
	inicioNivel1 := time.Now()

	var mapeamentoEleicao map[uuid.UUID]uuid.UUID
	var mapeamentoUE map[uuid.UUID]uuid.UUID
	var mapeamentoPartido map[uuid.UUID]uuid.UUID
	var err error

	if len(dados.Eleicoes) > 0 {
		mapeamentoEleicao, err = repo.InserirEleicoesComRetorno(ctx, tx, valores(dados.Eleicoes), lote, resultado)
		if err != nil {
			return err
		}
		qtdOperacoes1++
	}
	if len(dados.UnidadesEleitorais) > 0 {
		mapeamentoUE, err = repo.InserirUnidadesEleitoraisComRetorno(ctx, tx, valores(dados.UnidadesEleitorais), lote, resultado)
		if err != nil {
			return err
		}
		qtdOperacoes1++
	}
	if len(dados.Partidos) > 0 {
		mapeamentoPartido, err = repo.InserirPartidosComRetorno(ctx, tx, valores(dados.Partidos), lote, resultado)
		if err != nil {
			return err
		}
		qtdOperacoes1++
	}

	remapearEleicaoIDs(dados, mapeamentoEleicao)
	remapearUnidadeEleitoralIDs(dados, mapeamentoUE)
	remapearPartidoIDs(dados, mapeamentoPartido)
	log.Debug("nivel 1 persistido", "duracao", time.Since(inicio).String())

	durCopy1 := resultado.TempoCOPY - copyAntes1
	durMerge1 := resultado.TempoMerge - mergeAntes1
	durTotal1 := time.Since(inicioNivel1)
	durParse1 := durTotal1 - durCopy1 - durMerge1
	if durParse1 < 0 {
		durParse1 = 0
	}
	qtdRegistros1 := int64(len(dados.Eleicoes) + len(dados.UnidadesEleitorais) + len(dados.Partidos))
	resultado.RegistrarNivel("n1_dimensoes", repositorios.NivelTiming{
		Copia:     durCopy1,
		Mesclar:   durMerge1,
		Parse:     durParse1,
		Total:     durTotal1,
		Registros: qtdRegistros1,
		Operacoes: qtdOperacoes1,
	})
	if resultado.SetEntidade != nil {
		resultado.SetEntidade("")
	}

	// Nivel 2: Candidatos (serial, depende dos IDs do nivel 1)
	if len(dados.Candidatos) > 0 {
		if resultado.SetEntidade != nil {
			resultado.SetEntidade("n2_candidatos")
		}
		qtdOperacoes2 := 0
		copyAntes2 := resultado.TempoCOPY
		mergeAntes2 := resultado.TempoMerge
		inicioNivel2 := time.Now()

		log.Debug("persistindo", "nivel", "2", "entidade", "Candidatos", "registros", len(dados.Candidatos))
		mapeamentoCandidato, err := repo.InserirCandidatosComRetorno(ctx, tx, valores(dados.Candidatos), lote, resultado)
		if err != nil {
			return fmt.Errorf("candidato: %w", err)
		}
		qtdOperacoes2++
		remapearCandidatoIDs(dados, mapeamentoCandidato)
		dados.CandidatosPorID = make(map[uuid.UUID]*types.Candidato, len(dados.Candidatos))
		for _, c := range dados.Candidatos {
			dados.CandidatosPorID[c.ID] = c
		}
		log.Debug("nivel 2 persistido", "duracao", time.Since(inicio).String())

		durCopy2 := resultado.TempoCOPY - copyAntes2
		durMerge2 := resultado.TempoMerge - mergeAntes2
		durTotal2 := time.Since(inicioNivel2)
		durParse2 := durTotal2 - durCopy2 - durMerge2
		if durParse2 < 0 {
			durParse2 = 0
		}
		resultado.RegistrarNivel("n2_candidatos", repositorios.NivelTiming{
			Copia:     durCopy2,
			Mesclar:   durMerge2,
			Parse:     durParse2,
			Total:     durTotal2,
			Registros: int64(len(dados.Candidatos)),
			Operacoes: qtdOperacoes2,
		})
		if resultado.SetEntidade != nil {
			resultado.SetEntidade("")
		}
	}

	// Nivel 3: Fornecedores, Doadores (serial — pgx.Tx nao suporta uso concorrente)
	if resultado.SetEntidade != nil {
		resultado.SetEntidade("n3_fornecedores_doadores")
	}
	qtdOperacoes3 := 0
	copyAntes3 := resultado.TempoCOPY
	mergeAntes3 := resultado.TempoMerge
	inicioNivel3 := time.Now()

	var mapeamentoFornecedor map[uuid.UUID]uuid.UUID
	var mapeamentoDoador map[uuid.UUID]uuid.UUID

	if len(dados.Fornecedores) > 0 {
		mapeamentoFornecedor, err = repo.InserirFornecedoresComRetorno(ctx, tx, valores(dados.Fornecedores), lote, resultado)
		if err != nil {
			return err
		}
		qtdOperacoes3++
	}
	if len(dados.Doadores) > 0 {
		mapeamentoDoador, err = repo.InserirDoadoresComRetorno(ctx, tx, valores(dados.Doadores), lote, resultado)
		if err != nil {
			return err
		}
		qtdOperacoes3++
	}

	remapearFornecedorIDs(dados, mapeamentoFornecedor)
	remapearDoadorIDs(dados, mapeamentoDoador)
	log.Debug("nivel 3 persistido", "duracao", time.Since(inicio).String())

	durCopy3 := resultado.TempoCOPY - copyAntes3
	durMerge3 := resultado.TempoMerge - mergeAntes3
	durTotal3 := time.Since(inicioNivel3)
	durParse3 := durTotal3 - durCopy3 - durMerge3
	if durParse3 < 0 {
		durParse3 = 0
	}
	qtdRegistros3 := int64(len(dados.Fornecedores) + len(dados.Doadores))
	resultado.RegistrarNivel("n3_fornecedores_doadores", repositorios.NivelTiming{
		Copia:     durCopy3,
		Mesclar:   durMerge3,
		Parse:     durParse3,
		Total:     durTotal3,
		Registros: qtdRegistros3,
		Operacoes: qtdOperacoes3,
	})
	if resultado.SetEntidade != nil {
		resultado.SetEntidade("")
	}

	// Nivel 4: PrestacoesContas (serial, depende dos niveis 2 e 3)
	if len(dados.Prestacoes) > 0 {
		if resultado.SetEntidade != nil {
			resultado.SetEntidade("n4_prestacoes")
		}
		qtdOperacoes4 := 0
		copyAntes4 := resultado.TempoCOPY
		mergeAntes4 := resultado.TempoMerge
		inicioNivel4 := time.Now()

		log.Debug("persistindo", "nivel", "4", "entidade", "PrestacoesContas", "registros", len(dados.Prestacoes))

		// Reconstroi map apos remapeamento de dimensoes para evitar
		// duplicates no ON CONFLICT (SQLSTATE 21000) e
		// violacao de check constraint (SQLSTATE 23514)
		prestacoesDedup := make(map[string]*types.PrestacaoContas, len(dados.Prestacoes))
		for _, p := range dados.Prestacoes {
			chave := chavePrestacaoNatural(p.TipoPrestador, p.EleicaoID, p.SQPrestadorContas)
			if existente, ok := prestacoesDedup[chave]; ok {
				if existente.CandidatoID == nil && p.CandidatoID != nil {
					existente.CandidatoID = p.CandidatoID
				}
				if existente.PartidoID == nil && p.PartidoID != nil {
					existente.PartidoID = p.PartidoID
				}
				if existente.UFSigla == nil && p.UFSigla != nil {
					existente.UFSigla = p.UFSigla
				}
				if existente.UnidadeEleitoralID == nil && p.UnidadeEleitoralID != nil {
					existente.UnidadeEleitoralID = p.UnidadeEleitoralID
				}
				if existente.TipoPrestacao == "" && p.TipoPrestacao != "" {
					existente.TipoPrestacao = p.TipoPrestacao
				}
				if existente.DataPrestacao == nil && p.DataPrestacao != nil {
					existente.DataPrestacao = p.DataPrestacao
				}
				if existente.Turno == nil && p.Turno != nil {
					existente.Turno = p.Turno
				}
				if existente.CNPJPrestadorConta == "" && p.CNPJPrestadorConta != "" {
					existente.CNPJPrestadorConta = p.CNPJPrestadorConta
				}
			} else {
				prestacoesDedup[chave] = p
			}
		}
		dados.Prestacoes = prestacoesDedup

		mapeamentoCorrigido := make(map[uuid.UUID]uuid.UUID, len(dados.Prestacoes))
		for _, p := range dados.Prestacoes {
			uuidOriginal := p.ID
			p.ID = uuid.Must(uuid.NewV7())
			mapeamentoCorrigido[uuidOriginal] = p.ID
		}

		mapeamentoPrestacao, err := repo.InserirPrestacoesComRetorno(ctx, tx, valores(dados.Prestacoes), lote, resultado)
		if err != nil {
			return fmt.Errorf("prestacao_contas: %w", err)
		}
		qtdOperacoes4++

		for uuidOriginal, uuidNovo := range mapeamentoCorrigido {
			if uuidDB, ok := mapeamentoPrestacao[uuidNovo]; ok {
				mapeamentoCorrigido[uuidOriginal] = uuidDB
			}
		}

		dados.PrestacoesPorID = make(map[uuid.UUID]*types.PrestacaoContas, len(dados.Prestacoes))
		for _, p := range dados.Prestacoes {
			dados.PrestacoesPorID[p.ID] = p
		}

		if err := remapearPrestacaoIDsComPlaceholderPgCopy(ctx, tx, repo, dados, mapeamentoCorrigido); err != nil {
			return err
		}
		sincronizarDependenciasDePrestacao(dados)
		log.Debug("nivel 4 persistido", "duracao", time.Since(inicio).String())

		durCopy4 := resultado.TempoCOPY - copyAntes4
		durMerge4 := resultado.TempoMerge - mergeAntes4
		durTotal4 := time.Since(inicioNivel4)
		durParse4 := durTotal4 - durCopy4 - durMerge4
		if durParse4 < 0 {
			durParse4 = 0
		}
		resultado.RegistrarNivel("n4_prestacoes", repositorios.NivelTiming{
			Copia:     durCopy4,
			Mesclar:   durMerge4,
			Parse:     durParse4,
			Total:     durTotal4,
			Registros: int64(len(dados.Prestacoes)),
			Operacoes: qtdOperacoes4,
		})
		if resultado.SetEntidade != nil {
			resultado.SetEntidade("")
		}
	}

	// Nivel 5: Despesas e Receitas (serial — pgx.Tx nao suporta uso concorrente)
	if resultado.SetEntidade != nil {
		resultado.SetEntidade("n5_receitas_despesas")
	}
	qtdOperacoes5 := 0
	copyAntes5 := resultado.TempoCOPY
	mergeAntes5 := resultado.TempoMerge
	inicioNivel5 := time.Now()

	var mapeamentoReceitaCand map[uuid.UUID]uuid.UUID
	var mapeamentoReceitaOrgao map[uuid.UUID]uuid.UUID

	if len(dados.DespesasCandidato) > 0 {
		_, err = repo.InserirDespesasCandidato(ctx, tx, dados.DespesasCandidato, lote, resultado)
		if err != nil {
			return err
		}
		qtdOperacoes5++
	}
	if len(dados.DespesasOrgaoPartidario) > 0 {
		_, err = repo.InserirDespesasOrgaoPartidario(ctx, tx, dados.DespesasOrgaoPartidario, lote, resultado)
		if err != nil {
			return err
		}
		qtdOperacoes5++
	}
	if len(dados.ReceitasCandidato) > 0 {
		mapeamentoReceitaCand, err = repo.InserirReceitasCandidatoComRetorno(ctx, tx, dados.ReceitasCandidato, lote, resultado)
		if err != nil {
			return err
		}
		qtdOperacoes5++
	}
	if len(dados.ReceitasOrgaoPartidario) > 0 {
		mapeamentoReceitaOrgao, err = repo.InserirReceitasOrgaoComRetorno(ctx, tx, dados.ReceitasOrgaoPartidario, lote, resultado)
		if err != nil {
			return err
		}
		qtdOperacoes5++
	}

	remapearReceitaCandidatoIDs(dados, mapeamentoReceitaCand)
	remapearReceitaOrgaoPartidarioIDs(dados, mapeamentoReceitaOrgao)
	log.Debug("nivel 5 persistido", "duracao", time.Since(inicio).String())

	durCopy5 := resultado.TempoCOPY - copyAntes5
	durMerge5 := resultado.TempoMerge - mergeAntes5
	durTotal5 := time.Since(inicioNivel5)
	durParse5 := durTotal5 - durCopy5 - durMerge5
	if durParse5 < 0 {
		durParse5 = 0
	}
	qtdRegistros5 := int64(len(dados.DespesasCandidato) + len(dados.DespesasOrgaoPartidario) + len(dados.ReceitasCandidato) + len(dados.ReceitasOrgaoPartidario))
	resultado.RegistrarNivel("n5_receitas_despesas", repositorios.NivelTiming{
		Copia:     durCopy5,
		Mesclar:   durMerge5,
		Parse:     durParse5,
		Total:     durTotal5,
		Registros: qtdRegistros5,
		Operacoes: qtdOperacoes5,
	})
	if resultado.SetEntidade != nil {
		resultado.SetEntidade("")
	}

	// Nivel 6: ReceitasDoadorOriginario e BensCandidato (serial — pgx.Tx nao suporta uso concorrente)
	if resultado.SetEntidade != nil {
		resultado.SetEntidade("n6_bens")
	}
	qtdOperacoes6 := 0
	copyAntes6 := resultado.TempoCOPY
	mergeAntes6 := resultado.TempoMerge
	inicioNivel6 := time.Now()

	if len(dados.ReceitasDoadorOriginarioCandidato) > 0 {
		_, err = repo.InserirEmLote(ctx, tx, dados.ReceitasDoadorOriginarioCandidato, lote, resultado)
		if err != nil {
			return err
		}
		qtdOperacoes6++
	}
	if len(dados.ReceitasDoadorOriginarioOrgaoPartidario) > 0 {
		_, err = repo.InserirEmLote(ctx, tx, dados.ReceitasDoadorOriginarioOrgaoPartidario, lote, resultado)
		if err != nil {
			return err
		}
		qtdOperacoes6++
	}
	if len(dados.BensCandidato) > 0 {
		_, err = repo.InserirEmLote(ctx, tx, dados.BensCandidato, lote, resultado)
		if err != nil {
			return err
		}
		qtdOperacoes6++
	}
	log.Debug("nivel 6 persistido", "duracao", time.Since(inicio).String())

	durCopy6 := resultado.TempoCOPY - copyAntes6
	durMerge6 := resultado.TempoMerge - mergeAntes6
	durTotal6 := time.Since(inicioNivel6)
	durParse6 := durTotal6 - durCopy6 - durMerge6
	if durParse6 < 0 {
		durParse6 = 0
	}
	qtdRegistros6 := int64(len(dados.ReceitasDoadorOriginarioCandidato) + len(dados.ReceitasDoadorOriginarioOrgaoPartidario) + len(dados.BensCandidato))
	resultado.RegistrarNivel("n6_bens", repositorios.NivelTiming{
		Copia:     durCopy6,
		Mesclar:   durMerge6,
		Parse:     durParse6,
		Total:     durTotal6,
		Registros: qtdRegistros6,
		Operacoes: qtdOperacoes6,
	})
	if resultado.SetEntidade != nil {
		resultado.SetEntidade("")
	}

	resultado.Etapas["copy"] += resultado.TempoCOPY
	resultado.Etapas["merge"] += resultado.TempoMerge
	log.Info("persistencia concluida", "duracao_total", time.Since(inicio).String())
	return nil
}
