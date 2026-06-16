package parse

import "strings"

const (
	tipoArquivoConsultaCandidato             = "consulta_candidato"
	tipoArquivoBemCandidato                  = "bem_candidato"
	tipoArquivoDespesaContratadaCandidato    = "despesa_contratada_candidato"
	tipoArquivoDespesaPagaCandidato          = "despesa_paga_candidato"
	tipoArquivoReceitaCandidato              = "receita_candidato"
	tipoArquivoReceitaCandidatoDoadorOrigem  = "receita_candidato_doador_originario"
	tipoArquivoDespesaContratadaOrgaoPartido = "despesa_contratada_orgao_partidario"
	tipoArquivoDespesaPagaOrgaoPartido       = "despesa_paga_orgao_partidario"
	tipoArquivoReceitaOrgaoPartido           = "receita_orgao_partidario"
	tipoArquivoReceitaOrgaoPartidoDoadorOrig = "receita_orgao_partidario_doador_originario"
)

func IdentificarTipoArquivo(nomeArquivo string) (string, bool) {
	nome := strings.ToLower(nomeArquivo)

	switch {
	case strings.HasPrefix(nome, "consulta_cand_"):
		return tipoArquivoConsultaCandidato, true
	case strings.HasPrefix(nome, "bem_candidato_"):
		return tipoArquivoBemCandidato, true
	case strings.HasPrefix(nome, "despesas_contratadas_candidatos_"):
		return tipoArquivoDespesaContratadaCandidato, true
	case strings.HasPrefix(nome, "despesas_pagas_candidatos_"):
		return tipoArquivoDespesaPagaCandidato, true
	case strings.HasPrefix(nome, "receitas_candidatos_doador_originario_"):
		return tipoArquivoReceitaCandidatoDoadorOrigem, true
	case strings.HasPrefix(nome, "receitas_candidatos_"):
		return tipoArquivoReceitaCandidato, true
	case strings.HasPrefix(nome, "despesas_contratadas_orgaos_partidarios_"):
		return tipoArquivoDespesaContratadaOrgaoPartido, true
	case strings.HasPrefix(nome, "despesas_pagas_orgaos_partidarios_"):
		return tipoArquivoDespesaPagaOrgaoPartido, true
	case strings.HasPrefix(nome, "receitas_orgaos_partidarios_doador_originario_"):
		return tipoArquivoReceitaOrgaoPartidoDoadorOrig, true
	case strings.HasPrefix(nome, "receitas_orgaos_partidarios_"):
		return tipoArquivoReceitaOrgaoPartido, true
	default:
		return "", false
	}
}

func PrioridadeTipoArquivo(tipo string) int {
	switch tipo {
	case tipoArquivoConsultaCandidato:
		return 1
	case tipoArquivoBemCandidato:
		return 2
	case tipoArquivoDespesaContratadaCandidato:
		return 3
	case tipoArquivoReceitaCandidato:
		return 4
	case tipoArquivoReceitaCandidatoDoadorOrigem:
		return 5
	case tipoArquivoDespesaContratadaOrgaoPartido:
		return 6
	case tipoArquivoReceitaOrgaoPartido:
		return 7
	case tipoArquivoReceitaOrgaoPartidoDoadorOrig:
		return 8
	case tipoArquivoDespesaPagaCandidato:
		return 9
	case tipoArquivoDespesaPagaOrgaoPartido:
		return 10
	default:
		return 99
	}
}
