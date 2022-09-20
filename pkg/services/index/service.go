package index

import (
	"habiko-go/pkg/templates"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi"

	"habiko-go/pkg/services"
)

type Service interface {
	GetIndexPage() []byte
	ServeStatic(r chi.Router, path string) error
}

type service struct {
	providers *services.Providers
}

func (s *service) GetIndexPage() []byte {
	generatedHTML := templates.PageTemplate(&templates.IndexPage{
		BasePage: templates.BasePage{
			Title:         "Index",
			Description:   "",
			AssetsVersion: "1",
		},
	})

	return []byte(generatedHTML)
}

func (s *service) ServeStatic(r chi.Router, path string) error {
	workDir, err := os.Getwd()
	if err != nil {
		return err
	}

	filesDir := filepath.Join(workDir, "static/dist")
	root := http.Dir(filesDir)
	fs := http.StripPrefix("/", http.FileServer(root))

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "max-age=2592000")
		fs.ServeHTTP(w, r)
	}))
	return nil
}

func NewService(providers *services.Providers) Service {
	return &service{providers: providers}
}
