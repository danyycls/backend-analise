package senado

import "encoding/json"

type ListaParlamentarEmExercicio struct {
	ListaParlamentarEmExercicio struct {
		Parlamentares struct {
			Parlamentar []ParlamentarResumo `json:"Parlamentar"`
		} `json:"Parlamentares"`
	} `json:"ListaParlamentarEmExercicio"`
}

type ParlamentarResumo struct {
	IdentificacaoParlamentar IdentificacaoParlamentar `json:"IdentificacaoParlamentar"`
	Mandato                  MandatoResumo            `json:"Mandato"`
}

type IdentificacaoParlamentar struct {
	CodigoParlamentar       string     `json:"CodigoParlamentar"`
	CodigoPublicoNaLegAtual string     `json:"CodigoPublicoNaLegAtual"`
	NomeParlamentar         string     `json:"NomeParlamentar"`
	NomeCompletoParlamentar string     `json:"NomeCompletoParlamentar"`
	SexoParlamentar         string     `json:"SexoParlamentar"`
	FormaTratamento         string     `json:"FormaTratamento"`
	UrlFotoParlamentar      string     `json:"UrlFotoParlamentar"`
	UrlPaginaParlamentar    string     `json:"UrlPaginaParlamentar"`
	EmailParlamentar        string     `json:"EmailParlamentar"`
	SiglaPartidoParlamentar string     `json:"SiglaPartidoParlamentar"`
	UfParlamentar           string     `json:"UfParlamentar"`
	MembroMesa              string     `json:"MembroMesa"`
	MembroLideranca         string     `json:"MembroLideranca"`
	Bloco                   *BlocoInfo `json:"Bloco"`
}

type BlocoInfo struct {
	CodigoBloco string `json:"CodigoBloco"`
	NomeBloco   string `json:"NomeBloco"`
	NomeApelido string `json:"NomeApelido"`
	DataCriacao string `json:"DataCriacao"`
}

type MandatoResumo struct {
	CodigoMandato                string      `json:"CodigoMandato"`
	UfParlamentar                string      `json:"UfParlamentar"`
	PrimeiraLegislaturaDoMandato Legislatura `json:"PrimeiraLegislaturaDoMandato"`
	SegundaLegislaturaDoMandato  Legislatura `json:"SegundaLegislaturaDoMandato"`
	DescricaoParticipacao        string      `json:"DescricaoParticipacao"`
}

type Legislatura struct {
	NumeroLegislatura string `json:"NumeroLegislatura"`
	DataInicio        string `json:"DataInicio"`
	DataFim           string `json:"DataFim"`
}

type DetalheParlamentar struct {
	DetalheParlamentar struct {
		Parlamentar ParlamentarDetalhe `json:"Parlamentar"`
	} `json:"DetalheParlamentar"`
}

type ParlamentarDetalhe struct {
	IdentificacaoParlamentar IdentificacaoParlamentar `json:"IdentificacaoParlamentar"`
	DadosBasicosParlamentar  DadosBasicosParlamentar  `json:"DadosBasicosParlamentar"`
	Telefones                *TelefonesWrapper        `json:"Telefones"`
	OutrasInformacoes        *OutrasInformacoes       `json:"OutrasInformacoes"`
}

type TelefonesWrapper struct {
	Telefone []Telefone `json:"Telefone"`
}

type Telefone struct {
	NumeroTelefone  string `json:"NumeroTelefone"`
	OrdemPublicacao string `json:"OrdemPublicacao"`
	IndicadorFax    string `json:"IndicadorFax"`
}

type OutrasInformacoes struct {
	Servico []Servico `json:"Servico"`
}

type Servico struct {
	NomeServico      string `json:"NomeServico"`
	DescricaoServico string `json:"DescricaoServico"`
	UrlServico       string `json:"UrlServico"`
}

