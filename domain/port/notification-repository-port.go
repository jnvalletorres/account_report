package port

import (
	"account_report/domain/entity"
)

type NotificationRepositoryPort interface {
	SaveNotification(notification entity.NotificationEntity) (string, error)
	FindById(id string) (*entity.NotificationEntity, error)
	UpdateById(id string, notificationEntity entity.NotificationEntity) error
}
