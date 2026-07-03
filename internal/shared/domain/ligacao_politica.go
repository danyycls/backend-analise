package domain

import (
	tsetypes "github.com/danyele/podp/internal/esferas-brasileiras/tse/types"
	"github.com/danyele/podp/internal/shared/clients/portaltransparencia"
	"github.com/danyele/podp/internal/shared/clients/tcu"
	"github.com/danyele/podp/internal/shared/types"
)

type VinculoLicitacao struct {
	NumeroControlePncp string             `bson:"numero_controle_pncp" json:"numero_controle_pncp"`
	CpfCnpj            string             `bson:"cpf_cnpj" json:"cpf_cnpj"`
	Socios             []SocioOutput      `bson:"socios,omitempty" json:"socios,omitempty"`
	DocumentosVinculos []DocumentoVinculo `bson:"documentos_vinculos,omitempty" json:"documentos_vinculos,omitempty"`
	ValorGlobal        float64            `bson:"valor_global" json:"valor_global"`
	NomeEmpresa        string             `bson:"nome_empresa" json:"nome_empresa"`
}

type SocioOutput struct {
	Nome      string `bson:"nome" json:"nome"`
	Documento string `bson:"documento" json:"documento"`
}

type DocumentoVinculo struct {
	DocumentoInput       string    `bson:"documento_input" json:"documento_input"`
	DocumentoNormalizado string    `bson:"documento_normalizado" json:"documento_normalizado"`
	Nome                 string    `bson:"nome" json:"nome"`
	Parcial              bool      `bson:"parcial" json:"parcial"`
	Origem               string    `bson:"origem" json:"origem"`
	Vinculos             []Vinculo `bson:"vinculos,omitempty" json:"vinculos,omitempty"`
}

type Vinculo struct {
	Tipo      string           `bson:"tipo" json:"tipo"`
	Descricao string           `bson:"descricao" json:"descricao"`
	Detalhes  *VinculoDetalhes `bson:"detalhes,omitempty" json:"detalhes,omitempty"`
}

type VinculoDetalhes struct {
	Fornecedor              *tsetypes.FornecedorDetalhado               `bson:"fornecedor,omitempty" json:"fornecedor,omitempty"`
	Doador                  *types.Doador                               `bson:"doador,omitempty" json:"doador,omitempty"`
	ReceitasCandidato       []*tsetypes.ReceitaCandidatoDetalhada       `bson:"receitas_candidato,omitempty" json:"receitas_candidato,omitempty"`
	ReceitasOrgaoPartidario []*tsetypes.ReceitaOrgaoPartidarioDetalhada `bson:"receitas_orgao_partidario,omitempty" json:"receitas_orgao_partidario,omitempty"`
	ContasIrregulares       []tcu.ContasIrregulares                     `bson:"contas_irregulares,omitempty" json:"contas_irregulares,omitempty"`
	Inabilitados            []tcu.Sancoes                               `bson:"inabilitados,omitempty" json:"inabilitados,omitempty"`
	Inidoneos               []tcu.Sancoes                               `bson:"inidoneos,omitempty" json:"inidoneos,omitempty"`
	ServidoresPublicos      []portaltransparencia.CadastroServidor      `bson:"servidores_publicos,omitempty" json:"servidores_publicos,omitempty"`
	PessoasPublicas         []portaltransparencia.PEP                   `bson:"pessoas_publicas,omitempty" json:"pessoas_publicas,omitempty"`
	DispensaValorLimite     *DispensaValorLimiteDetalhes                `bson:"dispensa_valor_limite,omitempty" json:"dispensa_valor_limite,omitempty"`
}

type DispensaValorLimiteDetalhes struct {
	Modalidade  string  `bson:"modalidade" json:"modalidade"`
	Categoria   string  `bson:"categoria" json:"categoria"`
	ValorGlobal float64 `bson:"valor_global" json:"valor_global"`
	Limite      float64 `bson:"limite" json:"limite"`
	Excedente   float64 `bson:"excedente" json:"excedente"`
	Objeto      string  `bson:"objeto" json:"objeto"`
	Regra       string  `bson:"regra" json:"regra"`
}