type DadosBasicosParlamentar struct {
	DataNascimento      string `json:"DataNascimento"`
	Naturalidade        string `json:"Naturalidade"`
	UfNaturalidade      string `json:"UfNaturalidade"`
	EnderecoParlamentar string `json:"EnderecoParlamentar"`
}

type CargoParlamentar struct {
	CargoParlamentar struct {
		Parlamentar struct {
			Codigo string `json:"Codigo"`
			Nome   string `json:"Nome"`
			Cargos struct {
				Cargo []Cargo `json:"Cargo"`
			} `json:"Cargos"`
		} `json:"Parlamentar"`
	} `json:"CargoParlamentar"`
}

type Cargo struct {
	IdentificacaoComissao IdentificacaoComissao `json:"IdentificacaoComissao"`
	CodigoCargo           string                `json:"CodigoCargo"`
	DescricaoCargo        string                `json:"DescricaoCargo"`
	DataInicio            string                `json:"DataInicio"`
	DataFim               string                `json:"DataFim"`
}

type IdentificacaoComissao struct {
	CodigoComissao    string `json:"CodigoComissao"`
	SiglaComissao     string `json:"SiglaComissao"`
	NomeComissao      string `json:"NomeComissao"`
	SiglaCasaComissao string `json:"SiglaCasaComissao"`
}

type MembroComissaoParlamentar struct {
	MembroComissaoParlamentar struct {
		Parlamentar struct {
			Codigo          string `json:"Codigo"`
			Nome            string `json:"Nome"`
			MembroComissoes struct {
				Comissao []ComissaoMembro `json:"Comissao"`
			} `json:"MembroComissoes"`
		} `json:"Parlamentar"`
	} `json:"MembroComissaoParlamentar"`
}

type ComissaoMembro struct {
	IdentificacaoComissao IdentificacaoComissao `json:"IdentificacaoComissao"`
	DescricaoParticipacao string                `json:"DescricaoParticipacao"`
	DataInicio            string                `json:"DataInicio"`
	DataFim               string                `json:"DataFim"`
}

type MandatoParlamentar struct {
	MandatoParlamentar struct {
		Parlamentar struct {
			Codigo   string `json:"Codigo"`
			Nome     string `json:"Nome"`
			Mandatos struct {
				Mandato []MandatoDetalhe `json:"Mandato"`
			} `json:"Mandatos"`
		} `json:"Parlamentar"`
	} `json:"MandatoParlamentar"`
}

type MandatoDetalhe struct {
	CodigoMandato                string             `json:"CodigoMandato"`
	UfParlamentar                string             `json:"UfParlamentar"`
	PrimeiraLegislaturaDoMandato Legislatura        `json:"PrimeiraLegislaturaDoMandato"`
	SegundaLegislaturaDoMandato  Legislatura        `json:"SegundaLegislaturaDoMandato"`
	DescricaoParticipacao        string             `json:"DescricaoParticipacao"`
	Suplentes                    *SuplentesWrapper  `json:"Suplentes"`
	Exercicios                   *ExerciciosWrapper `json:"Exercicios"`
	Partidos                     *PartidosWrapper   `json:"Partidos"`
}

type SuplentesWrapper struct {
	Suplente []Suplente `json:"Suplente"`
}

type Suplente struct {
	DescricaoParticipacao string `json:"DescricaoParticipacao"`
	CodigoParlamentar     string `json:"CodigoParlamentar"`
	NomeParlamentar       string `json:"NomeParlamentar"`
}

type ExerciciosWrapper struct {
	Exercicio []Exercicio `json:"Exercicio"`
}

type Exercicio struct {
	CodigoExercicio string `json:"CodigoExercicio"`
	DataInicio      string `json:"DataInicio"`
}

type PartidosWrapper struct {
	Partido []PartidoMandato `json:"Partido"`
}

type PartidoMandato struct {
	CodigoPartido   string `json:"CodigoPartido"`
	Sigla           string `json:"Sigla"`
	Nome            string `json:"Nome"`
	DataFiliacao    string `json:"DataFiliacao"`
	DataDesfiliacao string `json:"DataDesfiliacao"`
}

