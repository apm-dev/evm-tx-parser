package common

import (
	"fmt"
	"math/big"
	"strings"
)

func HexToInt(hex string) int64 {
	hex = strings.TrimPrefix(hex, "0x")
	b := new(big.Int)
	b.SetString(hex, 16)
	return b.Int64()
}

func HexToStringInt(hex string) string {
	hex = strings.TrimPrefix(hex, "0x")
	b := new(big.Int)
	b.SetString(hex, 16)
	return b.String()
}

func IntToHex(num int) string {
	return fmt.Sprintf("0x%x", num)
}
