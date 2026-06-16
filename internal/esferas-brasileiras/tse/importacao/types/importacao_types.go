package types

type ArquivoProcessado struct {
	NomeArquivo string `json:"nome_arquivo"`
	Tipo        string `json:"tipo"`
	Registros   int    `json:"registros"`
}

type ArquivoImportacao struct {
	Caminho         string
	CaminhoRelativo string
	Diretorio       string
	Nome            string
	Tipo            string
}

type tipoInfo struct{}
