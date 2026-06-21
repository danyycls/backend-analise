package handler

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/esferas-brasileiras/municipal/usecase"
	"github.com/danyele/podp/internal/shared/types"
	ws "github.com/danyele/podp/internal/shared/websocket"
)

type EsferaMunicipalBuscarDetalhesWSHandler struct {
	useCase *usecase.EsferaMunicipalBuscarDetalhesUseCase
}

func NovoEsferaMunicipalBuscarDetalhesWSHandler(useCase *usecase.EsferaMunicipalBuscarDetalhesUseCase) *EsferaMunicipalBuscarDetalhesWSHandler {
	return &EsferaMunicipalBuscarDetalhesWSHandler{useCase: useCase}
}

type wsMsgMuni struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func (h *EsferaMunicipalBuscarDetalhesWSHandler) BuscarDetalhesMunicipioWS(c *gin.Context) {
	codigoIBGEStr := c.Param("codigoIBGE")
	codigoIBGE, err := strconv.Atoi(codigoIBGEStr)
	if err != nil || codigoIBGE <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "codigo IBGE invalido"})
		return
	}

	exercicioStr := c.DefaultQuery("exercicio", "0")
	exercicio, _ := strconv.ParseInt(exercicioStr, 10, 64)

	nome := c.Query("nome")
	uf := c.Query("uf")

	conn, err := ws.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	ctx := c.Request.Context()
	var wg sync.WaitGroup
	ch := make(chan wsMsgMuni, 20)

	wg.Add(1)
	go func() {
		defer wg.Done()
		despesa := h.useCase.BuscarDespesaPessoal(ctx, codigoIBGE, exercicio)
		ch <- wsMsgMuni{Type: "despesa_pessoal", Data: despesa}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		gastos, receitas := h.useCase.BuscarRREO(ctx, codigoIBGE, exercicio)
		ch <- wsMsgMuni{Type: "gastos_por_funcao", Data: map[string]interface{}{"dados": gastos}}
		ch <- wsMsgMuni{Type: "receitas", Data: map[string]interface{}{"dados": receitas}}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		recursos := h.useCase.BuscarRecursosFederais(ctx, codigoIBGE)
		ch <- wsMsgMuni{Type: "recursos_federais", Data: map[string]interface{}{"dados": recursos}}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		contratos := h.useCase.BuscarContratos(ctx, codigoIBGE)
		ch <- wsMsgMuni{Type: "contratos", Data: map[string]interface{}{"dados": contratos}}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		servidores := h.useCase.BuscarServidores(ctx, codigoIBGE)
		ch <- wsMsgMuni{Type: "servidores", Data: map[string]interface{}{"dados": servidores}}
	}()

	go func() {
		wg.Wait()
		resumo := &types.DetalhesMunicipioResponse{
			CodigoIBGE: codigoIBGE,
			Nome:       nome,
			UF:         uf,
		}
		if exercicio > 0 {
			resumo.Exercicio = int(exercicio)
		}
		ch <- wsMsgMuni{Type: "concluido", Data: resumo}
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
