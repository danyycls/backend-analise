package usecase

import (
	"context"
	"strconv"

	"github.com/google/uuid"

	repositorio "github.com/danyele/podp/internal/esferas-brasileiras/tse/repositorio"
	tsetypes "github.com/danyele/podp/internal/esferas-brasileiras/tse/types"
	"github.com/danyele/podp/internal/shared/database"
	"github.com/danyele/podp/internal/shared/types"
)

type ConsultarEntidadeUseCase struct {
	db database.DB
}

func NovoConsultarEntidadeUseCase(db database.DB) *ConsultarEntidadeUseCase {
	return &ConsultarEntidadeUseCase{db: db}
}

func (u *ConsultarEntidadeUseCase) Executar(ctx context.Context, req *tsetypes.ConsultaEntidadeRequest) *tsetypes.ConsultaEntidadeResponse {
	repo := repositorio.Novo(u.db)

	switch req.Tipo {
	case "candidato":
		return u.buscarCandidato(ctx, repo, req.Chave)
	case "fornecedor":
		return u.buscarFornecedor(ctx, repo, req.Chave)
	case "doador":
		return u.buscarDoador(ctx, repo, req.Chave)
	case "receita":
		return u.buscarReceita(ctx, repo, req.Chave)
	case "despesa":
		return u.buscarDespesa(ctx, repo, req.Chave)
	default:
		return &tsetypes.ConsultaEntidadeResponse{
			Tipo:  req.Tipo,
			Chave: req.Chave,
			Erro:  "tipo desconhecido: " + req.Tipo,
		}
	}
}

func (u *ConsultarEntidadeUseCase) buscarCandidato(ctx context.Context, repo *repositorio.Repositorio, chave string) *tsetypes.ConsultaEntidadeResponse {
	if sq, err := strconv.ParseInt(chave, 10, 64); err == nil {
		cand, err := repo.CandidatoBuscarPorSQCandidato(ctx, sq)
		if err == nil && cand != nil {
			return &tsetypes.ConsultaEntidadeResponse{Tipo: "candidato", Chave: chave, Dados: u.candidatoParaEntidadeDTO(ctx, repo, cand)}
		}
	}
	cands, err := repo.CandidatosBuscarPorCPF(ctx, chave)
	if err == nil && len(cands) > 0 {
		return &tsetypes.ConsultaEntidadeResponse{Tipo: "candidato", Chave: chave, Dados: u.candidatoParaEntidadeDTO(ctx, repo, cands[0])}
	}
	return &tsetypes.ConsultaEntidadeResponse{Tipo: "candidato", Chave: chave, Erro: "candidato nao encontrado"}
}

