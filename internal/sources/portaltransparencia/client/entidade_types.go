package portaltransparencia

type Orgao struct {
	Nome           string      `json:"descricao"`
	CodigoSIAFI    string      `json:"codigo"`
	CNPJ           string      `json:"cnpj"`
	Sigla          string      `json:"sigla"`
	DescricaoPoder string      `json:"descricaoPoder"`
	OrgaoMaximo    OrgaoMaximo `json:"orgaoMaximo"`
}

type OrgaoMaximo struct {
	Codigo string `json:"codigo"`
	Sigla  string `json:"sigla"`
	Nome   string `json:"nome"`
}

type OrgaoVinculado struct {
	CodigoSIAFI string `json:"codigoSIAFI"`
	CNPJ        string `json:"cnpj"`
	Sigla       string `json:"sigla"`
	Nome        string `json:"nome"`
}

type OrgaoServidor struct {
	Codigo               string `json:"codigo"`
	Nome                 string `json:"nome"`
	Sigla                string `json:"sigla"`
	CodigoOrgaoVinculado string `json:"codigoOrgaoVinculado"`
	NomeOrgaoVinculado   string `json:"nomeOrgaoVinculado"`
}

type UnidadeGestora struct {
	Codigo         string         `json:"codigo"`
	Nome           string         `json:"nome"`
	DescricaoPoder string         `json:"descricaoPoder"`
	OrgaoVinculado OrgaoVinculado `json:"orgaoVinculado"`
	OrgaoMaximo    OrgaoMaximo    `json:"orgaoMaximo"`
}

type PessoaFisica struct {
	CPF                             string `json:"cpf"`
	Nome                            string `json:"nome"`
	NIS                             string `json:"nis"`
	FavorecidoDespesas              bool   `json:"favorecidoDespesas"`
	Servidor                        bool   `json:"servidor"`
	BeneficiarioDiarias             bool   `json:"beneficiarioDiarias"`
	Permissionario                  bool   `json:"permissionario"`
	Contratado                      bool   `json:"contratado"`
	SancionadoCEIS                  bool   `json:"sancionadoCEIS"`
	SancionadoCNEP                  bool   `json:"sancionadoCNEP"`
	SancionadoCEAF                  bool   `json:"sancionadoCEAF"`
	PortadorCPDC                    bool   `json:"portadorCPDC"`
	PortadorCPGF                    bool   `json:"portadorCPGF"`
	FavorecidoBolsaFamilia          bool   `json:"favorecidoBolsaFamilia"`
	FavorecidoPeti                  bool   `json:"favorecidoPeti"`
	FavorecidoSafra                 bool   `json:"favorecidoSafra"`
	FavorecidoSeguroDefeso          bool   `json:"favorecidoSeguroDefeso"`
	FavorecidoBpc                   bool   `json:"favorecidoBpc"`
	FavorecidoTransferencias        bool   `json:"favorecidoTransferencias"`
	FavorecidoCPCC                  bool   `json:"favorecidoCPCC"`
	FavorecidoCPDC                  bool   `json:"favorecidoCPDC"`
	FavorecidoCPGF                  bool   `json:"favorecidoCPGF"`
	ParticipanteLicitacao           bool   `json:"participanteLicitacao"`
	ServidorInativo                 bool   `json:"servidorInativo"`
	PensionistaOuRepresentanteLegal bool   `json:"pensionistaOuRepresentanteLegal"`
	InstituidorPensao               bool   `json:"instituidorPensao"`
	AuxilioEmergencial              bool   `json:"auxilioEmergencial"`
	FavorecidoAuxilioBrasil         bool   `json:"favorecidoAuxilioBrasil"`
	FavorecidoNovoBolsaFamilia      bool   `json:"favorecidoNovoBolsaFamilia"`
	FavorecidoAuxilioReconstrucao   bool   `json:"favorecidoAuxilioReconstrucao"`
}

type PessoaJuridica struct {
	CNPJ                      string `json:"cnpj"`
	RazaoSocial               string `json:"razaoSocial"`
	NomeFantasia              string `json:"nomeFantasia"`
	FavorecidoDespesas        bool   `json:"favorecidoDespesas"`
	PossuiContratacao         bool   `json:"possuiContratacao"`
	Convenios                 bool   `json:"convenios"`
	FavorecidoTransferencias  bool   `json:"favorecidoTransferencias"`
	SancionadoCEPIM           bool   `json:"sancionadoCEPIM"`
	SancionadoCEIS            bool   `json:"sancionadoCEIS"`
	SancionadoCNEP            bool   `json:"sancionadoCNEP"`
	SancionadoCEAF            bool   `json:"sancionadoCEAF"`
	ParticipanteLicitacao     bool   `json:"participanteLicitacao"`
	EmitiuNFe                 bool   `json:"emitiuNFe"`
	BeneficiadoRenunciaFiscal bool   `json:"beneficiadoRenunciaFiscal"`
	IsentoImuneRenunciaFiscal bool   `json:"isentoImuneRenunciaFiscal"`
	HabilitadoRenunciaFiscal  bool   `json:"habilitadoRenunciaFiscal"`
}

