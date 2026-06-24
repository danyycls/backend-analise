# Ligação Política

## Descrição da rota

A rota **`POST /busca/contexto`** cruza documentos de licitações e contratos públicos (tipicamente vindo do PNCP) com dados eleitorais do TSE e sanções do TCU, para descobrir **vínculos entre fornecedores/licitantes do Estado e o financiamento de campanhas eleitorais**.

### Importância da análise

A análise detecta padrões clássicos de **conflito de interesse e captura do Estado**: um fornecedor que venceu uma licitação pública e, ao mesmo tempo, financiou — direta ou indiretamente — a campanha do agente público que aprovou o contrato. É um indicador forte de corrupção eleitoral-licitatória.

Combinando os dados eleitorais com o TCU, a rota também sinaliza **red flags regulatórias**: se o mesmo documento (CPF ou CNPJ) é inidôneo (impedido de licitar), inabilitado (impedido de cargo em comissão) ou tem contas julgadas irregulares. Esses flags podem impedir a contratação ou justificar auditoria.

### Como é montado

A entrada é uma lista de licitações, cada uma com o `cpf_cnpj` vencedor e seu quadro societário (`socios`). O fluxo da análise é:

```
POST /busca/contexto (JSON)
        │
        ▼
┌───────────────┐
│   Handler     │  1. Bind JSON (valida array `licitacoes`)
│               │  2. Consulta cache Redis (chave por md5 do body)
│               │  3. Cache HIT  → 200 direto
└──────┬────────┘  4. Cache MISS → delega ao UseCase
       ▼
┌───────────────┐
│   UseCase     │  Para cada licitação:
│               │   - normaliza documento (remove não-dígitos, detecta `*` = parcial)
│               │   - processa sócios (origem="socio") e CNPJ principal (origem="principal")
└─┬───┬───┬─────┘
  │   │   │  Banco TSE (PostgreSQL):
  │   │   │   - fornecedor (com despesas_candidato / despesa_orgao_partidario)
  │   │   │   - doador + receita_candidato + receita_orgao_partidario
  ▼   ▼   ▼
┌───────────────┐
│  Enriquece    │  Em paralelo (max 5 goroutines):
│  Fornecedores │   - OpenCNPJ.Buscar(cnpj) → QSA, situação RFB, capital social
└──────┬────────┘   (apenas para vínculos "fornecedor" com CNPJ de 14 dígitos)
       ▼
┌───────────────┐
│  Enriquece    │  Em paralelo (max 5 goroutines):
│  TCU          │   - BuscarContasIrregulares (CPF e/ou CNPJ)
│               │   - BuscarInabilitados (só CPF)
│               │   - BuscarInidoneos (CPF e/ou CNPJ)
└──────┬────────┘   (apenas se retornar > 0 registros)
       ▼
┌───────────────────┐
│  Enriquece        │  Em paralelo (max 5 goroutines):
│  Servidor Público │   - ListarServidores (só CPF, via Portal da Transparência)
└────────┬──────────┘   (apenas se retornar > 0 registros)
         ▼
┌───────────────┐
│   Redis       │  Grava resposta no cache (TTL 30 dias)
│               │  Chave: podp:cache:ligacao-politica:<md5(body)>
└───────────────┘
       ▼
   200 OK + AnalisarLigacaoPoliticaResponse
```

Passos resumidos:

1. **Normalização**: o documento é limpo (remove pontuação); se contém `*`, vira busca parcial por nome.
2. **Busca no TSE (PostgreSQL)**: consulta as tabelas `fornecedor`, `doador`, `receita_candidato`, `receita_orgao_partidario` e tabelas dependentes. Para CPFs, busca também a variante com prefixo `000` (formato usado pelo TSE para CPF em CNPJ).
3. **Enriquecimento OpenCNPJ**: para cada vínculo de fornecedor com CNPJ de 14 dígitos, busca dados cadastrais (QSA, razão social, situação na RFB, capital social).
4. **Enriquecimento TCU**: para todos os documentos (principal, sócios, fornecedores e doadores encontrados), consulta sanções em paralelo. Só adiciona vínculo se houver registro.
5. **Enriquecimento Servidor Público**: para cada CPF encontrado (principal, sócios, fornecedores e doadores), consulta o **Portal da Transparência** (`/api-de-dados/servidores`) para verificar se a pessoa é servidor público federal. Só adiciona vínculo se houver registro.
6. **Cache Redis**: a resposta final é gravada em cache (chave derivada do body) com TTL de 30 dias. Falha no Redis é tratada como *warning* e **não aborta** a requisição.

### O que é analisado

Cada documento processado recebe um array de `vinculos`, cada um com um `tipo` indicando a ligação política encontrada:

