package entity

import "time"

type TransactionEntity struct {
	Id       int       `json:"id"`
	Date     time.Time `json:"date"`
	CardType string    `json:"card_type"`
	Amount   float64   `json:"amount"`
}
