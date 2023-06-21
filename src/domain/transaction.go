package domain

type Transaction struct {
	Hash        string      `json:"hash"`
	Value       string      `json:"value"`
	From        string      `json:"from"`
	To          string      `json:"to"`
	Nonce       int         `json:"nonce"`
	BlockNumber int         `json:"block_number"`
	Direction   TxDirection `json:"direction,omitempty"`
}

type TxDirection string

const (
	Incoming TxDirection = "incoming"
	Outgoing TxDirection = "outgoing"
)

type Block struct {
	Number       int           `json:"number"`
	Hash         string        `json:"hash"`
	ParentHash   string        `json:"parent_hash"`
	Transactions []Transaction `json:"transactions"`
}

type Blocks []Block

// implement sort.Interface
func (b Blocks) Len() int           { return len(b) }
func (b Blocks) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b Blocks) Less(i, j int) bool { return b[i].Number < b[j].Number }

type TransactionRepo interface {
	SaveMany(txs []Transaction) error
	FindByAddress(address string) ([]Transaction, error)
}
