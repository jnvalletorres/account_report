package model

import (
	"account_report/domain/entity"
	"fmt"
)

type TransactionResumeModel struct {
	Date              string
	NoTotal           string
	AmountDebitTotal  string
	AmountCreditTotal string
}

func ToTransactionResumeModel(e entity.TransactionResumeEntity) TransactionResumeModel {
	const DATE_FORMAT = "2006-Jan"

	return TransactionResumeModel{
		Date:              e.Date.Format(DATE_FORMAT),
		NoTotal:           fmt.Sprintf("%.2f", e.NoTotal),
		AmountDebitTotal:  fmt.Sprintf("$%.2f", e.AmountDebitAverage),
		AmountCreditTotal: fmt.Sprintf("$%.2f", e.AmountCreditAverage),
	}
}
