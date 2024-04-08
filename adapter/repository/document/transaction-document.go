package document

import (
	"account_report/domain/entity"
	"time"
)

type TransactionDocument struct {
	Id       int       `json:"id" bson:"id"`
	Date     time.Time `json:"date" bson:"date"`
	CardType string    `json:"card_type" bson:"card_type"`
	Amount   float64   `json:"amount" bson:"amount"`
}

func ToTransactionDocument(transactionEntity entity.TransactionEntity) TransactionDocument {
	return TransactionDocument{
		Id:       transactionEntity.Id,
		Date:     transactionEntity.Date,
		CardType: transactionEntity.CardType,
		Amount:   transactionEntity.Amount,
	}
}

func ToTransactionEntity(transactionDocument TransactionDocument) entity.TransactionEntity {
	return entity.TransactionEntity{
		Id:       transactionDocument.Id,
		Date:     transactionDocument.Date,
		CardType: transactionDocument.CardType,
		Amount:   transactionDocument.Amount,
	}
}
