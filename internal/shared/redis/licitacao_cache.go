package redis

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/danyele/podp/internal/shared/clients/pncp"
	"github.com/danyele/podp/internal/shared/logger"
)

const (
	chaveContrato   = "podp:licitacao:%s"
	chaveIdxUF      = "podp:idx:licitacao:uf:%s"
	chaveIdxMun     = "podp:idx:licitacao:municipio:%s"
	chaveIdxMes     = "podp:idx:licitacao:mes:%s"
	chaveIdxOrgao   = "podp:idx:licitacao:orgao:%s"
	chaveTempSUniao = "podp:tmp:sunicao:%d_%d"
	ttlCache        = 30 * 24 * time.Hour
	ttlTemp         = 30 * time.Second
)

type LicitacaoCache struct {
	cache Cache
}

func NovoLicitacaoCache(cache Cache) *LicitacaoCache {
	return &LicitacaoCache{cache: cache}
}

func (lc *LicitacaoCache) IndexarContratos(ctx context.Context, contratos []pncp.Contrato) error {
	log := logger.New("LicitacaoCache: IndexarContratos")

	for i := range contratos {
		c := &contratos[i]

		id := ""
		if c.NumeroControlePNCP != nil {
			id = *c.NumeroControlePNCP
		}
		if id == "" {
			continue
		}

		chave := fmt.Sprintf(chaveContrato, id)
		if err := lc.cache.SetEx(ctx, chave, c, ttlCache); err != nil {
			log.Warn("erro ao salvar contrato no cache", "id", id, "erro", err)
			continue
		}

		if uf := extrairUF(c); uf != "" {
			chaveUF := fmt.Sprintf(chaveIdxUF, uf)
			if _, err := lc.cache.SAdd(ctx, chaveUF, id); err != nil {
				log.Warn("erro ao indexar UF", "uf", uf, "erro", err)
			}
		}

		if codMun := extrairCodigoMunicipio(c); codMun != "" {
			chaveMun := fmt.Sprintf(chaveIdxMun, codMun)
			if _, err := lc.cache.SAdd(ctx, chaveMun, id); err != nil {
				log.Warn("erro ao indexar municipio", "codigo", codMun, "erro", err)
			}
		}

		if mes := extrairMes(c); mes != "" {
			chaveMes := fmt.Sprintf(chaveIdxMes, mes)
			if _, err := lc.cache.SAdd(ctx, chaveMes, id); err != nil {
				log.Warn("erro ao indexar mes", "mes", mes, "erro", err)
			}
		}

		if cnpjOrgao := extrairCNPJOrgao(c); cnpjOrgao != "" {
			chaveOrgao := fmt.Sprintf(chaveIdxOrgao, cnpjOrgao)
			if _, err := lc.cache.SAdd(ctx, chaveOrgao, id); err != nil {
				log.Warn("erro ao indexar orgao", "cnpj", cnpjOrgao, "erro", err)
			}
		}
	}

	return nil
}

func (lc *LicitacaoCache) BuscarPorFiltros(ctx context.Context, tipo, valor, dataInicial, dataFinal string) ([]pncp.Contrato, bool, error) {
	log := logger.New("LicitacaoCache: BuscarPorFiltros")

	meses := gerarMeses(dataInicial, dataFinal)
	if len(meses) == 0 {
		return nil, false, nil
	}

	chavesMes := make([]string, len(meses))
	for i, m := range meses {
		chavesMes[i] = fmt.Sprintf(chaveIdxMes, m)
	}

	var chaveFiltro string
	switch tipo {
	case "uf":
		chaveFiltro = fmt.Sprintf(chaveIdxUF, strings.ToUpper(valor))
	case "municipio":
		chaveFiltro = fmt.Sprintf(chaveIdxMun, valor)
	case "orgao":
		chaveFiltro = fmt.Sprintf(chaveIdxOrgao, valor)
	default:
		log.Warn("tipo de filtro desconhecido", "tipo", tipo)
		return nil, false, nil
	}

	existeFiltro, err := lc.cache.Exists(ctx, chaveFiltro)
	if err != nil || existeFiltro == 0 {
		return nil, false, nil
	}

	todosExistem := true
	for _, chave := range chavesMes {
		ok, err := lc.cache.Exists(ctx, chave)
		if err != nil || ok == 0 {
			todosExistem = false
			break
		}
	}
	if !todosExistem {
		log.Info("cache parcial — nem todos os meses estao indexados")
		return nil, false, nil
	}

	chaveTemp := fmt.Sprintf(chaveTempSUniao, time.Now().UnixNano(), rand.Int63())
	defer func() {
		if _, err := lc.cache.Del(ctx, chaveTemp); err != nil {
			log.Warn("erro ao limpar chave temporaria", "erro", err)
		}
	}()

	if _, err := lc.cache.SUnionStore(ctx, chaveTemp, chavesMes...); err != nil {
		return nil, false, fmt.Errorf("erro ao unir meses: %w", err)
	}

	ids, err := lc.cache.SInter(ctx, chaveTemp, chaveFiltro)
	if err != nil {
		return nil, false, fmt.Errorf("erro ao intersectar indices: %w", err)
	}

	if len(ids) == 0 {
		return nil, false, nil
	}

	contratos := make([]pncp.Contrato, 0, len(ids))
	for _, id := range ids {
		var c pncp.Contrato
		ok, err := lc.cache.Get(ctx, fmt.Sprintf(chaveContrato, id), &c)
		if err != nil || !ok {
			continue
		}
		contratos = append(contratos, c)
	}

	log.Info("cache hit", "tipo", tipo, "valor", valor, "contratos", len(contratos))
	return contratos, true, nil
}

