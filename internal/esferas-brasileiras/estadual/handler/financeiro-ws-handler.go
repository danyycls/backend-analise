package handler

import (
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/esferas-brasileiras/estadual/usecase/dadosfinanceiros"
	"github.com/danyele/podp/internal/shared/types"
	ws "github.com/danyele/podp/internal/shared/websocket"
)

type EsferaEstadualBuscarFinanceiroWSHandler struct {
	despesaPessoalUC   *dadosfinanceiros.EsferaEstadualBuscarDespesaPessoalUseCase
	despesaCategoriaUC *dadosfinanceiros.EsferaEstadualBuscarDespesaCategoriaUseCase
	rreoUC             *dadosfinanceiros.EsferaEstadualBuscarRREOUseCase
	recursosFederaisUC *dadosfinanceiros.EsferaEstadualBuscarRecursosFederaisUseCase
}

func NovoEsferaEstadualBuscarFinanceiroWSHandler(
	despesaPessoalUC *dadosfinanceiros.EsferaEstadualBuscarDespesaPessoalUseCase,
	despesaCategoriaUC *dadosfinanceiros.EsferaEstadualBuscarDespesaCategoriaUseCase,
	rreoUC *dadosfinanceiros.EsferaEstadualBuscarRREOUseCase,
	recursosFederaisUC *dadosfinanceiros.EsferaEstadualBuscarRecursosFederaisUseCase,
) *EsferaEstadualBuscarFinanceiroWSHandler {
	return &EsferaEstadualBuscarFinanceiroWSHandler{
		despesaPessoalUC:   despesaPessoalUC,
		despesaCategoriaUC: despesaCategoriaUC,
		rreoUC:             rreoUC,
		recursosFederaisUC: recursosFederaisUC,
	}
}

type wsMsg struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func (h *EsferaEstadualBuscarFinanceiroWSHandler) BuscarFinanceiroWS(c *gin.Context) {
	uf := strings.ToUpper(c.Param("uf"))
	if uf == "" || len(uf) != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "UF invalida"})
		return
	}

	conn, err := ws.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	ctx := c.Request.Context()
	var wg sync.WaitGroup
	ch := make(chan wsMsg, 20)

	wg.Add(1)
	go func() {
		defer wg.Done()
		resp, _ := h.despesaPessoalUC.Executar(ctx, &dadosfinanceiros.EsferaEstadualBuscarDespesaPessoalRequest{UF: uf})
		ch <- wsMsg{Type: "despesa_pessoal", Data: resp.Dados}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		resp, _ := h.despesaCategoriaUC.Executar(ctx, &dadosfinanceiros.EsferaEstadualBuscarDespesaCategoriaRequest{UF: uf})
		ch <- wsMsg{Type: "despesa_categoria", Data: map[string]interface{}{"dados": resp.Dados}}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		resp, _ := h.rreoUC.Executar(ctx, &dadosfinanceiros.EsferaEstadualBuscarRREORequest{UF: uf})
		ch <- wsMsg{Type: "gastos_por_funcao", Data: map[string]interface{}{"dados": resp.Gastos}}
		ch <- wsMsg{Type: "receitas", Data: map[string]interface{}{"dados": resp.Receitas}}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		resp, _ := h.recursosFederaisUC.Executar(ctx, &dadosfinanceiros.EsferaEstadualBuscarRecursosFederaisRequest{UF: uf})
		ch <- wsMsg{Type: "recursos_federais", Data: map[string]interface{}{"dados": resp.Dados}}
	}()

	go func() {
		wg.Wait()
		resumo := &types.DadosEstadoFinanceiroResumo{
			UF: uf,
		}
		ch <- wsMsg{Type: "concluido", Data: resumo}
		close(ch)
	}()

	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				return
			}
			ws.WriteJSON(conn, msg)
		case <-ctx.Done():
			return
		}
	}
}
