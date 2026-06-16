CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE eleicao (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    codigo_tse INTEGER NOT NULL,
    ano SMALLINT NOT NULL,
    codigo_tipo_eleicao INTEGER,
    nome_tipo_eleicao VARCHAR(100),
    descricao VARCHAR(255) NOT NULL,
    data_eleicao DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_eleicao_codigo_tse UNIQUE (codigo_tse)
);

CREATE TABLE unidade_eleitoral (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sg_uf VARCHAR(2) NOT NULL,
    codigo_tse VARCHAR(16) NOT NULL,
    nome VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_unidade_eleitoral_uf_codigo UNIQUE (sg_uf, codigo_tse)
);

CREATE TABLE partido (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    numero SMALLINT NOT NULL,
    sigla VARCHAR(20) NOT NULL,
    nome VARCHAR(255) NOT NULL,
    federacao_codigo_tse BIGINT,
    federacao_sigla VARCHAR(50),
    federacao_nome VARCHAR(255),
    coligacao_codigo_tse BIGINT,
    coligacao_nome VARCHAR(255),
    coligacao_composicao TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_partido_numero UNIQUE (numero),
    CONSTRAINT uq_partido_sigla UNIQUE (sigla)
);

CREATE TABLE candidato (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sq_candidato BIGINT NOT NULL,
    eleicao_id UUID NOT NULL REFERENCES eleicao(id),
    sg_uf VARCHAR(2) NOT NULL,
    partido_id UUID REFERENCES partido(id),
    cargo_codigo INTEGER,
    cargo_nome VARCHAR(100),
    genero_descricao VARCHAR(100),
    cor_raca_descricao VARCHAR(100),
    estado_civil_nome VARCHAR(100),
    grau_instrucao_nome VARCHAR(150),
    ocupacao_codigo INTEGER,
    ocupacao_nome VARCHAR(255),
    numero_candidato INTEGER,
    cpf VARCHAR(11),
    cpf_vice VARCHAR(11),
    nome_completo VARCHAR(255) NOT NULL,
    nome_urna VARCHAR(255),
    nome_social VARCHAR(255),
    data_nascimento DATE,
    situacao_totalizacao_descricao VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_candidato_sq_candidato UNIQUE (sq_candidato)
);

CREATE TABLE bem_candidato (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    candidato_id UUID NOT NULL REFERENCES candidato(id),
    tipo_bem_codigo INTEGER,
    tipo_bem_nome VARCHAR(255),
    numero_ordem INTEGER NOT NULL,
    descricao TEXT NOT NULL,
    valor NUMERIC(18, 2) NOT NULL,
    data_ultima_atualizacao DATE,
    hora_ultima_atualizacao TIME,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_bem_candidato_ordem UNIQUE (candidato_id, numero_ordem)
);

CREATE TABLE fornecedor (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cpf_cnpj VARCHAR(14) NOT NULL,
    nome VARCHAR(255) NOT NULL,
    nome_rfb VARCHAR(255),
    tipo_fornecedor_codigo INTEGER,
    tipo_fornecedor_descricao VARCHAR(100),
    cnae_codigo VARCHAR(20),
    cnae_descricao VARCHAR(255),
    esfera_partidaria_codigo VARCHAR(10),
    esfera_partidaria_descricao VARCHAR(100),
    sg_uf VARCHAR(2),
    municipio_nome VARCHAR(255),
    sq_candidato_relacionado BIGINT,
    numero_candidato_relacionado INTEGER,
    cargo_codigo_relacionado INTEGER,
    cargo_descricao_relacionada VARCHAR(100),
    partido_numero_relacionado SMALLINT,
    partido_sigla_relacionado VARCHAR(20),
    partido_nome_relacionado VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_fornecedor_cpf_cnpj UNIQUE (cpf_cnpj)
);

CREATE TABLE doador (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cpf_cnpj VARCHAR(14) NOT NULL,
    nome VARCHAR(255) NOT NULL,
    nome_rfb VARCHAR(255),
    cnae_codigo VARCHAR(20),
    cnae_descricao VARCHAR(255),
    esfera_partidaria_codigo VARCHAR(10),
    esfera_partidaria_descricao VARCHAR(100),
    sg_uf VARCHAR(2),
    municipio_nome VARCHAR(255),
    sq_candidato_relacionado BIGINT,
    numero_candidato_relacionado INTEGER,
    cargo_codigo_relacionado INTEGER,
    cargo_descricao_relacionada VARCHAR(100),
    partido_numero_relacionado SMALLINT,
    partido_sigla_relacionado VARCHAR(20),
    partido_nome_relacionado VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_doador_cpf_cnpj UNIQUE (cpf_cnpj)
);

