package server

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gorilla/mux"
	"grmpkg.com/grmpkg/internal/model"
)

var tpl string = `
<!DOCTYPE html>
<html>
	<head>
		<title>{{.Name}}</title>
		<meta name="go-import" content="{{.Name}} mod https://{{.Info.Hostname}}">
	</head>
	<body>
		Download with

		<pre>
			go get {{.Name}}
		</pre>
	</body>
</html>
`

func (s *Server) PackageHTMLHandler() http.HandlerFunc {
	ddb := dynamodb.New(s.session)

	return func(rw http.ResponseWriter, r *http.Request) {
		s.logger.Infow("Getting package", "package", mux.Vars(r)["package"])
		packageName := fmt.Sprintf("%s/%s", r.Host, mux.Vars(r)["package"])

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

		temp := template.New("body")
		tempRender, _ := temp.Parse(tpl)
		tempRender.Execute(rw, pack)
	}
}