type Pessoa struct {
	ID                    int    `json:"id"`
	CPFFormatado          string `json:"cpfFormatado"`
	CNPJFormatado         string `json:"cnpjFormatado"`
	NumeroInscricaoSocial string `json:"numeroInscricaoSocial"`
	Nome                  string `json:"nome"`
	RazaoSocialReceita    string `json:"razaoSocialReceita"`
	NomeFantasiaReceita   string `json:"nomeFantasiaReceita"`
	Tipo                  string `json:"tipo"`
}

type Cartao struct {
	ID              int               `json:"id"`
	MesExtrato      string            `json:"mesExtrato"`
	DataTransacao   string            `json:"dataTransacao"`
	ValorTransacao  string            `json:"valorTransacao"`
	TipoCartao      IdCodigoDescricao `json:"tipoCartao"`
	Estabelecimento Pessoa            `json:"estabelecimento"`
	UnidadeGestora  UnidadeGestora    `json:"unidadeGestora"`
	Portador        Beneficiario      `json:"portador"`
}

type IdCodigoDescricao struct {
	ID        int    `json:"id"`
	Codigo    string `json:"codigo"`
	Descricao string `json:"descricao"`
}

type Beneficiario struct {
	CPFFormatado string `json:"cpfFormatado"`
	NIS          string `json:"nis"`
	Nome         string `json:"nome"`
}

type Servidor struct {
	ID                              int            `json:"id"`
	IDservidorAposentadoPensionista int            `json:"idServidorAposentadoPensionista"`
	Pessoa                          Pessoa         `json:"pessoa"`
	Situacao                        string         `json:"situacao"`
	OrgaoServidorLotacao            OrgaoServidor  `json:"orgaoServidorLotacao"`
	OrgaoServidorExercicio          OrgaoServidor  `json:"orgaoServidorExercicio"`
	EstadoExercicio                 UF             `json:"estadoExercicio"`
	TipoServidor                    string         `json:"tipoServidor"`
	Funcao                          FuncaoServidor `json:"funcao"`
	CodigoMatriculaFormatado        string         `json:"codigoMatriculaFormatado"`
	FlagAfastado                    int            `json:"flagAfastado"`
}

type CadastroServidor struct {
	Servidor              ServidorAposentadoPensionista `json:"servidor"`
	FichasCargoEfetivo    []FichaCargoEfetivo           `json:"fichasCargoEfetivo"`
	FichasFuncao          []FichaFuncao                 `json:"fichasFuncao"`
	FichasMilitar         []FichaMilitar                `json:"fichasMilitar"`
	FichasDemaisSituacoes []FichaServidorCivil          `json:"fichasDemaisSituacoes"`
	FichasAposentadoria   []FichaAposentadoria          `json:"fichasAposentadoria"`
	FichasReformado       []FichaReformado              `json:"fichasReformado"`
	FichasPensaoCivil     []FichaPensaoCivil            `json:"fichasPensaoCivil"`
	FichasPensaoMilitar   []FichaPensaoMilitar          `json:"fichasPensaoMilitar"`
}

type ServidorAposentadoPensionista struct {
	ID                              int                      `json:"id"`
	IDservidorAposentadoPensionista int                      `json:"idServidorAposentadoPensionista"`
	Pessoa                          Pessoa                   `json:"pessoa"`
	Situacao                        string                   `json:"situacao"`
	OrgaoServidorLotacao            OrgaoServidor            `json:"orgaoServidorLotacao"`
	OrgaoServidorExercicio          OrgaoServidor            `json:"orgaoServidorExercicio"`
	EstadoExercicio                 UF                       `json:"estadoExercicio"`
	TipoServidor                    string                   `json:"tipoServidor"`
	Funcao                          FuncaoServidor           `json:"funcao"`
	ServidorInativoInstuidorPensao  ServidorInativo          `json:"servidorInativoInstuidorPensao"`
	PensionistaRepresentante        PensionistaRepresentante `json:"pensionistaRepresentante"`
	CodigoMatriculaFormatado        string                   `json:"codigoMatriculaFormatado"`
	FlagAfastado                    int                      `json:"flagAfastado"`
}

type UF struct {
	Sigla string `json:"sigla"`
	Nome  string `json:"nome"`
}

type FuncaoServidor struct {
	CodigoFuncaoCargo    string `json:"codigoFuncaoCargo"`
	DescricaoFuncaoCargo string `json:"descricaoFuncaoCargo"`
}

type ServidorInativo struct {
	ID           int    `json:"id"`
	CPFFormatado string `json:"cpfFormatado"`
	Nome         string `json:"nome"`
}

type PensionistaRepresentante struct {
	ID           int    `json:"id"`
	CPFFormatado string `json:"cpfFormatado"`
	Nome         string `json:"nome"`
}
