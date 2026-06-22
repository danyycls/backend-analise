# TCU

**Nome cliente:** TCU — Tribunal de Contas da União

**Descrição:** Cliente para consulta ao cadastro de sanções e condenações do Tribunal de Contas da União. Permite pesquisar responsáveis com contas julgadas irregulares (CADIRREG), contas com possível implicação eleitoral, responsáveis inabilitados para cargo em comissão ou função de confiança, e licitantes inidôneos.

## Doc Client

**Documentação de integração client:** <https://sites.tcu.gov.br/dados-abertos/webservices-tcu/#sancoes-e-condenacoes/>

**Base URL:** `https://certidoes.apps.tcu.gov.br/api/publico`

## APIs Integradas

### Contas Irregulares

| Método | URL | Input | Output | Descrição |
|--------|-----|-------|--------|-----------|
| `BuscarContasIrregulares` | `POST /responsaveis-contas-irregulares` | `TCUQueryParams` | `[]ContasIrregulares` | Webservice REST que retorna as pessoas que, por decisão do TCU, tiveram suas contas julgadas irregulares com decisão transitada em julgado. |

### Fins Eleitorais

| Método | URL | Input | Output | Descrição |
|--------|-----|-------|--------|-----------|
| `BuscarFinsEleitorais` | `POST /responsaveis-fins-eleitorais` | `TCUQueryParams` | `[]FinsEleitorais` | Webservice REST que retorna as pessoas que, por decisão do TCU, tiveram suas contas julgadas irregulares com imputação de débito e com decisão transitada em julgado nos últimos 8 anos. |

### Inabilitados

| Método | URL | Input | Output | Descrição |
|--------|-----|-------|--------|-----------|
| `BuscarInabilitados` | `POST /responsaveis-inabilitados` | `TCUQueryParams` | `[]Sancoes` | Webservice REST que consolida informações de pessoas que, por decisão do TCU, são consideradas inabilitadas e não podem exercer cargo em comissão ou função de confiança no âmbito da Administração Pública por um período de 5 a 8 anos. |

### Inidôneos

| Método | URL | Input | Output | Descrição |
|--------|-----|-------|--------|-----------|
| `BuscarInidoneos` | `POST /responsaveis-inidoneos` | `TCUQueryParams` | `[]Sancoes` | Webservice REST que consolida informações de pessoas físicas ou jurídicas que, por decisão do TCU, são consideradas inidôneas e estão impedidas de participar de licitações no âmbito da Administração Pública. |
