package utils

import (
	"log"
	"time"
)

type ILogger interface {
	Date() time.Time
	Error(variables ...interface{})
	Log(variables ...interface{})
	Fatal(variables ...interface{})
}

const (
	iError = "ERROR"
	ilog   = "LOG"
	iFatal = "FATAL"
)

type logger struct {
	time time.Time
}

func NewLogger(timezone string) (ILogger, error) {
	if utc, err := dateTimeToUTC(timezone); err != nil {
		log.Fatalf("failed to instantiate logger %v", err)
		return &logger{}, err
	} else {
		return &logger{time: utc}, nil
	}
}

func dateTimeToUTC(timezone string) (time.Time, error) {
	if timezone == "" {
		timezone = "America/Toronto"
	}
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, err
	}
	format := "2006-01-02 15:04:05"
	dt, err := time.Parse(format, time.Now().In(location).Format(format))
	if err != nil {
		return time.Time{}, err
	}
	return dt, nil
}

func (l *logger) Date() time.Time {
	return l.time
}

func (l *logger) Error(variables ...interface{}) {
	log.Printf("%s %s %v\n", l.time, iError, variables)
}

func (l *logger) Log(variables ...interface{}) {
	log.Printf("%s %s %v\n", l.time, ilog, variables)
}

func (l *logger) Fatal(variables ...interface{}) {
	log.Fatalf("%s %s %v\n", l.time, iFatal, variables)
}
