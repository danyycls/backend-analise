package repositorios

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/danyele/podp/internal/shared/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func setID(v any, id uuid.UUID) {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	f := rv.FieldByName("ID")
	if f.IsValid() && f.CanSet() {
		f.Set(reflect.ValueOf(id))
	}
}

func getID(v any) uuid.UUID {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	f := rv.FieldByName("ID")
	if f.IsValid() {
		return f.Interface().(uuid.UUID)
	}
	return uuid.Nil
}

type tableConfig struct {
	stagingName string
	targetName  string
}

func tabelaConfig(target string) tableConfig {
	return tableConfig{
		stagingName: "stg_" + target,
		targetName:  target,
	}
}

func garantirTabelaStaging(ctx context.Context, tx pgx.Tx, cfg tableConfig) error {
	sql := fmt.Sprintf(
		`CREATE TEMP TABLE IF NOT EXISTS %s (LIKE %s INCLUDING DEFAULTS)`,
		cfg.stagingName, cfg.targetName,
	)
	if _, err := tx.Exec(ctx, sql); err != nil {
		return fmt.Errorf("criar staging %s: %w", cfg.stagingName, err)
	}
	if _, err := tx.Exec(ctx, fmt.Sprintf(`TRUNCATE TABLE %s`, cfg.stagingName)); err != nil {
		return fmt.Errorf("truncate staging %s: %w", cfg.stagingName, err)
	}
	return nil
}

func strNil(s string) any {
	if s == "" {
		return nil
	}
	return s
}

func formatKey(v any) string {
	switch val := v.(type) {
	case [16]byte:
		return uuid.UUID(val).String()
	case string:
		return val
	case int64:
		return fmt.Sprintf("%d", val)
	case int32:
		return fmt.Sprintf("%d", val)
	case int16:
		return fmt.Sprintf("%d", val)
	case int:
		return fmt.Sprintf("%d", val)
	default:
		return fmt.Sprintf("%v", val)
	}
}

func copyInsertReturning[T any](
	ctx context.Context, tx pgx.Tx,
	valores []*T, lote int,
	nomeTabela string,
	columns []string,
	conflictTarget string,
	setClause string,
	returningColumns []string,
	extractValues func(v *T) []any,
	keyFunc func(v *T) string,
) (map[uuid.UUID]uuid.UUID, error) {
	log := logger.New("LeitorCSV: Repositorio: copyInsertReturning")
	resultado := make(map[uuid.UUID]uuid.UUID)
	if len(valores) == 0 {
		return resultado, nil
	}

	cfg := tabelaConfig(nomeTabela)
	if err := garantirTabelaStaging(ctx, tx, cfg); err != nil {
		return nil, err
	}

	colList := strings.Join(columns, ", ")

	for i := 0; i < len(valores); i += lote {
		end := i + lote
		if end > len(valores) {
			end = len(valores)
		}
		batch := valores[i:end]

		rows := make([][]any, len(batch))
		for j, v := range batch {
			rows[j] = extractValues(v)
		}

		startCopy := time.Now()
		copied, err := tx.CopyFrom(ctx, pgx.Identifier{cfg.stagingName}, columns, pgx.CopyFromRows(rows))
		durCopy := time.Since(startCopy)
		if err != nil {
			return nil, fmt.Errorf("copy %s: %w", nomeTabela, err)
		}
		log.Info("pgcopy batch copiado",
			"tabela", nomeTabela, "inicio", i, "fim", end, "copiados", copied, "duracao", durCopy.String())
		if copied != int64(len(batch)) {
			return nil, fmt.Errorf("copy %s: esperado %d, copiado %d", nomeTabela, len(batch), copied)
		}

		returningList := strings.Join(returningColumns, ", ")

		keys := strings.Trim(conflictTarget, "() ")
		selectFromStaging := cfg.stagingName
		if keys != "" {
			selectFromStaging = fmt.Sprintf("(SELECT DISTINCT ON (%s) %s FROM %s ORDER BY %s, id) AS %s", keys, colList, cfg.stagingName, keys, cfg.stagingName)
		}

		mergeSQL := fmt.Sprintf(
			`INSERT INTO %s (%s) SELECT %s FROM %s ON CONFLICT %s DO UPDATE SET %s RETURNING %s`,
			cfg.targetName, colList, colList, selectFromStaging,
			conflictTarget, setClause, returningList,
		)

		startMerge := time.Now()
		rowsResult, err := tx.Query(ctx, mergeSQL)
		durMerge := time.Since(startMerge)
		if err != nil {
			return nil, fmt.Errorf("merge %s: %w", nomeTabela, err)
		}
		log.Info("pgcopy merge executado",
			"tabela", nomeTabela, "inicio", i, "fim", end, "merge_duracao", durMerge.String())

		keyToList := make(map[string][]*T, len(batch))
		lookup := make(map[string]*T, len(batch))
		origIDs := make(map[*T]uuid.UUID, len(batch))
		for _, v := range batch {
			k := keyFunc(v)
			keyToList[k] = append(keyToList[k], v)
			lookup[k] = v
			origIDs[v] = getID(v)
		}

		countReturning := 0
		for rowsResult.Next() {
			vals, err := rowsResult.Values()
			if err != nil {
				rowsResult.Close()
				return nil, fmt.Errorf("scan values %s: %w", nomeTabela, err)
			}
			if len(vals) < 2 {
				continue
			}

			idRaw := vals[0].([16]byte)
			idUUID := uuid.UUID(idRaw)

			var chave string
			if len(vals) == 2 {
				chave = formatKey(vals[1])
			} else {
				parts := make([]string, len(vals)-1)
				for idx, v := range vals[1:] {
					parts[idx] = formatKey(v)
				}
				chave = strings.Join(parts, "|")
			}

			if list, ok := keyToList[chave]; ok {
				for _, entry := range list {
					antigoID := origIDs[entry]
					resultado[antigoID] = idUUID
					setID(entry, idUUID)
					countReturning++
				}
			} else if v, ok := lookup[chave]; ok {
				antigoID := getID(v)
				resultado[antigoID] = idUUID
				setID(v, idUUID)
				countReturning++
			}
		}
		rowsResult.Close()
		if countReturning == 0 {
			log.Info("pgcopy merge sem returning",
				"tabela", nomeTabela, "inicio", i, "fim", end, "returning_rows", 0)
		} else {
			log.Info("pgcopy merge returning",
				"tabela", nomeTabela, "inicio", i, "fim", end, "returning_rows", countReturning)
		}
		if err := rowsResult.Err(); err != nil {
			return nil, fmt.Errorf("iteracao %s: %w", nomeTabela, err)
		}

		if _, err := tx.Exec(ctx, fmt.Sprintf(`TRUNCATE TABLE %s`, cfg.stagingName)); err != nil {
			return nil, fmt.Errorf("truncate staging %s: %w", nomeTabela, err)
		}
	}

	return resultado, nil
}

