package auth

import (
	"github.com/iTchTheRightSpot/erp-golang/pkg/models"
)

type MockJwtService struct {
	JwtResponse       *models.JwtResponse
	StaffJwtObj       *models.JwtObj
	GenerateJwtCalled bool
	GenerateJwtError  error
	ValidateJwtCalled bool
	ValidateJwtError  error
}

func (m *MockJwtService) GenerateJwt(*models.JwtObj, int) (*models.JwtResponse, error) {
	m.GenerateJwtCalled = true
	return m.JwtResponse, m.GenerateJwtError
}

func (m *MockJwtService) ValidateJwt(string) (*models.JwtObj, error) {
	m.ValidateJwtCalled = true
	return m.StaffJwtObj, m.ValidateJwtError
}
