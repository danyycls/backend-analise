package deputados

type Partido struct {
	ID    int    `json:"id"`
	Sigla string `json:"sigla"`
	Nome  string `json:"nome"`
	URI   string `json:"uri"`
}
