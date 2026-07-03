package usecase

import (
	"context"
	"time"

	"github.com/danyele/podp/internal/shared/domain"
	"github.com/danyele/podp/internal/shared/logger"
	"github.com/danyele/podp/internal/shared/mongodb"
	anomalia "github.com/danyele/podp/internal/worker/anomalia"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AnomaliaConsultaUseCase struct {
	mongo mongodb.Client
}

func NovoAnomaliaConsultaUseCase(mongo mongodb.Client) *AnomaliaConsultaUseCase {
	return &AnomaliaConsultaUseCase{mongo: mongo}
}

type AnomaliaFiltro struct {
	Documento string
	Uf        string
	Municipio string
	Tag       string
	Categoria string
	Pagina    int
	PorPagina int
}

type ListarAnomaliasResultado struct {
	Total     int
	Anomalias []anomalia.AnomaliaDocumento
}

func (u *AnomaliaConsultaUseCase) Listar(ctx context.Context, filtro AnomaliaFiltro) (*ListarAnomaliasResultado, error) {
	log := logger.New("Worker: Anomalia: ConsultaUseCase: Listar")

	filter := bson.M{}

	if filtro.Documento != "" {
		filter["documento_fornecedor_pncp"] = filtro.Documento
	}
	if filtro.Uf != "" {
		filter["uf"] = filtro.Uf
	}
	if filtro.Municipio != "" {
		filter["municipio"] = filtro.Municipio
	}
	switch {
	case filtro.Tag != "" && filtro.Categoria != "":
		filter["$and"] = []bson.M{
			{"tags": filtro.Tag},
			{"tags": primitive.Regex{Pattern: "-" + filtro.Categoria + "$"}},
		}
	case filtro.Tag != "":
		filter["tags"] = filtro.Tag
	case filtro.Categoria != "":
		filter["tags"] = primitive.Regex{Pattern: "-" + filtro.Categoria + "$"}
	}

	opts := options.Find()

	if filtro.Pagina > 0 && filtro.PorPagina > 0 {
		skip := int64((filtro.Pagina - 1) * filtro.PorPagina)
		limit := int64(filtro.PorPagina)
		opts.SetSkip(skip)
		opts.SetLimit(limit)
	}

	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	results, err := u.mongo.Find(ctx, "anomalias", filter, opts)
	if err != nil {
		log.Error("erro ao buscar anomalias", "erro", err)
		return nil, err
	}

	anomalias := make([]anomalia.AnomaliaDocumento, 0, len(results))
	for _, r := range results {
		a := anomalia.AnomaliaDocumento{}
		if id, ok := r["_id"]; ok {
			a.ID = id.(primitive.ObjectID)
		}
		if v, ok := r["job_id"]; ok {
			a.JobID = v.(string)
		}
		if v, ok := r["documento_fornecedor_pncp"]; ok {
			a.DocumentoFornecedorPNCP = v.(string)
		}
		if v, ok := r["nome_fornecedor_pncp"]; ok {
			a.NomeFornecedorPNCP = v.(string)
		}
		if v, ok := r["numero_controle_pncp"]; ok {
			a.NumeroControlePncp = v.(string)
		}
		if v, ok := r["orgao_cnpj"]; ok {
			a.OrgaoCNPJ = v.(string)
		}
		if v, ok := r["orgao_nome"]; ok {
			a.OrgaoNome = v.(string)
		}
		if v, ok := r["uf"]; ok {
			a.Uf = v.(string)
		}
		if v, ok := r["municipio"]; ok {
			a.Municipio = v.(string)
		}
		if v, ok := r["titulo"]; ok {
			a.Titulo = v.(string)
		}
		if v, ok := r["tags"]; ok {
			tagsRaw := v.(bson.A)
			a.Tags = make([]string, len(tagsRaw))
			for i, t := range tagsRaw {
				a.Tags[i] = t.(string)
			}
		}
		if v, ok := r["documentos_vinculos"]; ok {
			documentosBSON, err := bson.Marshal(bson.M{"documentos_vinculos": v})
			if err == nil {
				var wrapper struct {
					DocumentosVinculos []domain.DocumentoVinculo `bson:"documentos_vinculos"`
				}
				_ = bson.Unmarshal(documentosBSON, &wrapper)
				a.DocumentosVinculos = wrapper.DocumentosVinculos
			}
		}
		if v, ok := r["created_at"]; ok {
			a.CreatedAt = v.(primitive.DateTime).Time()
		}
		anomalias = append(anomalias, a)
	}

	return &ListarAnomaliasResultado{
		Total:     len(anomalias),
		Anomalias: anomalias,
	}, nil
}

func (u *AnomaliaConsultaUseCase) Inserir(ctx context.Context, anomalia anomalia.AnomaliaDocumento) error {
	log := logger.New("Worker: Anomalia: ConsultaUseCase: Inserir")
	anomalia.CreatedAt = time.Now()

	_, err := u.mongo.InsertOne(ctx, "anomalias", anomalia)
	if err != nil {
		log.Error("erro ao inserir anomalia", "erro", err)
		return err
	}
	return nil
}
