package opencnpj

//go:generate mockgen -source=interface.go -destination=mock.go -package=opencnpj

import (
	"context"

	"github.com/danyele/laceu/internal/shared/types"
)

type Client interface {
	Buscar(ctx context.Context, cnpj string) (*types.OpenCNPJResponse, error)
}

var _ Client = (*OpenCNPJClient)(nil)
