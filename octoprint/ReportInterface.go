package octoprint

import (
	"fmt"
	"github.com/jdevelop/golang-rpi-extras/lcd_hd44780"
)

type Report interface {
	Render(printerState TPrinterStatus, progress *TProgress)
}

type localConsole int

const format = "%dm of %dm"

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
		fmt.Printf(format+"\n", progress.PrintTime/60, progress.PrintTimeLeft/60)
	}
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
		l.l.Print(fmt.Sprintf("Printing %.1f%%\n", progress.Completion))
	default:
		l.l.Print("Unknown")
	}
	l.l.SetCursor(1, 0)
	if progress != nil {
		str := fmt.Sprintf(format, progress.PrintTime/60, progress.PrintTimeLeft/60)
		l.l.Print(str)
	}
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
