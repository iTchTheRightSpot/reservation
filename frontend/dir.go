package frontend

import (
	"context"
	"embed"
	"github.com/iTchTheRightSpot/utility/utils"
	"io/fs"
	"net/http"
)

//go:embed dist/frontend/browser
var frontend embed.FS

type FrontendStruct struct {
	Logger utils.ILogger
}

func (dep *FrontendStruct) FileSystem() (http.FileSystem, error) {
	build, err := fs.Sub(frontend, "dist/frontend/browser")
	if err != nil {
		dep.Logger.Error(context.Background(), err.Error())
		return nil, &utils.ServerError{}
	}
	return http.FS(build), nil
}
