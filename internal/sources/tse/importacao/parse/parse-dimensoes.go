// Pacote service contém funções auxiliares para a importação de dados
// eleitorais. Este arquivo implementa o parser das entidades dimensionais
// (eleições, UFs, unidades eleitorais) e funções auxiliares de geração de
// chaves naturais.
package parse

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"github.com/danyele/podp/internal/shared/types"
)

// -----------------------------------------------------------------------------
// Funções de geração de chaves naturais (usadas para deduplicação).
// -----------------------------------------------------------------------------

// chaveUnidadeEleitoral compõe a chave natural de uma unidade eleitoral
// no formato "UF|codigoTSE".
func chaveUnidadeEleitoral(ufSigla, codigoTSE string) string {
	return fmt.Sprintf("%s|%s", ufSigla, codigoTSE)
}

// chavePrestacaoNatural compõe a chave natural de uma prestação de contas
// no formato "tipoPrestador|eleicaoID|sqPrestador".
func chavePrestacaoNatural(tipoPrestador string, eleicaoID uuid.UUID, sqPrestador int64) string {
	buf := make([]byte, 0, len(tipoPrestador)+1+36+1+20)
	buf = append(buf, tipoPrestador...)
	buf = append(buf, '|')
	buf = append(buf, eleicaoID.String()...)
	buf = append(buf, '|')
	buf = strconv.AppendInt(buf, sqPrestador, 10)
	return string(buf)
}

// -----------------------------------------------------------------------------
// Garantia de entidades dimensionais.
// -----------------------------------------------------------------------------

// garantirEleicao busca uma eleição pelo código TSE ou cria um novo
// registro com os dados informados. Usa fallback para descrição vazia.
func (p *ProcessadorLeitorCSV) garantirEleicao(_ context.Context, codigoTexto, anoTexto, codigoTipoTexto, nomeTipo, descricao, dataTexto string) (uuid.UUID, error) {
	codigo := inteiroOpcional(codigoTexto)
	if codigo == nil {
		return uuid.Nil, fmt.Errorf("CD_ELEICAO obrigatorio e invalido: %q", codigoTexto)
	}

	if existente, ok := p.dados.Eleicoes[*codigo]; ok {
		return existente.ID, nil
	}

	ano := inteiroOpcional(anoTexto)
	if ano == nil {
		return uuid.Nil, fmt.Errorf("ANO_ELEICAO obrigatorio e invalido: %q", anoTexto)
	}

	e := &types.Eleicao{
		Ano:               int16(*ano),
		CodigoTSE:         *codigo,
		CodigoTipoEleicao: inteiroOpcional(codigoTipoTexto),
		NomeTipoEleicao:   textoOpcional(nomeTipo),
		Descricao:         textoComFallback(descricao, "Eleicao sem descricao"),
		DataEleicao:       dataOpcional(dataTexto),
	}
	e.ID = uuid.Must(uuid.NewV7())
	p.dados.Eleicoes[*codigo] = e
	return e.ID, nil
}

// garantirUF valida e retorna a sigla da UF em maiúsculas. Retorna erro
// se a sigla estiver vazia, pois UF é obrigatória na maioria dos registros.
func (p *ProcessadorLeitorCSV) garantirUF(_ context.Context, sigla string) (string, error) {
	valor := textoOpcional(sigla)
	if valor == "" {
		return "", fmt.Errorf("SG_UF obrigatoria")
	}
	return strings.ToUpper(valor), nil
}

// garantirUFOpcional retorna a sigla da UF em maiúsculas ou nil se a
// coluna estiver vazia. Usada quando UF é opcional (ex.: órgãos partidários).
func (p *ProcessadorLeitorCSV) garantirUFOpcional(_ context.Context, sigla string) (*string, error) {
	valor := textoOpcional(sigla)
	if valor == "" {
		return nil, nil
	}
	valor = strings.ToUpper(valor)
	return &valor, nil
}

// garantirUnidadeEleitoral busca ou cria uma unidade eleitoral compondo a
// chave UF|codigoTSE. Retorna erro se a unidade for obrigatória e não
// puder ser determinada.
func (p *ProcessadorLeitorCSV) garantirUnidadeEleitoral(ctx context.Context, ufSigla, codigoTSE, nome string) (uuid.UUID, error) {
	unidadeID, err := p.garantirUnidadeEleitoralOpcional(ctx, &ufSigla, codigoTSE, nome)
	if err != nil {
		return uuid.Nil, err
	}
	if unidadeID == nil {
		return uuid.Nil, fmt.Errorf("unidade eleitoral obrigatoria")
	}
	return *unidadeID, nil
}

// garantirUnidadeEleitoralOpcional busca ou cria uma unidade eleitoral.
// Retorna nil se UF ou código TSE não forem informados.
func (p *ProcessadorLeitorCSV) garantirUnidadeEleitoralOpcional(_ context.Context, ufSigla *string, codigoTSE, nome string) (*uuid.UUID, error) {
	if ufSigla == nil {
		return nil, nil
	}

	codigo := textoOpcional(codigoTSE)
	if codigo == "" {
		return nil, nil
	}

	chave := chaveUnidadeEleitoral(*ufSigla, codigo)
	if existente, ok := p.dados.UnidadesEleitorais[chave]; ok {
		id := existente.ID
		return &id, nil
	}

	entidade := &types.UnidadeEleitoral{
		UFSigla:   *ufSigla,
		CodigoTSE: codigo,
		Nome:      textoComFallback(nome, codigo),
	}
	entidade.ID = uuid.Must(uuid.NewV7())
	p.dados.UnidadesEleitorais[chave] = entidade
	id := entidade.ID
	return &id, nil
}

// -----------------------------------------------------------------------------
// Função utilitária de formatação de erro.
// -----------------------------------------------------------------------------

// erroLinha formata uma mensagem de erro incluindo o número da linha do CSV.
func erroLinha(numeroLinha int, err error) error {
	return fmt.Errorf("linha %d: %w", numeroLinha, err)
}
