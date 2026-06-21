package usecase

import (
	"context"
	"errors"

	repositorio "github.com/danyele/podp/internal/esferas-brasileiras/tse/repositorio"
	tsetypes "github.com/danyele/podp/internal/esferas-brasileiras/tse/types"
	"github.com/danyele/podp/internal/shared/database"

	"github.com/jackc/pgx/v5"
)

type BuscarDoadorUseCase struct {
	db database.DB
}

func NovoBuscarDoadorUseCase(db database.DB) *BuscarDoadorUseCase {
	return &BuscarDoadorUseCase{db: db}
}

func (u *BuscarDoadorUseCase) Executar(ctx context.Context, documento string) (*tsetypes.DoadorRelacoesResponse, error) {
	repo := repositorio.Novo(u.db)
	doador, err := repo.DoadorBuscarPorCNPJExato(ctx, documento)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}
	dto := &tsetypes.DoadorRelacoesResponse{
		Doador: &tsetypes.DoadorRelacoes{
			CPFCNPJ: doador.CPFCNPJ,
			Nome:    doador.Nome,
		},
	}
	recsCand, _ := repo.ReceitasCandidatoBuscarPorDoadorIDComPrestacao(ctx, doador.ID.String())
	for _, r := range recsCand {
		dto.Receitas = append(dto.Receitas, receitaCandidatoParaDTO(ctx, repo, r))
	}
	recsPart, _ := repo.ReceitasOrgaoBuscarPorDoadorIDComPrestacao(ctx, doador.ID.String())
	for _, r := range recsPart {
		dto.Receitas = append(dto.Receitas, receitaPartidoParaDTO(ctx, repo, r))
	}
	dto.TotalReceitas = len(dto.Receitas)
	return dto, nil
}
