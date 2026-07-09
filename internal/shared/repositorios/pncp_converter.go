package repositorios

import (
	"encoding/json"
	"time"

	"github.com/danyele/podp/internal/shared/clients/pncp"
	"github.com/danyele/podp/internal/shared/types"
)

func ContratoParaPersistido(c pncp.Contrato) ContratoPersistido {
	out := ContratoPersistido{
		NumeroControlePNCP: strVal(c.NumeroControlePNCP),
	}

	if c.CNPJOrgao != nil {
		out.CNPJOrgao = *c.CNPJOrgao
	} else if c.OrgaoEntidade != nil && c.OrgaoEntidade.CNPJ != nil {
		out.CNPJOrgao = *c.OrgaoEntidade.CNPJ
	}

	if c.UG != nil {
		out.UGUFSigla = strVal(c.UG.UFSigla)
		out.UGCodigoIbge = strVal(c.UG.CodigoIbge)
		out.CodigoUG = c.UG.CodigoUnidade
		out.NomeUG = c.UG.NomeUnidade
		out.UGMunicipioNome = c.UG.MunicipioNome
		out.UGUFNome = c.UG.UFNome
	} else if c.OrgaoVinculado != nil {
		out.UGUFSigla = strVal(c.OrgaoVinculado.UFSigla)
		out.UGCodigoIbge = strVal(c.OrgaoVinculado.CodigoIbge)
	}

	out.DataPublicacaoPncp = parseData(c.DataPublicacao)
	out.DataAssinatura = parseData(c.DataAssinatura)
	out.DataInicioVigencia = parseData(c.DataInicioVigencia)
	out.DataTerminoVigencia = parseData(c.DataTerminoVigencia)

	out.ValorGlobal = c.ValorGlobal
	out.ValorInicial = c.ValorInicial
	out.ValorTotalEstimado = c.ValorTotalEstimado
	out.ValorTotalHomologado = c.ValorTotalHomologado

	out.NIFornecedor = c.NIFornecedor
	out.NomeRazaoSocialFornecedor = c.NomeRazaoSocialFornecedor

	if c.AmparoLegal != nil && c.AmparoLegal.Codigo != nil {
		out.CodigoAmparoLegal = c.AmparoLegal.Codigo
	}

	out.NumeroContrato = c.NumeroContrato
	out.CodigoContrato = c.CodigoContrato
	out.CodigoTipoContrato = c.CodigoTipoContrato
	if c.TipoContrato != nil {
		out.TipoContratoNome = c.TipoContrato.Nome
	}
	out.ModalidadeNome = c.ModalidadeNome

	out.CodigoOrgao = c.CodigoOrgao
	out.NomeOrgao = c.NomeOrgao
	out.NomeOrgaoSub = c.NomeOrgaoSub
	out.ObjetoContrato = c.ObjetoCompra
	out.NumeroLicitacao = c.NumeroLicitação
	out.OrigemLicitacao = c.OrigemLicitação
	out.Produto = c.Produto
	out.SubtipoContrato = c.SubtipoContrato
	out.AnoContrato = c.AnoContrato

	raw, _ := json.Marshal(c)
	out.DadosCompletos = raw

	return out
}

func PersistidoParaContrato(cp ContratoPersistido) pncp.Contrato {
	var c pncp.Contrato
	if len(cp.DadosCompletos) > 0 {
		_ = json.Unmarshal(cp.DadosCompletos, &c)
	}
	if c.CNPJOrgao == nil || *c.CNPJOrgao == "" {
		c.CNPJOrgao = pncp.StrPtr(cp.CNPJOrgao)
	}
	if c.NumeroControlePNCP == nil || *c.NumeroControlePNCP == "" {
		c.NumeroControlePNCP = pncp.StrPtr(cp.NumeroControlePNCP)
	}
	if c.NIFornecedor == nil || *c.NIFornecedor == "" {
		c.NIFornecedor = cp.NIFornecedor
	}
	if c.NomeRazaoSocialFornecedor == nil || *c.NomeRazaoSocialFornecedor == "" {
		c.NomeRazaoSocialFornecedor = cp.NomeRazaoSocialFornecedor
	}
	return c
}

func FornecedorParaPersistido(f types.FornecedorOpenCNPJ) FornecedorPersistido {
	out := FornecedorPersistido{}
	if f.CNPJ != nil {
		out.CNPJ = *f.CNPJ
	}
	if f.RazaoSocial != nil {
		out.RazaoSocial = *f.RazaoSocial
	}
	raw, _ := json.Marshal(f)
	out.DadosCompletos = raw
	return out
}

func PersistidoParaFornecedor(fp FornecedorPersistido) *types.FornecedorOpenCNPJ {
	var f types.FornecedorOpenCNPJ
	if len(fp.DadosCompletos) > 0 {
		_ = json.Unmarshal(fp.DadosCompletos, &f)
	}
	if f.CNPJ == nil || *f.CNPJ == "" {
		f.CNPJ = pncp.StrPtr(fp.CNPJ)
	}
	if f.RazaoSocial == nil || *f.RazaoSocial == "" {
		f.RazaoSocial = pncp.StrPtr(fp.RazaoSocial)
	}
	return &f
}

func SocioParaPersistido(s types.Socio) SocioPersistido {
	return SocioPersistido{
		CNPJCPFSocio: strVal(s.CNPJCPFSocio),
		NomeSocio:    s.NomeSocio,
	}
}

func SocioParaFornecedorSocio(cnpjFornecedor, socioID string, s types.Socio) FornecedorSocioPersistido {
	return FornecedorSocioPersistido{
		CNPJFornecedor:            cnpjFornecedor,
		SocioID:                   socioID,
		DataEntradaSociedade:      s.DataEntradaSociedade,
		IdentificadorSocio:        s.IdentificadorSocio,
		NomeSocio:                 s.NomeSocio,
		QualificacaoSocio:         s.QualificacaoSocio,
		NomeRepresentante:         s.NomeRepresentante,
		QualificacaoRepresentante: qualificacaoDescricao(s.QualificacaoRepresentante),
		RepresentanteLegal:        s.RepresentanteLegal,
		FaixaEtaria:               s.FaixaEtaria,
		PaisCodigo:                paisCodigo(s.Pais),
		PaisDescricao:             paisDescricao(s.Pais),
	}
}

func ExtrairAmparoLegal(a *pncp.AmparoLegal) *AmparoLegalPersistido {
	if a == nil || a.Codigo == nil {
		return nil
	}
	return &AmparoLegalPersistido{
		Codigo:    *a.Codigo,
		Nome:      strVal(a.Nome),
		Descricao: a.Descricao,
	}
}

func strVal(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func parseData(s *string) *time.Time {
	if s == nil || *s == "" {
		return nil
	}
	v := *s
	formats := []string{
		"2006-01-02T15:04:05",
		"2006-01-02",
		"02/01/2006",
	}
	for _, f := range formats {
		t, err := time.Parse(f, v)
		if err == nil {
			return &t
		}
	}
	return nil
}

func paisCodigo(p *types.PaisInfo) *string {
	if p == nil {
		return nil
	}
	return p.Codigo
}

func paisDescricao(p *types.PaisInfo) *string {
	if p == nil {
		return nil
	}
	return p.Descricao
}

func qualificacaoDescricao(q *types.QualificacaoInfo) *string {
	if q == nil {
		return nil
	}
	return q.Descricao
}
