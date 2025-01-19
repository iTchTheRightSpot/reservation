package auth

import (
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

type IPasswordService interface {
	PasswordRegex(str string) (bool, error)
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

func (dep *passwordService) PasswordRegex(str string) (bool, error) {
	// https://stackoverflow.com/questions/19605150/regex-for-password-must-contain-at-least-eight-characters-at-least-one-number-a
	regex := `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[#$@!%&*?])[A-Za-z\d#$@!%&*?]{8,15}$`
	return regexp.MatchString(regex, str)
}

func (dep *passwordService) Encode(str string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
}

func (dep *passwordService) Validate(hash, str []byte) error {
	return bcrypt.CompareHashAndPassword(hash, str)
}
