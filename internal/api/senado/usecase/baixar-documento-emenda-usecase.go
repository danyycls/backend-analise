package usecase

import (
	"context"

	senado "github.com/danyele/podp/internal/sources/senado/client"
)

type BaixarDocumentoEmendaUseCase struct {
	client *senado.SenadoClient
}

func NovoBaixarDocumentoEmendaUseCase(c *senado.SenadoClient) *BaixarDocumentoEmendaUseCase {
	return &BaixarDocumentoEmendaUseCase{client: c}
}

func (u *BaixarDocumentoEmendaUseCase) Executar(ctx context.Context, idDocumento int) ([]byte, string, error) {
	return u.client.BaixarDocumentoEmenda(ctx, idDocumento)
}
