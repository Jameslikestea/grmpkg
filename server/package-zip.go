package server

import (
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gorilla/mux"
)

func (s *Server) PackageZipHandler() http.HandlerFunc {
	boto := s3.New(s.session)

	return func(rw http.ResponseWriter, r *http.Request) {
		s.logger.Infow("Getting package zip", "package", mux.Vars(r)["package"], "version", mux.Vars(r)["version"])

		object, err := boto.GetObject(&s3.GetObjectInput{
			Bucket: aws.String("grmpkg"),
			Key:    aws.String(mux.Vars(r)["package"] + "@" + mux.Vars(r)["version"] + ".zip"),
		})
		if err != nil {
			rw.Header().Set("package-manager", "grm-pkg")
			rw.WriteHeader(404)
			rw.Write([]byte("404 Not Found"))
			return
		}

		rw.Header().Set("package-manager", "grm-pkg")
		rw.Header().Set("content-type", "application/zip")
		rw.WriteHeader(200)
		io.Copy(rw, object.Body)
	}
}
