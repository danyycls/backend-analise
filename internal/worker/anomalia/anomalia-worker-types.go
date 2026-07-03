package anomalia

type IniciarWorkerRequest struct {
	Licitacoes []LicitacaoInput `json:"licitacoes" binding:"required"`
}

type LicitacaoInput struct {
	NumeroControlePncp string              `json:"numero_controle_pncp"`
	CpfCnpj            string              `json:"cpf_cnpj"`
	Socios             []SocioInput        `json:"socios"`
	OrgaoCnpj          string              `json:"orgao_cnpj"`
	OrgaoNome          string              `json:"orgao_nome"`
	Uf                 string              `json:"uf"`
	Municipio          string              `json:"municipio"`
	ValorGlobal        float64             `json:"valor_global"`
	NomeEmpresa        string              `json:"nome_empresa"`
	Anormalidades      []AnormalidadeInput `json:"anormalidades,omitempty"`
}

type SocioInput struct {
	Nome      string `json:"nome"`
	Documento string `json:"documento"`
}

type AnormalidadeInput struct {
	Tipo      string                    `json:"tipo"`
	Descricao string                    `json:"descricao"`
	Detalhes  AnormalidadeDetalhesInput `json:"detalhes"`
}

type AnormalidadeDetalhesInput struct {
	DispensaValorLimite *DispensaValorLimiteInput `json:"dispensa_valor_limite,omitempty"`
}

type DispensaValorLimiteInput struct {
	Modalidade  string  `json:"modalidade"`
	Categoria   string  `json:"categoria"`
	ValorGlobal float64 `json:"valor_global"`
	Limite      float64 `json:"limite"`
	Excedente   float64 `json:"excedente"`
	Objeto      string  `json:"objeto"`
	Regra       string  `json:"regra"`
}

type WorkerProgressoResponse struct {
	JobID                string `json:"job_id"`
	Status               string `json:"status"`
	Total                int    `json:"total"`
	Processed            int    `json:"processed"`
	Success              int    `json:"success"`
	Errors               int    `json:"errors"`
	AnomaliasEncontradas int    `json:"anomalias_encontradas"`
	Message              string `json:"message,omitempty"`
	EtapaAtual           string `json:"etapa_atual"`
}

type WorkerEvento struct {
	Type                 string `json:"type"`
	Message              string `json:"message,omitempty"`
	Processed            int    `json:"processed"`
	Total                int    `json:"total"`
	Success              int    `json:"success"`
	Errors               int    `json:"errors"`
	AnomaliasEncontradas int    `json:"anomalias_encontradas"`
	EtapaAtual           string `json:"etapa_atual"`
	Documento            string `json:"documento,omitempty"`
}

type ListarAnomaliasResponse struct {
	Total     int                 `json:"total"`
	Pagina    int                 `json:"pagina"`
	PorPagina int                 `json:"por_pagina"`
	Anomalias []AnomaliaDocumento `json:"anomalias"`
}
