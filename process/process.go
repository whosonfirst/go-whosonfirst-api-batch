package process

import (
	"encoding/json"
	"errors"
	"github.com/whosonfirst/go-whosonfirst-api"
	"github.com/whosonfirst/go-whosonfirst-api-batch"
	"github.com/whosonfirst/go-whosonfirst-api-batch/parse"
	"github.com/whosonfirst/go-whosonfirst-api/client"
	"github.com/whosonfirst/go-whosonfirst-api/endpoint"
	"github.com/whosonfirst/go-whosonfirst-hash"
	"log"
	"net/url"
	"time"
)

// see this function signature? it _will_ change... (20171004/thisisaaronland)

func ProcessBatch(input []byte, api_key string) (*batch.BatchResponseSet, error) {

	// please wrap all of this in a library somewhere...

	hasher, err := hash.NewWOFHash()

	if err != nil {
		return nil, err
	}

	input_hash, err := hasher.HashBytes(input)

	if err != nil {
		return nil, err
	}

	request_key := batch.BatchRequestKey{
		APIKey:    api_key,
		InputHash: input_hash,
	}

	// check to see request_key isn't already being processed

	parse_opts := parse.NewDefaultParseRequestOptions()

	requests, err := parse.ParseRequest(input, parse_opts)

	if err != nil {
		return nil, err
	}

	request_set := batch.BatchRequestSet{
		APIKey:     api_key,
		Requests:   requests,
		RequestKey: request_key,
	}

	// log api_key + "#" + hash here - it would be nice to all of this using
	// BatchRequestSet but that means always parsing body first...

	// see notes above wrt a timeout context (as in: it does not exist yet)

	return ProcessRequestSet(request_set)
}

// please for to be passing in a timeout context here and to make
// sure it bubbles down to any individual requests being processed

func ProcessRequestSet(rs batch.BatchRequestSet) (*batch.BatchResponseSet, error) {

	// maybe just pass the client as part of the function call?

	e, err := endpoint.NewMapzenAPIEndpoint(rs.APIKey)

	if err != nil {
		return nil, err
	}

	cl, err := client.NewHTTPClient(e)

	if err != nil {
		return nil, err
	}

	response_ch := make(chan batch.BatchResponse)
	error_ch := make(chan error)
	done_ch := make(chan bool)

	complete_ch := make(chan bool)

	pending := len(rs.Requests)
	responses := make([]interface{}, pending) // see below because []interface{} is a stop-gap

	go func() {

		for pending > 0 {

			select {
			case rsp := <-response_ch:

				// see this? it's basically just so that we can easily serialize the reponses
				// in the main handler - it does not take in to account any of the work that's
				// been done around SPRs or non-JSON formatted responses... it should, although
				// please don't ask me what that means yet... (20171004/thisisaaronland)

				var i interface{}
				json.Unmarshal(rsp.APIResponse.Raw(), &i)
				responses[rsp.Index] = i

			case err := <-error_ch:
				// please to be inserting an error in responses[idx] here
				log.Println(err)
			case <-done_ch:
				pending -= 1
			}
		}

		complete_ch <- true
	}()

	t1 := time.Now()

	for idx, req := range rs.Requests {

		// please for to be rate-limiting here...

		go ProcessRequest(cl, idx, req, response_ch, error_ch, done_ch)
	}

	<-complete_ch

	t2 := time.Since(t1)

	response_set := batch.BatchResponseSet{
		RequestKey: rs.RequestKey,
		Responses:  responses,
		Timing:     t2,
	}

	return &response_set, nil
}

func ProcessRequest(cl api.APIClient, idx int, req batch.BatchRequest, response_ch chan batch.BatchResponse, error_ch chan error, done_ch chan bool) {

	defer func() {
		done_ch <- true
	}()

	t1 := time.Now()

	cb := func(rsp api.APIResponse) error {

		t2 := time.Since(t1)

		response := batch.BatchResponse{
			Index:       idx,
			APIResponse: rsp,
			Timing:      t2,
		}

		response_ch <- response
		return nil
	}

	method := ""
	args := url.Values{}

	for k, v := range req {

		if k == "method" {
			method = v
			continue
		}

		args.Set(k, v)
	}

	if method == "" {
		error_ch <- errors.New("Missing API method")
		return
	}

	err := cl.ExecuteMethodWithCallback(method, &args, cb)

	if err != nil {
		error_ch <- err
	}
}
