package entity

type InputEntity struct {
	Accout       AccountEntity       `json:"account"`
	Transactions []TransactionEntity `json:"transactions"`
	Sources      []string            `json:"sources"`
}
