package config

import (
	"cmp"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"os"
)

type SecretVariables struct {
	Profile            string
	Address            string
	DbConnectionString string
	PrivateKeyPath     string
	PublicKeyPath      string
	CookieParam        *utils.CookieParam
	FrontEnd           string
}

func (dep *SecretVariables) Config() *SecretVariables {
	profile := cmp.Or(os.Getenv("PROFILE"), "development")
	return &SecretVariables{
		Address:            cmp.Or(os.Getenv("PORT"), ":8080"),
		Profile:            profile,
		DbConnectionString: cmp.Or(os.Getenv("DB_CONN"), "postgres://reservation:reservation@localhost:5432/reservation_db?sslmode=disable"),
		PrivateKeyPath:     cmp.Or(os.Getenv("JWT_PRIV_PATH"), "./keys/private.key"),
		PublicKeyPath:      cmp.Or(os.Getenv("JWT_PUB_PATH"), "./keys/public.key"),
		FrontEnd:           cmp.Or(os.Getenv("FRONTEND"), "http://localhost:4200"),
		CookieParam: &utils.CookieParam{
			CookieName:   cmp.Or(os.Getenv("COOKIENAME"), "RESERVATION_COOKIE"),
			CookieSecure: profile == "production",
			CookieDomain: cmp.Or(os.Getenv("COOKIEDOMAIN"), "localhost"),
		},
	}
}
