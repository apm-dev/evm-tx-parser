package domain

type Address struct {
	ID      uint
	Address string
}

type AddressRepo interface {
	Save(address string) (id int, err error)
	Exist(address string) bool
}
