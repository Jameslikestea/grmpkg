package server

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gorilla/mux"
	"golang.org/x/mod/modfile"
	"grmpkg.com/grmpkg/internal/model"
	"grmpkg.com/grmpkg/internal/validator"
)

func (s *Server) UploadHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		packageName := mux.Vars(r)["package"]
		packageVersion := mux.Vars(r)["version"]

		s.logger.Infow("Uploading package", "name", packageName, "version", packageVersion)
		s.logger.Infow("package size", "size", r.ContentLength)

		bs, err := ioutil.ReadAll(r.Body)
		s.logger.Infow("body content", "body", string(bs[:10]))
		if err != nil {
			s.logger.Fatalw("Error reading body", "error", err)
			return
		}

		// Although the regex should not pass a valid version, we're still going to check here.
		if !validator.ValidateVersion(packageVersion) {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("Version Not Valid"))
			return
		}

		reader, err := zip.NewReader(bytes.NewReader(bs), int64(len(bs)))
		if err != nil {
			s.logger.Errorw("could not decode zip", "error", err)
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("Could not decode zip"))
			return
		}

		foundMod := false

		for _, file := range reader.File {
			s.logger.Infow("found file", "package", packageName, "version", packageVersion, "file", file.Name)

			if strings.HasSuffix(file.Name, "go.mod") {
				reader, err := file.Open()
				if err != nil {
					s.logger.Errorw("cannot open go.mod")
				}
				bytes, err := ioutil.ReadAll(reader)
				if err != nil {
					s.logger.Errorw("cannot open go.mod")
				}
				if modfile.ModulePath(bytes) != packageName {
					s.logger.Errorw("module name does not match")
					rw.WriteHeader(http.StatusBadRequest)
					rw.Write([]byte("package names do not match"))
					return
				}
				s.uploadToS3(fmt.Sprintf("%s@%s.mod", packageName, packageVersion), bytes)
				foundMod = true
				break
			}
		}

		if !foundMod {
			s.logger.Errorw("no go.mod file found")
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("no go.mod file"))
			return
		}

		s.uploadToS3(fmt.Sprintf("%s@%s.zip", packageName, packageVersion), bs)
		s.updatePackage(packageName, packageVersion)

		rw.Header().Set("content-type", "application/zip")
		rw.WriteHeader(200)
		rw.Write([]byte("200 OK"))
	}
}

func (s *Server) uploadToS3(key string, data []byte) error {
	s3s := s3.New(s.session)

	_, err := s3s.PutObject(&s3.PutObjectInput{
		Bucket: aws.String("grmpkg"),
		Key:    &key,
		Body:   bytes.NewReader(data),
	})

	return err
}

func (s *Server) updatePackage(packageName, packageVersion string) error {
	ddb := dynamodb.New(s.session)

	obj, err := ddb.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("grmpkg"),
		Key: map[string]*dynamodb.AttributeValue{
			"name": {
				S: &packageName,
			},
		},
	})
	if err != nil {
		return err
	}

	var pack model.Package

	err = dynamodbattribute.UnmarshalMap(obj.Item, &pack)

	if err != nil {
		return err
	}

	return nil
}
