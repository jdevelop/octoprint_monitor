package octoprint

type TPrinterStatus int

const (
	PrinterOk TPrinterStatus = iota
	Printing
	PrinterFailed
)

type PrinterStatus interface {
	GetPrinterStatus() (TPrinterStatus, error)
	GetProgress() (*TProgress, error)
}
