package deputados

type Bloco struct {
	ID            string `json:"id"`
	Nome          string `json:"nome"`
	URI           string `json:"uri"`
	IDLegislatura string `json:"idLegislatura"`
	Federacao     bool   `json:"federacao"`
}
