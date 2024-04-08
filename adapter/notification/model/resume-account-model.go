package model

import (
	"account_report/domain/entity"
	"fmt"
	"time"
)

type ResumeModel struct {
	Title               string
	Name                string
	Now                 string
	BalanceTotal        string
	TraansactionsResume []TransactionResumeModel
}

func ToResumeModel(e entity.ResumeEntity) ResumeModel {

	transactions := make([]TransactionResumeModel, 0)

	for _, v := range e.TraansactionsResume {
		transactions = append(transactions, ToTransactionResumeModel(v))
	}

	return ResumeModel{
		Name:                e.Name,
		Now:                 e.Now.Format(time.RFC822),
		BalanceTotal:        fmt.Sprintf("$%.2f", e.BalanceTotal),
		TraansactionsResume: transactions,
	}
}
