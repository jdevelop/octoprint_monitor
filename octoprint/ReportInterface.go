package octoprint

import (
	"fmt"
	"github.com/jdevelop/golang-rpi-extras/lcd_hd44780"
	"time"
)

type Report interface {
	Render(printerState TPrinterStatus, progress *TProgress)
	Welcome(v *ApiVersion)
}

type localConsole int

const (
	format             = "%s / %s"
	IdleStatus         = "Idle"
	DisconnectedStatus = "Disconnected"
	UnknownStatus      = "Unknown"
	PrintingStatus     = "Printing %.1f%%"
)

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return fmt.Sprintf("%02d:%02d", h, m)
}

func (c *localConsole) Render(printerState TPrinterStatus, progress *TProgress) {
	switch printerState {
	case PrinterOk:
		fmt.Println(IdleStatus)
	case PrinterFailed:
		fmt.Println(DisconnectedStatus)
	case Printing:
		fmt.Printf(PrintingStatus+"\n", progress.Completion)
	default:
		fmt.Println(UnknownStatus)
	}
	if progress != nil {
		fmt.Printf(format+"\n",
			fmtDuration(time.Duration(progress.PrintTime)*time.Second),
			fmtDuration(time.Duration(progress.PrintTimeLeft)*time.Second),
		)
	}
}

func (c *localConsole) Welcome(v *ApiVersion) {
	if v == nil {
		fmt.Println("    OctoPrint   ")
	} else {
		fmt.Printf("OctoPrint %s\n", v.ServerVersion)
	}
	fmt.Println(" Status Monitor ")
}

type localLCD struct {
	l lcd_hd44780.PiLCD4
}

func (l *localLCD) Render(printerState TPrinterStatus, progress *TProgress) {
	l.l.Cls()
	l.l.SetCursor(0, 0)
	switch printerState {
	case PrinterOk:
		l.l.Print(IdleStatus)
	case PrinterFailed:
		l.l.Print(DisconnectedStatus)
	case Printing:
		l.l.Print(fmt.Sprintf(PrintingStatus, progress.Completion))
	default:
		l.l.Print(UnknownStatus)
	}
	l.l.SetCursor(1, 0)
	if progress != nil {
		str := fmt.Sprintf(format,
			fmtDuration(time.Duration(progress.PrintTime)*time.Second),
			fmtDuration(time.Duration(progress.PrintTimeLeft)*time.Second),
		)
		l.l.Print(str)
	}
}

func (l *localLCD) Welcome(v *ApiVersion) {
	l.l.Cls()
	l.l.SetCursor(0, 0)
	if v != nil {
		l.l.Print(fmt.Sprintf("OctoPrint %s", v.ServerVersion))
	} else {
		l.l.Print("    OctoPrint")
	}
	l.l.SetCursor(1, 0)
	l.l.Print(" Status Monitor ")
}

func MakeLCD(dataPins []int, resetPin int, strobePin int) (r Report, err error) {
	lcd, err := lcd_hd44780.NewLCD4(dataPins, resetPin, strobePin)
	if err != nil {
		return
	}
	lcd.Init()
	r = &localLCD{
		l: lcd,
	}
	return
}

var lc = localConsole(0)

func MakeConsole() (r Report, err error) {
	r = &lc
	return
}
