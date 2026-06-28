package domain

import (
	tsetypes "github.com/danyele/podp/internal/esferas-brasileiras/tse/types"
	"github.com/danyele/podp/internal/shared/clients/portaltransparencia"
	"github.com/danyele/podp/internal/shared/clients/tcu"
	"github.com/danyele/podp/internal/shared/types"
)

type VinculoLicitacao struct {
	NumeroControlePncp string             `json:"numero_controle_pncp"`
	CpfCnpj            string             `json:"cpf_cnpj"`
	Socios             []SocioOutput      `json:"socios,omitempty"`
	Documentos         []DocumentoVinculo `json:"documentos,omitempty"`
}

type SocioOutput struct {
	Nome      string `json:"nome"`
	Documento string `json:"documento"`
}

type DocumentoVinculo struct {
	DocumentoInput       string    `json:"documento_input"`
	DocumentoNormalizado string    `json:"documento_normalizado"`
	Nome                 string    `json:"nome"`
	Parcial              bool      `json:"parcial"`
	Origem               string    `json:"origem"`
	Vinculos             []Vinculo `json:"vinculos,omitempty"`
}

type Vinculo struct {
	Tipo      string           `json:"tipo"`
	Descricao string           `json:"descricao"`
	Detalhes  *VinculoDetalhes `json:"detalhes,omitempty"`
}

type VinculoDetalhes struct {
	Fornecedor              *tsetypes.FornecedorDetalhado          `json:"fornecedor,omitempty"`
	Doador                  *types.Doador                          `json:"doador,omitempty"`
	ReceitasCandidato       []*types.ReceitaCandidato              `json:"receitas_candidato,omitempty"`
	ReceitasOrgaoPartidario []*types.ReceitaOrgaoPartidario        `json:"receitas_orgao_partidario,omitempty"`
	ContasIrregulares       []tcu.ContasIrregulares                `json:"contas_irregulares,omitempty"`
	Inabilitados            []tcu.Sancoes                          `json:"inabilitados,omitempty"`
	Inidoneos               []tcu.Sancoes                          `json:"inidoneos,omitempty"`
	ServidoresPublicos      []portaltransparencia.CadastroServidor `json:"servidores_publicos,omitempty"`
}
