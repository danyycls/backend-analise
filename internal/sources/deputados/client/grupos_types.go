package deputados

type UltimoStatusGrupo struct {
	CodTipoOrgao  string `json:"codTipoOrgao"`
	TipoOrgao     string `json:"tipoOrgao"`
	DataStatus    string `json:"dataStatus"`
	IDLegislatura string `json:"idLegislatura"`
	Situacao      string `json:"situacao"`
	NomeCasa      string `json:"nomeCasa"`
}

type Grupo struct {
	ID              int               `json:"id"`
	URI             string            `json:"uri"`
	Nome            string            `json:"nome"`
	AnoCriacao      string            `json:"anoCriacao"`
	Ativo           string            `json:"ativo"`
	GrupoMisto      string            `json:"grupoMisto"`
	Subvencionado   string            `json:"subvencionado"`
	ResolucaoTitulo string            `json:"resolucaoTitulo"`
	ResolucaoURI    string            `json:"resolucaoUri"`
	Observacoes     string            `json:"observacoes"`
	UltimoStatus    UltimoStatusGrupo `json:"ultimoStatus"`
}

type GrupoDetalhe struct {
	ID              int               `json:"id"`
	URI             string            `json:"uri"`
	Nome            string            `json:"nome"`
	AnoCriacao      string            `json:"anoCriacao"`
	Ativo           string            `json:"ativo"`
	GrupoMisto      string            `json:"grupoMisto"`
	Subvencionado   string            `json:"subvencionado"`
	ResolucaoTitulo string            `json:"resolucaoTitulo"`
	ResolucaoURI    string            `json:"resolucaoUri"`
	ProjetoTitulo   string            `json:"projetoTitulo"`
	ProjetoURI      string            `json:"projetoUri"`
	Observacoes     string            `json:"observacoes"`
	UltimoStatus    UltimoStatusGrupo `json:"ultimoStatus"`
}

type MembroGrupo struct {
	IDLegislatura string `json:"idLegislatura"`
	Nome          string `json:"nome"`
	URI           string `json:"uri"`
	Cargo         string `json:"cargo"`
	Tipo          string `json:"tipo"`
	DataInicio    string `json:"dataInicio"`
	DataFim       string `json:"dataFim"`
	OrdemEntrada  int    `json:"ordemEntrada"`
}

type HistoricoGrupo struct {
	IDLegislatura          string `json:"idLegislatura"`
	DataStatus             string `json:"dataStatus"`
	Situacao               string `json:"situacao"`
	Observacao             string `json:"observacao"`
	Presidente             string `json:"presidente"`
	PresidenteURI          string `json:"presidenteUri"`
	OficioTitulo           string `json:"oficioTitulo"`
	OficioAutor            string `json:"oficioAutor"`
	OficioAutorTipo        string `json:"oficioAutorTipo"`
	OficioAutorURI         string `json:"oficioAutorUri"`
	OficioDataApresentacao string `json:"oficioDataApresentacao"`
	OficioDataPublicacao   string `json:"oficioDataPublicacao"`
	DocumentoSGM           string `json:"documentoSGM"`
}