| `tipo` | Significado | Detalhe anexado em `detalhes` |
|--------|-------------|-------------------------------|
| `fornecedor` | Documento é fornecedor de campanha eleitoral | `FornecedorDetalhado` (com despesas de candidato/partido + enriquecimento OpenCNPJ) |
| `doador` | Documento é doador registrado no TSE | `Doador` |
| `receita_candidato` | Doação efetiva a um candidato | `[]ReceitaCandidato` (valor, data, descrição, `sq_receita`) |
| `receita_orgao_partidario` | Doação efetiva a um partido/órgão partidário | `[]ReceitaOrgaoPartidario` |
| `tcu_contas_irregulares` | Contas julgadas irregulares no TCU | `[]ContasIrregulares` |
| `tcu_inabilitado` | Inabilitado para cargo em comissão (TCU) | `[]Sancoes` |
| `tcu_inidoneo` | Licitante inidôneo (TCU) | `[]Sancoes` |
| `servidor_publico` | CPF registrado como servidor público federal no Portal da Transparência | `[]CadastroServidor` |

Cada vínculo traz ainda uma `descricao` textual legível por humano (ex.: `"11222333000181 é fornecedor de campanha com 3 despesa(s) de candidato e 1 despesa(s) partidária(s)"`), gerada em `analisar-ligacao-politica-usecase.go:401-423`.

O campo `origem` do documento indica a proveniência na licitação:

| `origem` | Quando |
|----------|--------|
| `principal` | O `cpf_cnpj` vencedor da licitação |
| `socio` | O documento de um sócio listado em `socios` |

### Clients consultados e análises que fazem

| Client | Tipo | Métodos usados | Análise feita |
|--------|------|----------------|---------------|
| **TSE (PostgreSQL)** | Repositório interno (não-HTTP) | `FornecedoresBuscarPorDocumento`, `DoadoresBuscarPorDocumento`, `ReceitasCandidatoBuscarPorDoadorID`, `ReceitasOrgaoBuscarPorDoadorID`, `DespesasCandidatoBuscarPorFornecedorID`, `DespesasPartidoBuscarPorFornecedorID` | Verifica se o documento aparece como fornecedor/doador de campanha e lista as doações/despesas correspondentes |
| **OpenCNPJ** | HTTP — `internal/shared/clients/opencnpj` | `Buscar(ctx, cnpj)` | Enriquece vínculos `fornecedor` com CNPJ de 14 dígitos: QSA, razão social, nome fantasia, situação cadastral RFB, capital social |
| **TCU** | HTTP — `internal/shared/clients/tcu` | `BuscarContasIrregulares`, `BuscarInabilitados`, `BuscarInidoneos` (não usa `BuscarFinsEleitorais` nesta rota) | Verifica sanções TCU para cada documento (principal, sócios, fornecedores e doadores) |
| **Portal da Transparência** | HTTP — `internal/shared/clients/portaltransparencia` | `ListarServidores` | Verifica se cada CPF encontrado é servidor público federal |
| **Redis** | Cache — `internal/shared/redis` | `Get`, `Set` (chave `podp:cache:ligacao-politica:<md5>`, TTL 30 dias) | Evita reprocessar o mesmo body; falha não aborta a requisição |

> Documentação dos clients: **[docs/clientes/tcu.md](./clientes/tcu.md)**, **[docs/clientes/opencnpj.md](./clientes/opencnpj.md)** (este último ainda como placeholder).
> Documentação do banco TSE: **[docs/db-tse.md](./db-tse.md)**.
> Documentação de como o banco é populado: **[docs/tse-importacao.md](./tse-importacao.md)**.

---

## URL da rota

```
POST /busca/contexto
Content-Type: application/json
```

Rota registrada em `internal/app/routes.go:22`. Handler: `AnalisarLigacaoPoliticaHandler.Analisar` em `internal/ligacao-politica/handler/analisar-ligacao-politica-handler.go:23`.

---

## Request

O body é um objeto com o array `licitacoes` (obrigatório, não pode ser vazio). Cada licitação traz o `numero_controle_pncp` (identificador PNCP do contrato), o `cpf_cnpj` vencedor e a lista de `socios`.

### Estrutura (binding)

```go
struct {
    Licitacoes []AnalisarLigacaoPoliticaRequest `json:"licitacoes" binding:"required"`
}

type AnalisarLigacaoPoliticaRequest struct {
    NumeroControlePncp string                                `json:"numero_controle_pncp"`
    CpfCnpj            string                                `json:"cpf_cnpj"`
    Socios             []AnalisarLigacaoPoliticaSocioRequest `json:"socios"`
}

type AnalisarLigacaoPoliticaSocioRequest struct {
    Nome      string `json:"nome"`
    Documento string `json:"documento"`
}
```

