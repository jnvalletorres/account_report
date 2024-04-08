package notification

import (
	"account_report/adapter/common"
	"account_report/adapter/notification/model"
	"account_report/domain/entity"
	"account_report/domain/port"
	"path/filepath"
)

type smtpEmailNotification struct {
	sender *common.Sender
	logger port.LoggerPort
}

func NewSmtpEmailNotification(sender *common.Sender, logger port.LoggerPort) port.EmailNotificationPort {
	return &smtpEmailNotification{sender: sender, logger: logger}
}

func (n *smtpEmailNotification) SendAccountReport(resumeEntity entity.ResumeEntity) error {
	const (
		DIRECTORY     = ".resources"
		TEMPLATE_NAME = "account-repport-template.html"
	)

	message, err := common.NewHtmlMessage(
		resumeEntity.Subject,
		model.ToResumeModel(resumeEntity),
		filepath.Join(DIRECTORY, TEMPLATE_NAME))

	if err != nil {
		n.logger.Error("Error to create smtp message", err)
		return err
	}

	message.To = resumeEntity.To
	message.Cc = resumeEntity.Cc
	message.Bcc = resumeEntity.Bcc

	for _, attachFile := range resumeEntity.AttachFiles {
		message.AddAttachFile(attachFile)
	}

	err = n.sender.Send(message)
	if err != nil {
		n.logger.Error("Error to send smtp message", err)
	}

	return nil
}
