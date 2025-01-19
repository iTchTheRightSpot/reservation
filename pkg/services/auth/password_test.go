package auth

import (
	"testing"
)

func TestProfileModel(t *testing.T) {
	ps := NewPasswordService(nil)

	t.Run("regex. password has no uppercase", func(t *testing.T) {
		b, err := ps.PasswordRegex("password123#!")
		if err != nil {
			t.Error(err.Error())
		}

		if b {
			t.Error("expect false, got true")
		}
	})

	t.Run("regex. password has no lowercase", func(t *testing.T) {
		b, err := ps.PasswordRegex("PASSWORD123#!")
		if err != nil {
			t.Error(err.Error())
		}

		if b {
			t.Error("expect false, got true")
		}
	})

	t.Run("regex. password has no number", func(t *testing.T) {
		b, err := ps.PasswordRegex("passworD#!")
		if err != nil {
			t.Error(err.Error())
		}

		if b {
			t.Error("expect false, got true")
		}
	})

	t.Run("regex. password has no special character", func(t *testing.T) {
		b, err := ps.PasswordRegex("passworD123")
		if err != nil {
			t.Error(err.Error())
		}

		if b {
			t.Error("expect false, got true")
		}
	})

	t.Run("regex. success", func(t *testing.T) {
		b, err := ps.PasswordRegex("pa(ssworD123#")
		if err != nil {
			t.Error(err.Error())
		}

		if !b {
			t.Error("expect true, got false")
		}
	})
}
