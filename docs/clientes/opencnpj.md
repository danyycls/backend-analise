# OpenCNPJ

**Nome cliente:** OpenCNPJ

**Descrição:** Cliente Go para consulta de dados cadastrais de pessoas jurídicas (CNPJ) via API pública OpenCNPJ, retornando razão social, nome fantasia, situação cadastral, capital social e quadro societário (QSA).

## Doc Client

**Documentação de integração client:** <https://opencnpj.org>

**Base URL:** `https://api.opencnpj.org/%s`

## APIs Integradas

### Consulta de CNPJ

| Método | URL | Input | Output | Descrição |
|--------|-----|-------|--------|-----------|
| `Buscar` | `GET /<cnpj>` | `ctx context.Context`, `cnpj string` | `*OpenCNPJResponse` | Consulta dados cadastrais de uma pessoa jurídica pelo CNPJ, retornando razão social, nome fantasia, situação cadastral, capital social e o quadro de sócios e administradores (QSA). |
