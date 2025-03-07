package utils

import (
	"log"
	"time"
)

type mockLogger struct {
	loc *time.Location
}

func NewDevLogger() ILogger {
	loc, err := dateInTimezone("UTC")
	if err != nil {
		return nil
	}
	return &mockLogger{loc: loc}
}

func (m *mockLogger) Timezone() *time.Location {
	return m.loc
}

func (m *mockLogger) Date() time.Time {
	dt, err := time.Parse(TimeFormat, time.Now().In(m.loc).Format(TimeFormat))
	if err != nil {
		log.Print(err.Error())
	}
	return dt
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
