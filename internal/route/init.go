package route


import (
	"embed"
	"errors"
	"github.com/gin-gonic/gin"
	"io/fs"
	"net/http"
	"path"
	"path/filepath"
	"strings"
)

type Resource struct {
	fs embed.FS
	path string
}

func NewResource() *Resource {
	return &Resource{
		fs: Static,
		path: "html",
	}
}

func (r *Resource) Open(name string) (fs.File, error) {
	if filepath.Separator != '/' && strings.ContainsRune(name, filepath.Separator) {
		return nil, errors.New("http: invalid character in file path")
	}
	fullName := filepath.Join(r.path, filepath.FromSlash(path.Clean("/static/" + name)))
	file, err := r.fs.Open(fullName)

	return file, err
}

func InitResource(engine *gin.Engine) *gin.Engine {
	engine.StaticFS("/static", http.FS(NewResource()))
	return engine
}