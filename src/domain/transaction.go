package domain

type Transaction struct {
	Hash        string `json:"hash"`
	Value       string `json:"value"`
	From        string `json:"from"`
	To          string `json:"to"`
	Nonce       int    `json:"nonce"`
	BlockNumber int    `json:"block_number"`
}

type Block struct {
	Number       int           `json:"number"`
	Hash         string        `json:"hash"`
	ParentHash   string        `json:"parent_hash"`
	Transactions []Transaction `json:"transactions"`
}

type TransactionRepository interface {
	SaveMany(txs []Transaction) error
	FindByAddress(address string) ([]Transaction, error)
}
