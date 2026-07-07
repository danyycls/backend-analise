package parse

import (
	"context"
	"fmt"

	"github.com/danyele/podp/internal/shared/logger"
	"github.com/danyele/podp/internal/shared/types"

	repositorios "github.com/danyele/podp/internal/esferas-brasileiras/tse/importacao/repositorios"
	tipos "github.com/danyele/podp/internal/esferas-brasileiras/tse/importacao/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const codigoSentinelaEleicao = -1

func garantirEleicaoSentinela(dados *tipos.DadosImportacao) uuid.UUID {
	if dados == nil || dados.Eleicoes == nil {
		return uuid.Nil
	}
	if existente, ok := dados.Eleicoes[codigoSentinelaEleicao]; ok {
		return existente.ID
	}
	e := &types.Eleicao{
		CodigoTSE: codigoSentinelaEleicao,
		Ano:       0,
		Descricao: "__ELEICAO_DESCONHECIDA",
	}
	e.ID = uuid.Must(uuid.NewV7())
	dados.Eleicoes[codigoSentinelaEleicao] = e
	return e.ID
}

type itemRef struct {
	setPrestacao func(uuid.UUID)
	prestacaoID  uuid.UUID
	tipo         string
	eleicaoID    uuid.UUID
	candidatoID  *uuid.UUID
	partidoID    *uuid.UUID
	sqPrestador  int64
	nomeEntidade string
	entidadeID   uuid.UUID
}

func remapearPrestacaoIDsComPlaceholderPgCopy(ctx context.Context, tx pgx.Tx, repo *repositorios.Repositorio, dados *tipos.DadosImportacao, mapeamento map[uuid.UUID]uuid.UUID) error {
	log := logger.New("LeitorCSV: Utils: remapearPrestacaoIDsComPlaceholderPgCopy")

	items := make([]itemRef, 0)

	for _, dc := range dados.DespesasCandidato {
		if novo, ok := mapeamento[dc.PrestacaoContasID]; ok {
			dc.PrestacaoContasID = novo
			continue
		}
		if obterPrestacaoPorID(dados, dc.PrestacaoContasID) != nil {
			continue
		}
		if dc.PrestacaoContasID == uuid.Nil {
			continue
		}
		eleicaoID := garantirEleicaoSentinela(dados)
		if c := obterCandidatoPorID(dados, dc.CandidatoID); c != nil {
			eleicaoID = c.EleicaoID
		}
		items = append(items, itemRef{
			setPrestacao: func(id uuid.UUID) { dc.PrestacaoContasID = id },
			prestacaoID:  dc.PrestacaoContasID,
			tipo:         "CANDIDATO",
			eleicaoID:    eleicaoID,
			candidatoID:  &dc.CandidatoID,
			sqPrestador:  dc.SQPrestadorContas,
			nomeEntidade: "despesa_candidato",
			entidadeID:   dc.ID,
		})
	}

	for _, d := range dados.DespesasOrgaoPartidario {
		if novo, ok := mapeamento[d.PrestacaoContasID]; ok {
			d.PrestacaoContasID = novo
			continue
		}
		if obterPrestacaoPorID(dados, d.PrestacaoContasID) != nil {
			continue
		}
		if d.PrestacaoContasID == uuid.Nil {
			continue
		}
		items = append(items, itemRef{
			setPrestacao: func(id uuid.UUID) { d.PrestacaoContasID = id },
			prestacaoID:  d.PrestacaoContasID,
			tipo:         "ORGAO_PARTIDARIO",
			eleicaoID:    garantirEleicaoSentinela(dados),
			partidoID:    &d.PartidoID,
			sqPrestador:  d.SQPrestadorContas,
			nomeEntidade: "despesa_orgao_partidario",
			entidadeID:   d.ID,
		})
	}

	for _, r := range dados.ReceitasCandidato {
		if novo, ok := mapeamento[r.PrestacaoContasID]; ok {
			r.PrestacaoContasID = novo
			continue
		}
		if obterPrestacaoPorID(dados, r.PrestacaoContasID) != nil {
			continue
		}
		if r.PrestacaoContasID == uuid.Nil {
			continue
		}
		eleicaoID := garantirEleicaoSentinela(dados)
		if c := obterCandidatoPorID(dados, r.CandidatoID); c != nil {
			eleicaoID = c.EleicaoID
		}
		items = append(items, itemRef{
			setPrestacao: func(id uuid.UUID) { r.PrestacaoContasID = id },
			prestacaoID:  r.PrestacaoContasID,
			tipo:         "CANDIDATO",
			eleicaoID:    eleicaoID,
			candidatoID:  &r.CandidatoID,
			sqPrestador:  r.SQPrestadorContas,
			nomeEntidade: "receita_candidato",
			entidadeID:   r.ID,
		})
	}

	for _, r := range dados.ReceitasOrgaoPartidario {
		if novo, ok := mapeamento[r.PrestacaoContasID]; ok {
			r.PrestacaoContasID = novo
			continue
		}
		if obterPrestacaoPorID(dados, r.PrestacaoContasID) != nil {
			continue
		}
		if r.PrestacaoContasID == uuid.Nil {
			continue
		}
		items = append(items, itemRef{
			setPrestacao: func(id uuid.UUID) { r.PrestacaoContasID = id },
			prestacaoID:  r.PrestacaoContasID,
			tipo:         "ORGAO_PARTIDARIO",
			eleicaoID:    garantirEleicaoSentinela(dados),
			partidoID:    &r.PartidoID,
			sqPrestador:  r.SQPrestadorContas,
			nomeEntidade: "receita_orgao_partidario",
			entidadeID:   r.ID,
		})
	}

	for _, r := range dados.ReceitasDoadorOriginarioCandidato {
		if novo, ok := mapeamento[r.PrestacaoContasID]; ok {
			r.PrestacaoContasID = novo
			continue
		}
		if obterPrestacaoPorID(dados, r.PrestacaoContasID) != nil {
			continue
		}
		if r.PrestacaoContasID == uuid.Nil {
			continue
		}
		eleicaoID := garantirEleicaoSentinela(dados)
		var candidatoID *uuid.UUID
		if rc := dados.ReceitasCandidatoPorSQ[r.SQReceita]; rc != nil {
			candidatoID = &rc.CandidatoID
			if c := obterCandidatoPorID(dados, rc.CandidatoID); c != nil {
				eleicaoID = c.EleicaoID
			}
		}
		items = append(items, itemRef{
			setPrestacao: func(id uuid.UUID) { r.PrestacaoContasID = id },
			prestacaoID:  r.PrestacaoContasID,
			tipo:         "CANDIDATO",
			eleicaoID:    eleicaoID,
			candidatoID:  candidatoID,
			sqPrestador:  r.SQPrestadorContas,
			nomeEntidade: "receita_doador_originario_candidato",
			entidadeID:   r.ID,
		})
	}

	for _, r := range dados.ReceitasDoadorOriginarioOrgaoPartidario {
		if novo, ok := mapeamento[r.PrestacaoContasID]; ok {
			r.PrestacaoContasID = novo
			continue
		}
		if obterPrestacaoPorID(dados, r.PrestacaoContasID) != nil {
			continue
		}
		if r.PrestacaoContasID == uuid.Nil {
			continue
		}
		var partidoID *uuid.UUID
		if ro := dados.ReceitasOrgaoPorSQ[r.SQReceita]; ro != nil {
			partidoID = &ro.PartidoID
		}
		items = append(items, itemRef{
			setPrestacao: func(id uuid.UUID) { r.PrestacaoContasID = id },
			prestacaoID:  r.PrestacaoContasID,
			tipo:         "ORGAO_PARTIDARIO",
			eleicaoID:    garantirEleicaoSentinela(dados),
			partidoID:    partidoID,
			sqPrestador:  r.SQPrestadorContas,
			nomeEntidade: "receita_doador_originario_orgao_partidario",
			entidadeID:   r.ID,
		})
	}

	if len(items) == 0 {
		return nil
	}

	uniqueIDs := make([]uuid.UUID, 0, len(items))
	seenUUID := make(map[uuid.UUID]struct{}, len(items))
	for _, ref := range items {
		if _, ok := seenUUID[ref.prestacaoID]; !ok {
			seenUUID[ref.prestacaoID] = struct{}{}
			uniqueIDs = append(uniqueIDs, ref.prestacaoID)
		}
	}

	existsInDB := make(map[uuid.UUID]struct{}, len(uniqueIDs))
	if len(uniqueIDs) > 0 {
		rows, err := tx.Query(ctx, `SELECT id FROM prestacao_contas WHERE id = ANY($1)`, uniqueIDs)
		if err != nil {
			return fmt.Errorf("batch check prestacao_contas: %w", err)
		}
		for rows.Next() {
			var id uuid.UUID
			if err := rows.Scan(&id); err != nil {
				rows.Close()
				return fmt.Errorf("scan prestacao_contas id: %w", err)
			}
			existsInDB[id] = struct{}{}
		}
		rows.Close()
	}

	var created int
	keyCache := make(map[string]uuid.UUID)
	sentinelID := garantirEleicaoSentinela(dados)
	groups := make(map[string][]itemRef)
	for _, ref := range items {
		if _, ok := existsInDB[ref.prestacaoID]; ok {
			ref.setPrestacao(ref.prestacaoID)
			continue
		}
		grupoKey := fmt.Sprintf("%s:%d", ref.tipo, ref.sqPrestador)
		groups[grupoKey] = append(groups[grupoKey], ref)
	}

	for _, group := range groups {
		if len(group) == 0 {
			continue
		}
		bestEleicaoID := group[0].eleicaoID
		bestCandidatoID := group[0].candidatoID
		bestPartidoID := group[0].partidoID
		for _, ref := range group {
			if ref.eleicaoID != sentinelID {
				bestEleicaoID = ref.eleicaoID
			}
			if ref.candidatoID != nil {
				bestCandidatoID = ref.candidatoID
			}
			if ref.partidoID != nil {
				bestPartidoID = ref.partidoID
			}
		}
		chave := fmt.Sprintf("%s:%s:%d", group[0].tipo, bestEleicaoID.String(), group[0].sqPrestador)
		if placeholderID, ok := keyCache[chave]; ok {
			for _, ref := range group {
				ref.setPrestacao(placeholderID)
			}
			continue
		}
		placeholder, err := repo.GarantirPrestacaoPlaceholder(ctx, tx, group[0].tipo, bestEleicaoID, bestCandidatoID, bestPartidoID, group[0].sqPrestador)
		if err != nil {
			return err
		}
		keyCache[chave] = placeholder
		for _, ref := range group {
			ref.setPrestacao(placeholder)
		}
		created++
		log.Info("prestacao nao encontrada - atribuindo placeholder",
			"tipo", group[0].tipo, "eleicao", bestEleicaoID, "sq_prestador_contas", group[0].sqPrestador, "placeholder", placeholder)
	}

	log.Info("placeholders criados", "total_placeholders", created)
	return nil
}
