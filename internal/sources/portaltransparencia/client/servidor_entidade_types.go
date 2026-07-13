package portaltransparencia

type FichaCargoEfetivo struct {
	Nome                                          string   `json:"nome"`
	CPFDescaracterizado                           string   `json:"cpfDescaracterizado"`
	MatriculaDescaracterizada                     string   `json:"matriculaDescaracterizada"`
	DataPublicacaoDocumentoIngressoServicoPublico string   `json:"dataPublicacaoDocumentoIngressoServicoPublico"`
	DiplomaLegal                                  string   `json:"diplomaLegal"`
	JornadaTrabalho                               string   `json:"jornadaTrabalho"`
	RegimeJuridico                                string   `json:"regimeJuridico"`
	SituacaoServidor                              string   `json:"situacaoServidor"`
	Afastamentos                                  []string `json:"afastamentos"`
	OrgaoSuperiorLotacao                          string   `json:"orgaoSuperiorLotacao"`
	OrgaoLotacao                                  string   `json:"orgaoLotacao"`
	UorgLotacao                                   string   `json:"uorgLotacao"`
	OrgaoServidorLotacao                          string   `json:"orgaoServidorLotacao"`
	DataIngressoOrgao                             string   `json:"dataIngressoOrgao"`
	DataIngressoServicoPublico                    string   `json:"dataIngressoServicoPublico"`
	OrgaoSuperiorExercicio                        string   `json:"orgaoSuperiorExercicio"`
	OrgaoExercicio                                string   `json:"orgaoExercicio"`
	OrgaoServidorExercicio                        string   `json:"orgaoServidorExercicio"`
	UorgExercicio                                 string   `json:"uorgExercicio"`
	Cargo                                         string   `json:"cargo"`
	ClasseCargo                                   string   `json:"classeCargo"`
	PadraoCargo                                   string   `json:"padraoCargo"`
	NivelCargo                                    string   `json:"nivelCargo"`
	DataIngressoCargo                             string   `json:"dataIngressoCargo"`
	FormaIngresso                                 string   `json:"formaIngresso"`
	UFExercicio                                   string   `json:"ufExercicio"`
}

type FichaFuncao struct {
	Nome                                          string   `json:"nome"`
	CPFDescaracterizado                           string   `json:"cpfDescaracterizado"`
	MatriculaDescaracterizada                     string   `json:"matriculaDescaracterizada"`
	DataPublicacaoDocumentoIngressoServicoPublico string   `json:"dataPublicacaoDocumentoIngressoServicoPublico"`
	DiplomaLegal                                  string   `json:"diplomaLegal"`
	JornadaTrabalho                               string   `json:"jornadaTrabalho"`
	RegimeJuridico                                string   `json:"regimeJuridico"`
	SituacaoServidor                              string   `json:"situacaoServidor"`
	Afastamentos                                  []string `json:"afastamentos"`
	OrgaoSuperiorLotacao                          string   `json:"orgaoSuperiorLotacao"`
	OrgaoLotacao                                  string   `json:"orgaoLotacao"`
	UorgLotacao                                   string   `json:"uorgLotacao"`
	OrgaoServidorLotacao                          string   `json:"orgaoServidorLotacao"`
	DataIngressoOrgao                             string   `json:"dataIngressoOrgao"`
	DataIngressoServicoPublico                    string   `json:"dataIngressoServicoPublico"`
	OrgaoSuperiorExercicio                        string   `json:"orgaoSuperiorExercicio"`
	OrgaoExercicio                                string   `json:"orgaoExercicio"`
	UorgExercicio                                 string   `json:"uorgExercicio"`
	OrgaoServidorExercicio                        string   `json:"orgaoServidorExercicio"`
	Funcao                                        string   `json:"funcao"`
	Atividade                                     string   `json:"atividade"`
	OpcaoFuncao                                   string   `json:"opcaoFuncao"`
	DataIngressoFuncao                            string   `json:"dataIngressoFuncao"`
	UFExercicio                                   string   `json:"ufExercicio"`
}

