package repositorios

import "context"

type PNCPRepository interface {
	SalvarContratos(ctx context.Context, contratos []ContratoPersistido) error
	SalvarFornecedores(ctx context.Context, fornecedores []FornecedorPersistido) error
	SalvarSocio(ctx context.Context, socio SocioPersistido) (string, error)
	SalvarFornecedorSocios(ctx context.Context, vinculos []FornecedorSocioPersistido) error
	SalvarAmparosLegais(ctx context.Context, amparos []AmparoLegalPersistido) error

	BuscarContratosPorFiltro(ctx context.Context, tipo, valor string, ano, mes int) ([]ContratoPersistido, error)
	BuscarContratosPorFiltroEPeriodo(ctx context.Context, tipo, valor, dataInicial, dataFinal string) ([]ContratoPersistido, error)
	BuscarContratoPorNumeroControle(ctx context.Context, numeroControle string) (*ContratoPersistido, error)

	BuscarFornecedor(ctx context.Context, cnpj string) (*FornecedorPersistido, error)
	BuscarSociosPorFornecedor(ctx context.Context, cnpj string) ([]FornecedorSocioPersistido, error)
	BuscarFornecedoresPorSocio(ctx context.Context, cnpjCpfSocio string) ([]FornecedorSocioPersistido, error)

	BuscaJaRealizada(ctx context.Context, tipo, valor string, ano, mes int) (bool, error)
	RegistrarBusca(ctx context.Context, controle BuscaControlePersistido) error
	AtualizarBusca(ctx context.Context, controle BuscaControlePersistido) error

	ComTransaction(ctx context.Context, fn func(PNCPRepository) error) error
}
