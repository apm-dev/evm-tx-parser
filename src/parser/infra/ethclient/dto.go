package ethclient

import (
	"encoding/json"

	"github.com/apm-dev/evm-tx-parser/src/common"
	"github.com/apm-dev/evm-tx-parser/src/domain"
)

type (
	RpcRequest struct {
		ID      int64       `json:"id"`
		JsonRpc string      `json:"jsonrpc"`
		Method  string      `json:"method"`
		Params  interface{} `json:"params,omitempty"`
	}

	RpcResponse struct {
		ID      int64           `json:"id"`
		JsonRpc string          `json:"jsonrpc"`
		Error   *RpcError       `json:"error,omitempty"`
		Result  json.RawMessage `json:"result,omitempty"`
	}

	RpcError struct {
		Code    int    `json:"code"`
		Data    string `json:"data"`
		Message string `json:"message"`
	}

	Transaction struct {
		BlockNumber string `json:"blockNumber"`
		From        string `json:"from"`
		To          string `json:"to"`
		Hash        string `json:"hash"`
		Nonce       string `json:"nonce"`
		Value       string `json:"value"`
	}

	Block struct {
		Hash         string        `json:"hash"`
		Number       string        `json:"number"`
		ParentHash   string        `json:"parentHash"`
		Transactions []Transaction `json:"transactions"`
	}
)

func (b *Block) ToEntity() *domain.Block {
	return &domain.Block{
		Number:       int(common.HexToInt(b.Number)),
		Hash:         b.Hash,
		ParentHash:   b.ParentHash,
		Transactions: TxsToEntity(b.Transactions),
	}
}

func (t *Transaction) ToEntity() *domain.Transaction {
	return &domain.Transaction{
		Hash:        t.Hash,
		Value:       common.HexToStringInt(t.Value),
		From:        t.From,
		To:          t.To,
		Nonce:       int(common.HexToInt(t.Nonce)),
		BlockNumber: int(common.HexToInt(t.BlockNumber)),
	}
}

func TxsToEntity(txs []Transaction) []domain.Transaction {
	result := make([]domain.Transaction, 0, len(txs))
	for _, tx := range txs {
		result = append(result, *tx.ToEntity())
	}
	return result
}