func copyInsertEmLote[T any](
	ctx context.Context, tx pgx.Tx,
	valores []*T, lote int,
	nomeTabela string,
	columns []string,
	conflictTarget string,
	extractValues func(v *T) []any,
) (int64, error) {
	if len(valores) == 0 {
		return 0, nil
	}

	cfg := tabelaConfig(nomeTabela)
	if err := garantirTabelaStaging(ctx, tx, cfg); err != nil {
		return 0, err
	}

	colList := strings.Join(columns, ", ")

	var totalInserido int64
	for i := 0; i < len(valores); i += lote {
		end := i + lote
		if end > len(valores) {
			end = len(valores)
		}
		batch := valores[i:end]

		rows := make([][]any, len(batch))
		for j, v := range batch {
			rows[j] = extractValues(v)
		}

		copied, err := tx.CopyFrom(ctx, pgx.Identifier{cfg.stagingName}, columns, pgx.CopyFromRows(rows))
		if err != nil {
			return totalInserido, fmt.Errorf("copy %s: %w", nomeTabela, err)
		}
		if copied != int64(len(batch)) {
			return totalInserido, fmt.Errorf("copy %s: esperado %d, copiado %d", nomeTabela, len(batch), copied)
		}

		keys := strings.Trim(conflictTarget, "() ")
		selectFromStaging := cfg.stagingName
		if keys != "" {
			selectFromStaging = fmt.Sprintf("(SELECT DISTINCT ON (%s) %s FROM %s ORDER BY %s, id) AS %s", keys, colList, cfg.stagingName, keys, cfg.stagingName)
		}

		mergeSQL := fmt.Sprintf(
			`INSERT INTO %s (%s) SELECT %s FROM %s ON CONFLICT %s DO NOTHING`,
			cfg.targetName, colList, colList, selectFromStaging,
			conflictTarget,
		)

		tag, err := tx.Exec(ctx, mergeSQL)
		if err != nil {
			return totalInserido, fmt.Errorf("merge %s: %w", nomeTabela, err)
		}
		totalInserido += tag.RowsAffected()

		if _, err := tx.Exec(ctx, fmt.Sprintf(`TRUNCATE TABLE %s`, cfg.stagingName)); err != nil {
			return totalInserido, fmt.Errorf("truncate staging %s: %w", nomeTabela, err)
		}
	}

	return totalInserido, nil
}
