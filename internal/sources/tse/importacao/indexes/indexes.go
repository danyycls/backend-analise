package indexes

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type IndexDef struct {
	Name string
	SQL  string
}

var SecondaryIndexes = []IndexDef{
	{Name: "idx_eleicao_codigo_tse", SQL: "CREATE INDEX IF NOT EXISTS idx_eleicao_codigo_tse ON eleicao (codigo_tse)"},
	{Name: "idx_unidade_eleitoral_sg_uf", SQL: "CREATE INDEX IF NOT EXISTS idx_unidade_eleitoral_sg_uf ON unidade_eleitoral (sg_uf)"},
	{Name: "idx_unidade_eleitoral_codigo_tse", SQL: "CREATE INDEX IF NOT EXISTS idx_unidade_eleitoral_codigo_tse ON unidade_eleitoral (codigo_tse)"},
	{Name: "idx_partido_numero", SQL: "CREATE INDEX IF NOT EXISTS idx_partido_numero ON partido (numero)"},
	{Name: "idx_candidato_sq_candidato", SQL: "CREATE INDEX IF NOT EXISTS idx_candidato_sq_candidato ON candidato (sq_candidato)"},
	{Name: "idx_candidato_cpf", SQL: "CREATE INDEX IF NOT EXISTS idx_candidato_cpf ON candidato (cpf)"},
	{Name: "idx_candidato_eleicao_uf", SQL: "CREATE INDEX IF NOT EXISTS idx_candidato_eleicao_uf ON candidato (eleicao_id, sg_uf)"},
	{Name: "idx_candidato_partido_id", SQL: "CREATE INDEX IF NOT EXISTS idx_candidato_partido_id ON candidato (partido_id)"},
	{Name: "idx_bem_candidato_candidato_id", SQL: "CREATE INDEX IF NOT EXISTS idx_bem_candidato_candidato_id ON bem_candidato (candidato_id)"},
	{Name: "idx_fornecedor_cpf_cnpj", SQL: "CREATE INDEX IF NOT EXISTS idx_fornecedor_cpf_cnpj ON fornecedor (cpf_cnpj)"},
	{Name: "idx_fornecedor_sg_uf", SQL: "CREATE INDEX IF NOT EXISTS idx_fornecedor_sg_uf ON fornecedor (sg_uf)"},
	{Name: "idx_doador_cpf_cnpj", SQL: "CREATE INDEX IF NOT EXISTS idx_doador_cpf_cnpj ON doador (cpf_cnpj)"},
	{Name: "idx_doador_sg_uf", SQL: "CREATE INDEX IF NOT EXISTS idx_doador_sg_uf ON doador (sg_uf)"},
	{Name: "idx_prestacao_contas_sq_prestador", SQL: "CREATE INDEX IF NOT EXISTS idx_prestacao_contas_sq_prestador ON prestacao_contas (sq_prestador_contas)"},
	{Name: "idx_prestacao_contas_eleicao_tipo", SQL: "CREATE INDEX IF NOT EXISTS idx_prestacao_contas_eleicao_tipo ON prestacao_contas (eleicao_id, tipo_prestador)"},
	{Name: "idx_despesa_candidato_sq_despesa", SQL: "CREATE INDEX IF NOT EXISTS idx_despesa_candidato_sq_despesa ON despesa_candidato (sq_despesa)"},
	{Name: "idx_despesa_candidato_candidato_id", SQL: "CREATE INDEX IF NOT EXISTS idx_despesa_candidato_candidato_id ON despesa_candidato (candidato_id)"},
	{Name: "idx_despesa_candidato_fornecedor_id", SQL: "CREATE INDEX IF NOT EXISTS idx_despesa_candidato_fornecedor_id ON despesa_candidato (fornecedor_id)"},
	{Name: "idx_despesa_candidato_prestacao_id", SQL: "CREATE INDEX IF NOT EXISTS idx_despesa_candidato_prestacao_id ON despesa_candidato (prestacao_contas_id)"},
	{Name: "idx_despesa_candidato_data", SQL: "CREATE INDEX IF NOT EXISTS idx_despesa_candidato_data ON despesa_candidato (data_despesa)"},
	{Name: "idx_despesa_orgao_partidario_sq_despesa", SQL: "CREATE INDEX IF NOT EXISTS idx_despesa_orgao_partidario_sq_despesa ON despesa_orgao_partidario (sq_despesa)"},
	{Name: "idx_despesa_orgao_partidario_partido_id", SQL: "CREATE INDEX IF NOT EXISTS idx_despesa_orgao_partidario_partido_id ON despesa_orgao_partidario (partido_id)"},
	{Name: "idx_despesa_orgao_partidario_fornecedor_id", SQL: "CREATE INDEX IF NOT EXISTS idx_despesa_orgao_partidario_fornecedor_id ON despesa_orgao_partidario (fornecedor_id)"},
	{Name: "idx_despesa_orgao_partidario_prestacao_id", SQL: "CREATE INDEX IF NOT EXISTS idx_despesa_orgao_partidario_prestacao_id ON despesa_orgao_partidario (prestacao_contas_id)"},
	{Name: "idx_despesa_orgao_partidario_data", SQL: "CREATE INDEX IF NOT EXISTS idx_despesa_orgao_partidario_data ON despesa_orgao_partidario (data_despesa)"},
	{Name: "idx_receita_candidato_sq_receita", SQL: "CREATE INDEX IF NOT EXISTS idx_receita_candidato_sq_receita ON receita_candidato (sq_receita)"},
	{Name: "idx_receita_candidato_candidato_id", SQL: "CREATE INDEX IF NOT EXISTS idx_receita_candidato_candidato_id ON receita_candidato (candidato_id)"},
	{Name: "idx_receita_candidato_doador_id", SQL: "CREATE INDEX IF NOT EXISTS idx_receita_candidato_doador_id ON receita_candidato (doador_id)"},
	{Name: "idx_receita_candidato_prestacao_id", SQL: "CREATE INDEX IF NOT EXISTS idx_receita_candidato_prestacao_id ON receita_candidato (prestacao_contas_id)"},
	{Name: "idx_receita_candidato_data", SQL: "CREATE INDEX IF NOT EXISTS idx_receita_candidato_data ON receita_candidato (data_receita)"},
	{Name: "idx_receita_orgao_partidario_sq_receita", SQL: "CREATE INDEX IF NOT EXISTS idx_receita_orgao_partidario_sq_receita ON receita_orgao_partidario (sq_receita)"},
	{Name: "idx_receita_orgao_partidario_partido_id", SQL: "CREATE INDEX IF NOT EXISTS idx_receita_orgao_partidario_partido_id ON receita_orgao_partidario (partido_id)"},
	{Name: "idx_receita_orgao_partidario_doador_id", SQL: "CREATE INDEX IF NOT EXISTS idx_receita_orgao_partidario_doador_id ON receita_orgao_partidario (doador_id)"},
	{Name: "idx_receita_orgao_partidario_prestacao_id", SQL: "CREATE INDEX IF NOT EXISTS idx_receita_orgao_partidario_prestacao_id ON receita_orgao_partidario (prestacao_contas_id)"},
	{Name: "idx_receita_orgao_partidario_data", SQL: "CREATE INDEX IF NOT EXISTS idx_receita_orgao_partidario_data ON receita_orgao_partidario (data_receita)"},
	{Name: "idx_receita_doador_originario_candidato_prestacao_id", SQL: "CREATE INDEX IF NOT EXISTS idx_receita_doador_originario_candidato_prestacao_id ON receita_doador_originario_candidato (prestacao_contas_id)"},
	{Name: "idx_receita_doador_originario_candidato_receita_id", SQL: "CREATE INDEX IF NOT EXISTS idx_receita_doador_originario_candidato_receita_id ON receita_doador_originario_candidato (receita_candidato_id)"},
	{Name: "idx_receita_doador_originario_candidato_sq", SQL: "CREATE INDEX IF NOT EXISTS idx_receita_doador_originario_candidato_sq ON receita_doador_originario_candidato (sq_receita)"},
	{Name: "idx_receita_doador_originario_candidato_doc", SQL: "CREATE INDEX IF NOT EXISTS idx_receita_doador_originario_candidato_doc ON receita_doador_originario_candidato (documento_doador)"},
	{Name: "idx_receita_doador_originario_candidato_data", SQL: "CREATE INDEX IF NOT EXISTS idx_receita_doador_originario_candidato_data ON receita_doador_originario_candidato (data_receita)"},
	{Name: "idx_receita_doador_originario_orgao_prestacao_id", SQL: "CREATE INDEX IF NOT EXISTS idx_receita_doador_originario_orgao_prestacao_id ON receita_doador_originario_orgao_partidario (prestacao_contas_id)"},
	{Name: "idx_receita_doador_originario_orgao_receita_id", SQL: "CREATE INDEX IF NOT EXISTS idx_receita_doador_originario_orgao_receita_id ON receita_doador_originario_orgao_partidario (receita_orgao_partidario_id)"},
	{Name: "idx_receita_doador_originario_orgao_sq", SQL: "CREATE INDEX IF NOT EXISTS idx_receita_doador_originario_orgao_sq ON receita_doador_originario_orgao_partidario (sq_receita)"},
	{Name: "idx_receita_doador_originario_orgao_doc", SQL: "CREATE INDEX IF NOT EXISTS idx_receita_doador_originario_orgao_doc ON receita_doador_originario_orgao_partidario (documento_doador)"},
	{Name: "idx_receita_doador_originario_orgao_data", SQL: "CREATE INDEX IF NOT EXISTS idx_receita_doador_originario_orgao_data ON receita_doador_originario_orgao_partidario (data_receita)"},
	{Name: "idx_arquivo_importado_tipo", SQL: "CREATE INDEX IF NOT EXISTS idx_arquivo_importado_tipo ON arquivo_importado (tipo)"},
	{Name: "idx_arquivo_importado_uf", SQL: "CREATE INDEX IF NOT EXISTS idx_arquivo_importado_uf ON arquivo_importado (uf)"},
	{Name: "idx_candidato_sq", SQL: "CREATE INDEX IF NOT EXISTS idx_candidato_sq ON candidato (sq_candidato)"},
	{Name: "idx_receita_candidato_sq", SQL: "CREATE INDEX IF NOT EXISTS idx_receita_candidato_sq ON receita_candidato (sq_receita)"},
	{Name: "idx_despesa_candidato_sq_tipo", SQL: "CREATE INDEX IF NOT EXISTS idx_despesa_candidato_sq_tipo ON despesa_candidato (sq_despesa, tipo_registro)"},
	{Name: "idx_despesa_orgao_partidario_sq_tipo", SQL: "CREATE INDEX IF NOT EXISTS idx_despesa_orgao_partidario_sq_tipo ON despesa_orgao_partidario (sq_despesa, tipo_registro)"},
	{Name: "idx_receita_candidato_prestacao", SQL: "CREATE INDEX IF NOT EXISTS idx_receita_candidato_prestacao ON receita_candidato (prestacao_contas_id)"},
	{Name: "idx_receita_orgao_partidario_prestacao", SQL: "CREATE INDEX IF NOT EXISTS idx_receita_orgao_partidario_prestacao ON receita_orgao_partidario (prestacao_contas_id)"},
	{Name: "idx_despesa_candidato_prestacao", SQL: "CREATE INDEX IF NOT EXISTS idx_despesa_candidato_prestacao ON despesa_candidato (prestacao_contas_id)"},
	{Name: "idx_despesa_orgao_partidario_prestacao", SQL: "CREATE INDEX IF NOT EXISTS idx_despesa_orgao_partidario_prestacao ON despesa_orgao_partidario (prestacao_contas_id)"},
	{Name: "idx_arquivo_importado_nome", SQL: "CREATE INDEX IF NOT EXISTS idx_arquivo_importado_nome ON arquivo_importado (nome)"},
}

