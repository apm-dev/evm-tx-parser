package http

import (
	"net/http"

	"github.com/apm-dev/evm-tx-parser/src/domain"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

type ParserHandler struct {
	parser    domain.Parser
	validator *validator.Validate
}

func RegisterParserHandlers(
	e *echo.Echo,
	parser domain.Parser,
) {
	handler := &ParserHandler{
		parser:    parser,
		validator: validator.New(),
	}
	e.GET("/block/last", handler.GetLastParsedBlock)
	e.POST("/address", handler.SubscribeAddress)
	e.GET("/address/:address/txs", handler.GetTxsOfAddress)
}

func (h *ParserHandler) GetLastParsedBlock(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{"blockNumber": h.parser.GetCurrentBlock()})
}

func (h *ParserHandler) SubscribeAddress(c echo.Context) error {
	var body echo.Map
	err := c.Bind(body)
	if err != nil {
		log.Errorf("failed to subscribe for address, err: %s", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "failed to subscribe for address"})
	}
	address := body["address"]
	err = h.validator.Var(address, "required,eth_addr")
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "'address' body param is required and has to be valid eth address"})
	}
	ok := h.parser.Subscribe(address.(string))
	return c.JSON(http.StatusOK, echo.Map{"message": ok})
}

func (h *ParserHandler) GetTxsOfAddress(c echo.Context) error {
	address := c.Param("address")
	err := h.validator.Var(address, "required,eth_addr")
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "'address' path param is required and has to be valid eth address"})
	}
	txs := h.parser.GetTransactions(address)
	return c.JSON(http.StatusOK, txs)
}
