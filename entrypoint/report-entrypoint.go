package entrypoint

import "account_report/usecase"

type ReportEntryPoint interface {
	SendAccoutTransactionReport() error
}

type reportEntryPoint struct {
	reportUseCase usecase.ReportUseCase
}

func NewReportEntrypoint(reportUseCase usecase.ReportUseCase) ReportEntryPoint {
	return &reportEntryPoint{
		reportUseCase: reportUseCase,
	}
}

func (rc *reportEntryPoint) SendAccoutTransactionReport() error {
	return rc.reportUseCase.SendAccoutTransactionReport()
}
