package usecase

//go:generate mockgen -source=interface.go -destination=mock.go -package=usecase

import "context"

type AnalisarLigacaoPoliticaUseCaseInterface interface {
	Executar(ctx context.Context, licitacoes []AnalisarLigacaoPoliticaRequest) (*AnalisarLigacaoPoliticaResponse, error)
}

var _ AnalisarLigacaoPoliticaUseCaseInterface = (*AnalisarLigacaoPoliticaUseCase)(nil)
