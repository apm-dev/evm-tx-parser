package common

import (
	"math/big"
	"strings"
)

func HexToInt(hex string) (string, error) {
	hex = strings.TrimPrefix(hex, "0x")
	b := new(big.Int)
	b.SetString(hex, 16)
	return b.String(), nil
}
