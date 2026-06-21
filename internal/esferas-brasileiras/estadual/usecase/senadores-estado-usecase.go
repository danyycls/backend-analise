package usecase

import (
	"context"
	"fmt"
	"strings"

	senadoClient "github.com/danyele/podp/internal/shared/clients/senado"
	"github.com/danyele/podp/internal/shared/types"
)

type EsferaEstadualBuscarSenadoresRequest struct {
	UF string
}

type EsferaEstadualBuscarSenadoresResponse struct {
	Senadores []types.SenadorUF
}

type EsferaEstadualBuscarSenadoresUseCase struct {
	senadoCli *senadoClient.SenadoClient
}

func NovoEsferaEstadualBuscarSenadoresUseCase(senado *senadoClient.SenadoClient) *EsferaEstadualBuscarSenadoresUseCase {
	return &EsferaEstadualBuscarSenadoresUseCase{
		senadoCli: senado,
	}
}

func (u *EsferaEstadualBuscarSenadoresUseCase) Executar(ctx context.Context, req *EsferaEstadualBuscarSenadoresRequest) (*EsferaEstadualBuscarSenadoresResponse, error) {
	senadoresAPI, err := u.senadoCli.ListarSenadores(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro buscar senadores: %w", err)
	}

	var result []types.SenadorUF
	for _, s := range senadoresAPI {
		if s.IdentificacaoParlamentar.UfParlamentar == req.UF {
			fotoURL := s.IdentificacaoParlamentar.UrlFotoParlamentar
			fotoURL = strings.Replace(fotoURL, "http://", "https://", 1)
			result = append(result, types.SenadorUF{
				Codigo:          s.IdentificacaoParlamentar.CodigoParlamentar,
				NomeParlamentar: s.IdentificacaoParlamentar.NomeParlamentar,
				NomeCompleto:    s.IdentificacaoParlamentar.NomeCompletoParlamentar,
				Uf:              s.IdentificacaoParlamentar.UfParlamentar,
				Partido:         s.IdentificacaoParlamentar.SiglaPartidoParlamentar,
				URLFoto:         fotoURL,
			})
		}
	}
	return &EsferaEstadualBuscarSenadoresResponse{Senadores: result}, nil
}
