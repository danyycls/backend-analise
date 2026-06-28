package services

import (
	"context"
	"fmt"
	"sync"

	tcu "github.com/danyele/podp/internal/shared/clients/tcu"
	"github.com/danyele/podp/internal/shared/domain"
)

type VerificarSancoesTCUInput struct {
	Documentos []string
}

type VerificarSancoesTCUOutput struct {
	Resultados map[string][]domain.Vinculo
}

type VerificarSancoesTCUService interface {
	Executar(ctx context.Context, input VerificarSancoesTCUInput) (*VerificarSancoesTCUOutput, error)
}

type verificarSancoesTCUServiceImpl struct {
	tcu tcu.Client
}

func NovoVerificarSancoesTCUService(tcu tcu.Client) VerificarSancoesTCUService {
	return &verificarSancoesTCUServiceImpl{tcu: tcu}
}

func (s *verificarSancoesTCUServiceImpl) Executar(ctx context.Context, input VerificarSancoesTCUInput) (*VerificarSancoesTCUOutput, error) {
	if len(input.Documentos) == 0 {
		return &VerificarSancoesTCUOutput{Resultados: make(map[string][]domain.Vinculo)}, nil
	}

	resultados := make(map[string][]domain.Vinculo, len(input.Documentos))
	var mu sync.Mutex

	const maxConcorrencia = 5
	sem := make(chan struct{}, maxConcorrencia)
	var wg sync.WaitGroup

	for _, doc := range input.Documentos {
		wg.Add(1)
		sem <- struct{}{}
		doc := doc
		go func() {
			defer wg.Done()
			defer func() { <-sem }()

			isCPF := len(doc) == 11
			filterCPF := tcu.TCUQueryParams{CPF: doc}
			filterCNPJ := tcu.TCUQueryParams{CNPJ: doc}

			var vinculos []domain.Vinculo

			if resp, err := s.tcu.BuscarContasIrregulares(ctx, filterCPF); err == nil && len(resp) > 0 {
				vinculos = append(vinculos, domain.Vinculo{
					Tipo:      "tcu_contas_irregulares",
					Descricao: fmt.Sprintf("%d registro(s) de contas julgadas irregulares no TCU", len(resp)),
					Detalhes: &domain.VinculoDetalhes{
						ContasIrregulares: resp,
					},
				})
			}

			if respCNPJ, err := s.tcu.BuscarContasIrregulares(ctx, filterCNPJ); err == nil && len(respCNPJ) > 0 {
				vinculos = append(vinculos, domain.Vinculo{
					Tipo:      "tcu_contas_irregulares",
					Descricao: fmt.Sprintf("%d registro(s) de contas julgadas irregulares no TCU", len(respCNPJ)),
					Detalhes: &domain.VinculoDetalhes{
						ContasIrregulares: respCNPJ,
					},
				})
			}

			if isCPF {
				if resp, err := s.tcu.BuscarInabilitados(ctx, filterCPF); err == nil && len(resp) > 0 {
					vinculos = append(vinculos, domain.Vinculo{
						Tipo:      "tcu_inabilitado",
						Descricao: fmt.Sprintf("%d registro(s) de inabilitado no TCU", len(resp)),
						Detalhes: &domain.VinculoDetalhes{
							Inabilitados: resp,
						},
					})
				}
			}

			if resp, err := s.tcu.BuscarInidoneos(ctx, filterCPF); err == nil && len(resp) > 0 {
				vinculos = append(vinculos, domain.Vinculo{
					Tipo:      "tcu_inidoneo",
					Descricao: fmt.Sprintf("%d registro(s) de inidôneo no TCU", len(resp)),
					Detalhes: &domain.VinculoDetalhes{
						Inidoneos: resp,
					},
				})
			}

			if !isCPF {
				if respCNPJ, err := s.tcu.BuscarInidoneos(ctx, filterCNPJ); err == nil && len(respCNPJ) > 0 {
					vinculos = append(vinculos, domain.Vinculo{
						Tipo:      "tcu_inidoneo",
						Descricao: fmt.Sprintf("%d registro(s) de inidôneo no TCU", len(respCNPJ)),
						Detalhes: &domain.VinculoDetalhes{
							Inidoneos: respCNPJ,
						},
					})
				}
			}

			if len(vinculos) > 0 {
				mu.Lock()
				resultados[doc] = vinculos
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	return &VerificarSancoesTCUOutput{Resultados: resultados}, nil
}
