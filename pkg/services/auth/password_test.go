package auth

import (
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"testing"
)

func TestProfileModel(t *testing.T) {
	ps := NewPasswordService(utils.NewMockLogger())

	t.Run("regex. password has no uppercase", func(t *testing.T) {
		if err := ps.PasswordRegex("password123#!"); err == nil {
			t.Error(err.Error())
		}
	})

	t.Run("regex. password has no lowercase", func(t *testing.T) {
		if err := ps.PasswordRegex("PASSWORD123#!"); err == nil {
			t.Error(err.Error())
		}
	})

	t.Run("regex. password has no number", func(t *testing.T) {
		if err := ps.PasswordRegex("passworD#!"); err == nil {
			t.Error(err.Error())
		}
	})

	t.Run("regex. password has no special character", func(t *testing.T) {
		if err := ps.PasswordRegex("passworD123"); err == nil {
			t.Error(err.Error())
		}
	})

	t.Run("regex. success", func(t *testing.T) {
		if err := ps.PasswordRegex("pa(ssworD123#"); err != nil {
			t.Error(err.Error())
		}
	})
}
