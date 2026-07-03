package usecase

import (
	"context"

	repositorio "github.com/danyele/podp/internal/esferas-brasileiras/tse/repositorio"
	"github.com/danyele/podp/internal/shared/database"
	"github.com/danyele/podp/internal/shared/domain"
	"github.com/danyele/podp/internal/shared/logger"
	"github.com/danyele/podp/internal/shared/services"
)

type AnalisarLigacaoPoliticaUseCase struct {
	db                      database.DB
	normalizarDocumentosSvc services.NormalizarDocumentosService
	buscarLigacaoTseSvc     services.BuscarLigacaoPoliticaTSEService
	verificarSancoesTcuSvc  services.VerificarSancoesTCUService
	verificarServidorSvc    services.VerificarServidorPublicoService
}

func NovoAnalisarLigacaoPoliticaUseCase(
	db database.DB,
	normalizarDocumentosSvc services.NormalizarDocumentosService,
	buscarLigacaoTseSvc services.BuscarLigacaoPoliticaTSEService,
	verificarSancoesTcuSvc services.VerificarSancoesTCUService,
	verificarServidorSvc services.VerificarServidorPublicoService,
) *AnalisarLigacaoPoliticaUseCase {
	return &AnalisarLigacaoPoliticaUseCase{
		db:                      db,
		normalizarDocumentosSvc: normalizarDocumentosSvc,
		buscarLigacaoTseSvc:     buscarLigacaoTseSvc,
		verificarSancoesTcuSvc:  verificarSancoesTcuSvc,
		verificarServidorSvc:    verificarServidorSvc,
	}
}

func (u *AnalisarLigacaoPoliticaUseCase) Executar(ctx context.Context, licitacoes []AnalisarLigacaoPoliticaRequest) (*AnalisarLigacaoPoliticaResponse, error) {
	log := logger.New("LigacaoPolitica: UseCase: Executar")
	log.Info("licitacoes recebidas", "quantidade", len(licitacoes))

	repo := repositorio.Novo(u.db)

	resultados := u.normalizarDocumentos(ctx, licitacoes)

	u.buscarLigacoesTSE(ctx, repo, resultados)

	u.verificarSancoesTCU(ctx, resultados)

	u.verificarServidores(ctx, resultados)

	return &AnalisarLigacaoPoliticaResponse{
		DocumentosProcessados: len(resultados),
		Resultados:            resultados,
	}, nil
}

func (u *AnalisarLigacaoPoliticaUseCase) normalizarDocumentos(
	ctx context.Context,
	licitacoes []AnalisarLigacaoPoliticaRequest,
) []domain.VinculoLicitacao {
	resultados := make([]domain.VinculoLicitacao, 0, len(licitacoes))

	for _, lic := range licitacoes {
		vl := domain.VinculoLicitacao{
			NumeroControlePncp: lic.NumeroControlePncp,
			CpfCnpj:            lic.CpfCnpj,
		}

		sociosOut := make([]domain.SocioOutput, 0, len(lic.Socios))
		for _, s := range lic.Socios {
			sociosOut = append(sociosOut, domain.SocioOutput(s))

			dv, _ := u.normalizarDocumentosSvc.Executar(ctx, services.NormalizarDocumentosInput{
				Documento: s.Documento,
				Nome:      s.Nome,
				Origem:    "socio",
			})
			if dv != nil {
				vl.DocumentosVinculos = append(vl.DocumentosVinculos, dv.DocumentoVinculo)
			}
		}
		vl.Socios = sociosOut

		dv, _ := u.normalizarDocumentosSvc.Executar(ctx, services.NormalizarDocumentosInput{
			Documento: lic.CpfCnpj,
			Nome:      "",
			Origem:    "principal",
		})
		if dv != nil {
			vl.DocumentosVinculos = append(vl.DocumentosVinculos, dv.DocumentoVinculo)
		}

		resultados = append(resultados, vl)
	}

	return resultados
}

func (u *AnalisarLigacaoPoliticaUseCase) buscarLigacoesTSE(
	ctx context.Context,
	repo *repositorio.Repositorio,
	resultados []domain.VinculoLicitacao,
) {
	for li := range resultados {
		for di := range resultados[li].DocumentosVinculos {
			doc := &resultados[li].DocumentosVinculos[di]
			vinculos, _ := u.buscarLigacaoTseSvc.Executar(ctx, repo, services.BuscarLigacaoPoliticaTSEInput{
				DocumentoNormalizado: doc.DocumentoNormalizado,
				Nome:                 doc.Nome,
				Parcial:              doc.Parcial,
			})
			if vinculos != nil {
				doc.Vinculos = vinculos.Vinculos
			}
		}
	}
}

func (u *AnalisarLigacaoPoliticaUseCase) verificarSancoesTCU(
	ctx context.Context,
	resultados []domain.VinculoLicitacao,
) {
	documentos := coletarDocumentosNormalizados(resultados)
	if len(documentos) == 0 {
		return
	}

	resultado, err := u.verificarSancoesTcuSvc.Executar(ctx, services.VerificarSancoesTCUInput{
		Documentos: documentos,
	})
	if err != nil || resultado == nil {
		return
	}

	anexarVinculosPorDocumento(resultados, resultado.Resultados)
}

func (u *AnalisarLigacaoPoliticaUseCase) verificarServidores(
	ctx context.Context,
	resultados []domain.VinculoLicitacao,
) {
	documentos := coletarDocumentosNormalizados(resultados)
	if len(documentos) == 0 {
		return
	}

	resultado, err := u.verificarServidorSvc.Executar(ctx, services.VerificarServidorInput{
		Documentos: documentos,
	})
	if err != nil || resultado == nil {
		return
	}

	anexarVinculosPorDocumento(resultados, resultado.Resultados)
}

func coletarDocumentosNormalizados(resultados []domain.VinculoLicitacao) []string {
	if len(resultados) == 0 {
		return nil
	}

	documentos := make([]string, 0, len(resultados)*2)
	for _, r := range resultados {
		for _, d := range r.DocumentosVinculos {
			if len(d.DocumentoNormalizado) >= 3 {
				documentos = append(documentos, d.DocumentoNormalizado)
			}
		}
	}
	return documentos
}

func anexarVinculosPorDocumento(
	resultados []domain.VinculoLicitacao,
	vinculos map[string][]domain.Vinculo,
) {
	if len(vinculos) == 0 {
		return
	}

	for li := range resultados {
		for di := range resultados[li].DocumentosVinculos {
			doc := resultados[li].DocumentosVinculos[di].DocumentoNormalizado
			if v, ok := vinculos[doc]; ok {
				resultados[li].DocumentosVinculos[di].Vinculos = append(
					resultados[li].DocumentosVinculos[di].Vinculos, v...,
				)
			}
		}
	}
}
