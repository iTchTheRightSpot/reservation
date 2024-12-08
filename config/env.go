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
}

func (dep *SecretVariables) Config() (*SecretVariables, error) {
	profile := cmp.Or(os.Getenv("PROFILE"), "development")
	return &SecretVariables{
		Address:            cmp.Or(os.Getenv("PORT"), ":8080"),
		Profile:            profile,
		DbConnectionString: cmp.Or(os.Getenv("DB_CONN"), "postgres://mrp:mrp@localhost:5432/mrp_db?sslmode=disable"),
		PrivateKeyPath:     cmp.Or(os.Getenv("JWT_PRIV_PATH"), "./keys/private.key"),
		PublicKeyPath:      cmp.Or(os.Getenv("JWT_PUB_PATH"), "./keys/public.key"),
		CookieParam: &utils.CookieParam{
			CookieName:   cmp.Or(os.Getenv("COOKIENAME"), "ERPCOOKIE"),
			CookieSecure: profile == "production",
			CookieDomain: cmp.Or(os.Getenv("COOKIEDOMAIN"), "localhost"),
		},
	}, nil
}