type FichaMilitar struct {
	Nome                                          string   `json:"nome"`
	CPFDescaracterizado                           string   `json:"cpfDescaracterizado"`
	MatriculaDescaracterizada                     string   `json:"matriculaDescaracterizada"`
	DataPublicacaoDocumentoIngressoServicoPublico string   `json:"dataPublicacaoDocumentoIngressoServicoPublico"`
	DiplomaLegal                                  string   `json:"diplomaLegal"`
	JornadaTrabalho                               string   `json:"jornadaTrabalho"`
	RegimeJuridico                                string   `json:"regimeJuridico"`
	SituacaoServidor                              string   `json:"situacaoServidor"`
	Afastamentos                                  []string `json:"afastamentos"`
	OrgaoSuperior                                 string   `json:"orgaoSuperior"`
	Orgao                                         string   `json:"orgao"`
	OrgaoServidorLotacao                          string   `json:"orgaoServidorLotacao"`
	Cargo                                         string   `json:"cargo"`
	DataIngressoOrgao                             string   `json:"dataIngressoOrgao"`
}

type FichaServidorCivil struct {
	Nome                                          string   `json:"nome"`
	CPFDescaracterizado                           string   `json:"cpfDescaracterizado"`
	MatriculaDescaracterizada                     string   `json:"matriculaDescaracterizada"`
	DataPublicacaoDocumentoIngressoServicoPublico string   `json:"dataPublicacaoDocumentoIngressoServicoPublico"`
	DiplomaLegal                                  string   `json:"diplomaLegal"`
	JornadaTrabalho                               string   `json:"jornadaTrabalho"`
	RegimeJuridico                                string   `json:"regimeJuridico"`
	SituacaoServidor                              string   `json:"situacaoServidor"`
	Afastamentos                                  []string `json:"afastamentos"`
	OrgaoSuperiorLotacao                          string   `json:"orgaoSuperiorLotacao"`
	OrgaoLotacao                                  string   `json:"orgaoLotacao"`
	UorgLotacao                                   string   `json:"uorgLotacao"`
	OrgaoServidorLotacao                          string   `json:"orgaoServidorLotacao"`
	DataIngressoOrgao                             string   `json:"dataIngressoOrgao"`
	DataIngressoServicoPublico                    string   `json:"dataIngressoServicoPublico"`
}

type FichaAposentadoria struct {
	Nome                                          string   `json:"nome"`
	CPFDescaracterizado                           string   `json:"cpfDescaracterizado"`
	MatriculaDescaracterizada                     string   `json:"matriculaDescaracterizada"`
	DataPublicacaoDocumentoIngressoServicoPublico string   `json:"dataPublicacaoDocumentoIngressoServicoPublico"`
	DiplomaLegal                                  string   `json:"diplomaLegal"`
	JornadaTrabalho                               string   `json:"jornadaTrabalho"`
	RegimeJuridico                                string   `json:"regimeJuridico"`
	SituacaoServidor                              string   `json:"situacaoServidor"`
	Afastamentos                                  []string `json:"afastamentos"`
	OrgaoSuperiorLotacao                          string   `json:"orgaoSuperiorLotacao"`
	OrgaoLotacao                                  string   `json:"orgaoLotacao"`
	UorgLotacao                                   string   `json:"uorgLotacao"`
	OrgaoServidorLotacao                          string   `json:"orgaoServidorLotacao"`
	DataIngressoOrgao                             string   `json:"dataIngressoOrgao"`
	DataIngressoServicoPublico                    string   `json:"dataIngressoServicoPublico"`
	FormaIngresso                                 string   `json:"formaIngresso"`
	DataIngressoCargo                             string   `json:"dataIngressoCargo"`
	Cargo                                         string   `json:"cargo"`
	TipoAposentadoria                             string   `json:"tipoAposentadoria"`
	FundamentacaoAposentadoria                    string   `json:"fundamentacaoAposentadoria"`
	DataAposentadoria                             string   `json:"dataAposentadoria"`
}

