// Pacote service contém funções auxiliares para a importação de dados
// eleitorais. Este arquivo implementa o parser da planilha de consulta de
// candidatos (consulta_cand_), que é a principal fonte de entidades do
// domínio: eleições, UFs, unidades eleitorais, partidos, federações,
// coligações, cargos, estados civis, graus de instrução, ocupações e
// candidatos.
package parse

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/uuid"

	"github.com/danyele/podp/internal/shared/types"
)

// -----------------------------------------------------------------------------
// Processamento da planilha de consulta de candidatos.
// -----------------------------------------------------------------------------

// processarConsultaCandidato percorre o CSV de consulta de candidatos linha a
// linha. Para cada registro, garante a existência (com deduplicação) de todas
// as entidades relacionadas (eleição, UF, unidade eleitoral, partido,
// federação, coligação, cargo, estado civil, grau de instrução, ocupação)
// antes de montar e armazenar o Candidato.
func (p *ProcessadorLeitorCSV) processarConsultaCandidato(ctx context.Context, caminho string) (int, error) {
	total := 0

	err := lerArquivoCSV(caminho, func(numeroLinha int, registro map[string]string) error {
		eleicaoID, err := p.garantirEleicao(ctx, registro["CD_ELEICAO"], registro["ANO_ELEICAO"], registro["CD_TIPO_ELEICAO"], registro["NM_TIPO_ELEICAO"], registro["DS_ELEICAO"], registro["DT_ELEICAO"])
		if err != nil {
			return erroLinha(numeroLinha, err)
		}

		ufSigla, err := p.garantirUF(ctx, registro["SG_UF"])
		if err != nil {
			return erroLinha(numeroLinha, err)
		}

		if _, err = p.garantirUnidadeEleitoral(ctx, ufSigla, registro["SG_UE"], registro["NM_UE"]); err != nil {
			return erroLinha(numeroLinha, err)
		}

		partidoID, err := p.garantirPartidoOpcional(ctx, registro["NR_PARTIDO"], registro["SG_PARTIDO"], registro["NM_PARTIDO"],
			registro["NR_FEDERACAO"], registro["SG_FEDERACAO"], registro["NM_FEDERACAO"],
			registro["SQ_COLIGACAO"], registro["NM_COLIGACAO"], registro["DS_COMPOSICAO_COLIGACAO"])
		if err != nil {
			return erroLinha(numeroLinha, err)
		}

		sqCandidato := inteiro64Opcional(registro["SQ_CANDIDATO"])

		nomeCandidato := textoOpcional(registro["NM_CANDIDATO"])

		candidato := &types.Candidato{
			SQCandidato:                  *sqCandidato,
			EleicaoID:                    eleicaoID,
			UFSigla:                      ufSigla,
			PartidoID:                    partidoID,
			CargoCodigo:                  inteiroOpcional(registro["CD_CARGO"]),
			CargoNome:                    textoOpcional(registro["DS_CARGO"]),
			GeneroDescricao:              textoOpcional(registro["DS_GENERO"]),
			CorRacaDescricao:             textoOpcional(registro["DS_COR_RACA"]),
			EstadoCivilNome:              textoOpcional(registro["DS_ESTADO_CIVIL"]),
			GrauInstrucaoNome:            textoOpcional(registro["DS_GRAU_INSTRUCAO"]),
			OcupacaoCodigo:               inteiroOpcional(registro["CD_OCUPACAO"]),
			OcupacaoNome:                 textoOpcional(registro["DS_OCUPACAO"]),
			NumeroCandidato:              inteiroOpcional(registro["NR_CANDIDATO"]),
			CPF:                          documentoOpcional(registro["NR_CPF_CANDIDATO"]),
			NomeCompleto:                 nomeCandidato,
			NomeUrna:                     textoOpcional(registro["NM_URNA_CANDIDATO"]),
			NomeSocial:                   textoOpcional(registro["NM_SOCIAL_CANDIDATO"]),
			DataNascimento:               dataOpcional(registro["DT_NASCIMENTO"]),
			SituacaoTotalizacaoDescricao: textoOpcional(registro["DS_SIT_TOT_TURNO"]),
		}
		candidato.ID = uuid.Must(uuid.NewV7())
		p.dados.Candidatos[*sqCandidato] = candidato
		total++
		return nil
	})
	if err != nil {
		return 0, err
	}

	return total, nil
}

// -----------------------------------------------------------------------------
// Garantia de partido (opcional, obrigatório e por número).
// -----------------------------------------------------------------------------

