package utils

import (
	"log"
	"time"
)

type mockLogger struct{}

func NewMockLogger() ILogger {
	return &mockLogger{}
}

func (m mockLogger) Date() time.Time {
	utc, _ := dateTimeToUTC("")
	return utc
}

func (m mockLogger) Error(variables ...interface{}) {
	log.Printf("%s", variables)
}

func (m mockLogger) Log(variables ...interface{}) {
	log.Printf("%s", variables)
}

func (m mockLogger) Fatal(variables ...interface{}) {
	log.Fatalf("%s", variables)
}
