package config

import (
	"cmp"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"net/http"
	"os"
	"regexp"
)

type SecretVariables struct {
	Profile            string
	Address            string
	DbConnectionString string
	PrivateKeyPath     string
	PublicKeyPath      string
	CookieParam        *utils.CookieParam
	FrontEnd           string
	ApiPrefix          string
	SymmetricKey       string
	Discord            string
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
		SymmetricKey:       cmp.Or(os.Getenv("JWT_SECRET"), "jwt-secret"),
		CookieParam: &utils.CookieParam{
			CookieName:   cmp.Or(os.Getenv("COOKIENAME"), "GSESSIONID"),
			CookieSecure: profile == "production",
			CookieDomain: cookieDomain(),
			SameSite:     http.SameSiteLaxMode,
		},
		ApiPrefix: "/api/v1/",
		Discord:   cmp.Or(os.Getenv("DISCORD_WEBHOOK"), "http://localhost:4200"),
	}
}

func cookieDomain() string {
	val := cmp.Or(os.Getenv("COOKIEDOMAIN"), "localhost")
	if val == "localhost" {
		return val
	}
	// regex https://docs.spring.io/spring-session/reference/guides/java-custom-cookie.html
	return regexp.MustCompile(`^.+?\.(\w+\.[a-z]+)$`).FindStringSubmatch(val)[1]
}
