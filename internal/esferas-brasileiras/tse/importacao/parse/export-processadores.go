package parse

import (
	"context"
	"fmt"
	tipos "github.com/danyele/podp/internal/esferas-brasileiras/tse/importacao/types"
)

func ProcessarArquivo(ctx context.Context, processador *ProcessadorLeitorCSV, arquivo tipos.ArquivoImportacao) (int, error) {
	switch arquivo.Tipo {
	case tipoArquivoConsultaCandidato:
		return processador.ProcessarConsultaCandidato(ctx, arquivo.Caminho)
	case tipoArquivoBemCandidato:
		return processador.ProcessarBemCandidato(ctx, arquivo.Caminho)
	case tipoArquivoDespesaContratadaCandidato:
		return processador.ProcessarDespesaContratadaCandidato(ctx, arquivo.Caminho)
	case tipoArquivoDespesaPagaCandidato:
		return processador.ProcessarDespesaPagaCandidato(ctx, arquivo.Caminho)
	case tipoArquivoReceitaCandidato:
		return processador.ProcessarReceitaCandidato(ctx, arquivo.Caminho)
	case tipoArquivoReceitaCandidatoDoadorOrigem:
		return processador.ProcessarReceitaCandidatoDoadorOriginario(ctx, arquivo.Caminho)
	case tipoArquivoDespesaContratadaOrgaoPartido:
		return processador.ProcessarDespesaContratadaOrgaoPartidario(ctx, arquivo.Caminho)
	case tipoArquivoDespesaPagaOrgaoPartido:
		return processador.ProcessarDespesaPagaOrgaoPartidario(ctx, arquivo.Caminho)
	case tipoArquivoReceitaOrgaoPartido:
		return processador.ProcessarReceitaOrgaoPartidario(ctx, arquivo.Caminho)
	case tipoArquivoReceitaOrgaoPartidoDoadorOrig:
		return processador.ProcessarReceitaOrgaoPartidarioDoadorOriginario(ctx, arquivo.Caminho)
	default:
		return 0, fmt.Errorf("tipo de arquivo nao suportado: %s", arquivo.Tipo)
	}
}

func (p *ProcessadorLeitorCSV) ProcessarConsultaCandidato(ctx context.Context, caminho string) (int, error) {
	return p.processarConsultaCandidato(ctx, caminho)
}

func (p *ProcessadorLeitorCSV) ProcessarBemCandidato(ctx context.Context, caminho string) (int, error) {
	return p.processarBemCandidato(ctx, caminho)
}

func (p *ProcessadorLeitorCSV) ProcessarDespesaContratadaCandidato(ctx context.Context, caminho string) (int, error) {
	return p.processarDespesaContratadaCandidato(ctx, caminho)
}

func (p *ProcessadorLeitorCSV) ProcessarDespesaPagaCandidato(ctx context.Context, caminho string) (int, error) {
	return p.processarDespesaPagaCandidato(ctx, caminho)
}

func (p *ProcessadorLeitorCSV) ProcessarReceitaCandidato(ctx context.Context, caminho string) (int, error) {
	return p.processarReceitaCandidato(ctx, caminho)
}

func (p *ProcessadorLeitorCSV) ProcessarReceitaCandidatoDoadorOriginario(ctx context.Context, caminho string) (int, error) {
	return p.processarReceitaCandidatoDoadorOriginario(ctx, caminho)
}

func (p *ProcessadorLeitorCSV) ProcessarDespesaContratadaOrgaoPartidario(ctx context.Context, caminho string) (int, error) {
	return p.processarDespesaContratadaOrgaoPartidario(ctx, caminho)
}

func (p *ProcessadorLeitorCSV) ProcessarDespesaPagaOrgaoPartidario(ctx context.Context, caminho string) (int, error) {
	return p.processarDespesaPagaOrgaoPartidario(ctx, caminho)
}

func (p *ProcessadorLeitorCSV) ProcessarReceitaOrgaoPartidario(ctx context.Context, caminho string) (int, error) {
	return p.processarReceitaOrgaoPartidario(ctx, caminho)
}

func (p *ProcessadorLeitorCSV) ProcessarReceitaOrgaoPartidarioDoadorOriginario(ctx context.Context, caminho string) (int, error) {
	return p.processarReceitaOrgaoPartidarioDoadorOriginario(ctx, caminho)
}
