package domain

import "context"

type Parser interface {
	// last parsed block
	GetCurrentBlock() int

	// add address to observer
	Subscribe(address string) bool

	// list of inbound or outbound transactions for an address
	GetTransactions(address string) []Transaction

	// start reading the network and parsing blocks
	Start(ctx context.Context)
}

type EthereumClient interface {
	GetNowBlockNumber() (int, error)
	GetBlocksByRange(from, to int) ([]Block, error)
}

type ParserRepo interface {
	GetLastParsedBlock() (int, string)
	UpdateLastParsedBlock(num int, hash string) error
}
