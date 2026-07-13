package usecase

import (
	"context"
	"errors"

	repositorio "github.com/danyele/podp/internal/api/tse/repositorio"
	tsetypes "github.com/danyele/podp/internal/api/tse/types"
	"github.com/danyele/podp/internal/shared/database"
	"github.com/danyele/podp/internal/shared/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/danyele/podp/internal/shared/types"
)

type BuscarRelacoesUseCase struct {
	db database.DB
}

func NovoBuscarRelacoesUseCase(db database.DB) *BuscarRelacoesUseCase {
	return &BuscarRelacoesUseCase{db: db}
}

func (u *BuscarRelacoesUseCase) Executar(ctx context.Context, cnpj string) (*tsetypes.RelacoesResponse, error) {
	repo := repositorio.Novo(u.db)

	resp := &tsetypes.RelacoesResponse{
		Despesas: []tsetypes.DespesaRelacao{},
		Receitas: []tsetypes.ReceitaRelacao{},
	}

	forn, err := repo.FornecedorBuscarPorCNPJExato(ctx, cnpj)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	if forn != nil {
		resp.Fornecedor = &tsetypes.FornecedorEnriquecido{
			CPFCNPJ: forn.CPFCNPJ,
			Nome:    forn.Nome,
			NomeRFB: forn.NomeRFB,
		}

		despsCand, _ := repo.DespesasCandidatoBuscarPorFornecedorIDComPrestacao(ctx, forn.ID.String())
		for _, d := range despsCand {
			dto := despesaCandidatoParaDTO(ctx, repo, d)
			resp.Despesas = append(resp.Despesas, dto)
		}

		despsPart, _ := repo.DespesasPartidoBuscarPorFornecedorIDComPrestacao(ctx, forn.ID.String())
		for _, d := range despsPart {
			dto := despesaPartidoParaDTO(ctx, repo, d)
			resp.Despesas = append(resp.Despesas, dto)
		}
	}

	doador, err := repo.DoadorBuscarPorCNPJExato(ctx, cnpj)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	if doador != nil {
		resp.Doador = &tsetypes.DoadorRelacoes{
			CPFCNPJ: doador.CPFCNPJ,
			Nome:    doador.Nome,
		}

		recsCand, _ := repo.ReceitasCandidatoBuscarPorDoadorIDComPrestacao(ctx, doador.ID.String())
		for _, r := range recsCand {
			dto := receitaCandidatoParaDTO(ctx, repo, r)
			resp.Receitas = append(resp.Receitas, dto)
		}

		recsPart, _ := repo.ReceitasOrgaoBuscarPorDoadorIDComPrestacao(ctx, doador.ID.String())
		for _, r := range recsPart {
			dto := receitaPartidoParaDTO(ctx, repo, r)
			resp.Receitas = append(resp.Receitas, dto)
		}
	}

	resp.TotalDespesas = len(resp.Despesas)
	resp.TotalReceitas = len(resp.Receitas)

	return resp, nil
}

func despesaCandidatoParaDTO(ctx context.Context, repo *repositorio.Repositorio, d *types.DespesaCandidato) tsetypes.DespesaRelacao {
	dto := tsetypes.DespesaRelacao{
		SQDespesa:              d.SQDespesa,
		Tipo:                   "candidato",
		TipoRegistro:           d.TipoRegistro,
		Descricao:              d.Descricao,
		Valor:                  d.Valor,
		OrigemDespesaDescricao: d.OrigemDespesaDescricao,
	}
	if d.DataDespesa != nil {
		s := d.DataDespesa.Format("2006-01-02")
		dto.DataDespesa = &s
	}
	if d.CandidatoID != uuid.Nil {
		cand, err := repo.CandidatoBuscarPorID(ctx, d.CandidatoID)
		if err == nil && cand != nil {
			dto.Candidato = montarCandidatoResumido(ctx, repo, cand)
		}
	}
	return dto
}

func despesaPartidoParaDTO(ctx context.Context, repo *repositorio.Repositorio, d *types.DespesaOrgaoPartidario) tsetypes.DespesaRelacao {
	dto := tsetypes.DespesaRelacao{
		SQDespesa:              d.SQDespesa,
		Tipo:                   "partido",
		TipoRegistro:           d.TipoRegistro,
		Descricao:              d.Descricao,
		Valor:                  d.Valor,
		OrigemDespesaDescricao: d.OrigemDespesaDescricao,
	}
	if d.DataDespesa != nil {
		s := d.DataDespesa.Format("2006-01-02")
		dto.DataDespesa = &s
	}
	if d.PartidoID != uuid.Nil {
		part, err := repo.PartidosBuscarPorID(ctx, d.PartidoID)
		if err == nil && part != nil {
			dto.Partido = &tsetypes.PartidoResumido{
				Numero: part.Numero,
				Sigla:  part.Sigla,
				Nome:   part.Nome,
			}
		}
	}
	return dto
}

func receitaCandidatoParaDTO(ctx context.Context, repo *repositorio.Repositorio, r *types.ReceitaCandidato) tsetypes.ReceitaRelacao {
	dto := tsetypes.ReceitaRelacao{
		SQReceita:              r.SQReceita,
		Tipo:                   "candidato",
		Descricao:              r.Descricao,
		Valor:                  r.Valor,
		OrigemReceitaDescricao: r.OrigemReceitaDescricao,
	}
	if r.DataReceita != nil {
		s := r.DataReceita.Format("2006-01-02")
		dto.DataReceita = &s
	}
	if r.CandidatoID != uuid.Nil {
		cand, err := repo.CandidatoBuscarPorID(ctx, r.CandidatoID)
		if err == nil && cand != nil {
			dto.Candidato = montarCandidatoResumido(ctx, repo, cand)
		}
	}
	return dto
}

func receitaPartidoParaDTO(ctx context.Context, repo *repositorio.Repositorio, r *types.ReceitaOrgaoPartidario) tsetypes.ReceitaRelacao {
	dto := tsetypes.ReceitaRelacao{
		SQReceita:              r.SQReceita,
		Tipo:                   "partido",
		Descricao:              r.Descricao,
		Valor:                  r.Valor,
		OrigemReceitaDescricao: r.OrigemReceitaDescricao,
	}
	if r.DataReceita != nil {
		s := r.DataReceita.Format("2006-01-02")
		dto.DataReceita = &s
	}
	if r.PartidoID != uuid.Nil {
		part, err := repo.PartidosBuscarPorID(ctx, r.PartidoID)
		if err == nil && part != nil {
			dto.Partido = &tsetypes.PartidoResumido{
				Numero: part.Numero,
				Sigla:  part.Sigla,
				Nome:   part.Nome,
			}
		}
	}
	return dto
}

func montarCandidatoResumido(ctx context.Context, repo *repositorio.Repositorio, cand *types.Candidato) *tsetypes.CandidatoResumido {
	dto := &tsetypes.CandidatoResumido{
		SQCandidato:  cand.SQCandidato,
		NomeCompleto: cand.NomeCompleto,
		CargoNome:    cand.CargoNome,
		UFSigla:      cand.UFSigla,
	}
	if cand.PartidoID != nil && *cand.PartidoID != uuid.Nil {
		part, err := repo.PartidosBuscarPorID(ctx, *cand.PartidoID)
		if err == nil && part != nil {
			dto.PartidoSigla = part.Sigla
			dto.PartidoNome = part.Nome
		}
	}
	return dto
}

func init() {
	log := logger.New("TSE: UseCase: BuscarRelacoesUseCase")
	log.Info("BuscarRelacoesUseCase carregado")
}
