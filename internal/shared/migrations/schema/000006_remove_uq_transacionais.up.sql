ALTER TABLE prestacao_contas                        DROP CONSTRAINT IF EXISTS uq_prestacao_contas_natural;
ALTER TABLE despesa_candidato                        DROP CONSTRAINT IF EXISTS uq_despesa_candidato_natural;
ALTER TABLE despesa_orgao_partidario                  DROP CONSTRAINT IF EXISTS uq_despesa_orgao_partidario_natural;
ALTER TABLE receita_candidato                         DROP CONSTRAINT IF EXISTS uq_receita_candidato_natural;
ALTER TABLE receita_orgao_partidario                  DROP CONSTRAINT IF EXISTS uq_receita_orgao_partidario_natural;
ALTER TABLE receita_doador_originario_candidato       DROP CONSTRAINT IF EXISTS uq_receita_doador_originario_candidato;
ALTER TABLE receita_doador_originario_orgao_partidario DROP CONSTRAINT IF EXISTS uq_receita_doador_originario_orgao;
ALTER TABLE convenio                                  DROP CONSTRAINT IF EXISTS uq_convenio_numero;
