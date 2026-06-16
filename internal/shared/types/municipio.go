package types

type CandidatoEleito struct {
	ID                           string `json:"id"`
	SQCandidato                  int64  `json:"sq_candidato"`
	NomeUrna                     string `json:"nome_urna"`
	NomeCompleto                 string `json:"nome_completo"`
	PartidoSigla                 string `json:"partido_sigla"`
	PartidoNome                  string `json:"partido_nome"`
	CargoNome                    string `json:"cargo_nome"`
	SituacaoTotalizacaoDescricao string `json:"situacao_totalizacao_descricao"`
	AnoEleicao                   int16  `json:"ano_eleicao"`
	NumeroCandidato              *int   `json:"numero_candidato,omitempty"`
	CPF                          string `json:"cpf,omitempty"`
	EleicaoDescricao             string `json:"eleicao_descricao,omitempty"`
	EleicaoData                  string `json:"eleicao_data,omitempty"`
	EleicaoTipo                  string `json:"eleicao_tipo,omitempty"`
}

type MunicipioComDados struct {
	ID        int    `json:"id"`
	Nome      string `json:"nome"`
	Populacao int64  `json:"populacao"`
}

type DadosEstadoConsolidado struct {
	UF            string              `json:"uf"`
	Nome          string              `json:"nome"`
	Populacao     int64               `json:"populacao"`
	Municipios    []MunicipioComDados `json:"municipios"`
	Prefeitos     []CandidatoEleito   `json:"prefeitos"`
	VicePrefeitos []CandidatoEleito   `json:"vice_prefeitos"`
	Vereadores    []CandidatoEleito   `json:"vereadores"`
	Senadores     []SenadorUF         `json:"senadores"`
	Deputados     []DeputadoUF        `json:"deputados"`
}

type DadosCandidatosEstado struct {
	Prefeitos     []CandidatoEleito `json:"prefeitos"`
	VicePrefeitos []CandidatoEleito `json:"vice_prefeitos"`
	Vereadores    []CandidatoEleito `json:"vereadores"`
}

type SenadorUF struct {
	Codigo          string `json:"codigo"`
	NomeParlamentar string `json:"nome_parlamentar"`
	NomeCompleto    string `json:"nome_completo"`
	Uf              string `json:"uf"`
	Partido         string `json:"partido"`
	URLFoto         string `json:"url_foto"`
}

type DeputadoUF struct {
	ID            int    `json:"id"`
	Nome          string `json:"nome"`
	SiglaPartido  string `json:"sigla_partido"`
	SiglaUF       string `json:"sigla_uf"`
	URLFoto       string `json:"url_foto"`
	Email         string `json:"email"`
	NomeEleitoral string `json:"nome_eleitoral"`
}

type DetalhesMunicipioResponse struct {
	CodigoIBGE       int                      `json:"codigo_ibge"`
	Nome             string                   `json:"nome"`
	UF               string                   `json:"uf"`
	Exercicio        int                      `json:"exercicio"`
	DespesaPessoal   *DespesaPessoalResumo    `json:"despesa_pessoal,omitempty"`
	GastosPorFuncao  []GastoPorFuncao         `json:"gastos_por_funcao,omitempty"`
	Receitas         []ReceitaResumo          `json:"receitas,omitempty"`
	RecursosFederais []RecursoFederalRecebido `json:"recursos_federais,omitempty"`
	Contratos        []ContratoPNCP           `json:"contratos,omitempty"`
	Servidores       []ServidorMunicipio      `json:"servidores,omitempty"`
}

type ProgressoBusca struct {
	Etapa    string `json:"etapa"`
	Status   string `json:"status"`
	Mensagem string `json:"mensagem"`
}

type DespesaPessoalResumo struct {
	ValorTotal    float64 `json:"valor_total"`
	PercentualRCL float64 `json:"percentual_rcl"`
	Poder         string  `json:"poder"`
	Periodo       string  `json:"periodo"`
}

type GastoPorFuncao struct {
	Funcao    string  `json:"funcao"`
	Empenhado float64 `json:"empenhado"`
	Liquidado float64 `json:"liquidado"`
	Pago      float64 `json:"pago"`
}

type ReceitaResumo struct {
	Conta     string  `json:"conta"`
	Coluna    string  `json:"coluna"`
	Valor     float64 `json:"valor"`
	Exercicio int64   `json:"exercicio"`
}

type RecursoFederalRecebido struct {
	NomePessoa        string  `json:"nome_pessoa"`
	TipoPessoa        string  `json:"tipo_pessoa"`
	NomeUG            string  `json:"nome_ug"`
	NomeOrgao         string  `json:"nome_orgao"`
	NomeOrgaoSuperior string  `json:"nome_orgao_superior"`
	Valor             float64 `json:"valor"`
	MesAno            int     `json:"mes_ano"`
}

type ContratoPNCP struct {
	Orgao                string  `json:"orgao"`
	Objeto               string  `json:"objeto"`
	Valor                float64 `json:"valor"`
	NomeRazaoSocial      string  `json:"nome_razao_social"`
	DataVigenciaInicio   string  `json:"data_vigencia_inicio"`
	DataVigenciaFim      string  `json:"data_vigencia_fim"`
	DataPublicacao       string  `json:"data_publicacao"`
	NumeroContrato       string  `json:"numero_contrato"`
	NumeroControlePNCP   string  `json:"numero_controle_pncp"`
	ModalidadeNome       string  `json:"modalidade_nome"`
	NumeroLicitacao      string  `json:"numero_licitacao"`
	CodigoContrato       string  `json:"codigo_contrato"`
	OrigemLicitacao      string  `json:"origem_licitacao"`
	TipoContratoNome     string  `json:"tipo_contrato_nome"`
	ValorGlobal          float64 `json:"valor_global"`
	ValorParcela         float64 `json:"valor_parcela"`
	ValorTotalEstimado   float64 `json:"valor_total_estimado"`
	ValorTotalHomologado float64 `json:"valor_total_homologado"`
	AnoContrato          int     `json:"ano_contrato"`
	DataAssinatura       string  `json:"data_assinatura"`
	AmpLegalDescricao    string  `json:"amp_legal_descricao"`
	Produto              string  `json:"produto"`
	SubtipoContrato      string  `json:"subtipo_contrato"`
}

type ServidorMunicipio struct {
	Categoria         string  `json:"categoria"`
	Quantidade        *int    `json:"quantidade,omitempty"`
	DespesaTotal      float64 `json:"despesa_total"`
	PercentualDespesa float64 `json:"percentual_despesa,omitempty"`
}

type DadosEstadoFinanceiroResumo struct {
	UF string `json:"uf"`
}
