package portaltransparencia

import (
	"context"
)

func (c *PortalTransparenciaClient) ListarPessoasFisicas(ctx context.Context, filtro PessoaFisicaQueryParams) (*PessoaFisica, error) {
	params := map[string]string{}
	if filtro.CPF != "" {
		params["cpf"] = filtro.CPF
	}
	if filtro.Nome != "" {
		params["nome"] = filtro.Nome
	}
	if filtro.NIS != "" {
		params["nis"] = filtro.NIS
	}
	if filtro.FavorecidoDespesas != "" {
		params["favorecidoDespesas"] = filtro.FavorecidoDespesas
	}
	if filtro.Servidor != "" {
		params["servidor"] = filtro.Servidor
	}
	if filtro.BeneficiarioDiarias != "" {
		params["beneficiarioDiarias"] = filtro.BeneficiarioDiarias
	}
	if filtro.Permissionario != "" {
		params["permissionario"] = filtro.Permissionario
	}
	if filtro.Contratado != "" {
		params["contratado"] = filtro.Contratado
	}
	if filtro.SancionadoCEIS != "" {
		params["sancionadoCEIS"] = filtro.SancionadoCEIS
	}
	if filtro.SancionadoCNEP != "" {
		params["sancionadoCNEP"] = filtro.SancionadoCNEP
	}
	if filtro.SancionadoCEPIM != "" {
		params["sancionadoCEPIM"] = filtro.SancionadoCEPIM
	}
	if filtro.SancionadoCEAF != "" {
		params["sancionadoCEAF"] = filtro.SancionadoCEAF
	}
	if filtro.SancionadoAcordoLeniencia != "" {
		params["sancionadoAcordoLeniencia"] = filtro.SancionadoAcordoLeniencia
	}
	if filtro.Ordenacao != "" {
		params["ordenacao"] = filtro.Ordenacao
	}
	if filtro.OrdenacaoDirecao != "" {
		params["ordenacaoDirecao"] = filtro.OrdenacaoDirecao
	}
	var result PessoaFisica
	err := c.doGet(ctx, "/api-de-dados/pessoa-fisica", params, &result)
	return &result, err
}

func (c *PortalTransparenciaClient) ListarPessoasJuridicas(ctx context.Context, filtro PessoaJuridicaQueryParams) (*PessoaJuridica, error) {
	params := map[string]string{}
	if filtro.CNPJ != "" {
		params["cnpj"] = filtro.CNPJ
	}
	if filtro.RazaoSocial != "" {
		params["razaoSocial"] = filtro.RazaoSocial
	}
	if filtro.NomeFantasia != "" {
		params["nomeFantasia"] = filtro.NomeFantasia
	}
	if filtro.FavorecidoDespesas != "" {
		params["favorecidoDespesas"] = filtro.FavorecidoDespesas
	}
	if filtro.PossuiContratacao != "" {
		params["possuiContratacao"] = filtro.PossuiContratacao
	}
	if filtro.Convenios != "" {
		params["convenios"] = filtro.Convenios
	}
	if filtro.FavorecidoTransferencias != "" {
		params["favorecidoTransferencias"] = filtro.FavorecidoTransferencias
	}
	if filtro.SancionadoCEPIM != "" {
		params["sancionadoCEPIM"] = filtro.SancionadoCEPIM
	}
	if filtro.SancionadoCEIS != "" {
		params["sancionadoCEIS"] = filtro.SancionadoCEIS
	}
	if filtro.SancionadoCNEP != "" {
		params["sancionadoCNEP"] = filtro.SancionadoCNEP
	}
	if filtro.SancionadoCEAF != "" {
		params["sancionadoCEAF"] = filtro.SancionadoCEAF
	}
	if filtro.SancionadoAcordoLeniencia != "" {
		params["sancionadoAcordoLeniencia"] = filtro.SancionadoAcordoLeniencia
	}
	var result PessoaJuridica
	err := c.doGet(ctx, "/api-de-dados/pessoa-juridica", params, &result)
	return &result, err
}
