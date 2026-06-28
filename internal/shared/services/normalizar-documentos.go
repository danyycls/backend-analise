package services

import (
	"context"

	"github.com/danyele/podp/internal/shared/domain"
	"github.com/danyele/podp/internal/shared/utils"
)

type NormalizarDocumentosInput struct {
	Documento string
	Nome      string
	Origem    string
}

type NormalizarDocumentosOutput struct {
	domain.DocumentoVinculo
}

type NormalizarDocumentosService interface {
	Executar(ctx context.Context, input NormalizarDocumentosInput) (*NormalizarDocumentosOutput, error)
}

type normalizarDocumentosServiceImpl struct{}

func NovoNormalizarDocumentosService() NormalizarDocumentosService {
	return &normalizarDocumentosServiceImpl{}
}

func (s *normalizarDocumentosServiceImpl) Executar(_ context.Context, input NormalizarDocumentosInput) (*NormalizarDocumentosOutput, error) {
	doc, parcial := utils.NormalizarDocumento(input.Documento)

	return &NormalizarDocumentosOutput{
		DocumentoVinculo: domain.DocumentoVinculo{
			DocumentoInput:       input.Documento,
			DocumentoNormalizado: doc,
			Nome:                 input.Nome,
			Parcial:              parcial,
			Origem:               input.Origem,
		},
	}, nil
}
