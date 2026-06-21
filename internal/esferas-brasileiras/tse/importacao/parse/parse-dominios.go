// Pacote service contém funções auxiliares para a importação de dados
// eleitorais. Este arquivo implementa o cache de prestações de contas.
package parse

import (
	"context"
	"fmt"

	"github.com/danyele/podp/internal/shared/types"
)

// -----------------------------------------------------------------------------
// Cache de prestação de contas.
// -----------------------------------------------------------------------------

// cachearPrestacao armazena uma prestação de contas no mapa interno,
// indexada pela chave natural composta por tipo de prestador, eleição
// e SQ do prestador.
func (p *ProcessadorLeitorCSV) cachearPrestacao(prestacao *types.PrestacaoContas) {
	chave := chavePrestacaoNatural(prestacao.TipoPrestador, prestacao.EleicaoID, prestacao.SQPrestadorContas)
	p.dados.Prestacoes[chave] = prestacao
	chaveTipoESQ := fmt.Sprintf("%s|%d", prestacao.TipoPrestador, prestacao.SQPrestadorContas)
	p.dados.PrestacoesPorTipoESQ[chaveTipoESQ] = prestacao
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
func (p *ProcessadorLeitorCSV) garantirPrestacaoPorTipoESQ(_ context.Context, tipoPrestador, sqTexto string) (*types.PrestacaoContas, error) {
	sq := inteiro64Opcional(sqTexto)
	chave := fmt.Sprintf("%s|%d", tipoPrestador, *sq)
	pr, ok := p.dados.PrestacoesPorTipoESQ[chave]
	if !ok {
		return nil, fmt.Errorf("prestacao tipo %s SQ %d nao encontrada (importe despesas contratadas antes das pagas)", tipoPrestador, sq)
	}
	return pr, nil
}