// garantirPartidoOpcional retorna o UUID do partido se o número for
// informado, ou nil se a coluna estiver vazia.
func (p *ProcessadorLeitorCSV) garantirPartidoOpcional(ctx context.Context, numeroTexto, sigla, nome, federacaoCodigo, federacaoSigla, federacaoNome, coligacaoCodigo, coligacaoNome, coligacaoComposicao string) (*uuid.UUID, error) {
	numero := inteiro16Opcional(numeroTexto)
	if numero == nil {
		return nil, nil
	}
	partidoID, err := p.garantirPartidoPorNumero(ctx, *numero, sigla, nome, federacaoCodigo, federacaoSigla, federacaoNome, coligacaoCodigo, coligacaoNome, coligacaoComposicao)
	if err != nil {
		return nil, err
	}
	return &partidoID, nil
}

// garantirPartidoObrigatorio exige que NR_PARTIDO esteja presente; retorna
// erro se a coluna estiver vazia.
func (p *ProcessadorLeitorCSV) garantirPartidoObrigatorio(ctx context.Context, numeroTexto, sigla, nome, federacaoCodigo, federacaoSigla, federacaoNome, coligacaoCodigo, coligacaoNome, coligacaoComposicao string) (uuid.UUID, error) {
	numero := inteiro16Opcional(numeroTexto)
	if numero == nil {
		return uuid.Nil, fmt.Errorf("NR_PARTIDO obrigatorio")
	}
	return p.garantirPartidoPorNumero(ctx, *numero, sigla, nome, federacaoCodigo, federacaoSigla, federacaoNome, coligacaoCodigo, coligacaoNome, coligacaoComposicao)
}

// garantirPartidoPorNumero busca ou cria um partido pelo número. Usa sigla
// ou o próprio número como fallback para nome/sigla vazios.
func (p *ProcessadorLeitorCSV) garantirPartidoPorNumero(_ context.Context, numero int16, sigla, nome, federacaoCodigo, federacaoSigla, federacaoNome, coligacaoCodigo, coligacaoNome, coligacaoComposicao string) (uuid.UUID, error) {
	if existente, ok := p.dados.Partidos[numero]; ok {
		return existente.ID, nil
	}

	partido := &types.Partido{
		Numero:              numero,
		Sigla:               textoComFallback(sigla, strconv.Itoa(int(numero))),
		Nome:                textoComFallback(nome, textoComFallback(sigla, strconv.Itoa(int(numero)))),
		FederacaoCodigoTSE:  inteiro64Opcional(federacaoCodigo),
		FederacaoSigla:      textoOpcional(federacaoSigla),
		FederacaoNome:       textoOpcional(federacaoNome),
		ColigacaoCodigoTSE:  inteiro64Opcional(coligacaoCodigo),
		ColigacaoNome:       textoOpcional(coligacaoNome),
		ColigacaoComposicao: textoOpcional(coligacaoComposicao),
	}
	partido.ID = uuid.Must(uuid.NewV7())
	p.dados.Partidos[numero] = partido
	return partido.ID, nil
}

// -----------------------------------------------------------------------------
// Garantia de candidato e fornecedor por chave natural.
// -----------------------------------------------------------------------------

// garantirIDCandidato localiza um candidato já importado pelo SQ_CANDIDATO.
// Exige que a planilha consulta_cand tenha sido processada antes.
func (p *ProcessadorLeitorCSV) garantirIDCandidato(ctx context.Context, sqTexto string) (uuid.UUID, error) {
	sq := inteiro64Opcional(sqTexto)
	if sq == nil {
		return uuid.Nil, fmt.Errorf("SQ_CANDIDATO invalido ou ausente")
	}

	if p.cacheCandidatos != nil {
		if candidato, ok := p.cacheCandidatos[*sq]; ok {
			return candidato.ID, nil
		}
	}
	if candidato, ok := p.dados.Candidatos[*sq]; ok {
		return candidato.ID, nil
	}
	if p.buscarCandidatoPorSQ != nil {
		id, err := p.buscarCandidatoPorSQ(ctx, *sq)
		if err == nil {
			p.dados.Candidatos[*sq] = &types.Candidato{ModeloBase: types.ModeloBase{ID: id}, SQCandidato: *sq}
			return id, nil
		}
	}
	return uuid.Nil, fmt.Errorf("candidato SQ %d nao encontrado; importe consulta_cand antes de bens/despesas", *sq)
}