type FichaReformado struct {
	Nome                                          string   `json:"nome"`
	CPFDescaracterizado                           string   `json:"cpfDescaracterizado"`
	MatriculaDescaracterizada                     string   `json:"matriculaDescaracterizada"`
	DataPublicacaoDocumentoIngressoServicoPublico string   `json:"dataPublicacaoDocumentoIngressoServicoPublico"`
	DiplomaLegal                                  string   `json:"diplomaLegal"`
	JornadaTrabalho                               string   `json:"jornadaTrabalho"`
	RegimeJuridico                                string   `json:"regimeJuridico"`
	SituacaoServidor                              string   `json:"situacaoServidor"`
	Afastamentos                                  []string `json:"afastamentos"`
	OrgaoSuperior                                 string   `json:"orgaoSuperior"`
	Orgao                                         string   `json:"orgao"`
	OrgaoServidorLotacao                          string   `json:"orgaoServidorLotacao"`
	Cargo                                         string   `json:"cargo"`
	DataIngressoOrgao                             string   `json:"dataIngressoOrgao"`
	TipoAposentadoria                             string   `json:"tipoAposentadoria"`
	FundamentacaoAposentadoria                    string   `json:"fundamentacaoAposentadoria"`
	DataReforma                                   string   `json:"dataReforma"`
}

type FichaPensaoCivil struct {
	Nome                                          string   `json:"nome"`
	CPFDescaracterizado                           string   `json:"cpfDescaracterizado"`
	MatriculaDescaracterizada                     string   `json:"matriculaDescaracterizada"`
	DataPublicacaoDocumentoIngressoServicoPublico string   `json:"dataPublicacaoDocumentoIngressoServicoPublico"`
	DiplomaLegal                                  string   `json:"diplomaLegal"`
	JornadaTrabalho                               string   `json:"jornadaTrabalho"`
	RegimeJuridico                                string   `json:"regimeJuridico"`
	SituacaoServidor                              string   `json:"situacaoServidor"`
	Afastamentos                                  []string `json:"afastamentos"`
	OrgaoSuperiorLotacao                          string   `json:"orgaoSuperiorLotacao"`
	OrgaoLotacao                                  string   `json:"orgaoLotacao"`
	UorgLotacao                                   string   `json:"uorgLotacao"`
	OrgaoServidorLotacao                          string   `json:"orgaoServidorLotacao"`
	DataIngressoOrgao                             string   `json:"dataIngressoOrgao"`
	DataIngressoServicoPublico                    string   `json:"dataIngressoServicoPublico"`
	FormaIngresso                                 string   `json:"formaIngresso"`
	DataIngressoCargo                             string   `json:"dataIngressoCargo"`
	Cargo                                         string   `json:"cargo"`
	TipoPensao                                    string   `json:"tipoPensao"`
	FundamentacaoPensao                           string   `json:"fundamentacaoPensao"`
	DataInicioPensao                              string   `json:"dataInicioPensao"`
	ProporcaoPensao                               string   `json:"proporcaoPensao"`
	RepresentanteLegal                            string   `json:"representanteLegal"`
	CPFRepresentanteLegal                         string   `json:"cpfRepresentanteLegal"`
	NomeInstituidor                               string   `json:"nomeInstituidor"`
	CPFInstituidor                                string   `json:"cpfInstituidor"`
}

