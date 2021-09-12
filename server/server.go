package server

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Server struct {
	logger  *zap.SugaredLogger
	session *session.Session
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
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1"),
	})
	if err != nil {
		logger.Fatalw("cannot start server", "error", err)
	}

	if err != nil {
		logger.Errorw("cannot create table", "error", err)
	}

	return &Server{
		logger:  logger,
		session: sess,
	}
}

func (s *Server) Start() {
	s.logger.Infow("starting server")

	r := mux.NewRouter()
	subRouter := r.PathPrefix("/{package:[A-Za-z0-9/.]+}/@v/").Subrouter()

	r.HandleFunc("/{package:[A-Za-z0-9/.]+}", s.PackageHTMLHandler()).Methods(http.MethodGet)
	r.HandleFunc("/{package:[A-Za-z0-9/.]+}/@latest", s.PackageLatestHandler()).Methods(http.MethodGet)

	r.HandleFunc("/{package:[A-Za-z0-9/.]+}/{version:v(?:[0-9]+)\\.(?:[0-9]+)\\.(?:[0-9]+)(?:-.+)?}", s.UploadHandler()).Methods(http.MethodPost)

	subRouter.HandleFunc("/list", s.ListVersionsHandler()).Methods(http.MethodGet)
	subRouter.HandleFunc("/{version:v(?:[0-9]+)\\.(?:[0-9]+)\\.(?:[0-9]+)(?:-.+)?}.info", s.PackageInfoHandler()).Methods(http.MethodGet)
	subRouter.HandleFunc("/{version:v(?:[0-9]+)\\.(?:[0-9]+)\\.(?:[0-9]+)(?:-.+)?}.mod", s.PackageModHandler()).Methods(http.MethodGet)
	subRouter.HandleFunc("/{version:v(?:[0-9]+)\\.(?:[0-9]+)\\.(?:[0-9]+)(?:-.+)?}.zip", s.PackageZipHandler()).Methods(http.MethodGet)

	http.ListenAndServe(":80", r)

	s.logger.Infow("stopping server")
}
