package services

import (
	"context"
	"fmt"
	"sync"

	"github.com/danyele/podp/internal/shared/clients/portaltransparencia"
	"github.com/danyele/podp/internal/shared/domain"
)

type VerificarPessoaPublicaInput struct {
	Documentos []string
}

type VerificarPessoaPublicaOutput struct {
	Resultados map[string][]domain.Vinculo
}

type VerificarPessoaPublicaService interface {
	Executar(ctx context.Context, input VerificarPessoaPublicaInput) (*VerificarPessoaPublicaOutput, error)
}

type verificarPessoaPublicaServiceImpl struct {
	portaltransparencia *portaltransparencia.PortalTransparenciaClient
}

func NovoVerificarPessoaPublicaService(pt *portaltransparencia.PortalTransparenciaClient) VerificarPessoaPublicaService {
	return &verificarPessoaPublicaServiceImpl{portaltransparencia: pt}
}

func (s *verificarPessoaPublicaServiceImpl) Executar(ctx context.Context, input VerificarPessoaPublicaInput) (*VerificarPessoaPublicaOutput, error) {
	if s.portaltransparencia == nil {
		return &VerificarPessoaPublicaOutput{Resultados: make(map[string][]domain.Vinculo)}, nil
	}

	resultados := make(map[string][]domain.Vinculo)
	var mu sync.Mutex

	const maxConcorrencia = 5
	sem := make(chan struct{}, maxConcorrencia)
	var wg sync.WaitGroup

	for _, doc := range input.Documentos {
		if len(doc) != 11 {
			continue
		}
		wg.Add(1)
		sem <- struct{}{}
		doc := doc
		go func() {
			defer wg.Done()
			defer func() { <-sem }()

			peps, err := s.portaltransparencia.ListarPEPs(ctx, portaltransparencia.PEPQueryParams{
				Pagina: 1,
				CPF:    doc,
			})
			if err != nil || len(peps) == 0 {
				return
			}

			vinculos := []domain.Vinculo{{
				Tipo:      "pessoa_publica",
				Descricao: fmt.Sprintf("%d registro(s) de pessoa exposta politicamente no Portal da Transparência", len(peps)),
				Detalhes: &domain.VinculoDetalhes{
					PessoasPublicas: peps,
				},
			}}

			mu.Lock()
			resultados[doc] = vinculos
			mu.Unlock()
		}()
	}
	wg.Wait()

	return &VerificarPessoaPublicaOutput{Resultados: resultados}, nil
}
