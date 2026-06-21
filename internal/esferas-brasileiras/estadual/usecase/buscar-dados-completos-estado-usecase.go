package usecase

import (
	"context"
	"sort"
	"strings"

	"github.com/google/uuid"

	estadual "github.com/danyele/podp/internal/esferas-brasileiras/estadual"
	repositorioTSE "github.com/danyele/podp/internal/esferas-brasileiras/tse/repositorio"
	deputados "github.com/danyele/podp/internal/shared/clients/deputados"
	"github.com/danyele/podp/internal/shared/clients/ibge"
	senadoClient "github.com/danyele/podp/internal/shared/clients/senado"
	"github.com/danyele/podp/internal/shared/database"
	"github.com/danyele/podp/internal/shared/logger"
	"github.com/danyele/podp/internal/shared/types"
)

type EsferaEstadualBuscarDadosCompletosEstadoRequest struct {
	UF string
}

type EsferaEstadualBuscarDadosCompletosEstadoResponse struct {
	Dados *types.DadosEstadoConsolidado
}

type EsferaEstadualBuscarDadosCompletosEstadoUseCase struct {
	db           database.DB
	ibgeClient   *ibge.IBGEClient
	deputadosCli *deputados.DeputadosClient
	senadoCli    *senadoClient.SenadoClient
}

func NovoEsferaEstadualBuscarDadosCompletosEstadoUseCase(
	db database.DB,
	ibge *ibge.IBGEClient,
	deputadosCli *deputados.DeputadosClient,
	senado *senadoClient.SenadoClient,
) *EsferaEstadualBuscarDadosCompletosEstadoUseCase {
	return &EsferaEstadualBuscarDadosCompletosEstadoUseCase{
		db:           db,
		ibgeClient:   ibge,
		deputadosCli: deputadosCli,
		senadoCli:    senado,
	}
}

func (u *EsferaEstadualBuscarDadosCompletosEstadoUseCase) Executar(
	ctx context.Context,
	req *EsferaEstadualBuscarDadosCompletosEstadoRequest,
) (*EsferaEstadualBuscarDadosCompletosEstadoResponse, error) {
	log := logger.New("Estadual: UseCase: BuscarDadosCompletos")

	nomeEstado := u.buscarNomeEstado(ctx, req.UF)

	result := &types.DadosEstadoConsolidado{
		UF:   req.UF,
		Nome: nomeEstado,
	}

	repo := repositorioTSE.Novo(u.db)

	u.buscarMunicipiosComPopulacao(ctx, req.UF, result)
	u.buscarCandidatosEleitos(ctx, repo, req.UF, result, log)
	u.buscarDeputadosAtivos(ctx, req.UF, result, log)
	u.buscarSenadoresPorUF(ctx, req.UF, result, log)

	return &EsferaEstadualBuscarDadosCompletosEstadoResponse{Dados: result}, nil
}

func (u *EsferaEstadualBuscarDadosCompletosEstadoUseCase) buscarNomeEstado(ctx context.Context, uf string) string {
	estados, err := u.ibgeClient.ListarEstados(ctx)
	if err != nil {
		return uf
	}
	for _, e := range estados {
		if e.Sigla == uf {
			return e.Nome
		}
	}
	return uf
}

func (u *EsferaEstadualBuscarDadosCompletosEstadoUseCase) buscarMunicipiosComPopulacao(
	ctx context.Context,
	uf string,
	result *types.DadosEstadoConsolidado,
) {
	log := logger.New("Estadual: UseCase: BuscarDadosCompletos")
	municipiosIBGE, err := u.ibgeClient.ListarMunicipios(ctx, uf)
	if err != nil {
		log.Error("erro ao buscar municipios IBGE", "uf", uf, "erro", err)
		return
	}

	ids := make([]int, len(municipiosIBGE))
	for i, m := range municipiosIBGE {
		ids[i] = m.ID
	}

	populacaoMap, _ := u.ibgeClient.BuscarPopulacao(ctx, ids)

	var totalPop int64
	for _, m := range municipiosIBGE {
		pop := populacaoMap[m.ID]
		totalPop += pop
		result.Municipios = append(result.Municipios, types.MunicipioComDados{
			ID:        m.ID,
			Nome:      m.Nome,
			Populacao: pop,
		})
	}
	result.Populacao = totalPop

	sort.Slice(result.Municipios, func(i, j int) bool {
		return result.Municipios[i].Populacao < result.Municipios[j].Populacao
	})
}

