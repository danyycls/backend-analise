package usecase

import (
	"github.com/danyele/podp/internal/shared/domain"
)

type AnalisarLigacaoPoliticaRequest struct {
	NumeroControlePncp string                                `json:"numero_controle_pncp"`
	CpfCnpj            string                                `json:"cpf_cnpj"`
	Socios             []AnalisarLigacaoPoliticaSocioRequest `json:"socios"`
}

type AnalisarLigacaoPoliticaSocioRequest struct {
	Nome      string `json:"nome"`
	Documento string `json:"documento"`
}

type AnalisarLigacaoPoliticaResponse struct {
	DocumentosProcessados int                       `json:"documentos_processados"`
	Resultados            []domain.VinculoLicitacao `json:"resultados"`
}
