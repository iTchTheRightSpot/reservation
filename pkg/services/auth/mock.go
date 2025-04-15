package auth

import (
	"context"
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

func (m *MockJwtService) Encode(context.Context, *models.JwtObj, int) (*models.JwtResponse, error) {
	m.GenerateJwtCalled = true
	return m.JwtResponse, m.GenerateJwtError
}

func (m *MockJwtService) Decode(context.Context, string) (*models.JwtObj, error) {
	m.ValidateJwtCalled = true
	return m.StaffJwtObj, m.ValidateJwtError
}