// garantirFornecedorOpcional busca ou cria um fornecedor pelo CPF/CNPJ.
// Se o documento não for informado, retorna nil sem erro.
func (p *ProcessadorLeitorCSV) garantirFornecedorOpcional(ctx context.Context, registro map[string]string) (*uuid.UUID, error) {
	documento := documentoOpcional(registro["NR_CPF_CNPJ_FORNECEDOR"])
	if documento == "" {
		return nil, nil
	}

	if existente, ok := p.dados.Fornecedores[documento]; ok {
		id := existente.ID
		return &id, nil
	}

	ufSigla, err := p.garantirUFOpcional(ctx, registro["SG_UF_FORNECEDOR"])
	if err != nil {
		return nil, err
	}

	nomeFornecedor := textoOpcional(registro["NM_FORNECEDOR"])

	fornecedor := &types.Fornecedor{
		CPFCNPJ:                    documento,
		Nome:                       nomeFornecedor,
		NomeRFB:                    textoOpcional(registro["NM_FORNECEDOR_RFB"]),
		TipoFornecedorCodigo:       inteiroOpcional(registro["CD_TIPO_FORNECEDOR"]),
		TipoFornecedorDescricao:    textoOpcional(registro["DS_TIPO_FORNECEDOR"]),
		CNAECodigo:                 textoOpcional(registro["CD_CNAE_FORNECEDOR"]),
		CNAEDescricao:              textoOpcional(registro["DS_CNAE_FORNECEDOR"]),
		EsferaPartidariaCodigo:     textoOpcional(registro["CD_ESFERA_PART_FORNECEDOR"]),
		EsferaPartidariaDescricao:  textoOpcional(registro["DS_ESFERA_PART_FORNECEDOR"]),
		UFSigla:                    ufSigla,
		MunicipioNome:              textoOpcional(registro["NM_MUNICIPIO_FORNECEDOR"]),
		SQCandidatoRelacionado:     inteiro64Opcional(registro["SQ_CANDIDATO_FORNECEDOR"]),
		NumeroCandidatoRelacionado: inteiroOpcional(registro["NR_CANDIDATO_FORNECEDOR"]),
		CargoCodigoRelacionado:     inteiroOpcional(registro["CD_CARGO_FORNECEDOR"]),
		CargoDescricaoRelacionada:  textoOpcional(registro["DS_CARGO_FORNECEDOR"]),
		PartidoNumeroRelacionado:   inteiro16Opcional(registro["NR_PARTIDO_FORNECEDOR"]),
		PartidoSiglaRelacionado:    textoOpcional(registro["SG_PARTIDO_FORNECEDOR"]),
		PartidoNomeRelacionado:     textoOpcional(registro["NM_PARTIDO_FORNECEDOR"]),
	}
	fornecedor.ID = uuid.Must(uuid.NewV7())
	p.dados.Fornecedores[documento] = fornecedor
	id := fornecedor.ID
	return &id, nil
}

func (p *ProcessadorLeitorCSV) garantirDoadorReceitaOpcional(ctx context.Context, registro map[string]string) (*uuid.UUID, error) {
	documento := documentoOpcional(registro["NR_CPF_CNPJ_DOADOR"])
	if documento == "" {
		return nil, nil
	}

	if existente, ok := p.dados.Doadores[documento]; ok {
		id := existente.ID
		return &id, nil
	}

	ufSigla, err := p.garantirUFOpcional(ctx, registro["SG_UF_DOADOR"])
	if err != nil {
		return nil, err
	}

	doador := &types.Doador{
		CPFCNPJ:                    documento,
		Nome:                       textoComFallback(registro["NM_DOADOR"], "DOADOR SEM NOME"),
		NomeRFB:                    textoOpcional(registro["NM_DOADOR_RFB"]),
		CNAECodigo:                 textoOpcional(registro["CD_CNAE_DOADOR"]),
		CNAEDescricao:              textoOpcional(registro["DS_CNAE_DOADOR"]),
		EsferaPartidariaCodigo:     textoOpcional(registro["CD_ESFERA_PARTIDARIA_DOADOR"]),
		EsferaPartidariaDescricao:  textoOpcional(registro["DS_ESFERA_PARTIDARIA_DOADOR"]),
		UFSigla:                    ufSigla,
		MunicipioNome:              textoOpcional(registro["NM_MUNICIPIO_DOADOR"]),
		SQCandidatoRelacionado:     inteiro64Opcional(registro["SQ_CANDIDATO_DOADOR"]),
		NumeroCandidatoRelacionado: inteiroOpcional(registro["NR_CANDIDATO_DOADOR"]),
		CargoCodigoRelacionado:     inteiroOpcional(registro["CD_CARGO_CANDIDATO_DOADOR"]),
		CargoDescricaoRelacionada:  textoOpcional(registro["DS_CARGO_CANDIDATO_DOADOR"]),
		PartidoNumeroRelacionado:   inteiro16Opcional(registro["NR_PARTIDO_DOADOR"]),
		PartidoSiglaRelacionado:    textoOpcional(registro["SG_PARTIDO_DOADOR"]),
		PartidoNomeRelacionado:     textoOpcional(registro["NM_PARTIDO_DOADOR"]),
	}
	doador.ID = uuid.Must(uuid.NewV7())
	p.dados.Doadores[documento] = doador
	id := doador.ID
	return &id, nil
}
