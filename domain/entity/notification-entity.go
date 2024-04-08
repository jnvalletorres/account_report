package entity

import "time"

type NotificationEntity struct {
	ID           string              `json:"id"`
	SendStatus   string              `json:"send_status"`
	Accout       AccountEntity       `json:"account"`
	Transactions []TransactionEntity `json:"transactions"`
	Sources      []string            `json:"source"`
	UpdateAt     time.Time           `json:"update_at"`
	CreateAt     time.Time           `json:"create_at"`
}