type FichaPensaoMilitar struct {
	Nome                                          string   `json:"nome"`
	CPFDescaracterizado                           string   `json:"cpfDescaracterizado"`
	MatriculaDescaracterizada                     string   `json:"matriculaDescaracterizada"`
	DataPublicacaoDocumentoIngressoServicoPublico string   `json:"dataPublicacaoDocumentoIngressoServicoPublico"`
	DiplomaLegal                                  string   `json:"diplomaLegal"`
	JornadaTrabalho                               string   `json:"jornadaTrabalho"`
	RegimeJuridico                                string   `json:"regimeJuridico"`
	SituacaoServidor                              string   `json:"situacaoServidor"`
	Afastamentos                                  []string `json:"afastamentos"`
	OrgaoSuperior                                 string   `json:"orgaoSuperior"`
	Orgao                                         string   `json:"orgao"`
	OrgaoServidorLotacao                          string   `json:"orgaoServidorLotacao"`
	Cargo                                         string   `json:"cargo"`
	DataIngressoOrgao                             string   `json:"dataIngressoOrgao"`
	TipoPensao                                    string   `json:"tipoPensao"`
	FundamentacaoPensao                           string   `json:"fundamentacaoPensao"`
	DataInicioPensao                              string   `json:"dataInicioPensao"`
	ProporcaoPensao                               string   `json:"proporcaoPensao"`
	RepresentanteLegal                            string   `json:"representanteLegal"`
	CPFRepresentanteLegal                         string   `json:"cpfRepresentanteLegal"`
	NomeInstituidor                               string   `json:"nomeInstituidor"`
	CPFInstituidor                                string   `json:"cpfInstituidor"`
}

type ServidorRemuneracao struct {
	Servidor     ServidorAposentadoPensionista `json:"servidor"`
	Remuneracoes []Remuneracao                 `json:"remuneracoesDTO"`
}

type Remuneracao struct {
	SkMesReferencia                        string                   `json:"skMesReferencia"`
	MesAno                                 string                   `json:"mesAno"`
	ValorTotalRemuneracaoAposDeducoes      string                   `json:"valorTotalRemuneracaoAposDeducoes"`
	ValorTotalRemuneracaoDolarAposDeducoes string                   `json:"valorTotalRemuneracaoDolarAposDeducoes"`
	ValorTotalJetons                       string                   `json:"valorTotalJetons"`
	ValorTotalHonorariosAdvocaticios       string                   `json:"valorTotalHonorariosAdvocaticios"`
	Rubricas                               []Rubrica                `json:"rubricas"`
	Jetons                                 []Jetom                  `json:"jetons"`
	HonorariosAdvocaticios                 []HonorariosAdvocaticios `json:"honorariosAdvocaticios"`
	Observacoes                            []string                 `json:"observacoes"`
	RemuneracaoBasicaBruta                 string                   `json:"remuneracaoBasicaBruta"`
	RemuneracaoBasicaBrutaDolar            string                   `json:"remuneracaoBasicaBrutaDolar"`
	AbateRemuneracaoBasicaBruta            string                   `json:"abateRemuneracaoBasicaBruta"`
	AbateRemuneracaoBasicaBrutaDolar       string                   `json:"abateRemuneracaoBasicaBrutaDolar"`
	GratificacaoNatalina                   string                   `json:"gratificacaoNatalina"`
	GratificacaoNatalinaDolar              string                   `json:"gratificacaoNatalinaDolar"`
	AbateGratificacaoNatalina              string                   `json:"abateGratificacaoNatalina"`
	AbateGratificacaoNatalinaDolar         string                   `json:"abateGratificacaoNatalinaDolar"`
	Ferias                                 string                   `json:"ferias"`
	FeriasDolar                            string                   `json:"feriasDolar"`
	OutrasRemuneracoesEventuais            string                   `json:"outrasRemuneracoesEventuais"`
	OutrasRemuneracoesEventuaisDolar       string                   `json:"outrasRemuneracoesEventuaisDolar"`
	ImpostoRetidoNaFonte                   string                   `json:"impostoRetidoNaFonte"`
	ImpostoRetidoNaFonteDolar              string                   `json:"impostoRetidoNaFonteDolar"`
	PrevidenciaOficial                     string                   `json:"previdenciaOficial"`
	PrevidenciaOficialDolar                string                   `json:"previdenciaOficialDolar"`
	OutrasDeducoesObrigatorias             string                   `json:"outrasDeducoesObrigatorias"`
	OutrasDeducoesObrigatoriasDolar        string                   `json:"outrasDeducoesObrigatoriasDolar"`
	PensaoMilitar                          string                   `json:"pensaoMilitar"`
	PensaoMilitarDolar                     string                   `json:"pensaoMilitarDolar"`
	FundoSaude                             string                   `json:"fundoSaude"`
	FundoSaudeDolar                        string                   `json:"fundoSaudeDolar"`
	TaxaOcupacaoImovelFuncional            string                   `json:"taxaOcupacaoImovelFuncional"`
	TaxaOcupacaoImovelFuncionalDolar       string                   `json:"taxaOcupacaoImovelFuncionalDolar"`
	VerbasIndenizatoriasCivil              string                   `json:"verbasIndenizatoriasCivil"`
	VerbasIndenizatoriasCivilDolar         string                   `json:"verbasIndenizatoriasCivilDolar"`
	VerbasIndenizatoriasMilitar            string                   `json:"verbasIndenizatoriasMilitar"`
	VerbasIndenizatoriasMilitarDolar       string                   `json:"verbasIndenizatoriasMilitarDolar"`
	VerbasIndenizatoriasReferentesPDV      string                   `json:"verbasIndenizatoriasReferentesPDV"`
	VerbasIndenizatoriasReferentesPDVDolar string                   `json:"verbasIndenizatoriasReferentesPDVDolar"`
	RemuneracaoEmpresaPublica              bool                     `json:"remuneracaoEmpresaPublica"`
	ExisteValorMes                         bool                     `json:"existeValorMes"`
	MesAnoPorExtenso                       string                   `json:"mesAnoPorExtenso"`
	VerbasIndenizatorias                   string                   `json:"verbasIndenizatorias"`
	VerbasIndenizatoriasDolar              string                   `json:"verbasIndenizatoriasDolar"`
}

