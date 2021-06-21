package server

import (
	"encoding/json"
	"html/template"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

type Server struct {
	logger *zap.SugaredLogger
}

type Package struct {
	Scope   string `json:"scope"`
	Package string `json:"package"`
}

type Info struct {
	Name    string `json:"Name"`
	Short   string `json:"Short"`
	Version string `json:"Version"`
	Time    string `json:"Time"`
}

func New(logger *zap.SugaredLogger) *Server {
	return &Server{
		logger: logger,
	}
}

const proxy string = `
	<html>
		<head>
			<meta name="go-import" content="local.grmpkg.com/{{.Scope}}/{{.Package}} mod http://local.grmpkg.com/" />
		</head>
		<body><a href="/{{.Scope}}/{{.Package}}">Package</a></body>
	</html>
`

func (s *Server) Start() {
	s.logger.Infow("starting server")

	r := chi.NewRouter()

	r.Route("/{scope}/{package}", func(r chi.Router) {
		r.Get("/", func(rw http.ResponseWriter, r *http.Request) {
			p := Package{
				Scope:   chi.URLParam(r, "scope"),
				Package: chi.URLParam(r, "package"),
			}
			s.logger.Infow("getting package", "packageInfo", p)
			tpl, _ := template.New("template").Parse(proxy)
			tpl.Execute(rw, p)
		})
	})

	r.Route("/{domain}/{scope}/{package}", func(r chi.Router) {
		r.Get("/@v/list", func(rw http.ResponseWriter, r *http.Request) {
			s.logger.Infow("getting version")
			rw.Write([]byte("v1.1.0\nv1.2.0\nv1.1.2"))
		})
		r.Get("/@v/{version:.*}info", func(rw http.ResponseWriter, r *http.Request) {
			version := chi.URLParam(r, "version")
			version = version[:len(version)-1]
			s.logger.Infow("getting info", "version", version, "path", r.URL.Path)
			json.NewEncoder(rw).Encode(Info{
				Name:    version,
				Short:   version,
				Version: version,
				Time:    time.Now().Format("2006-01-02T15:04:05Z07:00"),
			})
		})
		r.Get("/@v/{version:.*}mod", func(rw http.ResponseWriter, r *http.Request) {
			version := chi.URLParam(r, "version")
			version = version[:len(version)-1]
			s.logger.Infow("getting mod", "version", version, "path", r.URL.Path)
			json.NewEncoder(rw).Encode(Info{
				Name:    version,
				Short:   version,
				Version: version,
				Time:    time.Now().Format("2006-01-02T15:04:05Z07:00"),
			})
		})
		r.Get("/@v/{version:.*}zip", func(rw http.ResponseWriter, r *http.Request) {
			version := chi.URLParam(r, "version")
			version = version[:len(version)-1]
			s.logger.Infow("getting zip", "version", version, "path", r.URL.Path)
			json.NewEncoder(rw).Encode(Info{
				Name:    version,
				Short:   version,
				Version: version,
				Time:    time.Now().Format("2006-01-02T15:04:05Z07:00"),
			})
		})
	})

	http.ListenAndServe(":80", r)

	s.logger.Infow("stopping server")
}
