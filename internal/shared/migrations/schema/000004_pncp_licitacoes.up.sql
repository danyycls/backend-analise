CREATE TABLE IF NOT EXISTS amparo_legal (
    codigo INTEGER PRIMARY KEY,
    nome TEXT NOT NULL,
    descricao TEXT
);

CREATE TABLE IF NOT EXISTS licitacao_fornecedor (
    cnpj VARCHAR(14) PRIMARY KEY,
    razao_social TEXT NOT NULL,
    dados_completos JSONB
);

CREATE TABLE IF NOT EXISTS socio (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cnpj_cpf_socio VARCHAR(14) NOT NULL UNIQUE,
    nome_socio TEXT
);
CREATE INDEX IF NOT EXISTS idx_socio_cpf ON socio (cnpj_cpf_socio);

CREATE TABLE IF NOT EXISTS fornecedor_socio (
    cnpj_fornecedor VARCHAR(14) NOT NULL REFERENCES licitacao_fornecedor(cnpj) ON DELETE CASCADE,
    socio_id UUID NOT NULL REFERENCES socio(id) ON DELETE CASCADE,
    data_entrada_sociedade TEXT,
    identificador_socio TEXT,
    nome_socio TEXT,
    qualificacao_socio TEXT,
    nome_representante TEXT,
    qualificacao_representante TEXT,
    representante_legal TEXT,
    faixa_etaria TEXT,
    pais_codigo TEXT,
    pais_descricao TEXT,
    PRIMARY KEY (cnpj_fornecedor, socio_id)
);
CREATE INDEX IF NOT EXISTS idx_fornecedor_socio_socio ON fornecedor_socio (socio_id);

CREATE TABLE IF NOT EXISTS licitacao_contrato (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    numero_controle_pncp VARCHAR(255) UNIQUE NOT NULL,
    cnpj_orgao VARCHAR(14) NOT NULL,
    ug_uf_sigla VARCHAR(2),
    ug_codigo_ibge VARCHAR(7),
    data_publicacao_pncp DATE,
    data_assinatura DATE,
    data_inicio_vigencia DATE,
    data_termino_vigencia DATE,
    valor_global NUMERIC(18,2),
    valor_inicial NUMERIC(18,2),
    valor_total_estimado NUMERIC(18,2),
    valor_total_homologado NUMERIC(18,2),
    ni_fornecedor VARCHAR(14),
    codigo_amparo_legal INTEGER REFERENCES amparo_legal(codigo),
    numero_contrato VARCHAR(50),
    codigo_contrato VARCHAR(50),
    codigo_tipo_contrato INTEGER,
    tipo_contrato_nome TEXT,
    codigo_ug VARCHAR(20),
    nome_ug TEXT,
    ug_municipio_nome TEXT,
    ug_uf_nome TEXT,
    modalidade_nome TEXT,
    codigo_orgao VARCHAR(20),
    nome_orgao TEXT,
    nome_orgao_sub TEXT,
    objeto_contrato TEXT,
    numero_licitacao VARCHAR(50),
    origem_licitacao TEXT,
    produto TEXT,
    subtipo_contrato TEXT,
    ano_contrato INTEGER,
    nome_razao_social_fornecedor TEXT,
    dados_completos JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_contrato_cnpj_orgao ON licitacao_contrato (cnpj_orgao);
CREATE INDEX IF NOT EXISTS idx_contrato_uf ON licitacao_contrato (ug_uf_sigla);
CREATE INDEX IF NOT EXISTS idx_contrato_codigo_ibge ON licitacao_contrato (ug_codigo_ibge);
CREATE INDEX IF NOT EXISTS idx_contrato_data_publicacao ON licitacao_contrato (data_publicacao_pncp);
CREATE INDEX IF NOT EXISTS idx_contrato_ano ON licitacao_contrato (ano_contrato);
CREATE INDEX IF NOT EXISTS idx_contrato_uf_ibge ON licitacao_contrato (ug_uf_sigla, ug_codigo_ibge);

CREATE TABLE IF NOT EXISTS licitacao_busca_controle (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tipo_busca VARCHAR(20) NOT NULL,
    valor_busca VARCHAR(20) NOT NULL,
    ano INTEGER NOT NULL,
    mes INTEGER NOT NULL,
    data_inicial DATE NOT NULL,
    data_final DATE NOT NULL,
    total_contratos_encontrados INTEGER NOT NULL DEFAULT 0,
    ultima_atualizacao TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_busca UNIQUE (tipo_busca, valor_busca, ano, mes)
);

CREATE INDEX IF NOT EXISTS idx_busca_controle_data ON licitacao_busca_controle (ultima_atualizacao);
