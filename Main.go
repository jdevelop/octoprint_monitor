package main

import (
	"flag"
	"fmt"
	"github.com/jdevelop/octoprint-status/octoprint"
	"github.com/spf13/viper"
	"log"
	"os/user"
	"sync"
	"time"
)

func main() {

	configPath := flag.String("config", "", "")
	flag.Parse()

	if *configPath == "" {
		u, _ := user.Current()

		viper.SetConfigFile(fmt.Sprintf("%s/.octoprint-st-rc", u.HomeDir))
	} else {
		viper.SetConfigFile(*configPath)
	}
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Can not find the configuration file at "+viper.ConfigFileUsed(), err)
	}

	type config struct {
		ApiKey  string        `mapstructure:"apikey"`
		Url     string        `mapstructure:"url"`
		Refresh time.Duration `mapstructure:"refresh"`
		PiLCD struct {
			Rs   int   `mapstructure:"rs"`
			E    int   `mapstructure:"e"`
			Data []int `mapstructure:"data"`
		} `mapstructure:"lcd"`
	}

	c := config{}

	if err := viper.Unmarshal(&c); err != nil {
		log.Fatal("Can't unmarshal config", err)
	}

	//op, err := octoprint.ConnectOctoprint("5EE1807CD2CB4AADA4029AAFFE11B05A", "http://10.0.0.28:5000")
	op, err := octoprint.ConnectOctoprint(c.ApiKey, c.Url)
	if err != nil {
		log.Fatal("Can't access print data", err)
	}

	var report octoprint.Report

	if c.PiLCD.Data != nil {
		fmt.Printf("Using LCD interface at %v\n", c.PiLCD)
		report, err = octoprint.MakeLCD(c.PiLCD.Data, c.PiLCD.Rs, c.PiLCD.E)
	} else {
		fmt.Println("Console access")
		report, err = octoprint.MakeConsole()
	}

	if err != nil {
		log.Fatal(err)
	}

	if c.Refresh == 0 {
		fmt.Println("Using default refresh interval")
		c.Refresh = 30 * time.Second
	}

	ticker := time.Tick(c.Refresh)

	var wg sync.WaitGroup
	wg.Add(1)

	report.Welcome()

	onTick := func() {
		for range ticker {
			s, err1 := op.GetPrinterStatus()
			p, err2 := op.GetProgress()

			if err1 == nil && err2 == nil {
				report.Render(s, p)
			} else {
				fmt.Println(err1, err2)
			}

		}
		wg.Done()
	}

	go onTick()

	wg.Wait()

}
