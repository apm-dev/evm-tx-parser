package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	_addrRepo "github.com/apm-dev/evm-tx-parser/src/address/infra/repo"
	"github.com/apm-dev/evm-tx-parser/src/common"
	"github.com/apm-dev/evm-tx-parser/src/config"
	"github.com/apm-dev/evm-tx-parser/src/parser"
	_parserHandler "github.com/apm-dev/evm-tx-parser/src/parser/delivery/http"
	"github.com/apm-dev/evm-tx-parser/src/parser/infra/ethclient"
	_parserRepo "github.com/apm-dev/evm-tx-parser/src/parser/infra/repo"
	_txRepo "github.com/apm-dev/evm-tx-parser/src/transaction/infra/repo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

func main() {
	config := config.NewConfig()

	logLevel, err := logrus.ParseLevel(config.App.LogLevel)
	if err != nil {
		panic(err)
	}
	logrus.SetLevel(logLevel)

	e := echo.New()
	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	// Healthz
	e.GET("/healthz", func(c echo.Context) error {
		// For liveness check, simply return HTTP 200 OK
		// indicating that the application is live and healthy.

		// For readiness check, we can perform additional checks
		// such as database connectivity, external service dependencies,
		// or other necessary components, and return HTTP 200 OK only
		// when the application is ready to serve traffic.
		return c.NoContent(http.StatusOK)
	})

	// Repos
	addrRepo := _addrRepo.NewAddressRepo()
	txRepo := _txRepo.NewTransactionRepo()
	parserRepo := _parserRepo.NewParserRepo()

	ethClient := ethclient.NewEthClient(config, "https://rpc.ankr.com/eth")

	// Set parser's starting point (block)
	err = parserRepo.UpdateLastParsedBlock(config.App.DefaultStartingBlockNum, config.App.DefaultStartingBlockHash)
	if err != nil {
		panic(err)
	}

	// Services
	parser := parser.NewParser(config, parserRepo, ethClient, txRepo, addrRepo)

	// Start Parser
	ctx := context.Background()
	ctxParser, cancelParser := context.WithCancel(ctx)
	defer cancelParser()
	go parser.Start(ctxParser)

	// Start Server
	_parserHandler.RegisterParserHandlers(e, parser)

	go func() {
		if err := e.Start(fmt.Sprintf(":%d", config.App.WebPort)); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	common.WaitForSignal()

	ctxEcho, cancelEcho := context.WithTimeout(ctx, 10*time.Second)
	defer cancelEcho()
	if err := e.Shutdown(ctxEcho); err != nil {
		logrus.Fatal(err)
	}
}
