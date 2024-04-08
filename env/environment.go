package env

import "os"

type FileInput struct {
	SourceDirectory     string
	AccountFileName     string
	TransactionFileName string
}

type App struct {
	AppId string
	GoEnv string
}

type Smtp struct {
	Host       string
	UserName   string
	Password   string
	PortNumber string
}

type MongoDB struct {
	Url     string
	Name    string
	TimeOut string
}

type Env struct {
	FileInput FileInput
	App       App
	Smtp      Smtp
	MongoDB   MongoDB
}

var Environment = Env{
	FileInput: FileInput{
		AccountFileName:     os.Getenv("ACCOUNT_FILE_NAME"),
		TransactionFileName: os.Getenv("TRANSACTION_FILE_NAME"),
		SourceDirectory:     os.Getenv("SOURCE_DIRECTORY"),
	},
	App: App{
		AppId: os.Getenv("APP_ID"),
		GoEnv: os.Getenv("GO_ENV"),
	},
	Smtp: Smtp{
		Host:       os.Getenv("SMTP_HOST"),
		UserName:   os.Getenv("SMTP_USER_NAME"),
		Password:   os.Getenv("SMTP_PASSWORD"),
		PortNumber: os.Getenv("SMTP_PORT_NUMBER"),
	},
	MongoDB: MongoDB{
		Url:     os.Getenv("DB_URL"),
		Name:    os.Getenv("DATABASE_NAME"),
		TimeOut: os.Getenv("DATABASE_TIMEOUT"),
	},
}