> **Regras de input**:
> - `licitacoes` é obrigatório (`binding:"required"`) e não pode ser array vazio (validação explícita no handler).
> - Documento pode conter pontuação (será normalizado) ou `*` para busca parcial por nome (nesse caso `nome` deve ser informado).
> - Documentos com menos de 3 dígitos são ignorados (não geram `DocumentoVinculo`).
> - Para CPFs, o use case busca automaticamente também a variante com prefixo `000` (formato TSE para CPF em CNPJ).

### Exemplo

```json
{
  "licitacoes": [
    {
      "numero_controle_pncp": "1A2B3C-2024-SP-001",
      "cpf_cnpj": "11222333000181",
      "socios": [
        { "nome": "João da Silva", "documento": "11122233344" },
        { "nome": "Maria Souza", "documento": "00011122233344" }
      ]
    }
  ]
}
```

---

## Response

Retorna `AnalisarLigacaoPoliticaResponse` com a quantidade de documentos processados e, para cada licitação, os vínculos encontrados por documento.

### Estrutura

```go
type AnalisarLigacaoPoliticaResponse struct {
    DocumentosProcessados int                `json:"documentos_processados"`
    Resultados            []VinculoLicitacao `json:"resultados"`
}

type VinculoLicitacao struct {
    NumeroControlePncp string             `json:"numero_controle_pncp"`
    CpfCnpj            string             `json:"cpf_cnpj"`
    Socios             []SocioOutput      `json:"socios,omitempty"`
    Documentos         []DocumentoVinculo `json:"documentos,omitempty"`
}

type DocumentoVinculo struct {
    DocumentoInput       string    `json:"documento_input"`
    DocumentoNormalizado string    `json:"documento_normalizado"`
    Nome                 string    `json:"nome"`
    Parcial              bool      `json:"parcial"`
    Origem               string    `json:"origem"`
    Vinculos             []Vinculo `json:"vinculos,omitempty"`
}

type Vinculo struct {
    Tipo      string           `json:"tipo"`
    Descricao string           `json:"descricao"`
    Detalhes  *VinculoDetalhes `json:"detalhes,omitempty"`
}
```

### Códigos HTTP

| HTTP | Quando | Body |
|------|--------|------|
| **200** | Sucesso (cache hit ou execução completa) | `AnalisarLigacaoPoliticaResponse` |
| **400** | JSON inválido ou campo `licitacoes` ausente | `{"erro": "corpo inválido: <detalhe>"}` |
| **400** | `licitacoes` é array vazio | `{"erro": "array licitacoes é obrigatório"}` |
| **500** | Erro interno no UseCase | `{"erro": "<mensagem>"}` |

> Na prática o UseCase trata erros de repositório/clients de forma tolerante (engole erros e segue), então o 500 é raro — aparece principalmente em falhas catastróficas de banco.

### Exemplo de sucesso (200)

Cenário: CNPJ principal é fornecedor de campanha com despesas, tem sanção de inidôneo no TCU, e um sócio (CPF) é doador com receita a candidato e contas irregulares no TCU.

