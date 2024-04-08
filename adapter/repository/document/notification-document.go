package document

import (
	"account_report/domain/entity"
	"time"
)

type NotificationDocument struct {
	ID           string                `json:"id" bson:"_id,omitempty"`
	SendStatus   string                `json:"send_status" bson:"send_status"`
	Accout       AccountDocument       `json:"account" bson:"account"`
	Transactions []TransactionDocument `json:"transactions" bson:"transactions"`
	Sources      []string              `json:"sources" bson:"sources"`
	UpdateAt     time.Time             `json:"update_at" bson:"update_at,omitempty"`
	CreateAt     time.Time             `json:"create_at" bson:"create_at,omitempty"`
}

func ToCreateNotificationDocument(notificationEntity entity.NotificationEntity) *NotificationDocument {

	transactions := []TransactionDocument{}
	for _, t := range notificationEntity.Transactions {
		transactions = append(transactions, ToTransactionDocument(t))
	}

	return &NotificationDocument{
		SendStatus:   notificationEntity.SendStatus,
		Accout:       *ToAccountDocument(notificationEntity.Accout),
		Transactions: transactions,
		Sources:      notificationEntity.Sources,
		CreateAt:     time.Now().UTC(),
		UpdateAt:     time.Now().UTC(),
	}
}

func ToUpdateNotificationDocument(notificationEntity entity.NotificationEntity) *NotificationDocument {

	transactions := []TransactionDocument{}
	for _, t := range notificationEntity.Transactions {
		transactions = append(transactions, ToTransactionDocument(t))
	}

	return &NotificationDocument{
		SendStatus:   notificationEntity.SendStatus,
		Accout:       *ToAccountDocument(notificationEntity.Accout),
		Transactions: transactions,
		Sources:      notificationEntity.Sources,
		UpdateAt:     time.Now().UTC(),
	}
}

func ToNotificationEntity(notificationDocument NotificationDocument) *entity.NotificationEntity {

	transactions := []entity.TransactionEntity{}
	for _, t := range notificationDocument.Transactions {
		transactions = append(transactions, ToTransactionEntity(t))
	}

	return &entity.NotificationEntity{
		ID:           notificationDocument.ID,
		SendStatus:   notificationDocument.SendStatus,
		Accout:       *ToAccountEntity(notificationDocument.Accout),
		Transactions: transactions,
		Sources:      notificationDocument.Sources,
		CreateAt:     time.Now().UTC(),
		UpdateAt:     time.Now().UTC(),
	}
}
