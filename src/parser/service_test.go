package parser_test

import (
	"context"
	"testing"
	"time"

	"github.com/apm-dev/evm-tx-parser/src/config"
	"github.com/apm-dev/evm-tx-parser/src/domain"
	"github.com/apm-dev/evm-tx-parser/src/domain/mocks"
	"github.com/apm-dev/evm-tx-parser/src/parser"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

func mockDependencies() (
	*mocks.ParserRepo,
	*mocks.EthereumClient,
	*mocks.TransactionRepo,
	*mocks.AddressRepo,
) {
	return new(mocks.ParserRepo),
		new(mocks.EthereumClient),
		new(mocks.TransactionRepo),
		new(mocks.AddressRepo)
}

func Test_Start(t *testing.T) {
	config := config.NewConfig()
	config.App.GetBlocksBatchSize = 3

	logrus.SetLevel(logrus.WarnLevel)

	noTxGetBlocksResult := []domain.Block{
		{Number: 101, Hash: "0x101", ParentHash: "0x100"},
		{Number: 102, Hash: "0x102", ParentHash: "0x101"},
		{Number: 103, Hash: "0x103", ParentHash: "0x102"},
	}

	t.Run("should fetch as many blocks as max block batch size config in a batch request, even there are more blocks to fetch", func(t *testing.T) {
		mockParserRepo, mockEthClient, mockTxRepo, mockAddressRepo := mockDependencies()
		// Arrange
		mockParserRepo.On("GetLastParsedBlock").
			Return(100, "0x100")
		mockEthClient.On("GetNowBlockNumber").
			Return(200, nil)
		mockEthClient.On("GetBlocksByRange", 101, 100+config.App.GetBlocksBatchSize).
			Return(noTxGetBlocksResult, nil)
		mockTxRepo.On("SaveMany", []domain.Transaction{}).
			Return(nil)
		mockParserRepo.On("UpdateLastParsedBlock", 103, "0x103").
			Return(nil)

		parser := parser.NewParser(config, mockParserRepo, mockEthClient, mockTxRepo, mockAddressRepo)
		// Action
		ctx, cancel := context.WithCancel(context.Background())
		go parser.Start(ctx)
		// Wait for a brief moment to allow some iterations of the loop.
		time.Sleep(100 * time.Millisecond)
		// Cancel the context to stop the function.
		cancel()
		// Wait for the function to exit.
		time.Sleep(100 * time.Millisecond)
		// Assert
		mockParserRepo.AssertExpectations(t)
		mockEthClient.AssertExpectations(t)
		mockTxRepo.AssertExpectations(t)
	})

	t.Run("should sleep and wait when there is no block to fetch", func(t *testing.T) {
		// Arrange
		mockParserRepo, mockEthClient, mockTxRepo, mockAddressRepo := mockDependencies()

		mockParserRepo.On("GetLastParsedBlock").
			Return(100, "0x100")
		mockEthClient.On("GetNowBlockNumber").
			Return(100, nil)

		parser := parser.NewParser(config, mockParserRepo, mockEthClient, mockTxRepo, mockAddressRepo)
		// Action
		ctx, cancel := context.WithCancel(context.Background())
		go parser.Start(ctx)
		// Wait for a brief moment to allow some iterations of the loop.
		time.Sleep(100 * time.Millisecond)
		// Cancel the context to stop the function.
		cancel()
		// Wait for the function to exit.
		time.Sleep(100 * time.Millisecond)
		// Assert
		mockParserRepo.AssertExpectations(t)
		mockEthClient.AssertExpectations(t)
		mockTxRepo.AssertNumberOfCalls(t, "SaveMany", 0)
		mockEthClient.AssertNumberOfCalls(t, "GetBlocksByRange", 0)
		mockParserRepo.AssertNumberOfCalls(t, "UpdateLastParsedBlock", 0)
	})

	t.Run("should stop parsing when detect orphan blocks", func(t *testing.T) {
		// Arrange
		mockParserRepo, mockEthClient, mockTxRepo, mockAddressRepo := mockDependencies()

		getBlocksResultWithOrphan := []domain.Block{
			{Number: 101, Hash: "0x101", ParentHash: "0x4234123"},
			{Number: 102, Hash: "0x102", ParentHash: "0x101"},
			{Number: 103, Hash: "0x103", ParentHash: "0x102"},
		}
		mockParserRepo.On("GetLastParsedBlock").
			Return(100, "0x100")
		mockEthClient.On("GetNowBlockNumber").
			Return(200, nil)
		mockEthClient.On("GetBlocksByRange", mock.Anything, mock.Anything).
			Return(getBlocksResultWithOrphan, nil)

		parser := parser.NewParser(config, mockParserRepo, mockEthClient, mockTxRepo, mockAddressRepo)
		// Action
		ctx, cancel := context.WithCancel(context.Background())
		go parser.Start(ctx)
		// Wait for a brief moment to allow some iterations of the loop.
		time.Sleep(100 * time.Millisecond)
		// Cancel the context to stop the function.
		cancel()
		// Wait for the function to exit.
		time.Sleep(100 * time.Millisecond)
		// Assert
		mockParserRepo.AssertExpectations(t)
		mockEthClient.AssertExpectations(t)
		mockTxRepo.AssertNumberOfCalls(t, "SaveMany", 0)
		mockAddressRepo.AssertNumberOfCalls(t, "Exist", 0)
		mockParserRepo.AssertNumberOfCalls(t, "UpdateLastParsedBlock", 0)
	})

	// These are POC unit-tests, of course there are lot more to write ;)
}
