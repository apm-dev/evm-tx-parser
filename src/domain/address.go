package domain

import "strings"

type AddressRepo interface {
	Save(address string) error
	Exist(address string) bool
}

func NormalizeAddress(address string) string {
	return strings.ToLower(address)
}
