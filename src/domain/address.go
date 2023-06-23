package domain

type AddressRepo interface {
	Save(address string) error
	Exist(address string) bool
}
