package pkg

import (
	"context"
	"database/sql"
	_ "embed"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist/file-adapter"
	"log"
	"os"
)

//go:embed casbin_config/model.conf
var ModelConf string

//go:embed casbin_config/model.csv
var PolicyCSV string

var CasbinEnforcer *casbin.Enforcer

func init() {
	m, err := model.NewModelFromString(ModelConf)
	if err != nil {
		log.Fatalf("Failed to load model: %v", err)
	}

	tempPolicyFile, err := os.CreateTemp("", "policy.csv")
	if err != nil {
		log.Fatalf("Failed to create temporary file: %v", err)
	}

	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Println(err)
		}
	}(tempPolicyFile.Name())

	_, err = tempPolicyFile.Write([]byte(PolicyCSV))
	if err != nil {
		log.Fatalf("Failed to write policy to temporary file: %v", err)
	}

	if err := tempPolicyFile.Close(); err != nil {
		return
	}

	a := fileadapter.NewAdapter(tempPolicyFile.Name())

	e, err := casbin.NewEnforcer(m, a)
	if err != nil {
		log.Fatalf("%s", err)
	}
	CasbinEnforcer = e
}

type Db interface {
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}
