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

const (
	iError = "ERROR"
	iLog   = "LOG"
	iFatal = "FATAL"
)

type logger struct {
	time     time.Time
	location *time.Location
	client   http.Client
	webhook  string
}

func NewLogger(timezone, url string) (ILogger, error) {
	if utc, loc, err := dateInTimezone(timezone); err != nil {
		return nil, err
	} else {
		return &logger{time: utc, location: loc, client: http.Client{Timeout: 2 * time.Second}, webhook: url}, nil
	}
}

func dateInTimezone(timezone string) (time.Time, *time.Location, error) {
	if timezone == "" {
		timezone = "America/Toronto"
	}

	location, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, nil, err
	}

	dt, err := time.Parse(TimeFormat, time.Now().In(location).Format(TimeFormat))

	if err != nil {
		return time.Time{}, nil, err
	}

	return dt, location, nil
}

func (l *logger) post(mess string) {
	type content struct {
		C string `json:"content"`
	}

	m, err := json.Marshal(&content{C: mess})
	if err != nil {
		log.Printf("%s %s %v\n", l.time, iError, err.Error())
		return
	}

	if _, err = http.Post(l.webhook, "application/json", bytes.NewBuffer(m)); err != nil {
		log.Printf("%s %s %v\n", l.time, iError, err.Error())
	}
}

func (l *logger) Timezone() *time.Location {
	return l.location
}

func (l *logger) Date() time.Time {
	return l.time
}

func (l *logger) Error(variables ...interface{}) {
	log.Printf("%s %s %v\n", l.time, iError, variables)
	l.post(fmt.Sprintf("## __**%s**__ on %s @everyone\n### %s", iError, l.time, variables))
}

func (l *logger) Log(variables ...interface{}) {
	log.Printf("%s %s %v\n", l.time, iLog, variables)
	l.post(fmt.Sprintf("## __**%s**__ on %s @everyone\n### %s", iLog, l.time, variables))
}

func (l *logger) Fatal(variables ...interface{}) {
	log.Fatalf("%s %s %v\n", l.time, iFatal, variables)
}