func (u *EsferaEstadualBuscarDadosCompletosEstadoUseCase) buscarCandidatosEleitos(
	ctx context.Context,
	repo *repositorioTSE.Repositorio,
	uf string,
	result *types.DadosEstadoConsolidado,
	log *logger.Logger,
) {
	candidatos, err := repo.CandidatosEleitosPorUF(ctx, uf, []string{"PREFEITO", "VICE-PREFEITO", "VEREADOR"})
	if err != nil {
		log.Error("erro ao buscar candidatos", "uf", uf, "erro", err)
		return
	}

	idsPartido, idsEleicao := u.coletarIDsDistintos(candidatos)
	partidoMap, _ := repo.PartidosBuscarPorIDs(ctx, idsPartido)
	eleicaoMap, _ := repo.EleicoesBuscarPorIDs(ctx, idsEleicao)

	for _, c := range candidatos {
		dto := u.montarCandidatoDTO(c, partidoMap, eleicaoMap)
		u.classificarCandidatoPorCargo(dto, result)
	}

	u.ordenarResultadosPorAno(result)
}

func (u *EsferaEstadualBuscarDadosCompletosEstadoUseCase) coletarIDsDistintos(
	candidatos []*types.Candidato,
) ([]uuid.UUID, []uuid.UUID) {
	idsPartido := make([]uuid.UUID, 0)
	idsEleicao := make([]uuid.UUID, 0)
	seenPartido := make(map[uuid.UUID]bool)
	seenEleicao := make(map[uuid.UUID]bool)

	for _, c := range candidatos {
		if c.PartidoID != nil && !seenPartido[*c.PartidoID] {
			idsPartido = append(idsPartido, *c.PartidoID)
			seenPartido[*c.PartidoID] = true
		}
		if c.ID != [16]byte{} && !seenEleicao[c.EleicaoID] {
			idsEleicao = append(idsEleicao, c.EleicaoID)
			seenEleicao[c.EleicaoID] = true
		}
	}

	return idsPartido, idsEleicao
}

func (u *EsferaEstadualBuscarDadosCompletosEstadoUseCase) montarCandidatoDTO(
	c *types.Candidato,
	partidoMap map[uuid.UUID]*types.Partido,
	eleicaoMap map[uuid.UUID]*types.Eleicao,
) types.CandidatoEleito {
	partidoSigla := ""
	partidoNome := ""
	if c.PartidoID != nil && partidoMap != nil {
		if p, ok := partidoMap[*c.PartidoID]; ok {
			partidoSigla = p.Sigla
			partidoNome = p.Nome
		}
	}

	ano := int16(0)
	eleicaoDesc := ""
	eleicaoData := ""
	eleicaoTipo := ""
	if eleicaoMap != nil {
		if e, ok := eleicaoMap[c.EleicaoID]; ok {
			ano = e.Ano
			eleicaoDesc = e.Descricao
			if e.DataEleicao != nil {
				eleicaoData = e.DataEleicao.Format("02/01/2006")
			}
			eleicaoTipo = estadual.NormalizarTipoEleicao(e.NomeTipoEleicao)
		}
	}

	return types.CandidatoEleito{
		ID:                           c.ID.String(),
		SQCandidato:                  c.SQCandidato,
		NomeUrna:                     c.NomeUrna,
		NomeCompleto:                 c.NomeCompleto,
		PartidoSigla:                 partidoSigla,
		PartidoNome:                  partidoNome,
		CargoNome:                    c.CargoNome,
		SituacaoTotalizacaoDescricao: c.SituacaoTotalizacaoDescricao,
		AnoEleicao:                   ano,
		NumeroCandidato:              c.NumeroCandidato,
		CPF:                          c.CPF,
		EleicaoDescricao:             eleicaoDesc,
		EleicaoData:                  eleicaoData,
		EleicaoTipo:                  eleicaoTipo,
	}
}

