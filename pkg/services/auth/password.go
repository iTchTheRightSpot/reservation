package auth

import (
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

type IPasswordService interface {
	PasswordRegex(str string) error
	Encode(str string) ([]byte, error)
	Validate(hash, str []byte) error
}

type passwordService struct {
	logger utils.ILogger
}

// NewPasswordService https://pkg.go.dev/golang.org/x/crypto/bcrypt
func NewPasswordService(l utils.ILogger) IPasswordService {
	return &passwordService{logger: l}
}

func (dep *passwordService) PasswordRegex(str string) error {
	mess := "password must be 8-20 characters long and include at least one uppercase letter, one lowercase letter, one number, and one special character (#?!@$%^&*-)"

	// length check
	if len(str) < 8 || len(str) > 20 {
		return &utils.BadRequestError{Message: mess}
	}

	// check for at least one lowercase letter
	if !regexp.MustCompile(`[a-z]`).MatchString(str) {
		return &utils.BadRequestError{Message: mess}
	}

	// check for at least one uppercase letter
	if !regexp.MustCompile(`[A-Z]`).MatchString(str) {
		return &utils.BadRequestError{Message: mess}
	}

	// check for at least one digit
	if !regexp.MustCompile(`\d`).MatchString(str) {
		return &utils.BadRequestError{Message: mess}
	}

	// check for at least one special character
	if !regexp.MustCompile(`[#?!@$%^&*-]`).MatchString(str) {
		return &utils.BadRequestError{Message: mess}
	}
	return nil
}

func (dep *passwordService) Encode(str string) ([]byte, error) {
	ps, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		dep.logger.Error(err.Error())
		return nil, &utils.ServerError{Message: "error encoding password"}
	}
	return ps, nil
}

func (dep *passwordService) Validate(hash, str []byte) error {
	if err := bcrypt.CompareHashAndPassword(hash, str); err != nil {
		dep.logger.Error(err.Error())
		return &utils.AuthenticationError{Message: "invalid email or password"}
	}
	return nil
}
