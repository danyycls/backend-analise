package deputados

import "encoding/json"

type Link struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
	Type string `json:"type"`
}

type Resultado struct {
	Dados json.RawMessage `json:"dados"`
	Links []Link          `json:"links"`
}
