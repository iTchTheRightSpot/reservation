package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
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

func construct(date time.Time, tp logType, variables ...interface{}) ([]byte, error) {
	var sb strings.Builder
	for _, v := range variables {
		sb.WriteString(fmt.Sprintf("%s", v))
	}
	sb.WriteString("\n")
	obj := discord{
		Status: tp,
		Time:   date,
		Info:   sb.String(),
	}
	return json.MarshalIndent(obj, "", "  ")
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
	js, err := construct(l.Date(), iError, variables)
	if err != nil {
		obj := discord{
			Status: iError,
			Time:   l.Date(),
			Info:   "failed to marshal ERROR object" + err.Error(),
		}
		log.Print(obj)
		return
	}
	str := string(js)
	log.Print(str)
	l.post(str)
}

func (l *logger) Log(variables ...interface{}) {
	js, err := construct(l.Date(), iLog, variables)
	if err != nil {
		obj := discord{
			Status: iError,
			Time:   l.Date(),
			Info:   "failed to marshal LOG object" + err.Error(),
		}
		log.Print(obj)
		return
	}
	v := string(js)
	log.Print(v)
	l.post(v)
}

func (l *logger) Fatal(variables ...interface{}) {
	js, err := construct(l.Date(), iFatal, variables)
	if err != nil {
		obj := discord{
			Status: iFatal,
			Time:   l.Date(),
			Info:   "failed to marshal FATAL object" + err.Error(),
		}
		log.Print(obj)
		return
	}
	log.Print(string(js))
}
