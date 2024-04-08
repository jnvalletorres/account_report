package repository

import (
	"account_report/domain/entity"
	"account_report/domain/port"
	"account_report/env"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

type fileInputRepository struct {
	logger   port.LoggerPort
	localEnv env.FileInput
}

func NewFileInputRepository(logger port.LoggerPort, environment env.Env) port.InputRepositoryPort {
	return &fileInputRepository{
		logger:   logger,
		localEnv: environment.FileInput,
	}
}

func (i *fileInputRepository) GetInput() (entity.InputEntity, error) {
	accountFileName := filepath.Join(i.localEnv.SourceDirectory, i.localEnv.AccountFileName)
	account, err := i.getAccountData(accountFileName)
	if err != nil {
		return entity.InputEntity{}, err
	}

	transactionFileName := filepath.Join(i.localEnv.SourceDirectory, i.localEnv.TransactionFileName)
	transactions, err := i.getTransactionData(transactionFileName)
	if err != nil {
		return entity.InputEntity{}, err
	}

	input := entity.InputEntity{
		Accout:       account,
		Transactions: transactions,
		Sources:      []string{accountFileName, transactionFileName},
	}
	return input, nil
}

func (i *fileInputRepository) getAccountData(fileName string) (entity.AccountEntity, error) {
	var (
		strErr  string
		matched bool
	)
	accountData, err := i.loadDataFromCSV(fileName)
	if err != nil {
		strErr = fmt.Sprintf("Error to load account csv file: %s incorrect", fileName)
		i.logger.Error(strErr)
		return entity.AccountEntity{}, err
	}

	if len(accountData) < 2 {
		strErr = fmt.Sprintf("Error to load account csv file: %s incorrect", fileName)
		i.logger.Error(strErr)
		return entity.AccountEntity{}, errors.New(strErr)
	}

	accountRow := accountData[1]

	if len(accountRow) < 3 {
		strErr = fmt.Sprintf("Error to load account csv file: %s incorrect", fileName)
		i.logger.Error(strErr)
		return entity.AccountEntity{}, errors.New(strErr)
	}

	matched, err = i.match(`^[[:digit:]]{1,}$`, accountRow[0])
	if err != nil || !matched {
		strErr = fmt.Sprintf("Validation error in csv file,  field(ID): %s is incorrect", accountRow[0])
		i.logger.Error(strErr)
		return entity.AccountEntity{}, err
	}

	id, err := strconv.Atoi(accountRow[0])
	if err != nil {
		strErr = fmt.Sprintf("Error to load account csv file, id: %s incorrect", accountRow[0])
		i.logger.Error(strErr)
		return entity.AccountEntity{}, err
	}
	accountEntity := entity.AccountEntity{
		Id:    id,
		Name:  accountRow[1],
		Email: accountRow[2],
	}

	return accountEntity, nil
}

func (i *fileInputRepository) getTransactionData(fileName string) ([]entity.TransactionEntity, error) {
	var (
		strErr  string
		matched bool
	)
	transactions := make([]entity.TransactionEntity, 0)

	transactionData, err := i.loadDataFromCSV(fileName)
	if err != nil {
		strErr = fmt.Sprintf("Error to load transaction csv file: %s incorrect", fileName)
		i.logger.Error(strErr)
		return transactions, err
	}

	if len(transactionData) < 2 {
		strErr = fmt.Sprintf("Error to load transaction csv file: %s incorrect", fileName)
		i.logger.Error(strErr)
		return transactions, errors.New(strErr)
	}

	for index, transactionRow := range transactionData {

		if index == 0 {
			continue
		}

		if len(transactionRow) < 2 {
			strErr = fmt.Sprintf("Error to load transaction csv file: %s incorrect", fileName)
			i.logger.Error(strErr)
			return []entity.TransactionEntity{}, errors.New(strErr)
		}

		matched, err = i.match(`^[[:digit:]]{1,}$`, transactionRow[0])
		if err != nil || !matched {
			strErr = fmt.Sprintf("Validation error in csv file,  field(ID): %s is incorrect", transactionRow[0])
			i.logger.Error(strErr)
			return []entity.TransactionEntity{}, err
		}

		matched, err = i.match(`^[+|-]{1,}(([[:digit:]]{1,})|([[:digit:]]{1,}\.{1,1}[[:digit:]]{1,}))$`, transactionRow[2])
		if err != nil || !matched {
			strErr = fmt.Sprintf("Validation error in csv file, field(Amount): %s is incorrect", transactionRow[2])
			i.logger.Error(strErr)
			return []entity.TransactionEntity{}, err
		}

		id, err := strconv.Atoi(transactionRow[0])
		if err != nil {
			strErr = fmt.Sprintf("Error to load transaction csv file, id: %s is incorrect", transactionRow[0])
			i.logger.Error(strErr)
			return []entity.TransactionEntity{}, err
		}

		strSymbol, strAmount := i.parserAmount(transactionRow[2])

		amount, err := strconv.ParseFloat(strAmount, 64)
		if err != nil {
			strErr = fmt.Sprintf("Error to load transaction csv file, amount: %s is incorrect", transactionRow[0])
			i.logger.Error(strErr)
			return []entity.TransactionEntity{}, err
		}

		date, err := time.Parse(time.DateOnly, transactionRow[1])
		if err != nil {
			strErr = fmt.Sprintf("Error to load transaction csv file, date: %s is incorrect", transactionRow[0])
			i.logger.Error(strErr)
			return []entity.TransactionEntity{}, err
		}

		transactionEntity := entity.TransactionEntity{
			Id:       id,
			Date:     date,
			CardType: strSymbol,
			Amount:   amount,
		}

		transactions = append(transactions, transactionEntity)
	}

	return transactions, nil
}

func (i *fileInputRepository) parserAmount(str string) (string, string) {
	return str[0:1], str[1:]

}

func (i *fileInputRepository) match(pattern string, str string) (bool, error) {
	matched, err := regexp.Match(pattern, []byte(str))
	if err != nil {
		return false, err
	}
	return matched, nil

}

func (i *fileInputRepository) loadDataFromCSV(fileName string) ([][]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return [][]string{}, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	return reader.ReadAll()

}