type OrcamentoLista struct {
	ListaLoteEmendas struct {
		LotesEmendasOrcamento struct {
			LoteEmendasOrcamento []LoteEmendasOrcamento `json:"LoteEmendasOrcamento"`
		} `json:"LotesEmendasOrcamento"`
	} `json:"ListaLoteEmendas"`
}

type LoteEmendasOrcamento struct {
	NomeAutorOrcamento       string `json:"NomeAutorOrcamento"`
	IndicadorAtivo           string `json:"IndicadorAtivo"`
	EmailAutorOrcamento      string `json:"EmailAutorOrcamento"`
	CodigoAutorOrcamento     string `json:"CodigoAutorOrcamento"`
	DataOperacao             string `json:"DataOperacao"`
	QuantidadeEmendas        string `json:"QuantidadeEmendas"`
	AnoExecucao              string `json:"AnoExecucao"`
	NumeroMateria            string `json:"NumeroMateria"`
	AnoMateria               string `json:"AnoMateria"`
	SiglaTipoPlOrcamento     string `json:"SiglaTipoPlOrcamento"`
	DescricaoTipoPlOrcamento string `json:"DescricaoTipoPlOrcamento"`
}

type VotacaoItem struct {
	Ano                     int    `json:"ano"`
	CasaSessao              string `json:"casaSessao"`
	CodigoMateria           int    `json:"codigoMateria"`
	CodigoSessao            int    `json:"codigoSessao"`
	CodigoSessaoLegislativa int    `json:"codigoSessaoLegislativa"`
	CodigoSessaoVotacao     int    `json:"codigoSessaoVotacao"`
	CodigoVotacaoSve        int    `json:"codigoVotacaoSve"`
	DataApresentacao        string `json:"dataApresentacao"`
	DataSessao              string `json:"dataSessao"`
	DescricaoVotacao        string `json:"descricaoVotacao"`
	Ementa                  string `json:"ementa"`
	IdProcesso              int    `json:"idProcesso"`
	Identificacao           string `json:"identificacao"`
}

type ProcessoAssunto struct {
	Id                int    `json:"id"`
	AssuntoGeral      string `json:"assuntoGeral"`
	AssuntoEspecifico string `json:"assuntoEspecifico"`
	DataInicio        string `json:"dataInicio"`
	DataFim           string `json:"dataFim"`
}

type VotacaoComissaoWrapper struct {
	VotacoesComissao struct {
		Votacoes struct {
			Votacao []VotacaoComissao `json:"Votacao"`
		} `json:"Votacoes"`
	} `json:"VotacoesComissao"`
}

type VotacaoComissao struct {
	CodigoVotacao         string       `json:"CodigoVotacao"`
	SiglaCasaColegiado    string       `json:"SiglaCasaColegiado"`
	CodigoReuniao         string       `json:"CodigoReuniao"`
	DataHoraInicioReuniao string       `json:"DataHoraInicioReuniao"`
	CodigoColegiado       string       `json:"CodigoColegiado"`
	SiglaColegiado        string       `json:"SiglaColegiado"`
	NomeColegiado         string       `json:"NomeColegiado"`
	IdentificacaoMateria  string       `json:"IdentificacaoMateria"`
	DescricaoVotacao      string       `json:"DescricaoVotacao"`
	Votos                 VotosWrapper `json:"Votos"`
}

type VotosWrapper struct {
	Voto []Voto `json:"Voto"`
}

