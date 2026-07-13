package tcu

//go:generate mockgen -source=interface.go -destination=mock.go -package=tcu

import "context"

type Client interface {
	BuscarContasIrregulares(ctx context.Context, filter TCUQueryParams) ([]ContasIrregulares, error)
	BuscarInabilitados(ctx context.Context, filter TCUQueryParams) ([]Sancoes, error)
	BuscarInidoneos(ctx context.Context, filter TCUQueryParams) ([]Sancoes, error)
	BuscarFinsEleitorais(ctx context.Context, filter TCUQueryParams) ([]FinsEleitorais, error)
}

var _ Client = (*TCUClient)(nil)
