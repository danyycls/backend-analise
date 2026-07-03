CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE TABLE convenio (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    numero_convenio VARCHAR(50) NOT NULL,
    uf VARCHAR(2),
    codigo_siafi_municipio VARCHAR(20),
    nome_municipio VARCHAR(255),
    situacao_convenio VARCHAR(100),
    numero_original VARCHAR(100),
    numero_processo VARCHAR(100),
    objeto_convenio TEXT,
    codigo_orgao_superior VARCHAR(20),
    nome_orgao_superior VARCHAR(255),
    codigo_orgao_concedente VARCHAR(20),
    nome_orgao_concedente VARCHAR(255),
    codigo_ug_concedente VARCHAR(20),
    nome_ug_concedente VARCHAR(255),
    codigo_convenente VARCHAR(20),
    tipo_convenente VARCHAR(100),
    nome_convenente VARCHAR(255),
    tipo_ente_convenente VARCHAR(100),
    tipo_instrumento VARCHAR(100),
    valor_convenio NUMERIC(18, 2),
    valor_liberado NUMERIC(18, 2),
    data_publicacao DATE,
    data_inicio_vigencia DATE,
    data_final_vigencia DATE,
    valor_contrapartida NUMERIC(18, 2),
    data_ultima_liberacao DATE,
    valor_ultima_liberacao NUMERIC(18, 2),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_convenio_numero UNIQUE (numero_convenio)
);

CREATE INDEX idx_convenio_uf ON convenio (uf);
CREATE INDEX idx_convenio_nome_municipio ON convenio (nome_municipio);
CREATE INDEX idx_convenio_nome_convenente ON convenio (nome_convenente);
CREATE INDEX idx_convenio_tipo_instrumento ON convenio (tipo_instrumento);
CREATE INDEX idx_convenio_situacao ON convenio (situacao_convenio);
CREATE INDEX idx_convenio_objeto_trgm ON convenio USING gin (objeto_convenio gin_trgm_ops);
CREATE INDEX idx_convenio_valores ON convenio (valor_convenio, valor_liberado);
CREATE INDEX idx_convenio_datas ON convenio (data_publicacao, data_inicio_vigencia);
