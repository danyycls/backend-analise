// Pacote service contém funções auxiliares para a importação de dados
// eleitorais. Este arquivo implementa o cache de prestações de contas.
package parse

import (
	"context"
	"fmt"
	"strconv"

	"github.com/danyele/podp/internal/shared/types"
)

// -----------------------------------------------------------------------------
// Cache de prestação de contas.
// -----------------------------------------------------------------------------

// cachearPrestacao armazena uma prestação de contas no mapa interno,
// indexada pela chave natural composta por tipo de prestador, eleição
// e SQ do prestador.
func (p *ProcessadorLeitorCSV) cachearPrestacao(prestacao *types.PrestacaoContas) {
	if p.dados == nil {
		return
	}
	chave := chavePrestacaoNatural(prestacao.TipoPrestador, prestacao.EleicaoID, prestacao.SQPrestadorContas)
	if p.dados.Prestacoes != nil {
		p.dados.Prestacoes[chave] = prestacao
	}
	if p.dados.PrestacoesPorID != nil {
		p.dados.PrestacoesPorID[prestacao.ID] = prestacao
	}
	if p.dados.PrestacoesPorTipoESQ != nil {
		chaveTipoESQ := prestacao.TipoPrestador + "|" + prestacao.EleicaoID.String() + "|" + strconv.FormatInt(prestacao.SQPrestadorContas, 10)
		p.dados.PrestacoesPorTipoESQ[chaveTipoESQ] = prestacao
	}
}

// -----------------------------------------------------------------------------
// Garantia de entidades de domínio (busca ou criação com deduplicação).
// -----------------------------------------------------------------------------

// -----------------------------------------------------------------------------
// Busca de prestação de contas por tipo e SQ do prestador.
// -----------------------------------------------------------------------------

// garantirPrestacaoPorTipoESQ localiza uma prestação de contas já cacheada
// pelo tipo de prestador e SQ do prestador. Retorna erro se não encontrar
// ou se houver múltiplas prestações com a mesma chave.
func (p *ProcessadorLeitorCSV) garantirPrestacaoPorTipoESQ(_ context.Context, tipoPrestador, sqTexto, codigoEleicao string) (*types.PrestacaoContas, error) {
	sq := inteiro64Opcional(sqTexto)
	if sq == nil {
		return nil, fmt.Errorf("SQ_PRESTADOR_CONTAS invalido")
	}
	var chave string
	cod := inteiroOpcional(codigoEleicao)
	if cod != nil {
		if e, ok := p.dados.Eleicoes[*cod]; ok {
			chave = tipoPrestador + "|" + e.ID.String() + "|" + strconv.FormatInt(*sq, 10)
		} else {
			p.log.Warn("eleicao nao encontrada em dados.Eleicoes — usando sentinela",
				"codigo_eleicao", *cod, "sq_prestador_contas", *sq)
			chave = tipoPrestador + "|" + garantirEleicaoSentinela(p.dados).String() + "|" + strconv.FormatInt(*sq, 10)
		}
	} else {
		p.log.Warn("CD_ELEICAO ausente no CSV — chave sem eleicao_id pode colidir",
			"sq_prestador_contas", *sq)
		chave = tipoPrestador + "|" + strconv.FormatInt(*sq, 10)
	}
	pr, ok := p.dados.PrestacoesPorTipoESQ[chave]
	if !ok {
		return nil, fmt.Errorf("prestacao tipo %s SQ %d nao encontrada (importe despesas contratadas antes das pagas)", tipoPrestador, *sq)
	}
	return pr, nil
}