var ConstraintIndexes = []IndexDef{
	{Name: "idx_prestacao_contas_chave", SQL: "CREATE UNIQUE INDEX IF NOT EXISTS idx_prestacao_contas_chave ON prestacao_contas (tipo_prestador, eleicao_id, sq_prestador_contas)"},
}

func DropSecondaryIndexes(ctx context.Context, tx pgx.Tx) error {
	for _, idx := range SecondaryIndexes {
		sql := "DROP INDEX IF EXISTS " + idx.Name
		if _, err := tx.Exec(ctx, sql); err != nil {
			return fmt.Errorf("drop index %s: %w", idx.Name, err)
		}
	}
	return nil
}

func RecreateConstraintIndexes(ctx context.Context, tx pgx.Tx) error {
	for _, idx := range ConstraintIndexes {
		if _, err := tx.Exec(ctx, "DROP INDEX IF EXISTS "+idx.Name); err != nil {
			return fmt.Errorf("drop constraint index %s: %w", idx.Name, err)
		}
		if _, err := tx.Exec(ctx, idx.SQL); err != nil {
			return fmt.Errorf("create constraint index %s: %w", idx.Name, err)
		}
	}
	return nil
}

func RecreateSecondaryIndexes(ctx context.Context, tx pgx.Tx) error {
	for _, idx := range SecondaryIndexes {
		if _, err := tx.Exec(ctx, idx.SQL); err != nil {
			return fmt.Errorf("create index %s: %w", idx.Name, err)
		}
	}
	return RecreateConstraintIndexes(ctx, tx)
}

var TablesToAnalyze = []string{
	"eleicao", "unidade_eleitoral", "partido",
	"candidato", "bem_candidato",
	"fornecedor", "doador",
	"prestacao_contas",
	"despesa_candidato", "despesa_orgao_partidario",
	"receita_candidato", "receita_orgao_partidario",
	"receita_doador_originario_candidato", "receita_doador_originario_orgao_partidario",
	"arquivo_importado", "convenio",
}

func AnalyzeTables(ctx context.Context, tx pgx.Tx) error {
	for _, tbl := range TablesToAnalyze {
		if _, err := tx.Exec(ctx, "ANALYZE "+tbl); err != nil {
			return fmt.Errorf("analyze %s: %w", tbl, err)
		}
	}
	return nil
}
