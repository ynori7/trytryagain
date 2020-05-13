package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/ynori7/trytryagain"
)

func main() {
	// set up a dummy http server which will return an error twice and then a success
	errCount := 2
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if errCount == 0 {
			rw.Write([]byte("success!"))
		} else {
			http.Error(rw, "something went wrong", 500)
			errCount--
		}
	}))
	defer server.Close()

	// build the dummy request
	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		log.Fatal("failed to create request")
	}

	// prepare the business logic
	var resp *http.Response // here's where we'll put the final result

	onError := func(err error) {
		log.Printf("Got error: %s\n", err.Error())
	}
	action := func() (error, bool) {
		var err error
		resp, err = server.Client().Do(req)
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
		log.Printf("`Do` failed with error: %s\n", err)
	} else {
		responseBody, _ := ioutil.ReadAll(resp.Body)
		log.Println("Got success response: ", string(responseBody))
	}
}
