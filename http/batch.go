package http

import(
	"encoding/json"
	"github.com/whosonfirst/go-whosonfirst-api-batch"
	"io/ioutil"
	"log"
	gohttp "net/http"
	"strings"
)

func BatchHandler() (gohttp.Handler, error) {

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		query := req.URL.Query()

		api_key := query.Get("api_key")
		api_key = strings.Trim(api_key, " ")

		if api_key == "" {
			gohttp.Error(rsp, "Missing API key", gohttp.StatusBadRequest)
			return
		}

		input, err := ioutil.ReadAll(req.Body)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
			return
		}

		// hash body here and ensure that we don't already have a running
		// batch for api_key + "#" + hash

		requests, err := batch.Parse(input)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
			return
		}

		request_set := batch.BatchRequestSet{
			APIKey:   api_key,
			Requests: requests,
		}

		// log api_key + "#" + hash here - it would be nice to all of this using
		// BatchRequestSet but that means always parsing body first...
		
		// see notes above wrt a timeout context (as in: it does not exist yet)

		response_set, err := batch.ProcessBatch(request_set)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		log.Println("TIMING", response_set.Timing)

		// something something something non-JSON responses something something
		// something see above in parse_request for discussion about SPR...
		
		rsp.Header().Set("Content-Type", "application/json")
		rsp.Header().Set("Access-Control-Allow-Origin", "*")

		js, err := json.Marshal(response_set.Responses)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		rsp.Write(js)
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
