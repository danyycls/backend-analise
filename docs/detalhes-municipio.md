# Detalhes de Município

## Descrição da rota

A rota **`GET /municipio/:codigoIBGE/detalhes/stream`** agrega, em tempo real via **WebSocket**, um painel fiscal-financeiro completo de um município brasileiro identificado por seu código IBGE. Combina dados de **3 relatórios oficiais do SICONFI** (Tesouro Nacional) com as **contratações públicas publicadas no PNCP**, entregues ao cliente como frames JSON `{"type","data"}` conforme cada seção fica pronta.

### Importância da análise

A rota permite que um cidadão, jornalista ou auditor veja, numa única conexão, a **saúde fiscal**, o **cumprimento de limites constitucionais** (15% saúde, 25% educação), a **execução orçamentária** e os **contratos firmados** por uma prefeitura — sem precisar consultar manualmente o SICONFI e o PNCP separadamente. É o equivalente municipal da rota `/estado/:uf/financeiro/stream`.

Como os relatórios SICONFI são tabulares e rotulados (não estruturados), o UseCase faz **parsing heurístico** (`strings.Contains` em `Rotulo`/`Conta`/`Coluna`) para extrair valores específicos de cada anexo e montar os structs de resposta.

### Como é montado

O handler faz o upgrade WebSocket, dispara **10 goroutines em paralelo** (uma por seção do painel), cada uma chama um método do UseCase e envia o resultado como frame JSON em um canal buffer 20. Quando todas terminam, emite opcionalmente `erro` (se o SICONFI caiu) e sempre `concluido`.

```
GET /municipio/:codigoIBGE/detalhes/stream?exercicio=2023 (WebSocket)
        │
        ▼
┌───────────────┐
│   Handler     │  1. Valida codigoIBGE (int > 0)
│               │  2. Lê query opcional exercicio (default ano atual - 1)
│               │  3. Upgrade WebSocket (aceita qualquer origem)
└──────┬────────┘  4. 10 goroutines em paralelo (WaitGroup + canal buffer 20)
       ▼
┌─────────────────────────────────────────────────────────────┐
│ 10 goroutines (cada uma com recover() + envia frame WS):    │
│                                                             │
│  SICONFI (9 seções):                PNCP (1 seção):         │
│   - BuscarRGF  (Anexo 02, 05, 06)   - BuscarContratacoes    │
│   - BuscarRREO (Anexo 05, 07,         PorMunicipio          │
│                  08, 09, 10)          (página 1, tam 20)    │
│   - BuscarDCA  (Anexo I-AB)                                 │
└──────┬──────────────────────────────────────────────────────┘
       ▼
┌───────────────┐
│  Canal → WS   │  Frames {"type":"...","data":...} (TextMessage)
│               │  Ordem NÃO determinística para os 10 de dados
└──────┬────────┘
       ▼
┌───────────────┐
│   Final       │  wg.Wait():
│               │   - se SICONFIIndisponivel() → frame "erro"
│               │   - sempre → frame "concluido"
│               │   - close(ch) encerra o loop
└───────────────┘
```

Passos resumidos:

1. **Validação**: `codigoIBGE` deve ser inteiro > 0 (espera-se 7 dígitos IBGE; o handler não faz padding — responsabilidade do caller). `exercicio` default = ano atual - 1.
2. **Upgrade WS**: `gorilla/websocket` com `CheckOrigin` permissivo.
3. **10 buscas paralelas**: 9 ao SICONFI (RGF, RREO, DCA em vários anexos) + 1 ao PNCP. Cada uma envia seu frame ao canal assim que termina.
4. **Frames WS**: o loop principal lê do canal e escreve `ws.WriteJSON` em `conn`. Respeita `ctx.Done()` (cliente desconecta → handler encerra).
5. **Encerramento**: após `wg.Wait()`, emite `erro` (se SICONFI indisponível) e `concluido`, fecha o canal e a conexão.

### O que é analisado

Cada frame WS tem `type` e `data`. Os 12 tipos enviados:

