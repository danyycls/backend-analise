package anomalia

import (
	"time"

	"github.com/danyele/podp/internal/shared/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AnomaliaDocumento struct {
	ID                      primitive.ObjectID        `bson:"_id,omitempty" json:"id,omitempty"`
	JobID                   string                    `bson:"job_id" json:"job_id"`
	NomeFornecedorPNCP      string                    `bson:"nome_fornecedor_pncp" json:"nome_fornecedor_pncp"`
	DocumentoFornecedorPNCP string                    `bson:"documento_fornecedor_pncp" json:"documento_fornecedor_pncp"`
	NumeroControlePncp      string                    `bson:"numero_controle_pncp" json:"numero_controle_pncp"`
	OrgaoCNPJ               string                    `bson:"orgao_cnpj" json:"orgao_cnpj"`
	OrgaoNome               string                    `bson:"orgao_nome" json:"orgao_nome"`
	Uf                      string                    `bson:"uf" json:"uf"`
	Municipio               string                    `bson:"municipio" json:"municipio"`
	Titulo                  string                    `bson:"titulo" json:"titulo"`
	Tags                    []string                  `bson:"tags" json:"tags"`
	Socios                  []domain.SocioOutput      `bson:"socios,omitempty" json:"socios,omitempty"`
	DocumentosVinculos      []domain.DocumentoVinculo `bson:"documentos_vinculos,omitempty" json:"documentos_vinculos,omitempty"`
	CreatedAt               time.Time                 `bson:"created_at" json:"created_at"`
}
