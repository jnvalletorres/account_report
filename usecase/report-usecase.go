package usecase

import (
	"account_report/domain/entity"
	"account_report/domain/enum"
	"account_report/domain/port"
	"fmt"
	"slices"
	"time"
)

type ReportUseCase interface {
	SendAccoutTransactionReport() error
}

type reportUseCase struct {
	notificationRepositoryPort port.NotificationRepositoryPort
	inputRepositoryPort        port.InputRepositoryPort
	emailNotificationPort      port.EmailNotificationPort
	loggerPort                 port.LoggerPort
}

func NewReportUseCase(
	notificationRepositoryPort port.NotificationRepositoryPort,
	inputRepositoryPort port.InputRepositoryPort,
	emailManagerPort port.EmailNotificationPort,
	loggerPort port.LoggerPort,
) ReportUseCase {
	instance := &reportUseCase{
		notificationRepositoryPort: notificationRepositoryPort,
		inputRepositoryPort:        inputRepositoryPort,
		emailNotificationPort:      emailManagerPort,
		loggerPort:                 loggerPort,
	}

	return instance
}

func (r *reportUseCase) SendAccoutTransactionReport() error {

	inputEntity, err := r.inputRepositoryPort.GetInput()
	if err != nil {
		r.loggerPort.Error("Error in SendAccoutTransactionReport", "load input data", err)
		return err
	}
	r.loggerPort.Info("load report date successfuly!")

	notificationId, err := r.saveNotification(inputEntity)
	if err != nil {
		r.loggerPort.Error("Error in SendAccoutTransactionReport", "save notifiction to db", err)
		return err
	}
	r.loggerPort.Info("Save notification report successfuly!")

	resume, err := r.createResume(inputEntity)
	if err != nil {
		r.loggerPort.Error("Error in SendAccoutTransactionReport", "create transaction resume", err)
		return err
	}
	r.loggerPort.Info("Create resume report successfuly!")

	err = r.emailNotificationPort.SendAccountReport(resume)
	if err != nil {
		r.loggerPort.Error("Error in SendAccoutTransactionReport", "send account report", err)
		return err
	}
	r.loggerPort.Info("Send account transaction report successfuly!")

	err = r.updateNotificationStatus(notificationId, "SUCCESS")
	if err != nil {
		r.loggerPort.Error("Error in SendAccoutTransactionReport", "update notifiction to db", err)
		return err
	}
	r.loggerPort.Info("Update notification report successfuly!")

	return nil
}

func (r *reportUseCase) updateNotificationStatus(notificationId string, sendStatus string) error {

	notificationEntityFound, err := r.notificationRepositoryPort.FindById(notificationId)
	if err != nil {
		return err
	}

	notificationEntityFound.SendStatus = sendStatus
	notificationEntityFound.UpdateAt = time.Now().UTC()

	err = r.notificationRepositoryPort.UpdateById(notificationId, *notificationEntityFound)
	if err != nil {
		return err
	}

	return nil
}

func (r *reportUseCase) saveNotification(inputEntity entity.InputEntity) (string, error) {

	notificationEntity := entity.NotificationEntity{
		Accout:       inputEntity.Accout,
		Transactions: inputEntity.Transactions,
		SendStatus:   "PENDING",
		Sources:      inputEntity.Sources,
	}

	id, err := r.notificationRepositoryPort.SaveNotification(notificationEntity)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *reportUseCase) createResume(input entity.InputEntity) (entity.ResumeEntity, error) {
	const DATE_FORMAT = "2006-Jan"
	const EMAIL_SUBJECT = "Account Report"
	var balanceTotal float64

	mapTransactionsRemume := make(map[string]*entity.TransactionResumeEntity, 0)

	for _, transaction := range input.Transactions {
		dateKey := transaction.Date.Format(DATE_FORMAT)

		cardType, err := enum.GetValueEnum(transaction.CardType)
		if err != nil {
			continue
		}

		_, ok := mapTransactionsRemume[dateKey]
		if !ok {
			mapTransactionsRemume[dateKey] = &entity.TransactionResumeEntity{}
		}

		date, _ := time.Parse(DATE_FORMAT, transaction.Date.Format(DATE_FORMAT))

		mapTransactionsRemume[dateKey].Date = date

		switch cardType {
		case enum.Credit:
			mapTransactionsRemume[dateKey].NoTotal += 1
			mapTransactionsRemume[dateKey].AmountCreditTotal += transaction.Amount
			mapTransactionsRemume[dateKey].AmountCreditAverage = mapTransactionsRemume[dateKey].AmountCreditTotal / mapTransactionsRemume[dateKey].NoTotal
			balanceTotal += transaction.Amount
			mapTransactionsRemume[dateKey].AmountDebitAverage = mapTransactionsRemume[dateKey].AmountDebitTotal / mapTransactionsRemume[dateKey].NoTotal
		case enum.Debit:
			mapTransactionsRemume[dateKey].NoTotal += 1
			mapTransactionsRemume[dateKey].AmountDebitTotal += transaction.Amount
			mapTransactionsRemume[dateKey].AmountDebitAverage = mapTransactionsRemume[dateKey].AmountDebitTotal / mapTransactionsRemume[dateKey].NoTotal
			balanceTotal += transaction.Amount
			mapTransactionsRemume[dateKey].AmountCreditAverage = mapTransactionsRemume[dateKey].AmountCreditTotal / mapTransactionsRemume[dateKey].NoTotal
		default:
			continue
		}

	}

	traansactionsResume := make([]entity.TransactionResumeEntity, 0)
	for _, value := range mapTransactionsRemume {
		traansactionsResume = append(traansactionsResume, *value)
	}

	slices.SortFunc(traansactionsResume, func(a, b entity.TransactionResumeEntity) int { return b.Date.Compare(a.Date) })

	return entity.ResumeEntity{
		To:                  []string{input.Accout.Email},
		Subject:             fmt.Sprintf("Stori : %s", EMAIL_SUBJECT),
		AttachFiles:         input.Sources,
		Now:                 time.Now(),
		Name:                input.Accout.Name,
		BalanceTotal:        balanceTotal,
		TraansactionsResume: traansactionsResume,
	}, nil
}
