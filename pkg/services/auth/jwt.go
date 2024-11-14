package auth

import (
	"crypto/rsa"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/iTchTheRightSpot/erp-golang/config"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"os"
	"time"
)

func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	keyData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return privateKey, nil
}

func loadPublicKey(path string) (*rsa.PublicKey, error) {
	keyData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key file: %w", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(keyData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	return publicKey, nil
}

type IJwtService interface {
	GenerateJwt(o *models.StaffJwtObj, expirationInSeconds int) (*models.JwtResponse, error)
	ValidateJwt(str string) (*models.StaffJwtObj, error)
}

type jwtService struct {
	logger  utils.ILogger
	privKey *rsa.PrivateKey
	pubKey  *rsa.PublicKey
}

func NewJwtService(l utils.ILogger, env *config.SecretVariables) IJwtService {
	priv, err := loadPrivateKey(env.PrivateKeyPath)
	if err != nil {
		l.Fatal(err)
	}

	pub, err := loadPublicKey(env.PublicKeyPath)
	if err != nil {
		l.Fatal(err)
	}

	return &jwtService{logger: l, privKey: priv, pubKey: pub}
}

func (dep *jwtService) GenerateJwt(o *models.StaffJwtObj, expirationInSeconds int) (*models.JwtResponse, error) {
	exp := dep.logger.Date().Add(time.Duration(expirationInSeconds) * time.Second)

	claims := jwt.NewWithClaims(
		jwt.SigningMethodRS256,
		jwt.MapClaims{
			"sub": o.StaffUUID,
			"obj": o,
			"iss": "Landscape ERP",
			"exp": exp.Unix(),
			"iat": dep.logger.Date().Unix(),
		},
	)

	token, err := claims.SignedString(dep.privKey)
	if err != nil {
		dep.logger.Error(fmt.Sprintf("Error signing token: %v", err))
		return nil, err
	}

	return &models.JwtResponse{Token: token, ExpireAt: exp}, nil
}

func (dep *jwtService) ValidateJwt(str string) (*models.StaffJwtObj, error) {
	token, err := jwt.Parse(str, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			errMsg := fmt.Sprintf("unexpected signing method: %v", token.Header["alg"])
			dep.logger.Error(errMsg)
			return nil, fmt.Errorf(errMsg)
		}
		return dep.pubKey, nil
	})

	if err != nil {
		dep.logger.Error(fmt.Sprintf("failed to parse token: %v", err))
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		dep.logger.Error("failed to parse claims from token")
		return nil, fmt.Errorf("failed to parse claims")
	}

	obj, ok := claims["obj"].(map[string]interface{})
	if !ok {
		dep.logger.Error("invalid object format in claims")
		return nil, fmt.Errorf("invalid object format in claims")
	}

	staffObj := &models.StaffJwtObj{}

	if staffUUID, ok := obj["staff_uuid"].(string); ok {
		staffObj.StaffUUID = staffUUID
	} else {
		dep.logger.Error("missing or invalid StaffUUID in token claims")
		return nil, fmt.Errorf("invalid or missing StaffUUID")
	}

	if roles, ok := obj["roles"].([]interface{}); ok {
		parsedRoles := make([]models.RoleEnum, len(roles))
		for i, role := range roles {
			if roleStr, ok := role.(string); ok {
				parsedRoles[i] = models.RoleEnum(roleStr)
			} else {
				dep.logger.Error("invalid role format in token claims")
				return nil, fmt.Errorf("invalid role format in token claims")
			}
		}
		staffObj.Roles = parsedRoles
	} else {
		dep.logger.Error("missing or invalid roles in token claims")
		return nil, fmt.Errorf("invalid or missing roles")
	}

	dep.logger.Log("successfully validated jwt")
	return staffObj, nil
}
