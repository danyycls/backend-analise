package handler

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/esferas-brasileiras/municipal/usecase"
	"github.com/danyele/podp/internal/shared/logger"
	ws "github.com/danyele/podp/internal/shared/websocket"
)

const msgSICONFIIndisponivel = "API SICONFI (Tesouro Nacional) temporariamente indisponível. Os dados financeiros do município não puderam ser carregados. Tente novamente mais tarde."

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
		defer func() {
			if r := recover(); r != nil {
				logger.New("Municipal: WS Handler").Error("panic BuscarDividaConsolidada", "recover", r)
			}
			wg.Done()
		}()
		d := h.useCase.BuscarDividaConsolidada(ctx, codigoIBGE, exercicio)
		ch <- wsMsgMuni{Type: "divida_consolidada", Data: d}
	}()

	wg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.New("Municipal: WS Handler").Error("panic BuscarDisponibilidadeCaixa", "recover", r)
			}
			wg.Done()
		}()
		d := h.useCase.BuscarDisponibilidadeCaixa(ctx, codigoIBGE, exercicio)
		ch <- wsMsgMuni{Type: "disponibilidade_caixa", Data: d}
	}()

	wg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.New("Municipal: WS Handler").Error("panic BuscarRestosAPagar", "recover", r)
			}
			wg.Done()
		}()
		d := h.useCase.BuscarRestosAPagar(ctx, codigoIBGE, exercicio)
		ch <- wsMsgMuni{Type: "restos_a_pagar", Data: d}
	}()

	wg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.New("Municipal: WS Handler").Error("panic BuscarGastoSaude", "recover", r)
			}
			wg.Done()
		}()
		d := h.useCase.BuscarGastoSaude(ctx, codigoIBGE, exercicio)
		ch <- wsMsgMuni{Type: "gasto_saude", Data: d}
	}()

	wg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.New("Municipal: WS Handler").Error("panic BuscarGastoEducacao", "recover", r)
			}
			wg.Done()
		}()
		d := h.useCase.BuscarGastoEducacao(ctx, codigoIBGE, exercicio)
		ch <- wsMsgMuni{Type: "gasto_educacao", Data: d}
	}()

	wg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.New("Municipal: WS Handler").Error("panic BuscarFundeb", "recover", r)
			}
			wg.Done()
		}()
		d := h.useCase.BuscarFundeb(ctx, codigoIBGE, exercicio)
		ch <- wsMsgMuni{Type: "fundeb", Data: d}
	}()

	wg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.New("Municipal: WS Handler").Error("panic BuscarBalancoPatrimonial", "recover", r)
			}
			wg.Done()
		}()
		d := h.useCase.BuscarBalancoPatrimonial(ctx, codigoIBGE, exercicio)
		ch <- wsMsgMuni{Type: "balanco_patrimonial", Data: d}
	}()

	wg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.New("Municipal: WS Handler").Error("panic BuscarDespesasPorGrupo", "recover", r)
			}
			wg.Done()
		}()
		d := h.useCase.BuscarDespesasPorGrupo(ctx, codigoIBGE, exercicio)
		ch <- wsMsgMuni{Type: "despesas_por_grupo", Data: map[string]interface{}{"dados": d}}
	}()

	wg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.New("Municipal: WS Handler").Error("panic BuscarTransferencias", "recover", r)
			}
			wg.Done()
		}()
		d := h.useCase.BuscarTransferencias(ctx, codigoIBGE, exercicio)
		ch <- wsMsgMuni{Type: "transferencias", Data: map[string]interface{}{"dados": d}}
	}()

	wg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.New("Municipal: WS Handler").Error("panic BuscarContratos", "recover", r)
			}
			wg.Done()
		}()
		d := h.useCase.BuscarContratos(ctx, codigoIBGE, int(exercicio))
		ch <- wsMsgMuni{Type: "contratos", Data: map[string]interface{}{"dados": d}}
	}()

	go func() {
		wg.Wait()
		if h.useCase.SICONFIIndisponivel() {
			ch <- wsMsgMuni{Type: "erro", Data: map[string]string{"erro": msgSICONFIIndisponivel}}
		}
		ch <- wsMsgMuni{Type: "concluido", Data: nil}
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