func (u *ConsultarEntidadeUseCase) candidatoParaEntidadeDTO(ctx context.Context, repo *repositorio.Repositorio, cand *types.Candidato) *tsetypes.CandidatoEntidade {
	dto := &tsetypes.CandidatoEntidade{
		SQCandidato:  cand.SQCandidato,
		NomeCompleto: cand.NomeCompleto,
		CPF:          cand.CPF,
		CargoNome:    cand.CargoNome,
		UFSigla:      cand.UFSigla,
		Genero:       cand.GeneroDescricao,
		CorRaca:      cand.CorRacaDescricao,
		OcupacaoNome: cand.OcupacaoNome,
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

func (u *ConsultarEntidadeUseCase) buscarFornecedor(ctx context.Context, repo *repositorio.Repositorio, chave string) *tsetypes.ConsultaEntidadeResponse {
	forn, err := repo.FornecedorBuscarPorCNPJExato(ctx, chave)
	if err == nil && forn != nil {
		return &tsetypes.ConsultaEntidadeResponse{Tipo: "fornecedor", Chave: chave, Dados: u.fornecedorParaEntidadeDTO(forn)}
	}
	return &tsetypes.ConsultaEntidadeResponse{Tipo: "fornecedor", Chave: chave, Erro: "fornecedor nao encontrado"}
}

func (u *ConsultarEntidadeUseCase) fornecedorParaEntidadeDTO(forn *types.Fornecedor) *tsetypes.FornecedorEntidade {
	dto := &tsetypes.FornecedorEntidade{
		CPFCNPJ: forn.CPFCNPJ,
		Nome:    forn.Nome,
		NomeRFB: forn.NomeRFB,
		CNAE:    forn.CNAEDescricao,
	}
	if forn.UFSigla != nil {
		dto.UFSigla = *forn.UFSigla
	}
	return dto
}

func (u *ConsultarEntidadeUseCase) buscarDoador(ctx context.Context, repo *repositorio.Repositorio, chave string) *tsetypes.ConsultaEntidadeResponse {
	doador, err := repo.DoadorBuscarPorCNPJExato(ctx, chave)
	if err == nil && doador != nil {
		return &tsetypes.ConsultaEntidadeResponse{Tipo: "doador", Chave: chave, Dados: u.doadorParaEntidadeDTO(doador)}
	}
	return &tsetypes.ConsultaEntidadeResponse{Tipo: "doador", Chave: chave, Erro: "doador nao encontrado"}
}

func (u *ConsultarEntidadeUseCase) doadorParaEntidadeDTO(doador *types.Doador) *tsetypes.DoadorEntidade {
	dto := &tsetypes.DoadorEntidade{
		CPFCNPJ: doador.CPFCNPJ,
		Nome:    doador.Nome,
		NomeRFB: doador.NomeRFB,
		CNAE:    doador.CNAEDescricao,
	}
	if doador.UFSigla != nil {
		dto.UFSigla = *doador.UFSigla
	}
	return dto
}

func (u *ConsultarEntidadeUseCase) buscarReceita(ctx context.Context, repo *repositorio.Repositorio, chave string) *tsetypes.ConsultaEntidadeResponse {
	sq, err := strconv.ParseInt(chave, 10, 64)
	if err != nil {
		return &tsetypes.ConsultaEntidadeResponse{Tipo: "receita", Chave: chave, Erro: "chave invalida, esperado numero"}
	}

	rc, err := repo.ReceitaCandidatoBuscarPorSQ(ctx, sq)
	if err == nil && rc != nil {
		return &tsetypes.ConsultaEntidadeResponse{Tipo: "receita", Chave: chave, Dados: u.receitaCandidatoParaEntidadeDTO(ctx, repo, rc)}
	}

	ro, err := repo.ReceitaOrgaoBuscarPorSQ(ctx, sq)
	if err == nil && ro != nil {
		return &tsetypes.ConsultaEntidadeResponse{Tipo: "receita", Chave: chave, Dados: u.receitaOrgaoParaEntidadeDTO(ctx, repo, ro)}
	}

	return &tsetypes.ConsultaEntidadeResponse{Tipo: "receita", Chave: chave, Erro: "receita nao encontrada"}
}

func (u *ConsultarEntidadeUseCase) receitaCandidatoParaEntidadeDTO(ctx context.Context, repo *repositorio.Repositorio, r *types.ReceitaCandidato) *tsetypes.ReceitaEntidade {
	dto := &tsetypes.ReceitaEntidade{
		SQReceita: r.SQReceita,
		Tipo:      "candidato",
		Descricao: r.Descricao,
		Valor:     r.Valor,
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
	if r.DoadorID != nil {
		doador, err := repo.DoadorBuscarPorID(ctx, *r.DoadorID)
		if err == nil && doador != nil {
			dto.DoadorNome = doador.Nome
		}
	}
	return dto
}

func (u *ConsultarEntidadeUseCase) receitaOrgaoParaEntidadeDTO(ctx context.Context, repo *repositorio.Repositorio, r *types.ReceitaOrgaoPartidario) *tsetypes.ReceitaEntidade {
	dto := &tsetypes.ReceitaEntidade{
		SQReceita: r.SQReceita,
		Tipo:      "orgao_partidario",
		Descricao: r.Descricao,
		Valor:     r.Valor,
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
	if r.DoadorID != nil {
		doador, err := repo.DoadorBuscarPorID(ctx, *r.DoadorID)
		if err == nil && doador != nil {
			dto.DoadorNome = doador.Nome
		}
	}
	return dto
}

func (u *ConsultarEntidadeUseCase) buscarDespesa(ctx context.Context, repo *repositorio.Repositorio, chave string) *tsetypes.ConsultaEntidadeResponse {
	sq, err := strconv.ParseInt(chave, 10, 64)
	if err != nil {
		return &tsetypes.ConsultaEntidadeResponse{Tipo: "despesa", Chave: chave, Erro: "chave invalida, esperado numero"}
	}

	dc, err := repo.DespesaCandidatoBuscarPorSQ(ctx, sq)
	if err == nil && dc != nil {
		return &tsetypes.ConsultaEntidadeResponse{Tipo: "despesa", Chave: chave, Dados: u.despesaCandidatoParaEntidadeDTO(ctx, repo, dc)}
	}

	dop, err := repo.DespesaOrgaoBuscarPorSQ(ctx, sq)
	if err == nil && dop != nil {
		return &tsetypes.ConsultaEntidadeResponse{Tipo: "despesa", Chave: chave, Dados: u.despesaOrgaoParaEntidadeDTO(ctx, repo, dop)}
	}

	return &tsetypes.ConsultaEntidadeResponse{Tipo: "despesa", Chave: chave, Erro: "despesa nao encontrada"}
}

func (u *ConsultarEntidadeUseCase) despesaCandidatoParaEntidadeDTO(ctx context.Context, repo *repositorio.Repositorio, d *types.DespesaCandidato) *tsetypes.DespesaEntidade {
	dto := &tsetypes.DespesaEntidade{
		SQDespesa: d.SQDespesa,
		Tipo:      "candidato",
		Descricao: d.Descricao,
		Valor:     d.Valor,
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
	if d.FornecedorID != nil {
		forn, err := repo.FornecedorBuscarPorID(ctx, *d.FornecedorID)
		if err == nil && forn != nil {
			dto.FornecedorNome = forn.Nome
		}
	}
	return dto
}

func (u *ConsultarEntidadeUseCase) despesaOrgaoParaEntidadeDTO(ctx context.Context, repo *repositorio.Repositorio, d *types.DespesaOrgaoPartidario) *tsetypes.DespesaEntidade {
	dto := &tsetypes.DespesaEntidade{
		SQDespesa: d.SQDespesa,
		Tipo:      "orgao_partidario",
		Descricao: d.Descricao,
		Valor:     d.Valor,
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
	if d.FornecedorID != nil {
		forn, err := repo.FornecedorBuscarPorID(ctx, *d.FornecedorID)
		if err == nil && forn != nil {
			dto.FornecedorNome = forn.Nome
		}
	}
	return dto
}
