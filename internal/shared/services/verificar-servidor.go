package services

import (
	"context"
	"fmt"
	"sync"

	"github.com/danyele/podp/internal/shared/clients/portaltransparencia"
	"github.com/danyele/podp/internal/shared/domain"
)

type VerificarServidorInput struct {
	Documentos []string
}

type VerificarServidorOutput struct {
	Resultados map[string][]domain.Vinculo
}

type VerificarServidorPublicoService interface {
	Executar(ctx context.Context, input VerificarServidorInput) (*VerificarServidorOutput, error)
}

type verificarServidorPublicoServiceImpl struct {
	portaltransparencia *portaltransparencia.PortalTransparenciaClient
}

func NovoVerificarServidorPublicoService(pt *portaltransparencia.PortalTransparenciaClient) VerificarServidorPublicoService {
	return &verificarServidorPublicoServiceImpl{portaltransparencia: pt}
}

func (s *verificarServidorPublicoServiceImpl) Executar(ctx context.Context, input VerificarServidorInput) (*VerificarServidorOutput, error) {
	if s.portaltransparencia == nil {
		return &VerificarServidorOutput{Resultados: make(map[string][]domain.Vinculo)}, nil
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

			servidores, err := s.portaltransparencia.ListarServidores(ctx, portaltransparencia.ServidorQueryParams{
				Pagina: 1,
				CPF:    doc,
			})
			if err != nil || len(servidores) == 0 {
				return
			}

			vinculos := []domain.Vinculo{{
				Tipo:      "servidor_publico",
				Descricao: fmt.Sprintf("%d registro(s) de servidor público no Portal da Transparência", len(servidores)),
				Detalhes: &domain.VinculoDetalhes{
					ServidoresPublicos: servidores,
				},
			}}

			mu.Lock()
			resultados[doc] = vinculos
			mu.Unlock()
		}()
	}
	wg.Wait()

	return &VerificarServidorOutput{Resultados: resultados}, nil
}