type Rubrica struct {
	Codigo          string  `json:"codigo"`
	Descricao       string  `json:"descricao"`
	Valor           float64 `json:"valor"`
	SkMesReferencia string  `json:"skMesReferencia"`
	ValorDolar      float64 `json:"valorDolar"`
}

type Jetom struct {
	Descricao     string  `json:"descricao"`
	Valor         float64 `json:"valor"`
	MesReferencia string  `json:"mesReferencia"`
}

type HonorariosAdvocaticios struct {
	MesReferencia         string  `json:"mesReferencia"`
	Valor                 float64 `json:"valor"`
	ValorFormatado        string  `json:"valorFormatado"`
	MensagemMesReferencia string  `json:"mensagemMesReferencia"`
}

type ServidorPorOrgao struct {
	QntPessoas                     int    `json:"qntPessoas"`
	QntVinculos                    int    `json:"qntVinculos"`
	Situacao                       int    `json:"skSituacao"`
	DescSituacao                   string `json:"descSituacao"`
	TipoVinculo                    int    `json:"skTipoVinculo"`
	DescTipoVinculo                string `json:"descTipoVinculo"`
	TipoServidor                   int    `json:"skTipoServidor"`
	DescTipoServidor               string `json:"descTipoServidor"`
	Licenca                        int    `json:"licenca"`
	CodOrgaoExercicioSiape         string `json:"codOrgaoExercicioSiape"`
	NomOrgaoExercicioSiape         string `json:"nomOrgaoExercicioSiape"`
	CodOrgaoSuperiorExercicioSiape string `json:"codOrgaoSuperiorExercicioSiape"`
	NomOrgaoSuperiorExercicioSiape string `json:"nomOrgaoSuperiorExercicioSiape"`
}

type PEP struct {
	CPF                 string `json:"cpf"`
	Nome                string `json:"nome"`
	SiglaFuncao         string `json:"sigla_funcao"`
	DescricaoFuncao     string `json:"descricao_funcao"`
	NivelFuncao         string `json:"nivel_funcao"`
	CodOrgao            string `json:"cod_orgao"`
	NomeOrgao           string `json:"nome_orgao"`
	DataInicioExercicio string `json:"dt_inicio_exercicio"`
	DataFimExercicio    string `json:"dt_fim_exercicio"`
	DataFimCarencia     string `json:"dt_fim_carencia"`
}
