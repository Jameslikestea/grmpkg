package server

import (
	"bytes"
	"fmt"
	"net/http"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gorilla/mux"
	"grmpkg.com/grmpkg/internal/model"
)

func (s *Server) ListVersionsHandler() http.HandlerFunc {
	ddb := dynamodb.New(s.session)

	return func(rw http.ResponseWriter, r *http.Request) {
		packageName := mux.Vars(r)["package"]
		s.logger.Infow("Handling function", "package", packageName)

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

		sort.Sort(pack.Versions)

		s.logger.Infow("got versions", "versions", len(pack.Versions))

		rw.Header().Set("package-manager", "grm-pkg")
		rw.WriteHeader(200)
		bs := bytes.NewBuffer([]byte{})
		for _, ver := range pack.Versions {
			bs.WriteString(fmt.Sprintf("%s\n", ver.Version))
		}
		rw.Write(bs.Bytes())
	}
}