func (w *VotosWrapper) UnmarshalJSON(data []byte) error {
	var raw struct {
		Voto json.RawMessage `json:"Voto"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	// Voto pode ser array ou objeto único
	if len(raw.Voto) > 0 && raw.Voto[0] == '{' {
		var single Voto
		if err := json.Unmarshal(raw.Voto, &single); err != nil {
			return err
		}
		w.Voto = []Voto{single}
	} else {
		var slice []Voto
		if err := json.Unmarshal(raw.Voto, &slice); err != nil {
			return err
		}
		w.Voto = slice
	}
	return nil
}

type Voto struct {
	CodigoParlamentar       string `json:"CodigoParlamentar"`
	NomeParlamentar         string `json:"NomeParlamentar"`
	SiglaPartidoParlamentar string `json:"SiglaPartidoParlamentar"`
	QualidadeVoto           string `json:"QualidadeVoto"`
}

type MateriaTramitacao struct {
	Materials struct {
		Materia []MateriaItem `json:"Materia"`
	} `json:"Materials"`
}

type MateriaItem struct {
	Codigo            string `json:"Codigo"`
	Sigla             string `json:"Sigla"`
	Numero            string `json:"Numero"`
	Ano               string `json:"Ano"`
	Ementa            string `json:"Ementa"`
	Identificacao     string `json:"Identificacao"`
	DescricaoSituacao string `json:"DescricaoSituacao"`
}

type PlenarioAgendaDia struct {
	AgendaDia struct {
		Reunioes struct {
			Reuniao []Reuniao `json:"Reuniao"`
		} `json:"Reunioes"`
	} `json:"AgendaDia"`
}

type PlenarioAgendaMes struct {
	AgendaMes struct {
		Reunioes struct {
			Reuniao []Reuniao `json:"Reuniao"`
		} `json:"Reunioes"`
	} `json:"AgendaMes"`
}

type Reuniao struct {
	Codigo     string `json:"Codigo"`
	Descricao  string `json:"Descricao"`
	Situacao   string `json:"Situacao"`
	Data       string `json:"Data"`
	HoraInicio string `json:"HoraInicio"`
	SiglaCasa  string `json:"SiglaCasa"`
}

type PlenarioEncontro struct {
	Encontro struct {
		Codigo    string `json:"Codigo"`
		Descricao string `json:"Descricao"`
		Situacao  string `json:"Situacao"`
	} `json:"Encontro"`
}

type ProcessoItem struct {
	Id               int    `json:"id"`
	Identificacao    string `json:"identificacao"`
	Ementa           string `json:"ementa"`
	DataApresentacao string `json:"dataApresentacao"`
	Situacao         string `json:"situacao"`
}

type ProcessoEmenda struct {
	Id                       int    `json:"id"`
	Identificacao            string `json:"identificacao"`
	DataApresentacao         string `json:"dataApresentacao"`
	Autoria                  string `json:"autoria"`
	DescricaoDocumentoEmenda string `json:"descricaoDocumentoEmenda"`
	Numero                   string `json:"numero"`
	Tipo                     string `json:"tipo"`
	TurnoApresentacao        string `json:"turnoApresentacao"`
	Casa                     string `json:"casa"`
	CodigoColegiado          int    `json:"codigoColegiado"`
	SiglaColegiado           string `json:"siglaColegiado"`
	NomeColegiado            string `json:"nomeColegiado"`
	IdCiEmenda               int    `json:"idCiEmenda"`
	IdCiEmendado             int    `json:"idCiEmendado"`
	IdDocumentoEmenda        int    `json:"idDocumentoEmenda"`
	IdProcesso               int    `json:"idProcesso"`
	UrlDocumentoEmenda       string `json:"urlDocumentoEmenda"`
}

type VotacaoComissaoParlamentar struct {
	VotacoesParlamentar struct {
		Votacoes struct {
			Votacao []VotacaoComissao `json:"Votacao"`
		} `json:"Votacoes"`
	} `json:"VotacoesParlamentar"`
}

type SenadorDetalhado struct {
	Senador   *ParlamentarDetalhe `json:"senador"`
	Cargos    []Cargo             `json:"cargos"`
	Comissoes []ComissaoMembro    `json:"comissoes"`
	Mandatos  []MandatoDetalhe    `json:"mandatos"`
}
