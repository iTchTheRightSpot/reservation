package utils

import (
	"log"
	"time"
)

type mockLogger struct {
	date time.Time
	loc  *time.Location
}

func NewDevLogger() ILogger {
	utc, loc, err := dateInTimezone("UTC")
	if err != nil {
		return nil
	}
	return &mockLogger{date: utc, loc: loc}
}

func (m *mockLogger) Timezone() *time.Location {
	return m.loc
}

func (m *mockLogger) Date() time.Time {
	return m.date
}

func (m *mockLogger) Error(variables ...interface{}) {
	log.Printf("%s", variables)
}

func (m *mockLogger) Log(variables ...interface{}) {
	log.Printf("%s", variables)
}

func (m *mockLogger) Fatal(variables ...interface{}) {
	log.Fatalf("%s", variables)
}
