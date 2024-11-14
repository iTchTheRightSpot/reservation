package config

type SecretVariables struct {
	Profile            string
	Address            string
	DbConnectionString string
}

func (dep *SecretVariables) Config() (*SecretVariables, error) {
	return &SecretVariables{
		Address:            ":8080",
		Profile:            "development",
		DbConnectionString: "postgres://mrp:mrp@localhost:5432/mrp_db?sslmode=disable",
	}, nil
}
