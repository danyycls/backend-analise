package types

type ArquivoProcessado struct {
	NomeArquivo        string `json:"nome_arquivo"`
	Tipo               string `json:"tipo"`
	Registros          int    `json:"registros"`
	RegistrosIgnorados int    `json:"registros_ignorados"`
	HashSHA256         string `json:"hash_sha256,omitempty"`
}

type ArquivoImportacao struct {
	Caminho         string
	CaminhoRelativo string
	Diretorio       string
	DiretorioLower  string
	Nome            string
	NomeLower       string
	Tipo            string
}