func (u *EsferaEstadualBuscarDadosCompletosEstadoUseCase) classificarCandidatoPorCargo(
	dto types.CandidatoEleito,
	result *types.DadosEstadoConsolidado,
) {
	switch dto.CargoNome {
	case "PREFEITO":
		result.Prefeitos = append(result.Prefeitos, dto)
	case "VICE-PREFEITO":
		result.VicePrefeitos = append(result.VicePrefeitos, dto)
	case "VEREADOR":
		result.Vereadores = append(result.Vereadores, dto)
	}
}

func (u *EsferaEstadualBuscarDadosCompletosEstadoUseCase) ordenarResultadosPorAno(result *types.DadosEstadoConsolidado) {
	sort.Slice(result.Prefeitos, func(i, j int) bool {
		return result.Prefeitos[i].AnoEleicao > result.Prefeitos[j].AnoEleicao
	})
	sort.Slice(result.VicePrefeitos, func(i, j int) bool {
		return result.VicePrefeitos[i].AnoEleicao > result.VicePrefeitos[j].AnoEleicao
	})
	sort.Slice(result.Vereadores, func(i, j int) bool {
		return result.Vereadores[i].AnoEleicao > result.Vereadores[j].AnoEleicao
	})
}

func (u *EsferaEstadualBuscarDadosCompletosEstadoUseCase) buscarDeputadosAtivos(
	ctx context.Context,
	uf string,
	result *types.DadosEstadoConsolidado,
	log *logger.Logger,
) {
	depParams := map[string]string{"siglaUf": uf}
	deputadosAPI, err := u.deputadosCli.ListarInfoDeputadosAtivos(ctx, depParams)
	if err != nil {
		log.Error("erro ao buscar deputados", "uf", uf, "erro", err)
		return
	}

	for _, d := range deputadosAPI {
		result.Deputados = append(result.Deputados, types.DeputadoUF{
			ID:            d.ID,
			Nome:          d.Nome,
			SiglaPartido:  d.SiglaPartido,
			SiglaUF:       d.SiglaUF,
			URLFoto:       d.URLFoto,
			Email:         d.Email,
			NomeEleitoral: "",
		})
	}
}

func (u *EsferaEstadualBuscarDadosCompletosEstadoUseCase) buscarSenadoresPorUF(
	ctx context.Context,
	uf string,
	result *types.DadosEstadoConsolidado,
	log *logger.Logger,
) {
	senadoresAPI, err := u.senadoCli.ListarSenadores(ctx)
	if err != nil {
		log.Error("erro ao buscar senadores", "erro", err)
		return
	}

	for _, s := range senadoresAPI {
		if s.IdentificacaoParlamentar.UfParlamentar == uf {
			fotoURL := s.IdentificacaoParlamentar.UrlFotoParlamentar
			fotoURL = strings.Replace(fotoURL, "http://", "https://", 1)
			result.Senadores = append(result.Senadores, types.SenadorUF{
				Codigo:          s.IdentificacaoParlamentar.CodigoParlamentar,
				NomeParlamentar: s.IdentificacaoParlamentar.NomeParlamentar,
				NomeCompleto:    s.IdentificacaoParlamentar.NomeCompletoParlamentar,
				Uf:              s.IdentificacaoParlamentar.UfParlamentar,
				Partido:         s.IdentificacaoParlamentar.SiglaPartidoParlamentar,
				URLFoto:         fotoURL,
			})
		}
	}
}
