package structs

import (
	"sync"
)

type Client struct {
	ID        int
	Account   *BankAccount
	TransChan chan Transaction
	TransBool bool
}

type Transaction struct {
	Amount  int
	IsDebit bool // if true then - else +
}

type BankAccount struct {
	sync.Mutex
	Balance int
}
