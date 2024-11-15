package auth

import (
	"github.com/iTchTheRightSpot/erp-golang/config"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"os"
	"reflect"
	"testing"
)

func TestJwtService(t *testing.T) {
	if err := os.Chdir("../../../"); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	t.Run("should generate and validate jwt", func(t *testing.T) {
		env := &config.SecretVariables{}
		con, err := env.Config()
		if err != nil {
			t.Fatalf("%s", err)
		}

		s := NewJwtService(utils.NewMockLogger(), con)

		// method to test
		roles := []models.RoleEnum{models.STAFF, models.DEVELOPER, models.USER}
		o := &models.JwtObj{Roles: roles, UserUUID: "staff-uuid"}

		// method to test & assert
		obj, err := s.GenerateJwt(
			o,
			utils.TwoDaysInSeconds,
		)
		if err != nil {
			t.Errorf("exception generating jwt %s", err)
		}

		// method to test & assert
		v, err := s.ValidateJwt(obj.Token)
		if err != nil {
			t.Errorf("exception validating generated token jwt %s", err)
		}

		if !reflect.DeepEqual(o.UserUUID, v.UserUUID) {
			t.Errorf("expect %s to equal given %s", o.UserUUID, v.UserUUID)
		}

		if !reflect.DeepEqual(o.Roles, v.Roles) {
			t.Errorf("expect %s to equal given %s", o.Roles, v.Roles)
		}
	})
}
