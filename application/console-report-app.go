package application

import (
	"account_report/config"
	"account_report/entrypoint"
)

var (
	reportController entrypoint.ReportEntryPoint
)

func init() {
	configuration := config.NewConsoleReportConfig()
	reportController = entrypoint.NewReportEntrypoint(configuration.CreateReporUseCase())
}

func Launch() {
	reportController.SendAccoutTransactionReport()
}
