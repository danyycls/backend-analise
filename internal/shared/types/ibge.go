package types

type MunicipioIBGE struct {
	ID   int    `json:"id"`
	Nome string `json:"nome"`
}

type EstadoIBGE struct {
	ID    int    `json:"id"`
	Sigla string `json:"sigla"`
	Nome  string `json:"nome"`
}

type MunicipioDetalhadoIBGE struct {
	ID           int          `json:"id"`
	Nome         string       `json:"nome"`
	Microrregiao Microrregiao `json:"microrregiao"`
}

type Microrregiao struct {
	ID          int          `json:"id"`
	Nome        string       `json:"nome"`
	Mesorregiao Mesorregiao2 `json:"mesorregiao"`
}

type Mesorregiao2 struct {
	ID   int        `json:"id"`
	Nome string     `json:"nome"`
	UF   EstadoIBGE `json:"UF"`
}
