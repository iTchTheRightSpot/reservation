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
	Encode(o *models.JwtObj, expirationInSeconds int) (*models.JwtResponse, error)
	Decode(str string) (*models.JwtObj, error)
}

type jwtService struct {
	logger  utils.ILogger
	privKey *rsa.PrivateKey
	pubKey  *rsa.PublicKey
}

func NewJwtServiceAsymmetric(l utils.ILogger, env *config.SecretVariables) IJwtService {
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

func (dep *jwtService) Encode(o *models.JwtObj, expirationInSeconds int) (*models.JwtResponse, error) {
	exp := dep.logger.Date().Add(time.Duration(expirationInSeconds) * time.Second)

	claims := jwt.NewWithClaims(
		jwt.SigningMethodRS256,
		jwt.MapClaims{
			"sub": o.UserId,
			"obj": o,
			"iss": "Reservation application powered by S.EJ.U development",
			"exp": exp.Unix(),
			"iat": dep.logger.Date().Unix(),
		},
	)

	token, err := claims.SignedString(dep.privKey)
	if err != nil {
		dep.logger.Error(err.Error())
		return nil, &utils.ServerError{Message: "error encoding to jwt"}
	}

	return &models.JwtResponse{Token: token, ExpireAt: exp}, nil
}

func (dep *jwtService) Decode(str string) (*models.JwtObj, error) {
	token, err := jwt.Parse(str, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			errMsg := fmt.Sprintf("unexpected signing method: %v", token.Header["alg"])
			dep.logger.Error(errMsg)
			return nil, &utils.AuthenticationError{Message: errMsg}
		}
		return dep.pubKey, nil
	})

	if err != nil {
		dep.logger.Error(err.Error())
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		dep.logger.Error("failed to parse claims from token")
		return nil, &utils.AuthenticationError{}
	}

	exp, err := claims.GetExpirationTime()
	if err != nil {
		dep.logger.Error(err)
		return nil, &utils.AuthenticationError{}
	}

	obj, ok := claims["obj"].(map[string]interface{})
	if !ok {
		dep.logger.Error("invalid object format in claims")
		return nil, &utils.AuthenticationError{}
	}

	jwtObj := &models.JwtObj{ExpireAt: &exp.Time}
	if uuid, ok := obj["user_id"].(string); ok {
		jwtObj.UserId = uuid
	} else {
		dep.logger.Error("missing or invalid UserId in token claims")
		return nil, &utils.AuthenticationError{}
	}

	if credentials, ok := obj["access_controls"].([]interface{}); ok {
		parsedRoles := make([]models.RolePermissionEnum, len(credentials))

		for i, cred := range credentials {
			if rolePermission, ok := cred.(map[string]interface{}); ok {
				role := models.RoleEnum(rolePermission["role"].(string))
				var permissions []models.PermissionEnum

				if perms, ok := rolePermission["permissions"].([]interface{}); ok {
					for _, perm := range perms {
						permissions = append(permissions, models.PermissionEnum(perm.(string)))
					}
				}

				parsedRoles[i] = models.RolePermissionEnum{
					Role:        role,
					Permissions: permissions,
				}
			}
		}

		jwtObj.AccessControls = parsedRoles
	} else {
		dep.logger.Error("missing or invalid roles in token claims")
		return nil, &utils.AuthenticationError{}
	}

	return jwtObj, nil
}

type jwt1Service struct {
	logger       utils.ILogger
	symmetricKey string
}

func NewJwtServiceSymmetric(l utils.ILogger, env *config.SecretVariables) IJwtService {
	return &jwt1Service{logger: l, symmetricKey: env.SymmetricKey}
}

func (dep *jwt1Service) Encode(o *models.JwtObj, expirationInSeconds int) (*models.JwtResponse, error) {
	exp := dep.logger.Date().Add(time.Duration(expirationInSeconds) * time.Second)

	claims := jwt.NewWithClaims(
		jwt.SigningMethodRS256,
		jwt.MapClaims{
			"sub": o.UserId,
			"obj": o,
			"iss": "Reservation application powered by S.EJ.U development",
			"exp": exp.Unix(),
			"iat": dep.logger.Date().Unix(),
		},
	)

	token, err := claims.SignedString(dep.symmetricKey)
	if err != nil {
		dep.logger.Error(err.Error())
		return nil, &utils.ServerError{Message: "error encoding to jwt"}
	}

	return &models.JwtResponse{Token: token, ExpireAt: exp}, nil
}

func (dep *jwt1Service) Decode(str string) (*models.JwtObj, error) {
	token, err := jwt.Parse(str, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(dep.symmetricKey), nil
	})

	if err != nil {
		dep.logger.Error(err.Error())
		return nil, &utils.AuthenticationError{}
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		dep.logger.Error("failed to parse claims from token")
		return nil, &utils.AuthenticationError{}
	}

	exp, err := claims.GetExpirationTime()
	if err != nil {
		dep.logger.Error(err)
		return nil, &utils.AuthenticationError{}
	}

	obj, ok := claims["obj"].(map[string]interface{})
	if !ok {
		dep.logger.Error("invalid object format in claims")
		return nil, &utils.AuthenticationError{}
	}

	jwtObj := &models.JwtObj{ExpireAt: &exp.Time}
	if uuid, ok := obj["user_id"].(string); ok {
		jwtObj.UserId = uuid
	} else {
		dep.logger.Error("missing or invalid UserId in token claims")
		return nil, &utils.AuthenticationError{}
	}

	if credentials, ok := obj["access_controls"].([]interface{}); ok {
		parsedRoles := make([]models.RolePermissionEnum, len(credentials))

		for i, cred := range credentials {
			if rolePermission, ok := cred.(map[string]interface{}); ok {
				role := models.RoleEnum(rolePermission["role"].(string))
				var permissions []models.PermissionEnum

				if perms, ok := rolePermission["permissions"].([]interface{}); ok {
					for _, perm := range perms {
						permissions = append(permissions, models.PermissionEnum(perm.(string)))
					}
				}

				parsedRoles[i] = models.RolePermissionEnum{
					Role:        role,
					Permissions: permissions,
				}
			}
		}

		jwtObj.AccessControls = parsedRoles
	} else {
		dep.logger.Error("missing or invalid roles in token claims")
		return nil, &utils.AuthenticationError{}
	}

	return jwtObj, nil
}
