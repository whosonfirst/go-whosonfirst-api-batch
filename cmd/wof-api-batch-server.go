package main

// https://github.com/whosonfirst/whosonfirst-www-api/issues/99#issuecomment-333960724

import (
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-api-batch/http"
	"log"
	gohttp "net/http"
)

func main() {

	var host = flag.String("host", "localhost", "The hostname to listen for requests on")
	var port = flag.Int("port", 8080, "The port number to listen for requests on")

	flag.Parse()

	// fetch list of valid API methods from the api.spec.method and pass along
	// to handler for basic validation on all requests here...

	batch_handler, err := http.BatchHandler()

	if err != nil {
		log.Fatal(err)
	}

	mux := gohttp.NewServeMux()
	mux.Handle("/", batch_handler)

	endpoint := fmt.Sprintf("%s:%d", *host, *port)
	log.Printf("listening on %s\n", endpoint)

	err = gohttp.ListenAndServe(endpoint, mux)

	if err != nil {
		log.Fatal(err)
	}

}
