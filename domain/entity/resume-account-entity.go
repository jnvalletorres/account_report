package entity

import "time"

type ResumeEntity struct {
	To                  []string                  `json:"to"`
	Cc                  []string                  `json:"cc"`
	Bcc                 []string                  `json:"bcc"`
	Subject             string                    `json:"subject"`
	AttachFiles         []string                  `json:"attach_files"`
	Name                string                    `json:"name"`
	Now                 time.Time                 `json:"now"`
	BalanceTotal        float64                   `json:"balance_total"`
	TraansactionsResume []TransactionResumeEntity `json:"transactions_resume"`
}
