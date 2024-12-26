package auth

import (
	"github.com/iTchTheRightSpot/erp-golang/config"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"os"
	"testing"
)

func TestJwtService(t *testing.T) {
	t.Parallel()

	if err := os.Chdir("../../../"); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}
	env := &config.SecretVariables{}
	con := env.Config()

	t.Run("should generate and validate jwt", func(t *testing.T) {
		t.Parallel()

		s := NewJwtService(utils.NewMockLogger(), con)

		cred := make([]models.RolePermission, 2)
		cred[0] = models.RolePermission{
			Role:        models.STAFF,
			Permissions: []models.PermissionEnum{models.READ, models.DELETE},
		}
		cred[1] = models.RolePermission{
			Role:        models.DEVELOPER,
			Permissions: []models.PermissionEnum{models.READ, models.DELETE},
		}

		o := &models.JwtObj{AccessControls: cred, UserId: "staff-uuid"}

		// method to test & assert
		obj, err := s.GenerateJwt(o, utils.TwoDaysInSeconds)
		if err != nil {
			t.Errorf("exception generating jwt %s", err)
		}

		// method to test & assert
		v, err := s.ValidateJwt(obj.Token)
		if err != nil {
			t.Errorf("exception validating generated token jwt %s", err)
		}

		if !(o.UserId == v.UserId) {
			t.Errorf("expect %s to equal given %s", o.UserId, v.UserId)
		}
	})
}
