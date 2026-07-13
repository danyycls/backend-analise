package deputados

type Votacao struct {
	ID                 int    `json:"id"`
	URI                string `json:"uri"`
	Data               string `json:"data"`
	DataHoraRegistro   string `json:"dataHoraRegistro"`
	SiglaOrgao         string `json:"siglaOrgao"`
	URIOrgao           string `json:"uriOrgao"`
	URIEvento          string `json:"uriEvento"`
	DescricaoResultado string `json:"descricaoResultado"`
	ObjetoSiglaTipo    string `json:"objetoSiglaTipo"`
	ObjetoNumero       string `json:"objetoNumero"`
	ObjetoAno          string `json:"objetoAno"`
	ObjetoURI          string `json:"objetoURI"`
	ObjetoCodTipo      string `json:"objetoCodTipo"`
	ObjetoTitulo       string `json:"objetoTitulo"`
	Aprovacao          string `json:"aprovacao"`
}

type EfeitosRegistrados struct {
	DataHoraResultado                    string `json:"dataHoraResultado"`
	DescResultado                        string `json:"descResultado"`
	TituloProposicao                     string `json:"tituloProposicao"`
	URIProposicao                        string `json:"uriProposicao"`
	DataHoraUltimaAberturaVotacao        string `json:"dataHoraUltimaAberturaVotacao"`
	DescUltimaAberturaVotacao            string `json:"descUltimaAberturaVotacao"`
	TituloProposicaoCitada               string `json:"tituloProposicaoCitada"`
	URIProposicaoCitada                  string `json:"uriProposicaoCitada"`
	DataHoraUltimaApresentacaoProposicao string `json:"dataHoraUltimaApresentacaoProposicao"`
	DescUltimaApresentacaoProposicao     string `json:"descUltimaApresentacaoProposicao"`
}

type VotacaoDetalhe struct {
	ID               int                `json:"id"`
	URI              string             `json:"uri"`
	Data             string             `json:"data"`
	DataHoraRegistro string             `json:"dataHoraRegistro"`
	SiglaOrgao       string             `json:"siglaOrgao"`
	URIOrgao         string             `json:"uriOrgao"`
	URIEvento        string             `json:"uriEvento"`
	Proposicao       string             `json:"proposicao"`
	ObjetoSiglaTipo  string             `json:"objetoSiglaTipo"`
	ObjetoNumero     string             `json:"objetoNumero"`
	ObjetoAno        string             `json:"objetoAno"`
	ObjetoURI        string             `json:"objetoURI"`
	Efeitos          EfeitosRegistrados `json:"efeitos"`
}

type Orientacao struct {
	CodPartidoBloco   int    `json:"codPartidoBloco"`
	SiglaPartidoBloco string `json:"siglaPartidoBloco"`
	URIPartidoBloco   string `json:"uriPartidoBloco"`
	CodTipoLideranca  string `json:"codTipoLideranca"`
	OrientacaoVoto    string `json:"orientacaoVoto"`
}

type Voto struct {
	Deputado         Deputado `json:"deputado_"`
	TipoVoto         string   `json:"tipoVoto"`
	DataRegistroVoto string   `json:"dataRegistroVoto"`
}
