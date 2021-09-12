package server

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gorilla/mux"
	"golang.org/x/mod/semver"
	"grmpkg.com/grmpkg/internal/model"
)

func (s *Server) PackageInfoHandler() http.HandlerFunc {
	ddb := dynamodb.New(s.session)

	return func(rw http.ResponseWriter, r *http.Request) {
		s.logger.Infow("Getting package info", "package", mux.Vars(r)["package"], "version", mux.Vars(r)["version"])
		packageName := mux.Vars(r)["package"]
		packageVersion := mux.Vars(r)["version"]

		item, err := ddb.GetItem(&dynamodb.GetItemInput{
			TableName: aws.String("grmpkg"),
			Key: map[string]*dynamodb.AttributeValue{
				"name": {
					S: &packageName,
				},
			},
		})
		if err != nil {
			rw.WriteHeader(404)
			rw.Write([]byte("404 Not Found"))
			return
		}

		var pack model.Package
		dynamodbattribute.UnmarshalMap(item.Item, &pack)

		s.logger.Infow("package found", "package", pack, "name", packageName)

		v := model.PackageVersion{}
		for _, version := range pack.Versions {
			s.logger.Infow("checking version", "version", version.Version)
			if semver.Compare(version.Version, packageVersion) == 0 {
				s.logger.Infow("found version", "version", version.Version)
				v = version
				break
			}
		}

		rw.Header().Set("package-manager", "grm-pkg")
		rw.Header().Set("content-type", "application/json")
		if v.Name == "" {
			s.logger.Errorw("cannot find package version", "version", packageVersion)
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte("404 Not Found"))
			return
		}
		rw.WriteHeader(200)
		json.NewEncoder(rw).Encode(v)
	}
}
