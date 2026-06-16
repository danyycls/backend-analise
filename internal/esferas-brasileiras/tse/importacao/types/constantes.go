// Pacote service contém funções auxiliares para a importação de dados
// eleitorais. Este arquivo define constantes compartilhadas usadas durante
// o parse e processamento dos arquivos CSV.
package types

// -----------------------------------------------------------------------------
// Constantes de tipo de prestador, tipo de registro e tamanho de lote.
// -----------------------------------------------------------------------------
const (
	// TipoPrestadorCandidato identifica prestações de contas de candidatos.
	TipoPrestadorCandidato = "CANDIDATO"

	// TipoPrestadorOrgaoPartidario identifica prestações de contas de
	// órgãos partidários.
	TipoPrestadorOrgaoPartidario = "ORGAO_PARTIDARIO"

	// TipoRegistroDespesaContratada classifica despesas do tipo contratada
	// (valor estimado/futuro).
	TipoRegistroDespesaContratada = "CONTRATADA"

	// TipoRegistroDespesaPaga classifica despesas do tipo paga
	// (valor efetivamente pago).
	TipoRegistroDespesaPaga = "PAGA"

	// TamanhoLotePadraoImportacao é o número de registros por lote usado
	// quando nenhum tamanho é explicitamente informado.
	TamanhoLotePadraoImportacao = 2000
)
