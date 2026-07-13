package types

import (
	"os"
	"strconv"
)

const (
	EnvDirCSV        = "IMPORTACAO_DIRETORIO_CSV"
	EnvBatchSize     = "IMPORT_BATCH_SIZE"
	EnvMaxWorkers    = "IMPORT_MAX_WORKERS"
	EnvFilesPerBatch = "IMPORT_FILES_PER_BATCH"
)

func GetEnvInt(key string, defaultVal int) int {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return defaultVal
	}
	return n
}
