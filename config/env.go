package config

import (
	"cmp"
	"os"
)

type SecretVariables struct {
	Profile            string
	Address            string
	DbConnectionString string
	PrivateKeyPath     string
	PublicKeyPath      string
}

func (dep *SecretVariables) Config() (*SecretVariables, error) {
	return &SecretVariables{
		Address:            cmp.Or(os.Getenv("PORT"), ":8080"),
		Profile:            cmp.Or(os.Getenv("PROFILE"), "development"),
		DbConnectionString: cmp.Or(os.Getenv("DB_CONN"), "postgres://mrp:mrp@localhost:5432/mrp_db?sslmode=disable"),
		PrivateKeyPath:     cmp.Or(os.Getenv("JWT_PRIV_PATH"), "./keys/private.key"),
		PublicKeyPath:      cmp.Or(os.Getenv("JWT_PUB_PATH"), "./keys/public.key"),
	}, nil
}
