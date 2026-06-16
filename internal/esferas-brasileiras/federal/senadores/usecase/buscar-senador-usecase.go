package usecase

import (
	"context"

	senado "github.com/danyele/laceu/internal/shared/clients/senado"
)

type BuscarSenadorUseCase struct {
	client *senado.SenadoClient
}

func NovoBuscarSenadorUseCase(c *senado.SenadoClient) *BuscarSenadorUseCase {
	return &BuscarSenadorUseCase{client: c}
}

func (u *BuscarSenadorUseCase) Buscar(ctx context.Context, codigo string) (*senado.SenadorDetalhado, error) {
	senador, err := u.client.BuscarSenador(ctx, codigo)
	if err != nil {
		return nil, err
	}

	result := &senado.SenadorDetalhado{
		Senador:   senador,
		Cargos:    []senado.Cargo{},
		Comissoes: []senado.ComissaoMembro{},
		Mandatos:  []senado.MandatoDetalhe{},
	}

	if cargos, err := u.client.ListarCargos(ctx, codigo); err == nil {
		result.Cargos = cargos
	}

	if comissoes, err := u.client.ListarComissoes(ctx, codigo); err == nil {
		result.Comissoes = comissoes
	}

	if mandatos, err := u.client.ListarMandatos(ctx, codigo); err == nil {
		result.Mandatos = mandatos
	}

	return result, nil
}
