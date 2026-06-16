package repositorio

import (
	"context"

	"github.com/google/uuid"

	tsetypes "github.com/danyele/laceu/internal/esferas-brasileiras/tse/types"
	"github.com/danyele/laceu/internal/shared/types"
)

// ---------------------------------------------------------------------------
// Cache de busca para evitar queries repetidas no mesmo documento
// ---------------------------------------------------------------------------

// CacheBusca armazena resultados ja montados para evitar consultas repetidas ao banco
type CacheBusca struct {
	repo         *Repositorio
	candidatos   map[uuid.UUID]*tsetypes.CandidatoDetalhado
	candidatosSQ map[int64]*tsetypes.CandidatoDetalhado
	prestacoes   map[uuid.UUID]*tsetypes.PrestacaoDetalhada
	fornecedores map[uuid.UUID]*types.Fornecedor
}

// NovoCacheBusca cria um cache vazio para ser usado durante o processamento de um lote
func NovoCacheBusca(repo *Repositorio) *CacheBusca {
	return &CacheBusca{
		repo:         repo,
		candidatos:   make(map[uuid.UUID]*tsetypes.CandidatoDetalhado),
		candidatosSQ: make(map[int64]*tsetypes.CandidatoDetalhado),
		prestacoes:   make(map[uuid.UUID]*tsetypes.PrestacaoDetalhada),
		fornecedores: make(map[uuid.UUID]*types.Fornecedor),
	}
}

// Candidato busca ou monta o candidatoDTO, cacheando por ID e SQCandidato
func (c *CacheBusca) Candidato(ctx context.Context, m *types.Candidato) *tsetypes.CandidatoDetalhado {
	if existing, ok := c.candidatos[m.ID]; ok {
		return existing
	}
	dto := MontarCandidatoDTO(ctx, c.repo, m)
	c.candidatos[m.ID] = dto
	c.candidatosSQ[m.SQCandidato] = dto
	return dto
}

// Prestacao busca ou monta a prestacaoDTO, cacheando por ID
func (c *CacheBusca) Prestacao(ctx context.Context, p *types.PrestacaoContas) *tsetypes.PrestacaoDetalhada {
	if existing, ok := c.prestacoes[p.ID]; ok {
		return existing
	}
	dto := MontarPrestacaoDTO(ctx, c.repo, p)
	despCand, _ := c.repo.DespesasCandidatoBuscarPorPrestacaoID(ctx, p.ID)
	despPart, _ := c.repo.DespesasPartidoBuscarPorPrestacaoID(ctx, p.ID)
	dto.DespesasCandidato = despCand
	dto.DespesasOrgaoPartidario = despPart
	c.prestacoes[p.ID] = dto
	return dto
}

// tsetypes.MontarCandidatoDTO monta o DTO completo de candidato com todos os relacionamentos
func MontarCandidatoDTO(ctx context.Context, repo *Repositorio, c *types.Candidato) *tsetypes.CandidatoDetalhado {
	dto := &tsetypes.CandidatoDetalhado{Candidato: c}
	if eleicao, err := repo.EleicoesBuscarPorID(ctx, c.EleicaoID); err == nil {
		dto.Eleicao = eleicao
	}
	if c.PartidoID != nil {
		if p, err := repo.PartidosBuscarPorID(ctx, *c.PartidoID); err == nil {
			dto.Partido = p
		}
	}
	return dto
}

// MontarPrestacaoDTO monta o tsetypes.PrestacaoDetalhada com eleicao, unidade eleitoral e partido
func MontarPrestacaoDTO(ctx context.Context, repo *Repositorio, p *types.PrestacaoContas) *tsetypes.PrestacaoDetalhada {
	dto := &tsetypes.PrestacaoDetalhada{PrestacaoContas: p}
	if eleicao, err := repo.EleicoesBuscarPorID(ctx, p.EleicaoID); err == nil {
		dto.Eleicao = eleicao
	}
	if p.UnidadeEleitoralID != nil {
		if ue, err := repo.UnidadesEleitoraisBuscarPorID(ctx, *p.UnidadeEleitoralID); err == nil {
			dto.UnidadeEleitoral = ue
		}
	}
	if p.PartidoID != nil {
		if part, err := repo.PartidosBuscarPorID(ctx, *p.PartidoID); err == nil {
			dto.Partido = part
		}
	}
	return dto
}
