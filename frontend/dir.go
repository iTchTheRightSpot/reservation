package frontend

import (
	"embed"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"io/fs"
	"net/http"
)

//go:embed dist/frontend/browser
var frontend embed.FS

type FrontendStruct struct {
	Logger utils.ILogger
}

func (dep *FrontendStruct) FileSystem() (http.FileSystem, error) {
	build, err := fs.Sub(frontend, "dist/ui/browser")
	if err != nil {
		dep.Logger.Error(err.Error())
		return nil, &utils.ServerError{}
	}
	return http.FS(build), nil
}