```json
{
  "documentos_processados": 1,
  "resultados": [
    {
      "numero_controle_pncp": "1A2B3C-2024-SP-001",
      "cpf_cnpj": "11222333000181",
      "socios": [
        { "nome": "João da Silva", "documento": "11122233344" },
        { "nome": "Maria Souza", "documento": "00011122233344" }
      ],
      "documentos": [
        {
          "documento_input": "11222333000181",
          "documento_normalizado": "11222333000181",
          "nome": "",
          "parcial": false,
          "origem": "principal",
          "vinculos": [
            {
              "tipo": "fornecedor",
              "descricao": "11222333000181 é fornecedor de campanha com 1 despesa(s) de candidato e 1 despesa(s) partidária(s)",
              "detalhes": {
                "fornecedor": {
                  "fornecedor": {
                    "cpf_cnpj": "11222333000181",
                    "nome": "Fornecedor Teste Ltda",
                    "tipo_fornecedor_descricao": "JURÍDICA",
                    "sg_uf": "SP",
                    "enriquecimento": {
                      "cnpj": "11222333000181",
                      "razaoSocial": "Empresa Teste Ltda",
                      "nomeFantasia": "Teste Fantasia",
                      "situacaoCadastral": "ATIVA",
                      "capitalSocial": "10000.00",
                      "qsa": [
                        { "nome_socio": "João da Silva", "cnpj_cpf_socio": "11122233344" }
                      ]
                    }
                  },
                  "despesas_candidato": [
                    {
                      "despesa": { "sq_despesa": 9001, "valor": 5000.00, "descricao": "Doação estimada" },
                      "sq_candidato": 7000001
                    }
                  ],
                  "despesas_orgao_partidario": [
                    {
                      "despesa": { "sq_despesa": 8001, "valor": 2000.00 },
                      "partido_numero": 13,
                      "partido_nome": "PARTIDO X"
                    }
                  ]
                }
              }
            },
            {
              "tipo": "tcu_inidoneo",
              "descricao": "1 registro(s) de inidôneo no TCU",
              "detalhes": {
                "inidoneos": [
                  {
                    "numeroProcessoFormatado": "0005678-90.2021.7.00.0000",
                    "nome": "Fornecedor Teste Ltda",
                    "dataFinalSancao": "31/12/2028"
                  }
                ]
              }
            }
          ]
        },
        {
          "documento_input": "11122233344",
          "documento_normalizado": "11122233344",
          "nome": "João da Silva",
          "parcial": false,
          "origem": "socio",
          "vinculos": [
            {
              "tipo": "doador",
              "descricao": "Doador registrado: João da Silva",
              "detalhes": {
                "doador": {
                  "cpf_cnpj": "11122233344",
                  "nome": "João da Silva",
                  "sg_uf": "SP"
                }
              }
            },
            {
              "tipo": "receita_candidato",
              "descricao": "Doação (sq_receita=12345) a candidato",
              "detalhes": {
                "receitas_candidato": [
                  { "sq_receita": 12345, "valor": 2500.00, "descricao": "Doação de pessoa física", "data_receita": "2024-09-15T00:00:00Z" }
                ]
              }
            },
            {
              "tipo": "tcu_contas_irregulares",
              "descricao": "1 registro(s) de contas julgadas irregulares no TCU",
              "detalhes": {
                "contas_irregulares": [
                  {
                    "numeroProcessoFormatado": "0001234-56.2020.7.00.0000",
                    "nome": "João da Silva",
                    "uf": "SP",
                    "dataTransitoEmJulgado": "10/05/2022"
                  }
                ]
              }
            },
            {
              "tipo": "tcu_inabilitado",
              "descricao": "1 registro(s) de inabilitado no TCU",
              "detalhes": {
                "inabilitados": [
                  { "nome": "João da Silva", "uf": "SP" }
                ]
              }
            }
          ]
        }
      ]
    }
  ]
}
```

### Exemplo sem correspondência (200)

Documento sem nenhum vínculo no TSE/TCU retorna com `vinculos` ausente (omitido por `omitempty`):

```json
{
  "documentos_processados": 1,
  "resultados": [
    {
      "numero_controle_pncp": "pncp-003",
      "cpf_cnpj": "00000000000000",
      "documentos": [
        {
          "documento_input": "00000000000000",
          "documento_normalizado": "00000000000000",
          "nome": "",
          "parcial": false,
          "origem": "principal"
        }
      ]
    }
  ]
}
```

### Exemplos de erro

**400 — corpo inválido** (JSON malformado ou `licitacoes` ausente):
```json
{ "erro": "corpo inválido: Key: 'Licitacoes' Error:Field validation for 'Licitacoes' failed on the 'required' tag" }
```

**400 — array vazio:**
```json
{ "erro": "array licitacoes é obrigatório" }
```

**500 — erro interno:**
```json
{ "erro": "erro ao consultar banco de dados: ..." }
```

---

## Referências

- **Clients:** [docs/clientes/tcu.md](./clientes/tcu.md), [docs/clientes/opencnpj.md](./clientes/opencnpj.md)
- **Banco TSE:** [docs/db-tse.md](./db-tse.md) | **Importação:** [docs/tse-importacao.md](./tse-importacao.md)
- **Código-fonte:** `internal/ligacao-politica/`
  - `handler/analisar-ligacao-politica-handler.go` — endpoint Gin (bind, cache, delega ao UseCase)
  - `usecase/analisar-ligacao-politica-usecase.go` — orquestração (busca TSE, enriquece OpenCNPJ + TCU em paralelo)
  - `usecase/ligacao_politica_types.go` — DTOs de request/response
  - `usecase/interface.go` — interface do UseCase (fonte do mockgen)
  - `testutils/builder.go` — builder de fixtures para testes
- **Repositório TSE compartilhado:** `internal/esferas-brasileiras/tse/repositorio/`
- **Clients Go:** `internal/shared/clients/{opencnpj,tcu}/`
- **Cache:** `internal/shared/redis/` (chave `podp:cache:ligacao-politica:<md5>`, TTL 30 dias)
