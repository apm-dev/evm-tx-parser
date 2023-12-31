package parser

import (
	"context"
	"math"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/apm-dev/evm-tx-parser/src/config"
	"github.com/apm-dev/evm-tx-parser/src/domain"
	log "github.com/sirupsen/logrus"
)

type parser struct {
	config      *config.Config
	once        *sync.Once
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
		once:        &sync.Once{},
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
	err := s.addressRepo.Save(domain.NormalizeAddress(address))
	if err != nil {
		log.Errorf("failed to subscribe address '%s', error '%s'", address, err)
		return false
	}
	return true
}

// list of inbound or outbound transactions for an address
func (s *parser) GetTransactions(address string) []domain.Transaction {
	txs, err := s.txRepo.FindByAddress(domain.NormalizeAddress(address))
	if err != nil {
		log.Errorf("failed to get transactions of address '%s', error '%s'", address, err)
		return nil
	}
	return txs
}

func (s *parser) Start(ctx context.Context) {
	s.once.Do(func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				lastBlockNum, lastBlockHash := s.parserRepo.GetLastParsedBlock()
				currentBlock, err := s.ethClient.GetNowBlockNumber()
				if err != nil {
					log.Errorf("failed to get current block, error '%s'", err)
					s.sleepOneBlockTime()
					continue
				}
				batchSize := s.howManyBlocksShouldFetch(lastBlockNum, currentBlock)
				if batchSize <= 0 {
					s.sleepOneBlockTime()
					continue
				}
				// Fetch range of blocks, use batch request to:
				// Reduced network latency
				// Improved performance
				// Atomicity
				fromBlock := lastBlockNum + 1
				toBlock := lastBlockNum + batchSize
				blocksRange, err := s.ethClient.GetBlocksByRange(fromBlock, toBlock)
				if err != nil {
					log.Errorf("failed to get blocks '%d -> %d', error '%s'", fromBlock, toBlock, err)
					s.sleepOneBlockTime()
					continue
				}
				if len(blocksRange) == 0 {
					log.Warnf("no blocks to parse '%d -> %d'", fromBlock, toBlock)
					s.sleepOneBlockTime()
					continue
				}
				// Sort Blocks by their Number in ASC order
				blocks := domain.Blocks(blocksRange)
				sort.Sort(blocks)
				// Stop parsing in case of reorgs (orphan blocks)
				if s.isThereOrphanBlock(lastBlockNum, lastBlockHash, blocks[0]) {
					log.Errorf("detect Orphan blocks on block '%d':'%s'", lastBlockNum, lastBlockHash)
					return
				}
				// Parse blocks concurrently to extract related transactions
				relatedTxs := make([]domain.Transaction, 0)
				resultChan := make(chan []domain.Transaction, len(blocks))
				for _, block := range blocks {
					go s.extractRelatedTxs(block, resultChan)
				}
				for i := 0; i < len(blocks); i++ {
					relatedTxs = append(relatedTxs, <-resultChan...)
				}
				// Store Txs and update last parsed block
				// better to do both in one DB (ACID) transaction
				err = s.txRepo.SaveMany(relatedTxs)
				if err != nil {
					log.Errorf("failed to store transactions '%d' of blocks '%d -> %d', error '%s'", len(relatedTxs), fromBlock, toBlock, err)
					continue
				}
				lastBlock := blocks[len(blocks)-1]
				err = s.parserRepo.UpdateLastParsedBlock(lastBlock.Number, lastBlock.Hash)
				if err != nil {
					log.Errorf("failed to update last parsed block '%d':'%s', error '%s'", lastBlock.Number, lastBlock.Hash, err)
					continue
				}
				log.Infof("Parsed '%d -> %d' '%d' blocks, detect '%d' txs", fromBlock, toBlock, len(blocks), len(relatedTxs))
			}
		}
	})
}

func (s *parser) sleepOneBlockTime() {
	time.Sleep(s.config.App.NetworkBlockTime)
}

func (s *parser) howManyBlocksShouldFetch(lastBlock, currentBlock int) int {
	// Stay a few blocks behind the network's head to prevent facing with reorgs (orphan blocks)
	safeBlock := currentBlock - s.config.App.OrphanPreventionBlockCount
	diff := safeBlock - lastBlock
	if diff <= 0 {
		return 0
	}
	return int(math.Min(float64(diff), float64(s.config.App.GetBlocksBatchSize)))
}

func (s *parser) isThereOrphanBlock(lastBlockNum int, lastBlockHash string, nextBlock domain.Block) bool {
	return nextBlock.Number != lastBlockNum+1 || !strings.EqualFold(nextBlock.ParentHash, lastBlockHash)
}

func (s *parser) extractRelatedTxs(b domain.Block, result chan<- []domain.Transaction) {
	txs := make([]domain.Transaction, 0)
	for _, tx := range b.Transactions {
		if s.addressRepo.Exist(domain.NormalizeAddress(tx.From)) || s.addressRepo.Exist(domain.NormalizeAddress(tx.To)) {
			txs = append(txs, tx)
		}
	}
	result <- txs
}
