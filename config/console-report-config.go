package config

import (
	"account_report/adapter/common"
	"account_report/adapter/logger"
	"account_report/adapter/notification"
	"account_report/adapter/repository"
	"account_report/domain/port"
	"account_report/env"
	"account_report/usecase"
	"strconv"
)

type ConsoleReportConfig interface {
	CreateReporUseCase() usecase.ReportUseCase
}

type consoleReportConfig struct {
	notificationRepositoryPort port.NotificationRepositoryPort
	inputRepositoryPort        port.InputRepositoryPort
	emailNotificationPort      port.EmailNotificationPort
	loggerPort                 port.LoggerPort
}

func NewConsoleReportConfig() ConsoleReportConfig {
	var err error
	logger := logger.NewConsoleLogger(env.Environment)
	emailSender := common.NewSender(
		env.Environment.Smtp.Host,
		env.Environment.Smtp.UserName,
		env.Environment.Smtp.Password,
		env.Environment.Smtp.PortNumber,
	)
	timeOut, err := strconv.Atoi(env.Environment.MongoDB.TimeOut)
	if err != nil {
		timeOut = 5
	}
	mongoDb, err := repository.NewMongoNotificationRepository(
		env.Environment.MongoDB.Url,
		env.Environment.MongoDB.Name,
		timeOut,
		logger,
	)
	if err != nil {
		panic(err)
	}

	return &consoleReportConfig{
		inputRepositoryPort:        repository.NewFileInputRepository(logger, env.Environment),
		emailNotificationPort:      notification.NewSmtpEmailNotification(emailSender, logger),
		loggerPort:                 logger,
		notificationRepositoryPort: mongoDb,
	}
}

func (rc *consoleReportConfig) CreateReporUseCase() usecase.ReportUseCase {
	return usecase.NewReportUseCase(rc.notificationRepositoryPort, rc.inputRepositoryPort, rc.emailNotificationPort, rc.loggerPort)
}
