package usecase

import (
	"context"
	"errors"

	repositorio "github.com/danyele/podp/internal/api/tse/repositorio"
	tsetypes "github.com/danyele/podp/internal/api/tse/types"
	"github.com/danyele/podp/internal/shared/database"

	"github.com/jackc/pgx/v5"
)

type BuscarFornecedorUseCase struct {
	db database.DB
}

func NovoBuscarFornecedorUseCase(db database.DB) *BuscarFornecedorUseCase {
	return &BuscarFornecedorUseCase{db: db}
}

func (u *BuscarFornecedorUseCase) Executar(ctx context.Context, documento string) (*tsetypes.FornecedorRelacoesResponse, error) {
	repo := repositorio.Novo(u.db)
	forn, err := repo.FornecedorBuscarPorCNPJExato(ctx, documento)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}
	dto := &tsetypes.FornecedorRelacoesResponse{
		Fornecedor: &tsetypes.FornecedorEnriquecido{
			CPFCNPJ: forn.CPFCNPJ,
			Nome:    forn.Nome,
		},
	}
	despsCand, _ := repo.DespesasCandidatoBuscarPorFornecedorIDComPrestacao(ctx, forn.ID.String())
	for _, d := range despsCand {
		dto.Despesas = append(dto.Despesas, despesaCandidatoParaDTO(ctx, repo, d))
	}
	despsPart, _ := repo.DespesasPartidoBuscarPorFornecedorIDComPrestacao(ctx, forn.ID.String())
	for _, d := range despsPart {
		dto.Despesas = append(dto.Despesas, despesaPartidoParaDTO(ctx, repo, d))
	}
	dto.TotalDespesas = len(dto.Despesas)
	return dto, nil
}
