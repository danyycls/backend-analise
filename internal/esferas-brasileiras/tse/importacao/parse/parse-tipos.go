// Pacote service contém funções auxiliares para a importação de dados
// eleitorais. Este arquivo implementa funções de convertsão e normalização
// de tipos de dados extraídos das colunas dos CSVs: texto, inteiro,
// decimal, data e documento (CPF/CNPJ).
package parse

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// -----------------------------------------------------------------------------
// Funções de tratamento de texto.
// -----------------------------------------------------------------------------

// textoComFallback retorna o valor textual se não estiver vazio, ou o
// fallback fornecido em caso contrário.
func textoComFallback(valor, fallback string) string {
	texto := textoOpcional(valor)
	if texto == "" {
		return fallback
	}
	return texto
}

// textoOpcional normaliza o valor textual: remove espaços nas bordas e
// retorna vazio para valores considerados nulos pelo TSE (#NULO, #NE,
// "NÃO DIVULGÁVEL", etc.).
func textoOpcional(valor string) string {
	texto := strings.TrimSpace(valor)
	switch strings.ToUpper(texto) {
	case "", "#NULO", "#NULO#", "#NE", "NÃO DIVULGÁVEL", "NAO DIVULGAVEL":
		return ""
	default:
		return texto
	}
}

// -----------------------------------------------------------------------------
// Funções de convertsão de inteiros.
// -----------------------------------------------------------------------------

// inteiroOpcional converte string para *int. Remove separadores de milhar
// (ponto) e decimais (vírgula) antes da convertsão. Retorna nil para valores
// inválidos ou vazios.
func inteiroOpcional(valor string) *int {
	texto := textoOpcional(valor)
	if texto == "" {
		return nil
	}
	// Remove separadores de milhar e decimais para normalizar.
	normalizado := strings.ReplaceAll(texto, ".", "")
	normalizado = strings.ReplaceAll(normalizado, ",", "")
	inteiro, err := strconv.Atoi(normalizado)
	if err != nil || inteiro < 0 {
		return nil
	}
	return &inteiro
}

// inteiro16Opcional converte string para *int16 (usado para números pequenos
// como ano, turno, número do partido).
func inteiro16Opcional(valor string) *int16 {
	inteiro := inteiroOpcional(valor)
	if inteiro == nil {
		return nil
	}
	valor16 := int16(*inteiro)
	return &valor16
}

// inteiro64Opcional converte string para *int64 (usado para SQs e códigos
// grandes do TSE).
func inteiro64Opcional(valor string) *int64 {
	texto := textoOpcional(valor)
	if texto == "" {
		return nil
	}
	inteiro, err := strconv.ParseInt(texto, 10, 64)
	if err != nil || inteiro < 0 {
		return nil
	}
	return &inteiro
}

// inteiro64OuZero converte string para int64, retornando 0 se inválido.
func inteiro64OuZero(valor string) int64 {
	inteiro := inteiro64Opcional(valor)
	if inteiro == nil {
		return 0
	}
	return *inteiro
}

// -----------------------------------------------------------------------------
// Funções de convertsão de decimais (valores monetários).
// -----------------------------------------------------------------------------

// decimalOpcional converte string para *float64. O formato brasileiro usa
// ponto como separador de milhar e vírgula como separador decimal.
// Exemplo: "1.234,56" → 1234.56
func decimalOpcional(valor string) *float64 {
	texto := textoOpcional(valor)
	if texto == "" {
		return nil
	}
	normalizado := strings.ReplaceAll(texto, ".", "")
	normalizado = strings.ReplaceAll(normalizado, ",", ".")
	decimal, err := strconv.ParseFloat(normalizado, 64)
	if err != nil {
		return nil
	}
	return &decimal
}

// decimalOuZero converte string para float64, retornando 0 se inválido.
func decimalOuZero(valor string) float64 {
	decimal := decimalOpcional(valor)
	if decimal == nil {
		return 0
	}
	return *decimal
}

// -----------------------------------------------------------------------------
// Funções de convertsão de data.
// -----------------------------------------------------------------------------

// dataOpcional converte string no formato "dd/mm/aaaa" para *time.Time.
// Retorna nil se a string estiver vazia ou for inválida.
func dataOpcional(valor string) *time.Time {
	texto := textoOpcional(valor)
	if texto == "" {
		return nil
	}
	data, err := time.Parse("02/01/2006", texto)
	if err != nil {
		return nil
	}
	return &data
}

// -----------------------------------------------------------------------------
// Funções de normalização de documentos (CPF/CNPJ).
// -----------------------------------------------------------------------------

var padroesPlaceholderDocumento = []string{
	"-4",
	"-1",
	"00000000000",
	"00000000000000",
}

func documentoOpcional(valor string) string {
	texto := textoOpcional(valor)
	if texto == "" {
		return ""
	}

	normalizado := strings.TrimSpace(texto)
	for _, p := range padroesPlaceholderDocumento {
		if normalizado == p {
			return ""
		}
	}

	var digitos strings.Builder
	for _, caractere := range texto {
		if caractere >= '0' && caractere <= '9' {
			digitos.WriteRune(caractere)
		}
	}
	documento := digitos.String()
	if documento == "" {
		return ""
	}
	if strings.Trim(documento, "0") == "" {
		return ""
	}
	return documento
}

// -----------------------------------------------------------------------------
// Função de validação de UUID opcional.
// -----------------------------------------------------------------------------

// uuidOpcional valida se um ponteiro de UUID não é nil nem Nil. Usado para
// verificar campos obrigatórios de chave estrangeira antes de montar a
// entidade.
func uuidOpcional(valor *uuid.UUID, campo string) *uuid.UUID {
	if valor == nil || *valor == uuid.Nil {
		fmt.Printf("uuuid %s obrigatorio", campo)
		return nil
	}

	return valor
}