| `type` | `data` | Origem | Descrição |
|--------|--------|--------|-----------|
| `divida_consolidada` | `*DividaConsolidada` ou `null` | SICONFI RGF-Anexo 02 | Dívida Consolidada Líquida (DCL), % RCL, limite legal LRF |
| `disponibilidade_caixa` | `*DisponibilidadeCaixa` ou `null` | SICONFI RGF-Anexo 05 | Caixa vinculada vs. não vinculada |
| `restos_a_pagar` | `*RestosAPagar` ou `null` | SICONFI RGF-Anexo 06 | Inscritos, pagos, cancelados |
| `gasto_saude` | `*GastoSaude` ou `null` | SICONFI RREO-Anexo 09 | Valor total e % aplicado (limite constitucional 15%) |
| `gasto_educacao` | `*GastoEducacao` ou `null` | SICONFI RREO-Anexo 10 | Valor total e % aplicado (limite constitucional 25%) |
| `fundeb` | `*FundebResumo` ou `null` | SICONFI RREO-Anexo 08 | Receita e despesa totais do FUNDEB |
| `balanco_patrimonial` | `*BalancoPatrimonial` ou `null` | SICONFI DCA-Anexo I-AB | Ativo/passivo (circ/não circ), patrimônio líquido |
| `despesas_por_grupo` | `{"dados": []DespesaPorGrupoItem}` | SICONFI RREO-Anexo 05 | Pessoal, Juros, Investimentos, Corrente, Capital — empenhado/liquidado/pago |
| `transferencias` | `{"dados": []TransferenciaItem}` | SICONFI RREO-Anexo 07 | Transferências recebidas por órgão/origem (FUNDEB, FPM, SUS, etc.) |
| `contratos` | `{"dados": []pncp.Contrato}` | PNCP `/contratacoes/publicacao` | Contratações publicadas pelo município no ano |
| `erro` | `{"erro": "..."}` | — | Condicional: só se SICONFI indisponível |
| `concluido` | `null` | — | Sempre o último frame |

> **Assimetria**: as 7 primeiras seções enviam o struct direto em `data`; as 3 últimas (`despesas_por_grupo`, `transferencias`, `contratos`) embrulham em `{"dados": ...}`. Seção sem dados envia `data: null` (ou `{"dados": null}`).
>
> **Ordem**: não-determinística para os 10 frames de dados (concorrência). Garantido: `erro` (se houver) e `concluido` vêm **por último**, nessa ordem.

### Clients consultados e análises que fazem

| Client | Tipo | Métodos usados | Análise feita |
|--------|------|----------------|---------------|
| **SICONFI** | HTTP — `internal/shared/clients/siconfi` | `BuscarRGF`, `BuscarRREO`, `BuscarDCA` | Extrai indicadores fiscais (DCL, caixa, restos a pagar) e execução orçamentária (despesas por grupo, transferências, FUNDEB, gastos em saúde/educação, balanço patrimonial) dos relatórios RGF, RREO e DCA |
| **PNCP** | HTTP — `internal/shared/clients/pncp` | `BuscarContratacoesPorMunicipio` | Lista contratações publicadas pelo município no ano (página 1, tamanho 20, modalidade default `"8"`) |

> Documentação dos clients: **[docs/clientes/siconfi.md](./clientes/siconfi.md)**, **[docs/clientes/pncp.md](./clientes/pncp.md)**.
> Esta rota **não usa PostgreSQL nem Redis** — é stateless e 100% via APIs externas.

---

## URL da rota

```
GET /municipio/:codigoIBGE/detalhes/stream?exercicio=2023
Upgrade: websocket
Connection: Upgrade
```

Rota registrada em `internal/app/routes.go:146-148`. Handler: `EsferaMunicipalBuscarDetalhesWSHandler.BuscarDetalhesMunicipioWS` em `internal/esferas-brasileiras/municipal/handler/detalhes-municipio-ws-handler.go:30`.

---

## Request

A rota usa **path param** + 1 **query param opcional**. Não há body. O header WebSocket é o padrão (`Upgrade: websocket`, `Connection: Upgrade`, `Sec-WebSocket-Key`, `Sec-WebSocket-Version: 13`); o `Upgrader` aceita qualquer `Origin`.

| Parâmetro | Local | Obrigatório | Descrição |
|-----------|-------|-------------|-----------|
| `codigoIBGE` | path | sim | Código IBGE do município (7 dígitos, ex: `3550308` = São Paulo/SP). Não há normalização/padding — o caller deve enviar os 7 dígitos corretos |
| `exercicio` | query | não | Ano-base (ex: `2023`). Se ausente/`0`/negativo, default = `ano atual - 1` |

> O handler só valida que `codigoIBGE` é inteiro > 0. Zeros à esquerda funcionam via `strconv.Atoi` (ex: `0520870` → `520870`).

### Exemplo

```
GET /municipio/3550308/detalhes/stream?exercicio=2023
```

---

## Response (frames WebSocket)

Após o upgrade (HTTP **101 Switching Protocols**), o servidor envia frames `TextMessage` no formato:

```go
type wsMsgMuni struct {
    Type string      `json:"type"`
    Data interface{} `json:"data"`
}
```

O fluxo é unidirecional servidor→cliente (o handler não lê mensagens do cliente). O loop principal respeita `ctx.Done()` — se o cliente desconectar, o handler encerra.

### Códigos HTTP

