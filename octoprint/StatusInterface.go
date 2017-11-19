package octoprint

type TPrinterStatus int

const (
	PrinterOk TPrinterStatus = iota
	Printing
	PrinterFailed
)

type PrinterStatus interface {
	GetVersionInfo() (*ApiVersion, error)
	GetPrinterStatus() (TPrinterStatus, error)
	GetProgress() (*TProgress, error)
}
