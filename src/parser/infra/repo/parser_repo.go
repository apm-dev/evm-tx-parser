package repo

import (
	"github.com/apm-dev/evm-tx-parser/src/domain"
)

type parserRepo struct {
	data *lastParsedBlock
}

func NewParserRepo() domain.ParserRepo {
	return &parserRepo{
		data: &lastParsedBlock{},
	}
}

type lastParsedBlock struct {
	num  int
	hash string
}

func (r *parserRepo) GetLastParsedBlock() (int, string) {
	return r.data.num, r.data.hash
}

func (r *parserRepo) UpdateLastParsedBlock(num int, hash string) error {
	r.data.num, r.data.hash = num, hash
	return nil
}
