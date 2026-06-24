package stream

import (
	"context"
	"encoding/json"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	gorilla "github.com/gorilla/websocket"

	dadosfinanceiros "github.com/danyele/podp/internal/esferas-brasileiras/estadual/usecase/dadosfinanceiros"
	handlerPNCP "github.com/danyele/podp/internal/esferas-brasileiras/federal/pncp/handler"
	usecaseMunicipal "github.com/danyele/podp/internal/esferas-brasileiras/municipal/usecase"
	pncp "github.com/danyele/podp/internal/shared/clients/pncp"
	"github.com/danyele/podp/internal/shared/logger"
	"github.com/danyele/podp/internal/shared/types"
	ws "github.com/danyele/podp/internal/shared/websocket"
)

const (
	msgSICONFIIndisponivelEstado = "API SICONFI (Tesouro Nacional) temporariamente indisponível. Os dados financeiros do estado não puderam ser carregados. Tente novamente mais tarde."
	msgSICONFIIndisponivelMuni   = "API SICONFI (Tesouro Nacional) temporariamente indisponível. Os dados financeiros do município não puderam ser carregados. Tente novamente mais tarde."
)

type StreamMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type ClientMessage struct {
	Channel    string `json:"channel"`
	JobID      string `json:"job_id,omitempty"`
	UF         string `json:"uf,omitempty"`
	CodigoIBGE int    `json:"codigo_ibge,omitempty"`
	Exercicio  int64  `json:"exercicio,omitempty"`
}

type Hub struct {
	orgaoHandler      *handlerPNCP.AnaliseOrgaoPNCPHandler
	publicacaoHandler *handlerPNCP.AnalisePublicacaoHandler

	despesaPessoalUC   *dadosfinanceiros.EsferaEstadualBuscarDespesaPessoalUseCase
	despesaCategoriaUC *dadosfinanceiros.EsferaEstadualBuscarDespesaCategoriaUseCase
	rreoUC             *dadosfinanceiros.EsferaEstadualBuscarRREOUseCase
	recursosFederaisUC *dadosfinanceiros.EsferaEstadualBuscarRecursosFederaisUseCase

	municipioUC *usecaseMunicipal.EsferaMunicipalBuscarDetalhesUseCase
}

func NewHub(
	orgaoHandler *handlerPNCP.AnaliseOrgaoPNCPHandler,
	publicacaoHandler *handlerPNCP.AnalisePublicacaoHandler,
	despesaPessoalUC *dadosfinanceiros.EsferaEstadualBuscarDespesaPessoalUseCase,
	despesaCategoriaUC *dadosfinanceiros.EsferaEstadualBuscarDespesaCategoriaUseCase,
	rreoUC *dadosfinanceiros.EsferaEstadualBuscarRREOUseCase,
	recursosFederaisUC *dadosfinanceiros.EsferaEstadualBuscarRecursosFederaisUseCase,
	municipioUC *usecaseMunicipal.EsferaMunicipalBuscarDetalhesUseCase,
) *Hub {
	return &Hub{
		orgaoHandler:       orgaoHandler,
		publicacaoHandler:  publicacaoHandler,
		despesaPessoalUC:   despesaPessoalUC,
		despesaCategoriaUC: despesaCategoriaUC,
		rreoUC:             rreoUC,
		recursosFederaisUC: recursosFederaisUC,
		municipioUC:        municipioUC,
	}
}