CREATE TABLE prestacao_contas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sq_prestador_contas BIGINT NOT NULL,
    eleicao_id UUID NOT NULL REFERENCES eleicao(id),
    candidato_id UUID REFERENCES candidato(id),
    partido_id UUID REFERENCES partido(id),
    sg_uf VARCHAR(2),
    unidade_eleitoral_id UUID REFERENCES unidade_eleitoral(id),
    tipo_prestador VARCHAR(30) NOT NULL,
    tipo_prestacao VARCHAR(30),
    data_prestacao DATE,
    turno SMALLINT,
    cnpj_prestador_conta VARCHAR(14),
    esfera_partidaria_codigo VARCHAR(10),
    esfera_partidaria_descricao VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_prestacao_contas_natural UNIQUE (tipo_prestador, eleicao_id, sq_prestador_contas),
    CONSTRAINT ck_prestacao_contas_tipo_prestador CHECK (tipo_prestador IN ('CANDIDATO', 'ORGAO_PARTIDARIO')),
    CONSTRAINT ck_prestacao_contas_titularidade CHECK (
        (tipo_prestador = 'CANDIDATO' AND candidato_id IS NOT NULL AND partido_id IS NULL) OR
        (tipo_prestador = 'ORGAO_PARTIDARIO' AND partido_id IS NOT NULL)
    )
);

