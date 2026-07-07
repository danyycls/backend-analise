ALTER TABLE prestacao_contas                        ADD CONSTRAINT uq_prestacao_contas_natural UNIQUE (tipo_prestador, eleicao_id, sq_prestador_contas);
ALTER TABLE despesa_candidato                        ADD CONSTRAINT uq_despesa_candidato_natural UNIQUE (sq_despesa, tipo_registro);
ALTER TABLE despesa_orgao_partidario                  ADD CONSTRAINT uq_despesa_orgao_partidario_natural UNIQUE (sq_despesa, tipo_registro);
ALTER TABLE receita_candidato                         ADD CONSTRAINT uq_receita_candidato_natural UNIQUE (sq_receita);
ALTER TABLE receita_orgao_partidario                  ADD CONSTRAINT uq_receita_orgao_partidario_natural UNIQUE (sq_receita);
ALTER TABLE receita_doador_originario_candidato       ADD CONSTRAINT uq_receita_doador_originario_candidato UNIQUE (sq_receita, documento_doador, nome_doador);
ALTER TABLE receita_doador_originario_orgao_partidario ADD CONSTRAINT uq_receita_doador_originario_orgao UNIQUE (sq_receita, documento_doador, nome_doador);
ALTER TABLE convenio                                  ADD CONSTRAINT uq_convenio_numero UNIQUE (numero_convenio);
