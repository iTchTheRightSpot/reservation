package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

var TimeFormat = time.RFC3339

type ILogger interface {
	Timezone() *time.Location
	Date() time.Time
	Error(variables ...interface{})
	Log(variables ...interface{})
	Fatal(variables ...interface{})
}

type logType string

const (
	iError logType = "ERROR"
	iLog   logType = "LOG"
	iFatal logType = "FATAL"
)

type discord struct {
	Status logType   `json:"status"`
	Time   time.Time `json:"time"`
	Info   string    `json:"info"`
}

type logger struct {
	location *time.Location
	client   http.Client
	webhook  string
}

func NewLogger(timezone, url string) (ILogger, error) {
	if loc, err := dateInTimezone(timezone); err != nil {
		return nil, err
	} else {
		return &logger{location: loc, client: http.Client{Timeout: 2 * time.Second}, webhook: url}, nil
	}
}

func dateInTimezone(timezone string) (*time.Location, error) {
	if timezone == "" {
		timezone = "America/Toronto"
	}

	location, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, err
	}

	return location, nil
}

func (l *logger) post(mess string) {
	type content struct {
		C string `json:"content"`
	}

	m, err := json.Marshal(&content{C: "@everyone\n" + mess})
	if err != nil {
		log.Printf("%s %s %v\n", l.Date(), iError, err.Error())
		return
	}

	if _, err = http.Post(l.webhook, "application/json", bytes.NewBuffer(m)); err != nil {
		log.Printf("%s %s %v\n", l.Date(), iError, err.Error())
	}
}

func (l *logger) Timezone() *time.Location {
	return l.location
}

func (l *logger) Date() time.Time {
	dt, err := time.Parse(TimeFormat, time.Now().In(l.location).Format(TimeFormat))
	if err != nil {
		log.Print(err.Error())
	}
	return dt
}

func (l *logger) Error(variables ...interface{}) {
	obj := discord{
		Status: iError,
		Time:   l.Date(),
		Info:   fmt.Sprintf("%v", variables),
	}

	js, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		log.Printf("%s %s Failed to marshal log object: %v\n", l.Date(), iLog, err.Error())
		return
	}
	log.Printf("%s\n", string(js))
	l.post(string(js))
}

func (l *logger) Log(variables ...interface{}) {
	obj := discord{
		Status: iLog,
		Time:   l.Date(),
		Info:   fmt.Sprintf("%v", variables),
	}

	js, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		log.Printf("%s %s Failed to marshal log object: %v\n", l.Date(), iLog, err.Error())
		return
	}
	log.Printf("%s\n", string(js))
	l.post(string(js))
}

func (l *logger) Fatal(variables ...interface{}) {
	obj := discord{
		Status: iFatal,
		Time:   l.Date(),
		Info:   fmt.Sprintf("%v", variables),
	}

	js, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		log.Printf("%s %s Failed to marshal log object: %v\n", l.Date(), iLog, err.Error())
		return
	}
	log.Printf("%s\n", string(js))
}