func (h *Hub) Handle(c *gin.Context) {
	conn, err := ws.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	_, msgBytes, err := conn.ReadMessage()
	if err != nil {
		return
	}

	var msg ClientMessage
	if err := json.Unmarshal(msgBytes, &msg); err != nil {
		return
	}

	ctx := c.Request.Context()

	switch msg.Channel {
	case "orgao_analise":
		h.streamOrgao(ctx, conn, msg.JobID)
	case "publicacao_analise":
		h.streamPublicacao(ctx, conn, msg.JobID)
	case "estado_financeiro":
		if msg.UF == "" || len(msg.UF) != 2 {
			ws.WriteJSON(conn, StreamMessage{Type: "erro", Data: map[string]string{"erro": "UF invalida"}})
			return
		}
		h.streamEstadoFinanceiro(ctx, conn, strings.ToUpper(msg.UF), msg.Exercicio)
	case "municipio_detalhes":
		if msg.CodigoIBGE <= 0 {
			ws.WriteJSON(conn, StreamMessage{Type: "erro", Data: map[string]string{"erro": "codigo IBGE invalido"}})
			return
		}
		h.streamMunicipioDetalhes(ctx, conn, msg.CodigoIBGE, msg.Exercicio)
	}
}

func (h *Hub) streamOrgao(ctx context.Context, conn *gorilla.Conn, jobID string) {
	eventChan, exists := h.orgaoHandler.GetJobChan(jobID)
	if !exists {
		ws.WriteJSON(conn, pncp.EventoAnalise{Type: "error", Message: "job nao encontrado"})
		return
	}
	for {
		select {
		case event, ok := <-eventChan:
			if !ok {
				ws.WriteJSON(conn, pncp.EventoAnalise{Type: "done"})
				return
			}
			if err := ws.WriteJSON(conn, event); err != nil {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func (h *Hub) streamPublicacao(ctx context.Context, conn *gorilla.Conn, jobID string) {
	eventChan, exists := h.publicacaoHandler.GetJobChan(jobID)
	if !exists {
		ws.WriteJSON(conn, pncp.EventoAnalise{Type: "error", Message: "job nao encontrado"})
		return
	}
	for {
		select {
		case event, ok := <-eventChan:
			if !ok {
				ws.WriteJSON(conn, pncp.EventoAnalise{Type: "done"})
				return
			}
			if err := ws.WriteJSON(conn, event); err != nil {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func (h *Hub) streamEstadoFinanceiro(ctx context.Context, conn *gorilla.Conn, uf string, exercicio int64) {
	var wg sync.WaitGroup
	ch := make(chan StreamMessage, 20)

	wg.Add(1)
	go func() {
		defer wg.Done()
		resp, _ := h.despesaPessoalUC.Executar(ctx, &dadosfinanceiros.EsferaEstadualBuscarDespesaPessoalRequest{UF: uf, Exercicio: exercicio})
		ch <- StreamMessage{Type: "despesa_pessoal", Data: resp.Dados}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		resp, _ := h.despesaCategoriaUC.Executar(ctx, &dadosfinanceiros.EsferaEstadualBuscarDespesaCategoriaRequest{UF: uf, Exercicio: exercicio})
		ch <- StreamMessage{Type: "despesa_categoria", Data: map[string]interface{}{"dados": resp.Dados}}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		resp, _ := h.rreoUC.Executar(ctx, &dadosfinanceiros.EsferaEstadualBuscarRREORequest{UF: uf, Exercicio: exercicio})
		ch <- StreamMessage{Type: "gastos_por_funcao", Data: map[string]interface{}{"dados": resp.Gastos}}
		ch <- StreamMessage{Type: "receitas", Data: map[string]interface{}{"dados": resp.Receitas}}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		resp, _ := h.recursosFederaisUC.Executar(ctx, &dadosfinanceiros.EsferaEstadualBuscarRecursosFederaisRequest{UF: uf, Exercicio: exercicio})
		ch <- StreamMessage{Type: "recursos_federais", Data: map[string]interface{}{"dados": resp.Dados}}
	}()

	go func() {
		wg.Wait()
		if h.despesaPessoalUC.SICONFIIndisponivel() {
			ch <- StreamMessage{Type: "erro", Data: map[string]string{"erro": msgSICONFIIndisponivelEstado}}
		}
		resumo := &types.DadosEstadoFinanceiroResumo{
			UF: uf,
		}
		ch <- StreamMessage{Type: "concluido", Data: resumo}
		close(ch)
	}()

	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				return
			}
			if err := ws.WriteJSON(conn, msg); err != nil {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func (h *Hub) streamMunicipioDetalhes(ctx context.Context, conn *gorilla.Conn, codigoIBGE int, exercicio int64) {
	var wg sync.WaitGroup
	ch := make(chan StreamMessage, 20)
	log := logger.New("Stream: Hub: MunicipioDetalhes")

	wg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error("panic BuscarDividaConsolidada", "recover", r)
			}
			wg.Done()
		}()
		d := h.municipioUC.BuscarDividaConsolidada(ctx, codigoIBGE, exercicio)
		ch <- StreamMessage{Type: "divida_consolidada", Data: d}
	}()

	wg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error("panic BuscarDisponibilidadeCaixa", "recover", r)
			}
			wg.Done()
		}()
		d := h.municipioUC.BuscarDisponibilidadeCaixa(ctx, codigoIBGE, exercicio)
		ch <- StreamMessage{Type: "disponibilidade_caixa", Data: d}
	}()

	wg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error("panic BuscarRestosAPagar", "recover", r)
			}
			wg.Done()
		}()
		d := h.municipioUC.BuscarRestosAPagar(ctx, codigoIBGE, exercicio)
		ch <- StreamMessage{Type: "restos_a_pagar", Data: d}
	}()

	wg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error("panic BuscarGastoSaude", "recover", r)
			}
			wg.Done()
		}()
		d := h.municipioUC.BuscarGastoSaude(ctx, codigoIBGE, exercicio)
		ch <- StreamMessage{Type: "gasto_saude", Data: d}
	}()

	wg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error("panic BuscarGastoEducacao", "recover", r)
			}
			wg.Done()
		}()
		d := h.municipioUC.BuscarGastoEducacao(ctx, codigoIBGE, exercicio)
		ch <- StreamMessage{Type: "gasto_educacao", Data: d}
	}()

	wg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error("panic BuscarFundeb", "recover", r)
			}
			wg.Done()
		}()
		d := h.municipioUC.BuscarFundeb(ctx, codigoIBGE, exercicio)
		ch <- StreamMessage{Type: "fundeb", Data: d}
	}()

	wg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error("panic BuscarBalancoPatrimonial", "recover", r)
			}
			wg.Done()
		}()
		d := h.municipioUC.BuscarBalancoPatrimonial(ctx, codigoIBGE, exercicio)
		ch <- StreamMessage{Type: "balanco_patrimonial", Data: d}
	}()

	wg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error("panic BuscarDespesasPorGrupo", "recover", r)
			}
			wg.Done()
		}()
		d := h.municipioUC.BuscarDespesasPorGrupo(ctx, codigoIBGE, exercicio)
		ch <- StreamMessage{Type: "despesas_por_grupo", Data: map[string]interface{}{"dados": d}}
	}()

	wg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error("panic BuscarTransferencias", "recover", r)
			}
			wg.Done()
		}()
		d := h.municipioUC.BuscarTransferencias(ctx, codigoIBGE, exercicio)
		ch <- StreamMessage{Type: "transferencias", Data: map[string]interface{}{"dados": d}}
	}()

	wg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error("panic BuscarContratos", "recover", r)
			}
			wg.Done()
		}()
		d := h.municipioUC.BuscarContratos(ctx, codigoIBGE, int(exercicio))
		ch <- StreamMessage{Type: "contratos", Data: map[string]interface{}{"dados": d}}
	}()

	go func() {
		wg.Wait()
		if h.municipioUC.SICONFIIndisponivel() {
			ch <- StreamMessage{Type: "erro", Data: map[string]string{"erro": msgSICONFIIndisponivelMuni}}
		}
		ch <- StreamMessage{Type: "concluido", Data: nil}
		close(ch)
	}()

	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				return
			}
			if err := ws.WriteJSON(conn, msg); err != nil {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}
