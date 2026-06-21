package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/danyele/laceu/internal/ligacao-politica/handler"
	"github.com/danyele/laceu/internal/ligacao-politica/usecase"
	opencnpjPkg "github.com/danyele/laceu/internal/shared/clients/opencnpj"
	tcuPkg "github.com/danyele/laceu/internal/shared/clients/tcu"
	redisPkg "github.com/danyele/laceu/internal/shared/redis"
	"github.com/danyele/laceu/internal/shared/testkit"
	"github.com/danyele/laceu/internal/shared/types"
)

type analisarTestCase struct {
	testkit.IntegrationTestCase
	body  any
	mocks func(ocMock *opencnpjPkg.MockClient, tcuMock *tcuPkg.MockClient, redisMock *redisPkg.MockCache)
}

func TestAnalisarLigacaoPolitica_Integration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	pgc, pool, db, err := testkit.StartPostgresContainer(ctx)
	require.NoError(t, err)
	defer pgc.Terminate(ctx)
	defer pool.Close()

	cases := []analisarTestCase{
		{
			IntegrationTestCase: testkit.IntegrationTestCase{
				Name: "CNPJ de fornecedor encontrado",
				Fixtures: func(ctx context.Context, t *testing.T, pool *pgxpool.Pool) {
					testkit.InsertFornecedor(ctx, t, pool, "11222333000181", "Fornecedor Teste Ltda")
				},
				Assert: func(t *testing.T, w *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusOK, w.Code)
					var resp usecase.AnalisarLigacaoPoliticaResponse
					require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
					require.Equal(t, 1, resp.DocumentosProcessados)
					require.Len(t, resp.Resultados, 1)
					r := resp.Resultados[0]
					require.Equal(t, "pncp-001", r.NumeroControlePncp)
					require.Len(t, r.Documentos, 1)
					d := r.Documentos[0]
					require.Equal(t, "11222333000181", d.DocumentoNormalizado)
					require.Len(t, d.Vinculos, 1)
					require.Equal(t, "fornecedor", d.Vinculos[0].Tipo)
				},
			},
			body: map[string][]usecase.AnalisarLigacaoPoliticaRequest{
				"licitacoes": {
					{NumeroControlePncp: "pncp-001", CpfCnpj: "11222333000181"},
				},
			},
			mocks: mocksVazios,
		},
		{
			IntegrationTestCase: testkit.IntegrationTestCase{
				Name: "CNPJ de fornecedor encontrado com prefixo 000 para CPF nao se aplica",
				Fixtures: func(ctx context.Context, t *testing.T, pool *pgxpool.Pool) {
					testkit.InsertFornecedor(ctx, t, pool, "11222333000181", "Fornecedor Teste Ltda")
				},
				Assert: func(t *testing.T, w *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusOK, w.Code)
					var resp usecase.AnalisarLigacaoPoliticaResponse
					require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
					require.Len(t, resp.Resultados[0].Documentos, 1)
					require.Empty(t, resp.Resultados[0].Documentos[0].Vinculos)
				},
			},
			body: map[string][]usecase.AnalisarLigacaoPoliticaRequest{
				"licitacoes": {
					{NumeroControlePncp: "pncp-011", CpfCnpj: "00011222333000181"},
				},
			},
			mocks: mocksVazios,
		},
		{
			IntegrationTestCase: testkit.IntegrationTestCase{
				Name: "CPF de doador encontrado",
				Fixtures: func(ctx context.Context, t *testing.T, pool *pgxpool.Pool) {
					testkit.InsertDoador(ctx, t, pool, "11122233344", "Doador Teste")
				},
				Assert: func(t *testing.T, w *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusOK, w.Code)
					var resp usecase.AnalisarLigacaoPoliticaResponse
					require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
					require.Len(t, resp.Resultados[0].Documentos, 1)
					d := resp.Resultados[0].Documentos[0]
					require.Equal(t, "11122233344", d.DocumentoNormalizado)
					require.Len(t, d.Vinculos, 1)
					require.Equal(t, "doador", d.Vinculos[0].Tipo)
				},
			},
			body: map[string][]usecase.AnalisarLigacaoPoliticaRequest{
				"licitacoes": {
					{NumeroControlePncp: "pncp-002", CpfCnpj: "11122233344"},
				},
			},
			mocks: mocksVazios,
		},
		{
			IntegrationTestCase: testkit.IntegrationTestCase{
				Name: "CPF de doador encontrado com prefixo 000",
				Fixtures: func(ctx context.Context, t *testing.T, pool *pgxpool.Pool) {
					testkit.InsertDoador(ctx, t, pool, "11122233344", "Doador Teste")
				},
				Assert: func(t *testing.T, w *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusOK, w.Code)
					var resp usecase.AnalisarLigacaoPoliticaResponse
					require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
					require.Len(t, resp.Resultados[0].Documentos, 1)
					d := resp.Resultados[0].Documentos[0]
					require.Equal(t, "00011122233344", d.DocumentoInput)
					require.Equal(t, "11122233344", d.DocumentoNormalizado)
					require.Len(t, d.Vinculos, 1)
				},
			},
			body: map[string][]usecase.AnalisarLigacaoPoliticaRequest{
				"licitacoes": {
					{NumeroControlePncp: "pncp-012", CpfCnpj: "00011122233344"},
				},
			},
			mocks: mocksVazios,
		},
		{
			IntegrationTestCase: testkit.IntegrationTestCase{
				Name: "Documento sem correspondencia no banco",
				Assert: func(t *testing.T, w *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusOK, w.Code)
					var resp usecase.AnalisarLigacaoPoliticaResponse
					require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
					require.Equal(t, 1, resp.DocumentosProcessados)
					require.Len(t, resp.Resultados[0].Documentos, 1)
					require.Empty(t, resp.Resultados[0].Documentos[0].Vinculos)
				},
			},
			body: map[string][]usecase.AnalisarLigacaoPoliticaRequest{
				"licitacoes": {
					{NumeroControlePncp: "pncp-003", CpfCnpj: "00000000000000"},
				},
			},
			mocks: mocksVazios,
		},
		{
			IntegrationTestCase: testkit.IntegrationTestCase{
				Name: "Documento invalido (len < 3) ignorado",
				Assert: func(t *testing.T, w *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusOK, w.Code)
					var resp usecase.AnalisarLigacaoPoliticaResponse
					require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
					require.Equal(t, 1, resp.DocumentosProcessados)
					require.Empty(t, resp.Resultados[0].Documentos)
				},
			},
			body: map[string][]usecase.AnalisarLigacaoPoliticaRequest{
				"licitacoes": {
					{NumeroControlePncp: "pncp-004", CpfCnpj: "12"},
				},
			},
			mocks: mocksVazios,
		},
		{
			IntegrationTestCase: testkit.IntegrationTestCase{
				Name: "Multiplas licitacoes (fornecedor + doador + sem dados)",
				Fixtures: func(ctx context.Context, t *testing.T, pool *pgxpool.Pool) {
					testkit.InsertFornecedor(ctx, t, pool, "11222333000181", "Fornecedor Teste Ltda")
					testkit.InsertDoador(ctx, t, pool, "11122233344", "Doador Teste")
				},
				Assert: func(t *testing.T, w *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusOK, w.Code)
					var resp usecase.AnalisarLigacaoPoliticaResponse
					require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
					require.Equal(t, 3, resp.DocumentosProcessados)
					require.Len(t, resp.Resultados, 3)
				},
			},
			body: map[string][]usecase.AnalisarLigacaoPoliticaRequest{
				"licitacoes": {
					{NumeroControlePncp: "pncp-005", CpfCnpj: "11222333000181"},
					{NumeroControlePncp: "pncp-006", CpfCnpj: "11122233344"},
					{NumeroControlePncp: "pncp-007", CpfCnpj: "00000000000000"},
				},
			},
			mocks: mocksVazios,
		},
		{
			IntegrationTestCase: testkit.IntegrationTestCase{
				Name: "Enriquecimento OpenCNPJ para fornecedor CNPJ",
				Fixtures: func(ctx context.Context, t *testing.T, pool *pgxpool.Pool) {
					testkit.InsertFornecedor(ctx, t, pool, "11222333000181", "Fornecedor Teste Ltda")
				},
				Assert: func(t *testing.T, w *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusOK, w.Code)
					var resp usecase.AnalisarLigacaoPoliticaResponse
					require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
					require.Len(t, resp.Resultados[0].Documentos, 1)
					v := resp.Resultados[0].Documentos[0].Vinculos[0]
					require.Equal(t, "fornecedor", v.Tipo)
					require.NotNil(t, v.Detalhes)
					require.NotNil(t, v.Detalhes.Fornecedor)
					require.NotNil(t, v.Detalhes.Fornecedor.Fornecedor.Enriquecimento)
					require.Equal(t, "Empresa Teste Ltda", *v.Detalhes.Fornecedor.Fornecedor.Enriquecimento.RazaoSocial)
					require.Equal(t, "ATIVA", *v.Detalhes.Fornecedor.Fornecedor.Enriquecimento.SituacaoCadastral)
				},
			},
			body: map[string][]usecase.AnalisarLigacaoPoliticaRequest{
				"licitacoes": {
					{NumeroControlePncp: "pncp-008", CpfCnpj: "11222333000181"},
				},
			},
			mocks: func(ocMock *opencnpjPkg.MockClient, tcuMock *tcuPkg.MockClient, redisMock *redisPkg.MockCache) {
				ocMock.EXPECT().Buscar(gomock.Any(), "11222333000181").Return(
					&types.OpenCNPJResponse{
						CNPJ:              "11222333000181",
						RazaoSocial:       "Empresa Teste Ltda",
						NomeFantasia:      "Teste Fantasia",
						SituacaoCadastral: "ATIVA",
						CapitalSocial:     "10000.00",
					}, nil)
				ocMock.EXPECT().Buscar(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
				mocksTcuVazio(tcuMock)
				mocksCacheMiss(redisMock)
			},
		},
		{
			IntegrationTestCase: testkit.IntegrationTestCase{
				Name: "Enriquecimento TCU para documento com dados",
				Fixtures: func(ctx context.Context, t *testing.T, pool *pgxpool.Pool) {
					testkit.InsertDoador(ctx, t, pool, "11122233344", "Doador Teste")
				},
				Assert: func(t *testing.T, w *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusOK, w.Code)
					var resp usecase.AnalisarLigacaoPoliticaResponse
					require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
					v := resp.Resultados[0].Documentos[0].Vinculos
					hasTCU := false
					for _, vv := range v {
						if vv.Tipo == "tcu_contas_irregulares" {
							hasTCU = true
							require.NotNil(t, vv.Detalhes)
							require.Len(t, vv.Detalhes.ContasIrregulares, 1)
							require.Equal(t, "Doador Teste", vv.Detalhes.ContasIrregulares[0].Nome)
						}
					}
					require.True(t, hasTCU, "deveria ter vinculo TCU de contas irregulares")
				},
			},
			body: map[string][]usecase.AnalisarLigacaoPoliticaRequest{
				"licitacoes": {
					{NumeroControlePncp: "pncp-009", CpfCnpj: "11122233344"},
				},
			},
			mocks: func(ocMock *opencnpjPkg.MockClient, tcuMock *tcuPkg.MockClient, redisMock *redisPkg.MockCache) {
				ocMock.EXPECT().Buscar(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
				tcuMock.EXPECT().BuscarContasIrregulares(gomock.Any(), gomock.Any()).Return(
					[]tcuPkg.ContasIrregulares{
						{
							NumeroProcessoFormatado:    "0001234-56.2020.7.00.0000",
							Nome:                       "Doador Teste",
							TipoRegistro:               "CPF",
							NumeroRegistro:             "11122233344",
							Municipio:                  "Brasília",
							UF:                         "DF",
							DataTransitoEmJulgado:      "01/01/2023",
							LinkDeliberacoesProcesso:   "http://example.com/delib",
							LinkAcompanhamentoProcesso: "http://example.com/acomp",
						},
					}, nil)
				tcuMock.EXPECT().BuscarInidoneos(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
				tcuMock.EXPECT().BuscarInabilitados(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
				tcuMock.EXPECT().BuscarContasIrregulares(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
				mocksCacheMiss(redisMock)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ocMock := opencnpjPkg.NewMockClient(ctrl)
			tcuMock := tcuPkg.NewMockClient(ctrl)
			redisMock := redisPkg.NewMockCache(ctrl)

			tc.mocks(ocMock, tcuMock, redisMock)

			testkit.CleanAllTables(ctx, t, pool)
			if tc.Fixtures != nil {
				tc.Fixtures(ctx, t, pool)
			}

			uc := usecase.NovoAnalisarLigacaoPoliticaUseCase(db, ocMock, tcuMock)
			h := handler.NovoAnalisarLigacaoPoliticaHandler(uc, redisMock)

			r := testkit.NewGinEngine(func(r *gin.Engine) {
				r.POST("/analisar", h.Analisar)
			})
			req := testkit.NewRequest(http.MethodPost, "/analisar", tc.body)
			w := testkit.ExecRequest(r, req)

			tc.Assert(t, w)
		})
	}
}

func mocksVazios(ocMock *opencnpjPkg.MockClient, tcuMock *tcuPkg.MockClient, redisMock *redisPkg.MockCache) {
	ocMock.EXPECT().Buscar(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
	mocksTcuVazio(tcuMock)
	mocksCacheMiss(redisMock)
}

func mocksTcuVazio(tcuMock *tcuPkg.MockClient) {
	tcuMock.EXPECT().BuscarContasIrregulares(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
	tcuMock.EXPECT().BuscarInidoneos(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
	tcuMock.EXPECT().BuscarInabilitados(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
}

func mocksCacheMiss(redisMock *redisPkg.MockCache) {
	redisMock.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, nil)
	redisMock.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
}
