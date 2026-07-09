package pncp

type EventoAnalise struct {
	Type                string              `json:"type"`
	CNPJ                string              `json:"cnpj,omitempty"`
	Orgao               string              `json:"orgao,omitempty"`
	TotalContratos      int                 `json:"totalContratos,omitempty"`
	ValorTotalContratos float64             `json:"valorTotalContratos,omitempty"`
	Message             string              `json:"message,omitempty"`
	Processed           int                 `json:"processed"`
	Total               int                 `json:"total"`
	Success             int                 `json:"success"`
	Errors              int                 `json:"errors"`
	Results             []*AnaliseResultado `json:"results,omitempty"`
}

type AnaliseContratoOrgaoRequest struct {
	CNPJs       []string `json:"cnpjs" binding:"required"`
	DataInicial string   `json:"dataInicial" binding:"required"`
	DataFinal   string   `json:"dataFinal" binding:"required"`
}

func StrPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
