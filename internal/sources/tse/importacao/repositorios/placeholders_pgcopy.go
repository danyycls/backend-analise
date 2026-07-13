package repositorios

import (
	"context"
	"errors"
	"fmt"

	"github.com/danyele/podp/internal/shared/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const maxPlaceholderCache = 1000

// GarantirPrestacaoPlaceholder garante que exista uma prestacao placeholder
// usando o SQ real do prestador (nao sentinela -1), permitindo que quando
// a prestacao real chegar em outro arquivo CSV, o ON CONFLICT atualize os dados.
func (r *Repositorio) GarantirPrestacaoPlaceholder(ctx context.Context, tx pgx.Tx, tipoPrestador string, eleicaoID uuid.UUID, candidatoID *uuid.UUID, partidoID *uuid.UUID, sqPrestadorContas int64) (uuid.UUID, error) {
	log := logger.New("LeitorCSV: Repositorio: GarantirPrestacaoPlaceholder")
	chave := fmt.Sprintf("%s:%s:%d", tipoPrestador, eleicaoID.String(), sqPrestadorContas)

	r.placeholderMu.Lock()
	if id, ok := r.prestacaoPlaceholderCache[chave]; ok {
		r.placeholderMu.Unlock()
		return id, nil
	}
	r.placeholderMu.Unlock()

	if candidatoID != nil && *candidatoID == uuid.Nil {
		candidatoID = nil
	}
	if partidoID != nil && *partidoID == uuid.Nil {
		partidoID = nil
	}

	if tipoPrestador == "CANDIDATO" && candidatoID == nil {
		var candID uuid.UUID
		candQuery := `INSERT INTO candidato (sq_candidato, eleicao_id, nome_completo, created_at, updated_at) VALUES (-1, $1, '__PLACEHOLDER', NOW(), NOW()) ON CONFLICT (sq_candidato) DO UPDATE SET nome_completo = EXCLUDED.nome_completo RETURNING id`
		if err := tx.QueryRow(ctx, candQuery, eleicaoID).Scan(&candID); err != nil {
			return uuid.Nil, fmt.Errorf("garantir placeholder candidato pgcopy: %w", err)
		}
		candidatoID = &candID
	}
	if tipoPrestador == "ORGAO_PARTIDARIO" && partidoID == nil {
		var pid uuid.UUID
		partQuery := `INSERT INTO partido (numero, sigla, nome, created_at, updated_at) VALUES (0, '__PH', '__PLACEHOLDER', NOW(), NOW()) ON CONFLICT (numero) DO NOTHING RETURNING id`
		if err := tx.QueryRow(ctx, partQuery).Scan(&pid); err != nil {
			sel := `SELECT id FROM partido WHERE numero = 0 LIMIT 1`
			if err2 := tx.QueryRow(ctx, sel).Scan(&pid); err2 != nil {
				return uuid.Nil, fmt.Errorf("garantir placeholder partido pgcopy: %w", err2)
			}
		}
		partidoID = &pid
	}

	var cand interface{}
	var part interface{}
	if candidatoID != nil {
		cand = *candidatoID
	} else {
		cand = nil
	}
	if partidoID != nil {
		part = *partidoID
	} else {
		part = nil
	}

	sqVal := sqPrestadorContas
	if sqVal <= 0 {
		sqVal = -1
	}

	sel := `SELECT id FROM prestacao_contas WHERE tipo_prestador = $1 AND eleicao_id = $2 AND sq_prestador_contas = $3 LIMIT 1`
	var existing uuid.UUID
	err := tx.QueryRow(ctx, sel, tipoPrestador, eleicaoID, sqVal).Scan(&existing)
	if err == nil {
		r.placeholderMu.Lock()
		r.prestacaoPlaceholderCache[chave] = existing
		if len(r.prestacaoPlaceholderCache) > maxPlaceholderCache {
			toRemove := len(r.prestacaoPlaceholderCache) / 2
			for k := range r.prestacaoPlaceholderCache {
				delete(r.prestacaoPlaceholderCache, k)
				toRemove--
				if toRemove <= 0 {
					break
				}
			}
		}
		r.placeholderMu.Unlock()
		return existing, nil
	}

	// Verifica se a eleicao jah existe no banco (persistida em lote anterior
	// ou no passo 1 do lote atual). Se existir, a FK jah esta satisfeita e
	// reaproveitamos o ID persistido. Se nao, inserimos a sentinela com
	// ON CONFLICT (codigo_tse) para evitar duplicates de codigo_tse=-1.
	var persistedID uuid.UUID
	if err := tx.QueryRow(ctx, `SELECT id FROM eleicao WHERE id = $1 LIMIT 1`, eleicaoID).Scan(&persistedID); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, fmt.Errorf("garantir eleicao placeholder pgcopy: %w", err)
		}
		if err := tx.QueryRow(ctx, `INSERT INTO eleicao (id, codigo_tse, ano, descricao, created_at, updated_at)
VALUES ($1, -1, 0, '__EL_PH', NOW(), NOW()) ON CONFLICT (codigo_tse) DO UPDATE SET updated_at = NOW() RETURNING id`, eleicaoID).Scan(&persistedID); err != nil {
			return uuid.Nil, fmt.Errorf("garantir eleicao placeholder pgcopy: %w", err)
		}
	}
	eleicaoID = persistedID

	insertSQL := `INSERT INTO prestacao_contas (sq_prestador_contas, eleicao_id, candidato_id, partido_id, tipo_prestador, tipo_prestacao, cnpj_prestador_conta, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, '__PLACEHOLDER', '__PLACEHOLDER', NOW(), NOW())
RETURNING id`

	var newID uuid.UUID
	if err := tx.QueryRow(ctx, insertSQL, sqVal, eleicaoID, cand, part, tipoPrestador).Scan(&newID); err != nil {
		return uuid.Nil, fmt.Errorf("garantir placeholder pgcopy: %w", err)
	}

	r.placeholderMu.Lock()
	r.prestacaoPlaceholderCache[chave] = newID
	if len(r.prestacaoPlaceholderCache) > maxPlaceholderCache {
		toRemove := len(r.prestacaoPlaceholderCache) / 2
		for k := range r.prestacaoPlaceholderCache {
			delete(r.prestacaoPlaceholderCache, k)
			toRemove--
			if toRemove <= 0 {
				break
			}
		}
	}
	r.placeholderMu.Unlock()

	log.Warn("placeholder_prestacao_pgcopy usado — dados reais ausentes",
		"id", newID, "tipo", tipoPrestador, "eleicao", eleicaoID, "sq_prestador_contas", sqVal)
	return newID, nil
}
