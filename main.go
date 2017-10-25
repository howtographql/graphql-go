package main

import (
	"log"
	"net/http"

	"github.com/neelance/graphql-go"
	"github.com/neelance/graphql-go/relay"
	"github.com/howtographql/graphql-go/resolvers"
	"fmt"
	"io/ioutil"
)

var schema *graphql.Schema

func init() {
	schemaFile, err := ioutil.ReadFile("schema.graphqls")
	if err != nil {
		panic(err)
	}

	schema = graphql.MustParseSchema(string(schemaFile), &resolvers.Resolver{})
}

func main() {
	http.Handle("/graphql", &relay.Handler{Schema: schema})

	fmt.Println("server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
