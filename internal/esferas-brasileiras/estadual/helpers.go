package estadual

import "strings"

func NormalizarTipoEleicao(nomeTipo string) string {
	upper := strings.ToUpper(nomeTipo)
	if strings.Contains(upper, "SUPLEMENTAR") {
		return "SUPLEMENTAR"
	}
	if strings.Contains(upper, "ORDINÁRIA") || strings.Contains(upper, "ORDINARIA") {
		return "ORDINÁRIA"
	}
	return nomeTipo
}