| HTTP | Quando | Body |
|------|--------|------|
| **400** | `codigoIBGE` não é inteiro ou `<= 0` | `{"erro": "codigo IBGE invalido"}` |
| **101** | Upgrade WS com sucesso | — (a partir daqui, frames WS) |

> Falha no upgrade é tratada pelo gorilla (escreve a resposta HTTP apropriada) e o handler simplesmente retorna.

### Frames de dados

**1. `divida_consolidada`** — Dívida Consolidada Líquida (SICONFI RGF-Anexo 02):
```json
{"type":"divida_consolidada","data":{"valor_dcl":125000000.00,"percentual_rcl":18.5,"limite_legal":120,"periodo":"2023"}}
```

**2. `disponibilidade_caixa`** — (SICONFI RGF-Anexo 05):
```json
{"type":"disponibilidade_caixa","data":{"vinculada":8500000.00,"nao_vinculada":3200000.00,"periodo":"2023"}}
```

**3. `restos_a_pagar`** — (SICONFI RGF-Anexo 06):
```json
{"type":"restos_a_pagar","data":{"inscritos":4200000.00,"pagos":1800000.00,"cancelados":300000.00,"periodo":"2023"}}
```

**4. `gasto_saude`** — (SICONFI RREO-Anexo 09, limite 15%):
```json
{"type":"gasto_saude","data":{"valor_total":67000000.00,"percentual_aplicado":21.2,"limite_constitutional":15,"periodo":"2023"}}
```

> O campo `limite_constitutional` (com "ti") é o JSON tag real do struct `GastoSaude` em `municipio.go:103` — manter como está.

**5. `gasto_educacao`** — (SICONFI RREO-Anexo 10, limite 25%):
```json
{"type":"gasto_educacao","data":{"valor_total":95000000.00,"percentual_aplicado":28.7,"limite_constitutional":25,"periodo":"2023"}}
```

**6. `fundeb`** — FUNDEB (SICONFI RREO-Anexo 08):
```json
{"type":"fundeb","data":{"receita_total":43000000.00,"despesa_total":41000000.00,"periodo":"2023"}}
```

**7. `balanco_patrimonial`** — (SICONFI DCA-Anexo I-AB):
```json
{"type":"balanco_patrimonial","data":{"ativo_circulante":12000000.00,"ativo_nao_circulante":450000000.00,"passivo_circulante":8000000.00,"passivo_nao_circulante":380000000.00,"patrimonio_liquido":74000000.00,"periodo":"2023"}}
```

**8. `despesas_por_grupo`** — embrulhado em `dados` (SICONFI RREO-Anexo 05):
```json
{"type":"despesas_por_grupo","data":{"dados":[
  {"grupo":"Pessoal","empenhado":90000000.00,"liquidado":88500000.00,"pago":88000000.00},
  {"grupo":"Juros e Encargos","empenhado":5000000.00,"liquidado":4950000.00,"pago":4900000.00},
  {"grupo":"Investimentos","empenhado":12000000.00,"liquidado":9000000.00,"pago":7000000.00},
  {"grupo":"Corrente","empenhado":150000000.00,"liquidado":145000000.00,"pago":143000000.00}
]}}
```

**9. `transferencias`** — embrulhado em `dados` (SICONFI RREO-Anexo 07):
```json
{"type":"transferencias","data":{"dados":[
  {"orgao":"Transferências do FUNDEB","valor":28000000.00},
  {"orgao":"Transferências da União - SUS","valor":15000000.00},
  {"orgao":"Cota-parte do FPM","valor":22000000.00}
]}}
```

**10. `contratos`** — embrulhado em `dados`, array `pncp.Contrato` (PNCP `/contratacoes/publicacao`):
```json
{"type":"contratos","data":{"dados":[
  {
    "anoContrato":2023,
    "numeroControlePNCP":"1A2B3C-2024-0001/0001-001",
    "numeroContrato":"12/2023",
    "objetoContrato":"Aquisição de medicamentos",
    "valorGlobal":450000.00,
    "valorParcela":75000.00,
    "nomeRazaoSocialFornecedor":"Farma LTDA",
    "niFornecedor":"11222333000181",
    "modalidadeNome":"Pregão",
    "dataPublicacaoPncp":"2023-03-15",
    "dataVigenciaInicio":"2023-04-01",
    "dataVigenciaFim":"2024-03-31",
    "orgaoEntidade":{"razaoSocial":"Município de São Paulo","cnpj":"12345678000190"}
  }
]}}
```

> O struct `pncp.Contrato` em `contratos_types.go:35-80` tem ~40 campos com ponteiros (`*int`/`*string`) — na prática muitos vêm `null`. O UseCase retorna `resp.Data` direto (sem transformar para `types.ContratoPNCP`), então o JSON é o formato bruto do PNCP.

### Frames de encerramento

