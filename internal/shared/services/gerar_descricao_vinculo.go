package services

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/danyele/podp/internal/shared/domain"
)

type GerarDescricaoVinculoInput struct {
	VinculoLicitacao domain.VinculoLicitacao
}

type GerarDescricaoVinculoOutput struct {
	Titulo string   `json:"titulo"`
	Tags   []string `json:"tags"`
}

type GerarDescricaoVinculoService interface {
	Executar(ctx context.Context, input GerarDescricaoVinculoInput) (*GerarDescricaoVinculoOutput, error)
}

type gerarDescricaoVinculoServiceImpl struct{}

func NovoGerarDescricaoVinculoService() GerarDescricaoVinculoService {
	return &gerarDescricaoVinculoServiceImpl{}
}

var tipoParaInfo = map[string]struct {
	Tag       string
	Fragmento string
}{
	"fornecedor":               {Tag: "Fornecedor-TSE", Fragmento: "fornecedor em campanha política"},
	"doador":                   {Tag: "Doador-TSE", Fragmento: "doador eleitoral"},
	"receita_candidato":        {Tag: "Doação Candidato-TSE", Fragmento: "doação a candidato"},
	"receita_orgao_partidario": {Tag: "Doação Partido-TSE", Fragmento: "doação a partido"},
	"tcu_contas_irregulares":   {Tag: "Contas Irregulares-TCU", Fragmento: "contas irregulares no TCU"},
	"tcu_inabilitado":          {Tag: "Inabilitado-TCU", Fragmento: "inabilitado pelo TCU"},
	"tcu_inidoneo":             {Tag: "Inidôneo-TCU", Fragmento: "inidôneo pelo TCU"},
	"servidor_publico":         {Tag: "Servidor Público-Portal Transparência", Fragmento: "servidor público federal"},
	"pessoa_publica":           {Tag: "Pessoa Exposta-Portal Transparência", Fragmento: "pessoa politicamente exposta"},
	"dispensa_valor_limite":    {Tag: "Dispensa acima do limite-Regra", Fragmento: "dispensa de licitação acima do limite legal de baixo valor"},
}

func (s *gerarDescricaoVinculoServiceImpl) Executar(
	ctx context.Context,
	input GerarDescricaoVinculoInput,
) (*GerarDescricaoVinculoOutput, error) {
	vl := input.VinculoLicitacao

	if len(vl.DocumentosVinculos) == 0 {
		return &GerarDescricaoVinculoOutput{}, nil
	}

	// Fixed prefix: "Empresa vencedora da licitação X no valor de R$ X, "
	empresa := vl.NomeEmpresa
	if empresa == "" {
		for _, doc := range vl.DocumentosVinculos {
			if doc.Origem == "principal" && doc.Nome != "" {
				empresa = doc.Nome
				break
			}
		}
	}

	var b strings.Builder
	b.WriteString(empresa)
	b.WriteString(" vencedora da licitação ")
	b.WriteString(vl.NumeroControlePncp)
	if vl.ValorGlobal > 0 {
		b.WriteString(fmt.Sprintf(" no valor de R$ %.2f", vl.ValorGlobal))
	}
	b.WriteString(", ")

	var fragmentos []string
	tagSet := make(map[string]struct{})

	for _, doc := range vl.DocumentosVinculos {
		for _, v := range doc.Vinculos {
			info, ok := tipoParaInfo[v.Tipo]
			if !ok {
				continue
			}
			tagSet[info.Tag] = struct{}{}

			frag := montarFragmento(v, doc)
			if frag != "" {
				fragmentos = append(fragmentos, frag)
			}
		}
	}

	if len(fragmentos) == 0 {
		return &GerarDescricaoVinculoOutput{}, nil
	}

	b.WriteString(strings.Join(fragmentos, ". "))

	tags := make([]string, 0, len(tagSet))
	for t := range tagSet {
		tags = append(tags, t)
	}
	sort.Strings(tags)

	return &GerarDescricaoVinculoOutput{
		Titulo: b.String(),
		Tags:   tags,
	}, nil
}

func montarFragmento(v domain.Vinculo, doc domain.DocumentoVinculo) string {
	d := v.Detalhes
	if d == nil {
		return fragmentoGenerico(v.Tipo, doc)
	}

	switch v.Tipo {
	case "doador":
		p := prefixoSocio(doc)
		if len(d.ReceitasCandidato) > 0 || len(d.ReceitasOrgaoPartidario) > 0 {
			var total float64
			for _, rc := range d.ReceitasCandidato {
				if rc != nil && rc.Receita != nil {
					total += rc.Receita.Valor
				}
			}
			for _, ro := range d.ReceitasOrgaoPartidario {
				if ro != nil && ro.Receita != nil {
					total += ro.Receita.Valor
				}
			}
			return fmt.Sprintf("%sfez doação(ões) (R$ %.2f)", p, total)
		}
		return fragmentoGenerico(v.Tipo, doc)
	case "receita_candidato":
		if len(d.ReceitasCandidato) > 0 && d.ReceitasCandidato[0] != nil && d.ReceitasCandidato[0].Receita != nil {
			rc := d.ReceitasCandidato[0].Receita
			p := prefixoSocio(doc)
			return fmt.Sprintf("%sfez doação a candidato (R$ %.2f)", p, rc.Valor)
		}
	case "receita_orgao_partidario":
		if len(d.ReceitasOrgaoPartidario) > 0 && d.ReceitasOrgaoPartidario[0] != nil && d.ReceitasOrgaoPartidario[0].Receita != nil {
			ro := d.ReceitasOrgaoPartidario[0].Receita
			p := prefixoSocio(doc)
			return fmt.Sprintf("%sfez doação a partido (R$ %.2f)", p, ro.Valor)
		}
	}

	return fragmentoGenerico(v.Tipo, doc)
}

func prefixoSocio(doc domain.DocumentoVinculo) string {
	if doc.Origem == "socio" && doc.Nome != "" {
		return "sócio " + doc.Nome + " "
	}
	return ""
}

func fragmentoGenerico(tipo string, doc domain.DocumentoVinculo) string {
	p := prefixoSocio(doc)

	switch tipo {
	case "fornecedor":
		if p != "" {
			return p + "é fornecedor em campanha política"
		}
		return "é fornecedor em campanha política"
	case "tcu_contas_irregulares":
		if p != "" {
			return p + "tem contas irregulares no TCU"
		}
		return "tem contas irregulares no TCU"
	case "tcu_inabilitado":
		if p != "" {
			return p + "é inabilitado pelo TCU"
		}
		return "é inabilitado pelo TCU"
	case "tcu_inidoneo":
		if p != "" {
			return p + "é inidôneo pelo TCU"
		}
		return "é inidôneo pelo TCU"
	case "servidor_publico":
		if p != "" {
			return p + "é servidor público federal"
		}
		return "é servidor público federal"
	case "dispensa_valor_limite":
		return "dispensa de licitação acima do limite legal de baixo valor"
	default:
		info, ok := tipoParaInfo[tipo]
		if !ok {
			return tipo
		}
		if p != "" {
			return p + info.Fragmento
		}
		return info.Fragmento
	}
}