CREATE TABLE despesa_candidato (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    prestacao_contas_id UUID NOT NULL REFERENCES prestacao_contas(id),
    candidato_id UUID NOT NULL REFERENCES candidato(id),
    fornecedor_id UUID REFERENCES fornecedor(id),
    sq_despesa BIGINT NOT NULL,
    tipo_registro VARCHAR(20) NOT NULL,
    tipo_documento VARCHAR(100),
    numero_documento VARCHAR(100),
    origem_despesa_codigo INTEGER,
    origem_despesa_descricao VARCHAR(255),
    fonte_despesa_codigo INTEGER,
    fonte_despesa_descricao VARCHAR(255),
    natureza_despesa_codigo INTEGER,
    natureza_despesa_descricao VARCHAR(255),
    especie_recurso_codigo INTEGER,
    especie_recurso_descricao VARCHAR(255),
    sq_parcelamento_despesa BIGINT,
    data_despesa DATE,
    descricao TEXT NOT NULL,
    valor NUMERIC(18, 2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_despesa_candidato_natural UNIQUE (sq_despesa, tipo_registro),
    CONSTRAINT ck_despesa_candidato_tipo_registro CHECK (tipo_registro IN ('CONTRATADA', 'PAGA'))
);

CREATE TABLE despesa_orgao_partidario (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    prestacao_contas_id UUID NOT NULL REFERENCES prestacao_contas(id),
    partido_id UUID NOT NULL REFERENCES partido(id),
    fornecedor_id UUID REFERENCES fornecedor(id),
    sq_despesa BIGINT NOT NULL,
    tipo_registro VARCHAR(20) NOT NULL,
    tipo_documento VARCHAR(100),
    numero_documento VARCHAR(100),
    origem_despesa_codigo INTEGER,
    origem_despesa_descricao VARCHAR(255),
    fonte_despesa_codigo INTEGER,
    fonte_despesa_descricao VARCHAR(255),
    natureza_despesa_codigo INTEGER,
    natureza_despesa_descricao VARCHAR(255),
    especie_recurso_codigo INTEGER,
    especie_recurso_descricao VARCHAR(255),
    sq_parcelamento_despesa BIGINT,
    data_despesa DATE,
    descricao TEXT NOT NULL,
    valor NUMERIC(18, 2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_despesa_orgao_partidario_natural UNIQUE (sq_despesa, tipo_registro),
    CONSTRAINT ck_despesa_orgao_partidario_tipo_registro CHECK (tipo_registro IN ('CONTRATADA', 'PAGA'))
);

CREATE TABLE receita_candidato (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    prestacao_contas_id UUID NOT NULL REFERENCES prestacao_contas(id),
    candidato_id UUID NOT NULL REFERENCES candidato(id),
    doador_id UUID REFERENCES doador(id),
    sq_receita BIGINT NOT NULL,
    fonte_receita_codigo INTEGER,
    fonte_receita_descricao VARCHAR(255),
    origem_receita_codigo INTEGER,
    origem_receita_descricao VARCHAR(255),
    natureza_receita_codigo INTEGER,
    natureza_receita_descricao VARCHAR(255),
    especie_receita_codigo INTEGER,
    especie_receita_descricao VARCHAR(255),
    numero_recibo_doacao VARCHAR(100),
    numero_documento_doacao VARCHAR(100),
    data_receita DATE,
    descricao TEXT NOT NULL,
    valor NUMERIC(18, 2) NOT NULL,
    natureza_recurso_estimavel TEXT,
    genero VARCHAR(100),
    cor_raca VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_receita_candidato_natural UNIQUE (sq_receita)
);

CREATE TABLE receita_orgao_partidario (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    prestacao_contas_id UUID NOT NULL REFERENCES prestacao_contas(id),
    partido_id UUID NOT NULL REFERENCES partido(id),
    doador_id UUID REFERENCES doador(id),
    sq_receita BIGINT NOT NULL,
    fonte_receita_codigo INTEGER,
    fonte_receita_descricao VARCHAR(255),
    origem_receita_codigo INTEGER,
    origem_receita_descricao VARCHAR(255),
    natureza_receita_codigo INTEGER,
    natureza_receita_descricao VARCHAR(255),
    especie_receita_codigo INTEGER,
    especie_receita_descricao VARCHAR(255),
    numero_recibo_doacao VARCHAR(100),
    numero_documento_doacao VARCHAR(100),
    data_receita DATE,
    descricao TEXT NOT NULL,
    valor NUMERIC(18, 2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_receita_orgao_partidario_natural UNIQUE (sq_receita)
);

CREATE TABLE receita_doador_originario_candidato (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    prestacao_contas_id UUID NOT NULL REFERENCES prestacao_contas(id),
    receita_candidato_id UUID REFERENCES receita_candidato(id),
    sq_receita BIGINT NOT NULL,
    documento_doador VARCHAR(14),
    nome_doador VARCHAR(255) NOT NULL,
    nome_doador_rfb VARCHAR(255),
    tipo_doador VARCHAR(100),
    cnae_codigo VARCHAR(20),
    cnae_descricao VARCHAR(255),
    data_receita DATE,
    descricao TEXT NOT NULL,
    valor NUMERIC(18, 2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_receita_doador_originario_candidato UNIQUE (sq_receita, documento_doador, nome_doador)
);

CREATE TABLE receita_doador_originario_orgao_partidario (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    prestacao_contas_id UUID NOT NULL REFERENCES prestacao_contas(id),
    receita_orgao_partidario_id UUID REFERENCES receita_orgao_partidario(id),
    sq_receita BIGINT NOT NULL,
    documento_doador VARCHAR(14),
    nome_doador VARCHAR(255) NOT NULL,
    nome_doador_rfb VARCHAR(255),
    tipo_doador VARCHAR(100),
    cnae_codigo VARCHAR(20),
    cnae_descricao VARCHAR(255),
    data_receita DATE,
    descricao TEXT NOT NULL,
    valor NUMERIC(18, 2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_receita_doador_originario_orgao UNIQUE (sq_receita, documento_doador, nome_doador)
);

CREATE TABLE arquivo_importado (
    caminho_relativo VARCHAR(500) PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    tipo VARCHAR(100) NOT NULL,
    uf VARCHAR(2) NOT NULL,
    total_registros INTEGER NOT NULL DEFAULT 0,
    criado_em TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes

CREATE INDEX idx_eleicao_codigo_tse ON eleicao (codigo_tse);
CREATE INDEX idx_unidade_eleitoral_sg_uf ON unidade_eleitoral (sg_uf);
CREATE INDEX idx_unidade_eleitoral_codigo_tse ON unidade_eleitoral (codigo_tse);
CREATE INDEX idx_partido_numero ON partido (numero);
CREATE INDEX idx_candidato_sq_candidato ON candidato (sq_candidato);
CREATE INDEX idx_candidato_cpf ON candidato (cpf);
CREATE INDEX idx_candidato_eleicao_uf ON candidato (eleicao_id, sg_uf);
CREATE INDEX idx_candidato_partido_id ON candidato (partido_id);
CREATE INDEX idx_bem_candidato_candidato_id ON bem_candidato (candidato_id);
CREATE INDEX idx_fornecedor_cpf_cnpj ON fornecedor (cpf_cnpj);
CREATE INDEX idx_fornecedor_sg_uf ON fornecedor (sg_uf);
CREATE INDEX idx_doador_cpf_cnpj ON doador (cpf_cnpj);
CREATE INDEX idx_doador_sg_uf ON doador (sg_uf);
CREATE INDEX idx_prestacao_contas_sq_prestador ON prestacao_contas (sq_prestador_contas);
CREATE INDEX idx_prestacao_contas_eleicao_tipo ON prestacao_contas (eleicao_id, tipo_prestador);
CREATE INDEX idx_despesa_candidato_sq_despesa ON despesa_candidato (sq_despesa);
CREATE INDEX idx_despesa_candidato_candidato_id ON despesa_candidato (candidato_id);
CREATE INDEX idx_despesa_candidato_fornecedor_id ON despesa_candidato (fornecedor_id);
CREATE INDEX idx_despesa_candidato_prestacao_id ON despesa_candidato (prestacao_contas_id);
CREATE INDEX idx_despesa_candidato_data ON despesa_candidato (data_despesa);
CREATE INDEX idx_despesa_orgao_partidario_sq_despesa ON despesa_orgao_partidario (sq_despesa);
CREATE INDEX idx_despesa_orgao_partidario_partido_id ON despesa_orgao_partidario (partido_id);
CREATE INDEX idx_despesa_orgao_partidario_fornecedor_id ON despesa_orgao_partidario (fornecedor_id);
CREATE INDEX idx_despesa_orgao_partidario_prestacao_id ON despesa_orgao_partidario (prestacao_contas_id);
CREATE INDEX idx_despesa_orgao_partidario_data ON despesa_orgao_partidario (data_despesa);
CREATE INDEX idx_receita_candidato_sq_receita ON receita_candidato (sq_receita);
CREATE INDEX idx_receita_candidato_candidato_id ON receita_candidato (candidato_id);
CREATE INDEX idx_receita_candidato_doador_id ON receita_candidato (doador_id);
CREATE INDEX idx_receita_candidato_prestacao_id ON receita_candidato (prestacao_contas_id);
CREATE INDEX idx_receita_candidato_data ON receita_candidato (data_receita);
CREATE INDEX idx_receita_orgao_partidario_sq_receita ON receita_orgao_partidario (sq_receita);
CREATE INDEX idx_receita_orgao_partidario_partido_id ON receita_orgao_partidario (partido_id);
CREATE INDEX idx_receita_orgao_partidario_doador_id ON receita_orgao_partidario (doador_id);
CREATE INDEX idx_receita_orgao_partidario_prestacao_id ON receita_orgao_partidario (prestacao_contas_id);
CREATE INDEX idx_receita_orgao_partidario_data ON receita_orgao_partidario (data_receita);
CREATE INDEX idx_receita_doador_originario_candidato_prestacao_id ON receita_doador_originario_candidato (prestacao_contas_id);
CREATE INDEX idx_receita_doador_originario_candidato_receita_id ON receita_doador_originario_candidato (receita_candidato_id);
CREATE INDEX idx_receita_doador_originario_candidato_sq ON receita_doador_originario_candidato (sq_receita);
CREATE INDEX idx_receita_doador_originario_candidato_doc ON receita_doador_originario_candidato (documento_doador);
CREATE INDEX idx_receita_doador_originario_candidato_data ON receita_doador_originario_candidato (data_receita);
CREATE INDEX idx_receita_doador_originario_orgao_prestacao_id ON receita_doador_originario_orgao_partidario (prestacao_contas_id);
CREATE INDEX idx_receita_doador_originario_orgao_receita_id ON receita_doador_originario_orgao_partidario (receita_orgao_partidario_id);
CREATE INDEX idx_receita_doador_originario_orgao_sq ON receita_doador_originario_orgao_partidario (sq_receita);
CREATE INDEX idx_receita_doador_originario_orgao_doc ON receita_doador_originario_orgao_partidario (documento_doador);
CREATE INDEX idx_receita_doador_originario_orgao_data ON receita_doador_originario_orgao_partidario (data_receita);
CREATE INDEX idx_arquivo_importado_tipo ON arquivo_importado (tipo);
CREATE INDEX idx_arquivo_importado_uf ON arquivo_importado (uf);
