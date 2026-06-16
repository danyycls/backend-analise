package utils

import (
	"github.com/danyele/laceu/internal/shared/types"
)

func BuildFornecedorDTO(data *types.OpenCNPJResponse) *types.FornecedorOpenCNPJ {
	if data == nil {
		return nil
	}

	socios := make([]types.Socio, len(data.Socios))
	for i, s := range data.Socios {
		var pais *types.PaisInfo
		if s.Pais != nil {
			pais = &types.PaisInfo{
				Codigo:    s.Pais.Codigo,
				Descricao: s.Pais.Descricao,
				Extra:     s.Pais.Extra,
			}
		}
		var qualificacao *types.QualificacaoInfo
		if s.QualificacaoRepresentante != nil {
			qualificacao = &types.QualificacaoInfo{
				Codigo:    s.QualificacaoRepresentante.Codigo,
				Descricao: s.QualificacaoRepresentante.Descricao,
				Extra:     s.QualificacaoRepresentante.Extra,
			}
		}
		socios[i] = types.Socio{
			CNPJCPFSocio:              s.CNPJCPFSocio,
			CodigoPais:                s.CodigoPais,
			DataEntradaSociedade:      s.DataEntradaSociedade,
			FaixaEtaria:               s.FaixaEtaria,
			IdentificadorSocio:        s.IdentificadorSocio,
			NomeRepresentante:         s.NomeRepresentante,
			NomeSocio:                 s.NomeSocio,
			Pais:                      pais,
			QualificacaoRepresentante: qualificacao,
			QualificacaoSocio:         s.QualificacaoSocio,
			RepresentanteLegal:        s.RepresentanteLegal,
			Extra:                     s.Extra,
		}
	}

	return &types.FornecedorOpenCNPJ{
		CapitalSocial:     strPtr(data.CapitalSocial),
		CNPJ:              strPtr(data.CNPJ),
		NomeFantasia:      strPtr(data.NomeFantasia),
		Socios:            socios,
		RazaoSocial:       strPtr(data.RazaoSocial),
		SituacaoCadastral: strPtr(data.SituacaoCadastral),
	}
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
