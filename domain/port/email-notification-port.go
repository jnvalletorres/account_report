package port

import "account_report/domain/entity"

type EmailNotificationPort interface {
	SendAccountReport(entity.ResumeEntity) error
}
