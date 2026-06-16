package parse

import (
	"context"

	"github.com/danyele/laceu/internal/shared/logger"

	repositorios "github.com/danyele/laceu/internal/esferas-brasileiras/tse/importacao/repositorios"
	tipos "github.com/danyele/laceu/internal/esferas-brasileiras/tse/importacao/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// remapearPrestacaoIDsComPlaceholderPgCopy aplica placeholders nas entradas
// que referenciam prestacoes inexistentes, usando o repositorio
func remapearPrestacaoIDsComPlaceholderPgCopy(ctx context.Context, tx pgx.Tx, repo *repositorios.Repositorio, dados *tipos.DadosImportacao, mapeamento map[uuid.UUID]uuid.UUID) error {
	log := logger.New("LeitorCSV: Utils: remapearPrestacaoIDsComPlaceholderPgCopy")
	created := 0
	nilCount := 0
	for _, r := range dados.ReceitasOrgaoPartidario {
		if r.PrestacaoContasID == uuid.Nil {
			nilCount++
		}
	}
	log.Info("inicio remapeamento de placeholders",
		"receitas_orgao_partidario", len(dados.ReceitasOrgaoPartidario), "prestacao_nil", nilCount)

	for _, despesa := range dados.DespesasCandidato {
		if novo, ok := mapeamento[despesa.PrestacaoContasID]; ok {
			despesa.PrestacaoContasID = novo
			continue
		}
		if obterPrestacaoPorID(dados, despesa.PrestacaoContasID) != nil {
			continue
		}
		tipo := "CANDIDATO"
		var eleicaoID uuid.UUID
		if candidato := obterCandidatoPorID(dados, despesa.CandidatoID); candidato != nil {
			eleicaoID = candidato.EleicaoID
		} else {
			for _, e := range dados.Eleicoes {
				eleicaoID = e.ID
				break
			}
		}
		placeholder, err := repo.GarantirPrestacaoPlaceholder(ctx, tx, tipo, eleicaoID, &despesa.CandidatoID, nil)
		if err != nil {
			return err
		}
		created++
		old := despesa.PrestacaoContasID
		despesa.PrestacaoContasID = placeholder
		log.Info("prestacao nao encontrada - atribuindo placeholder",
			"prestacao_antiga", old, "placeholder", placeholder, "entidade", "despesa_candidato", "id", despesa.ID)
	}

	for _, despesa := range dados.DespesasOrgaoPartidario {
		if novo, ok := mapeamento[despesa.PrestacaoContasID]; ok {
			despesa.PrestacaoContasID = novo
			continue
		}
		if obterPrestacaoPorID(dados, despesa.PrestacaoContasID) != nil {
			continue
		}
		tipo := "ORGAO_PARTIDARIO"
		var eleicaoID uuid.UUID
		for _, e := range dados.Eleicoes {
			eleicaoID = e.ID
			break
		}
		placeholder, err := repo.GarantirPrestacaoPlaceholder(ctx, tx, tipo, eleicaoID, nil, &despesa.PartidoID)
		if err != nil {
			return err
		}
		created++
		old := despesa.PrestacaoContasID
		despesa.PrestacaoContasID = placeholder
		log.Info("prestacao nao encontrada - atribuindo placeholder",
			"prestacao_antiga", old, "placeholder", placeholder, "entidade", "despesa_orgao_partidario", "id", despesa.ID)
	}

	for _, receita := range dados.ReceitasCandidato {
		if novo, ok := mapeamento[receita.PrestacaoContasID]; ok {
			receita.PrestacaoContasID = novo
			continue
		}
		if obterPrestacaoPorID(dados, receita.PrestacaoContasID) != nil {
			continue
		}
		tipo := "CANDIDATO"
		var eleicaoID uuid.UUID
		if candidato := obterCandidatoPorID(dados, receita.CandidatoID); candidato != nil {
			eleicaoID = candidato.EleicaoID
		} else {
			for _, e := range dados.Eleicoes {
				eleicaoID = e.ID
				break
			}
		}
		placeholder, err := repo.GarantirPrestacaoPlaceholder(ctx, tx, tipo, eleicaoID, &receita.CandidatoID, nil)
		if err != nil {
			return err
		}
		created++
		old := receita.PrestacaoContasID
		receita.PrestacaoContasID = placeholder
		log.Info("prestacao nao encontrada - atribuindo placeholder",
			"prestacao_antiga", old, "placeholder", placeholder, "entidade", "receita_candidato", "id", receita.ID)
	}

	for _, receita := range dados.ReceitasOrgaoPartidario {
		if novo, ok := mapeamento[receita.PrestacaoContasID]; ok {
			receita.PrestacaoContasID = novo
			continue
		}
		if obterPrestacaoPorID(dados, receita.PrestacaoContasID) != nil {
			continue
		}
		if receita.PrestacaoContasID != uuid.Nil {
			var exists int
			err := tx.QueryRow(ctx, `SELECT 1 FROM prestacao_contas WHERE id = $1`, receita.PrestacaoContasID).Scan(&exists)
			if err != nil {
			} else {
				continue
			}
		}
		tipo := "ORGAO_PARTIDARIO"
		var eleicaoID uuid.UUID
		for _, e := range dados.Eleicoes {
			eleicaoID = e.ID
			break
		}
		placeholder, err := repo.GarantirPrestacaoPlaceholder(ctx, tx, tipo, eleicaoID, nil, &receita.PartidoID)
		if err != nil {
			return err
		}
		created++
		old := receita.PrestacaoContasID
		receita.PrestacaoContasID = placeholder
		log.Info("prestacao nao encontrada - atribuindo placeholder",
			"prestacao_antiga", old, "placeholder", placeholder, "entidade", "receita_orgao_partidario", "id", receita.ID)
	}

	for _, origem := range dados.ReceitasDoadorOriginarioCandidato {
		if novo, ok := mapeamento[origem.PrestacaoContasID]; ok {
			origem.PrestacaoContasID = novo
			continue
		}
		if obterPrestacaoPorID(dados, origem.PrestacaoContasID) != nil {
			continue
		}
		tipo := "CANDIDATO"
		var eleicaoID uuid.UUID
		for _, e := range dados.Eleicoes {
			eleicaoID = e.ID
			break
		}
		placeholder, err := repo.GarantirPrestacaoPlaceholder(ctx, tx, tipo, eleicaoID, nil, nil)
		if err != nil {
			return err
		}
		created++
		old := origem.PrestacaoContasID
		origem.PrestacaoContasID = placeholder
		log.Info("prestacao nao encontrada - atribuindo placeholder",
			"prestacao_antiga", old, "placeholder", placeholder, "entidade", "receita_doador_originario_candidato", "id", origem.ID)
	}

	for _, origem := range dados.ReceitasDoadorOriginarioOrgaoPartidario {
		if novo, ok := mapeamento[origem.PrestacaoContasID]; ok {
			origem.PrestacaoContasID = novo
			continue
		}
		if obterPrestacaoPorID(dados, origem.PrestacaoContasID) != nil {
			continue
		}
		tipo := "ORGAO_PARTIDARIO"
		var eleicaoID uuid.UUID
		for _, e := range dados.Eleicoes {
			eleicaoID = e.ID
			break
		}
		placeholder, err := repo.GarantirPrestacaoPlaceholder(ctx, tx, tipo, eleicaoID, nil, nil)
		if err != nil {
			return err
		}
		created++
		old := origem.PrestacaoContasID
		origem.PrestacaoContasID = placeholder
		log.Info("prestacao nao encontrada - atribuindo placeholder",
			"prestacao_antiga", old, "placeholder", placeholder, "entidade", "receita_doador_originario_orgao_partidario", "id", origem.ID)
	}

	log.Info("placeholders criados",
		"total_placeholders", created)
	return nil
}
