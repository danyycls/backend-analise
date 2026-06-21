package repositorios

import (
	"context"
	"fmt"

	"github.com/danyele/podp/internal/shared/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// GarantirPrestacaoPlaceholder garante que exista uma prestacao placeholder
// similar ao comportamento do repositorio GORM.
func (r *Repositorio) GarantirPrestacaoPlaceholder(ctx context.Context, tx pgx.Tx, tipoPrestador string, eleicaoID uuid.UUID, candidatoID *uuid.UUID, partidoID *uuid.UUID) (uuid.UUID, error) {
	log := logger.New("LeitorCSV: Repositorio: GarantirPrestacaoPlaceholder")
	chave := fmt.Sprintf("%s:%s", tipoPrestador, eleicaoID.String())

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

	insertSQL := `INSERT INTO prestacao_contas (sq_prestador_contas, eleicao_id, candidato_id, partido_id, tipo_prestador, tipo_prestacao, cnpj_prestador_conta, created_at, updated_at)
VALUES (-1, $1, $2, $3, $4, '__PLACEHOLDER', '__PLACEHOLDER', NOW(), NOW())
ON CONFLICT (tipo_prestador, eleicao_id, sq_prestador_contas) DO NOTHING RETURNING id`

	var newID uuid.UUID
	row := tx.QueryRow(ctx, insertSQL, eleicaoID, cand, part, tipoPrestador)
	if err := row.Scan(&newID); err != nil {
		sel := `SELECT id FROM prestacao_contas WHERE tipo_prestador = $1 AND eleicao_id = $2 AND sq_prestador_contas = -1 LIMIT 1`
		var existing uuid.UUID
		if err2 := tx.QueryRow(ctx, sel, tipoPrestador, eleicaoID).Scan(&existing); err2 != nil {
			return uuid.Nil, fmt.Errorf("garantir placeholder pgcopy: falha ao selecionar placeholder: %w (insert err: %w)", err2, err)
		}
		newID = existing
	}

	r.placeholderMu.Lock()
	r.prestacaoPlaceholderCache[chave] = newID
	r.placeholderMu.Unlock()

	log.Info("placeholder_prestacao_pgcopy usado",
		"id", newID, "tipo", tipoPrestador, "eleicao", eleicaoID)
	return newID, nil
}
