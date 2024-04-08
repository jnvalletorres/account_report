package port

import "account_report/domain/entity"

type InputRepositoryPort interface {
	GetInput() (entity.InputEntity, error)
}
