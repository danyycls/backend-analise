package utils

import (
	"fmt"
	"time"
)

type AnoMes struct {
	Ano int
	Mes int
}

func ExtrairMeses(dataInicial, dataFinal string) []AnoMes {
	anoInicio, mesInicio := parseAnoMes(dataInicial)
	anoFim, mesFim := parseAnoMes(dataFinal)

	if anoInicio == 0 || mesInicio == 0 || anoFim == 0 || mesFim == 0 {
		return nil
	}

	if anoFim < anoInicio || (anoFim == anoInicio && mesFim < mesInicio) {
		return nil
	}

	var meses []AnoMes
	ano, mes := anoInicio, mesInicio
	for {
		meses = append(meses, AnoMes{Ano: ano, Mes: mes})
		if ano == anoFim && mes == mesFim {
			break
		}
		mes++
		if mes > 12 {
			mes = 1
			ano++
		}
	}
	return meses
}

func ExtrairMesesDoAno(ano int) []AnoMes {
	meses := make([]AnoMes, 12)
	for i := 1; i <= 12; i++ {
		meses[i-1] = AnoMes{Ano: ano, Mes: i}
	}
	return meses
}

func parseAnoMes(data string) (int, int) {
	clean := data
	if len(clean) >= 10 && clean[4] == '-' && clean[7] == '-' {
		ano := 0
		mes := 0
		for i := 0; i < 4; i++ {
			if clean[i] < '0' || clean[i] > '9' {
				return 0, 0
			}
			ano = ano*10 + int(clean[i]-'0')
		}
		for i := 5; i < 7; i++ {
			if clean[i] < '0' || clean[i] > '9' {
				return 0, 0
			}
			mes = mes*10 + int(clean[i]-'0')
		}
		return ano, mes
	}

	if len(clean) >= 8 {
		ano := 0
		mes := 0
		for i := 0; i < 4; i++ {
			if clean[i] < '0' || clean[i] > '9' {
				return 0, 0
			}
			ano = ano*10 + int(clean[i]-'0')
		}
		for i := 4; i < 6; i++ {
			if clean[i] < '0' || clean[i] > '9' {
				return 0, 0
			}
			mes = mes*10 + int(clean[i]-'0')
		}
		return ano, mes
	}

	return 0, 0
}

func FormatarMes(ano, mes int) string {
	return fmt.Sprintf("%04d%02d", ano, mes)
}

func FormatarPeriodoMes(ano, mes int) (dataInicial, dataFinal string) {
	dataInicial = fmt.Sprintf("%04d%02d01", ano, mes)
	lastDay := time.Date(ano, time.Month(mes), 1, 0, 0, 0, 0, time.UTC).AddDate(0, 1, -1)
	dataFinal = lastDay.Format("20060102")
	return
}