func extrairUF(c *pncp.Contrato) string {
	if c.UG != nil && c.UG.UFSigla != nil && *c.UG.UFSigla != "" {
		return strings.ToUpper(*c.UG.UFSigla)
	}
	if c.OrgaoVinculado != nil && c.OrgaoVinculado.UFSigla != nil && *c.OrgaoVinculado.UFSigla != "" {
		return strings.ToUpper(*c.OrgaoVinculado.UFSigla)
	}
	return ""
}

func extrairCodigoMunicipio(c *pncp.Contrato) string {
	if c.UG != nil && c.UG.CodigoIbge != nil && *c.UG.CodigoIbge != "" {
		return *c.UG.CodigoIbge
	}
	if c.OrgaoVinculado != nil && c.OrgaoVinculado.CodigoIbge != nil && *c.OrgaoVinculado.CodigoIbge != "" {
		return *c.OrgaoVinculado.CodigoIbge
	}
	return ""
}

func extrairCNPJOrgao(c *pncp.Contrato) string {
	if c.CNPJOrgao != nil && *c.CNPJOrgao != "" {
		return *c.CNPJOrgao
	}
	if c.OrgaoEntidade != nil && c.OrgaoEntidade.CNPJ != nil && *c.OrgaoEntidade.CNPJ != "" {
		return *c.OrgaoEntidade.CNPJ
	}
	return ""
}

func extrairMes(c *pncp.Contrato) string {
	data := ""
	if c.DataPublicacao != nil && *c.DataPublicacao != "" {
		data = *c.DataPublicacao
	} else if c.DataAssinatura != nil && *c.DataAssinatura != "" {
		data = *c.DataAssinatura
	} else if c.DataInicioVigencia != nil && *c.DataInicioVigencia != "" {
		data = *c.DataInicioVigencia
	}

	if data == "" {
		return ""
	}

	data = strings.TrimSpace(data)

	if len(data) >= 10 && data[4] == '-' && data[7] == '-' {
		return data[:7]
	}
	if len(data) >= 8 {
		return data[:6]
	}

	return ""
}

func gerarMeses(dataInicial, dataFinal string) []string {
	anoInicio, mesInicio := parseAnoMes(dataInicial)
	anoFim, mesFim := parseAnoMes(dataFinal)

	if anoInicio == 0 || mesInicio == 0 || anoFim == 0 || mesFim == 0 {
		return nil
	}

	if anoFim < anoInicio || (anoFim == anoInicio && mesFim < mesInicio) {
		return nil
	}

	var meses []string
	ano, mes := anoInicio, mesInicio
	for {
		meses = append(meses, fmt.Sprintf("%04d%02d", ano, mes))
		if ano == anoFim && mes == mesFim {
			break
		}
		mes++
		if mes > 12 {
			mes = 1
			ano++
		}
	}
	return meses
}

func parseAnoMes(data string) (int, int) {
	clean := strings.TrimSpace(data)

	if len(clean) >= 10 && clean[4] == '-' && clean[7] == '-' {
		ano := 0
		mes := 0
		for i := 0; i < 4; i++ {
			if clean[i] < '0' || clean[i] > '9' {
				return 0, 0
			}
			ano = ano*10 + int(clean[i]-'0')
		}
		for i := 5; i < 7; i++ {
			if clean[i] < '0' || clean[i] > '9' {
				return 0, 0
			}
			mes = mes*10 + int(clean[i]-'0')
		}
		return ano, mes
	}

	if len(clean) >= 8 {
		ano := 0
		mes := 0
		for i := 0; i < 4; i++ {
			if clean[i] < '0' || clean[i] > '9' {
				return 0, 0
			}
			ano = ano*10 + int(clean[i]-'0')
		}
		for i := 4; i < 6; i++ {
			if clean[i] < '0' || clean[i] > '9' {
				return 0, 0
			}
			mes = mes*10 + int(clean[i]-'0')
		}
		return ano, mes
	}

	return 0, 0
}
