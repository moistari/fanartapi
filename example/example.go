// example/example.go
package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/moistari/fanartapi"
)

func main() {
	api := flag.String("api", "", "api key")
	flag.Parse()
	cl := fanartapi.New(fanartapi.WithApiKey(*api))
	res, err := cl.Images(context.Background(), fanartapi.Movie, "tt0137523")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%q (%s)\n", res.Name, res.ID())
}
