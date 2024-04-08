package document

import "account_report/domain/entity"

type AccountDocument struct {
	Id    int    `json:"id" bson:"id"`
	Name  string `json:"name" bson:"name"`
	Email string `json:"email" bson:"email"`
}

func ToAccountDocument(accaountEntity entity.AccountEntity) *AccountDocument {
	return &AccountDocument{
		Id:    accaountEntity.Id,
		Name:  accaountEntity.Name,
		Email: accaountEntity.Email,
	}
}

func ToAccountEntity(accaountCocument AccountDocument) *entity.AccountEntity {
	return &entity.AccountEntity{
		Id:    accaountCocument.Id,
		Name:  accaountCocument.Name,
		Email: accaountCocument.Email,
	}
}
