package utils

import (
	"log"
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
	ilog   = "LOG"
	iFatal = "FATAL"
)

type logger struct {
	time     time.Time
	location *time.Location
}

func NewLogger(timezone string) (ILogger, error) {
	if utc, loc, err := dateInTimezone(timezone); err != nil {
		log.Fatalf("failed to instantiate logger %v", err)
		return nil, err
	} else {
		return &logger{time: utc, location: loc}, nil
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

func (l *logger) Timezone() *time.Location {
	return l.location
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