**11. `erro`** — condicional, só se SICONFI indisponível:
```json
{"type":"erro","data":{"erro":"API SICONFI (Tesouro Nacional) temporariamente indisponível. Os dados financeiros do município não puderam ser carregados. Tente novamente mais tarde."}}
```

**12. `concluido`** — sempre o último frame:
```json
{"type":"concluido","data":null}
```

### Exemplos de erro HTTP (antes do upgrade)

**400 — código IBGE não-numérico:**
```
GET /municipio/abc/detalhes/stream
→ 400 {"erro": "codigo IBGE invalido"}
```

**400 — código IBGE <= 0:**
```
GET /municipio/0/detalhes/stream
→ 400 {"erro": "codigo IBGE invalido"}
```

---

## Variáveis de ambiente

| Variável | Default | Descrição |
|----------|---------|-----------|
| `PNCP_BASE_URL` | `https://pncp.gov.br/pncp-consulta/v1` | Base URL da API de consulta do PNCP |
| `SICONFI_BASE_URL` | `https://apidatalake.tesouro.gov.br/ords/siconfi/tt` | Base URL da API de dados abertos do SICONFI |

Instanciados em `internal/app/app.go:178` (PNCP) e `app.go:184` (SICONFI).

---

## Notas técnicas

- **Concorrência**: 10 goroutines coordenadas por `sync.WaitGroup`, canal buffer 20. Sem semáforo adicional (10 fixas, sem worker pool).
- **Fallback de período/ano SICONFI** (para encontrar o relatório mais recente publicado):
  - **RGF** tenta `Q3` (3º quadrimestre) → `S2` (2º semestre) — ver `tentativasRGF` em `detalhes-municipio-usecase.go:113-122`.
  - **RREO** tenta bimestres `6` → `5`; se vazio e `ano > 2013`, decrementa `ano-1` e re-tenta `6` → `5`.
  - **DCA** tenta `ano` → `ano-1` → `ano-2` (loop em `buscarBalancoPatrimonial`).
- **SICONFI indisponível**: detectado quando a resposta HTTP não-200 contém `{"code":"AccountIsLocked"}` → o client retorna `ErrSICONFIIndisponivel` (`siconfi_client.go:17`). Qualquer goroutine que receber esse erro seta `atomic.Bool` `apiIndisponivel` no UseCase → handler emite frame `erro` ao final. **Não aborta as demais seções** — as que conseguem dados enviam normalmente.
- **Panic-safe**: cada goroutine tem `recover()` — um pânico (ex: nil pointer em item SICONFI malformado) é logado mas não derruba o handler. A seção afetada simplesmente não envia seu frame.
- **Sem retry/backoff** explícito: o "fallback" é tentar outro período/ano, não re-tentar a mesma chamada.
- **Sem timeout explícito** no handler: depende do `ctx` do request HTTP (Gin) e dos timeouts dos `http.Client` internos (SICONFI 30s, PNCP 300s). O loop respeita `ctx.Done()` (cliente desconecta → handler termina).
- **PNCP não pagina**: busca só página 1, tamanho 20. Municípios com > 20 contratações no ano verão só as 20 primeiras.
- **Parsing heurístico SICONFI**: cada método `buscar*` itera os itens retornados (cada item tem `Rotulo`, `Conta`, `Coluna`, `Valor`) e usa `strings.Contains` para casar rótulos/colunas específicos e acumular o `Valor`. Necessário porque o SICONFI retorna dados tabulares rotulados, não estruturados.
- **Sem Redis/PostgreSQL**: rota stateless, 100% via APIs externas (diferente de `/busca/contexto` que usa TSE+Redis, e `/estado/:uf/financeiro/stream` que usa Redis).

---

## Referências

- **Clients:** [docs/clientes/siconfi.md](./clientes/siconfi.md), [docs/clientes/pncp.md](./clientes/pncp.md)
- **Código-fonte:** `internal/esferas-brasileiras/municipal/`
  - `handler/detalhes-municipio-ws-handler.go` — endpoint Gin (valida IBGE, upgrade WS, 10 goroutines, loop de frames)
  - `usecase/detalhes-municipio-usecase.go` — orquestração (10 métodos públicos, fallback de período/ano, parsing heurístico SICONFI, flag `atomic.Bool` de indisponibilidade)
- **Tipos de resposta:** `internal/shared/types/municipio.go:63-139` (`DetalhesMunicipioResponse` e structs das 9 seções)
- **Helper WebSocket:** `internal/shared/websocket/websocket.go` (`Upgrader` com `CheckOrigin` permissivo, `WriteJSON`)
- **Clients Go:** `internal/shared/clients/{siconfi,pncp}/`
- **Wiring:** `internal/app/app.go:178,184` (clients), `app.go:265` (UseCase), `app.go:378` (handler)
