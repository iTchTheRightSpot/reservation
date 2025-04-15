package auth

import (
	"context"
	"github.com/iTchTheRightSpot/utility/utils"
	"testing"
)

func TestProfileModel(t *testing.T) {
	instance := NewPasswordService(utils.DevLogger("UTC"))

	t.Run("regex. password has no uppercase", func(t *testing.T) {
		if err := instance.PasswordRegex("password123#!"); err == nil {
			t.Error(err.Error())
		}
	})

	t.Run("regex. password has no lowercase", func(t *testing.T) {
		if err := instance.PasswordRegex("PASSWORD123#!"); err == nil {
			t.Error(err.Error())
		}
	})

	t.Run("regex. password has no number", func(t *testing.T) {
		if err := instance.PasswordRegex("passworD#!"); err == nil {
			t.Error(err.Error())
		}
	})

	t.Run("regex. password has no special character", func(t *testing.T) {
		if err := instance.PasswordRegex("passworD123"); err == nil {
			t.Error(err.Error())
		}
	})

	t.Run("regex. success", func(t *testing.T) {
		if err := instance.PasswordRegex("pa(ssworD123#"); err != nil {
			t.Error(err.Error())
		}
	})

	t.Run("should encode & validate password", func(t *testing.T) {
		// given
		str := "password"
		ctx := context.Background()

		// method to test & assert
		encode, err := instance.Encode(ctx, str)
		if err != nil {
			t.Error(err.Error())
		}

		// method to test & assert
		if err = instance.Validate(ctx, encode, []byte(str)); err != nil {
			t.Error(err.Error())
		}
	})
}