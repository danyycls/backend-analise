package utils

import (
	"strings"
)

func NormalizarCNPJ(c string) string {
	return strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, c)
}

func ApenasDigitos(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// NormalizarDocumento retorna apenas os digitos do documento e um booleano
// indicando se o documento original continha '*' (busca parcial).
func NormalizarDocumento(documento string) (doc string, parcial bool) {
	parcial = strings.Contains(documento, "*")
	return ApenasDigitos(documento), parcial
}
