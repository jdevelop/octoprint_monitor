package octoprint

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type octoprint struct {
	apiKey     string
	url        string
	httpClient *http.Client
}

type PrinterState struct {
	StatusText string `json:"text"`
	Flags      struct {
		Operational   bool `json:"operational"`
		Paused        bool `json:"paused"`
		Printing      bool `json:"printing"`
		SdReady       bool `json:"sdReady"`
		Error         bool `json:"error"`
		Ready         bool `json:"ready"`
		ClosedOnError bool `json:"closedOrError"`
	} `json:"flags"`
}

type Printer struct {
	State PrinterState `json:"state"`
}

type TemperatureData struct {
	Actual float32 `json:"actual"`
	Target float32 `json:"target"`
	Offset float32 `json:"offset"`
}

type Job struct {
	File struct {
		Date   uint32 `json:"date"`
		Name   string `json:"name"`
		Origin string `json:"origin"`
		Path   string `json:"path"`
		Size   uint32 `json:"size"`
	} `json:"file"`
	EstimatedPrintTime float32 `json:"estimatedPrintTime"`
	AveragePrintTime   float32 `json:"averagePrintTime"`
	LastPrintTime      float32 `json:"lastPrintTime"`
	Filament           struct {
		Length uint32 `json:"length"`
		Volume uint32 `json:"volume"`
	} `json:"filament"`
}

type TProgress struct {
	Completion          float32 `json:"completion"`
	FilePointer         uint32  `json:"filepos"`
	PrintTime           uint32  `json:"printTime"`
	PrintTimeLeft       uint32  `json:"printTimeLeft"`
	PrintTimeLeftOrigin string  `json:"printTimeLeftOrigin"`
}

type JobResponse struct {
	Job      Job       `json:"job"`
	Progress TProgress `json:"progress"`
	State    string    `json:"state"`
}

func (o *octoprint) request(url string) (req *http.Request, err error) {
	req, err = http.NewRequest("GET", o.url+url, nil)
	if err != nil {
		return
	}
	req.Header.Add("X-Api-Key", o.apiKey)
	return
}

func (o *octoprint) GetPrinterStatus() (s TPrinterStatus, err error) {
	req, err := o.request("/api/printer?history=false")
	resp, err := o.httpClient.Do(req)
	if err != nil {
		return
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if resp.StatusCode == 409 {
		s = PrinterFailed
		return
	}
	pState := Printer{}
	err = json.Unmarshal(data, &pState)
	if err != nil {
		return
	}
	if pState.State.Flags.ClosedOnError || pState.State.Flags.Error {
		s = PrinterFailed
	} else if pState.State.Flags.Printing {
		s = Printing
	} else {
		s = PrinterOk
	}
	return
}

func (o *octoprint) GetProgress() (p *TProgress, err error) {
	req, err := o.request("/api/job")
	if err != nil {
		return
	}
	resp, err := o.httpClient.Do(req)
	if err != nil {
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	jobState := JobResponse{}

	err = json.Unmarshal(data, &jobState)
	if err != nil {
		return
	}

	p = &jobState.Progress

	return
}

func ConnectOctoprint(api string, url string) (status PrinterStatus, err error) {
	c := http.Client{
		Timeout: 5 * time.Second,
	}
	status = &octoprint{
		apiKey:     api,
		url:        url,
		httpClient: &c,
	}
	return
}
