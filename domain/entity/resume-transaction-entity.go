package entity

import "time"

type TransactionResumeEntity struct {
	Date                time.Time `json:"date"`
	NoTotal             float64   `json:"no_total"`
	AmountDebitTotal    float64   `json:"amount_debit_total"`
	AmountCreditTotal   float64   `json:"amount_credit_total"`
	AmountDebitAverage  float64   `json:"amount_debit_average"`
	AmountCreditAverage float64   `json:"amount_credit_average"`
}
