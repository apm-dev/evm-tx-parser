package parser

import (
	"github.com/apm-dev/evm-tx-parser/src/config"
	"github.com/apm-dev/evm-tx-parser/src/domain"
	log "github.com/sirupsen/logrus"
)

type parser struct {
	config      *config.Config
	parserRepo  domain.ParserRepo
	ethClient   domain.EthereumClient
	txRepo      domain.TransactionRepo
	addressRepo domain.AddressRepo
}

func NewParser(
	config *config.Config,
	parserRepo domain.ParserRepo,
	ethClient domain.EthereumClient,
	txRepo domain.TransactionRepo,
	addressRepo domain.AddressRepo,
) domain.Parser {
	return &parser{
		config:      config,
		parserRepo:  parserRepo,
		ethClient:   ethClient,
		txRepo:      txRepo,
		addressRepo: addressRepo,
	}
}

// last parsed block
func (s *parser) GetCurrentBlock() int {
	num, _ := s.parserRepo.GetLastParsedBlock()
	return num
}

// add address to observer
func (s *parser) Subscribe(address string) bool {
	_, err := s.addressRepo.Save(address)
	if err != nil {
		log.Errorf("failed to subscribe address '%s', error '%s'", address, err)
		return false
	}
	return true
}

// list of inbound or outbound transactions for an address
func (s *parser) GetTransactions(address string) []domain.Transaction {
	txs, err := s.txRepo.FindByAddress(address)
	if err != nil {
		log.Errorf("failed to get transactions of address '%s', error '%s'", address, err)
		return nil
	}
	return txs
}
