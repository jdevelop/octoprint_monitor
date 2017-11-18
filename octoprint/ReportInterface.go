package octoprint

import (
	"fmt"
	"github.com/jdevelop/golang-rpi-extras/lcd_hd44780"
	"time"
)

type Report interface {
	Render(printerState TPrinterStatus, progress *TProgress)
	Welcome()
}

type localConsole int

const format = "%s / %s"

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
		fmt.Println("Idle")
	case PrinterFailed:
		fmt.Println("Dead")
	case Printing:
		fmt.Printf("Printing %.1f%%\n", progress.Completion)
	default:
		fmt.Println("Unknown")
	}
	if progress != nil {
		fmt.Printf(format+"\n",
			fmtDuration(time.Duration(progress.PrintTime)*time.Second),
			fmtDuration(time.Duration(progress.PrintTimeLeft)*time.Second),
		)
	}
}

func (c *localConsole) Welcome() {
	fmt.Println("    OctoPrint   ")
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
		l.l.Print("Idle")
	case PrinterFailed:
		l.l.Print("Dead")
	case Printing:
		l.l.Print(fmt.Sprintf("Printing %.1f%%", progress.Completion))
	default:
		l.l.Print("Unknown")
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

func (l *localLCD) Welcome() {
	l.l.Cls()
	l.l.SetCursor(0, 0)
	l.l.Print("    OctoPrint   ")
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
