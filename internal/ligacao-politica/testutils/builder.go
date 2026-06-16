package testutils

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/danyele/laceu/internal/ligacao-politica/usecase"
	"github.com/danyele/laceu/internal/shared/database"
)

type AnaliseRequestBuilder struct {
	licitacoes []usecase.AnalisarLigacaoPoliticaRequest
}

func NovoBuilder() *AnaliseRequestBuilder {
	return &AnaliseRequestBuilder{}
}

func (b *AnaliseRequestBuilder) WithLicitacaoCNPJ() *AnaliseRequestBuilder {
	cnpj := gofakeit.DigitN(14)
	lic := usecase.AnalisarLigacaoPoliticaRequest{
		NumeroControlePncp: gofakeit.UUID(),
		CpfCnpj:            cnpj,
	}
	qtd := gofakeit.Number(1, 3)
	for range qtd {
		lic.Socios = append(lic.Socios, usecase.AnalisarLigacaoPoliticaSocioRequest{
			Nome:      gofakeit.Name(),
			Documento: gofakeit.DigitN(11),
		})
	}
	b.licitacoes = append(b.licitacoes, lic)
	return b
}

func (b *AnaliseRequestBuilder) WithLicitacaoCPF() *AnaliseRequestBuilder {
	lic := usecase.AnalisarLigacaoPoliticaRequest{
		NumeroControlePncp: gofakeit.UUID(),
		CpfCnpj:            gofakeit.DigitN(11),
	}
	b.licitacoes = append(b.licitacoes, lic)
	return b
}

func (b *AnaliseRequestBuilder) WithLicitacaoDocInvalido() *AnaliseRequestBuilder {
	lic := usecase.AnalisarLigacaoPoliticaRequest{
		NumeroControlePncp: gofakeit.UUID(),
		CpfCnpj:            "12",
	}
	b.licitacoes = append(b.licitacoes, lic)
	return b
}

func (b *AnaliseRequestBuilder) Build() []usecase.AnalisarLigacaoPoliticaRequest {
	return b.licitacoes
}

func (b *AnaliseRequestBuilder) Persist(ctx context.Context, db database.DB) error {
	for _, lic := range b.licitacoes {
		doc := lic.CpfCnpj
		if len(doc) < 3 {
			continue
		}
		nome := "Fornecedor Fake"
		if len(doc) == 11 {
			nome = "Doador Fake"
		}
		tabela := "fornecedor"
		if len(doc) == 11 {
			tabela = "doador"
		}
		_, err := db.Exec(ctx,
			`INSERT INTO `+tabela+` (id, cpf_cnpj, nome, created_at, updated_at)
			 VALUES (gen_random_uuid(), $1, $2, NOW(), NOW())
			 ON CONFLICT (cpf_cnpj) DO NOTHING`, doc, nome)
		if err != nil {
			return err
		}
	}
	return nil
}
