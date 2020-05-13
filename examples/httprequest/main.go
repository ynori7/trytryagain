package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/ynori7/trytryagain"
)

func main() {
	// set up a dummy http server which always returns an error
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		http.Error(rw, "something went wrong", 500)
	}))
	defer server.Close()

	// build the dummy request
	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		log.Fatal("failed to create request")
	}

	// prepare the business logic
	onError := func(err error) {
		log.Printf("Got error: %s\n", err.Error())
	}
	action := func() (err error, retriable bool) {
		resp, err := server.Client().Do(req)
		if err != nil {
			return err, true
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf(resp.Status), resp.StatusCode >= 500
		}

		return nil, false
	}

	// prepare the retrier and do the action with retries
	err = trytryagain.NewRetrier(trytryagain.WithOnError(onError)).
		Do(context.Background(), action)

	if err != nil {
		log.Printf("Do failed with error: %s\n", err)
	} else {
		log.Println("Success!")
	}
}
