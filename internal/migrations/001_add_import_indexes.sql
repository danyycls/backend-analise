-- Índices para otimizar queries de importação de dados CSV
-- Executar antes da primeira importação para melhor performance

CREATE INDEX IF NOT EXISTS idx_candidato_sq ON candidato(sq_candidato);
CREATE INDEX IF NOT EXISTS idx_receita_candidato_sq ON receita_candidato(sq_receita);
CREATE INDEX IF NOT EXISTS idx_despesa_candidato_sq_tipo ON despesa_candidato(sq_despesa, tipo_registro);
CREATE INDEX IF NOT EXISTS idx_despesa_orgao_partidario_sq_tipo ON despesa_orgao_partidario(sq_despesa, tipo_registro);
CREATE INDEX IF NOT EXISTS idx_prestacao_contas_chave ON prestacao_contas(tipo_prestador, eleicao_id, sq_prestador_contas);
CREATE INDEX IF NOT EXISTS idx_receita_candidato_prestacao ON receita_candidato(prestacao_contas_id);
CREATE INDEX IF NOT EXISTS idx_receita_orgao_partidario_prestacao ON receita_orgao_partidario(prestacao_contas_id);
CREATE INDEX IF NOT EXISTS idx_despesa_candidato_prestacao ON despesa_candidato(prestacao_contas_id);
CREATE INDEX IF NOT EXISTS idx_despesa_orgao_partidario_prestacao ON despesa_orgao_partidario(prestacao_contas_id);
CREATE INDEX IF NOT EXISTS idx_arquivo_importado_nome ON arquivo_importado(nome);
