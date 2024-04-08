package enum

import (
	"errors"
	"fmt"
)

type CardTypeEnum string

const (
	Debit  CardTypeEnum = "Debit"
	Credit CardTypeEnum = "Credit"
)

func (ct CardTypeEnum) String() string {
	switch ct {
	case Debit:
		return "Debit(-)"
	case Credit:
		return "Credit(+)"
	default:
		return "unknown"
	}
}

func GetValueEnum(key string) (CardTypeEnum, error) {

	return func() (CardTypeEnum, error) {
		switch key {
		case "+":
			return Credit, nil
		case "-":
			return Debit, nil
		default:
			str := fmt.Sprintf("The key value: %s is incorrect", key)
			return "", errors.New(str)
		}
	}()
}
