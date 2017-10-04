package http

import (
	"encoding/json"
	"github.com/whosonfirst/go-whosonfirst-api-batch/process"
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

		// something something something check the size of req.Body before
		// reading it all in to memory or at least during the reading...

		input, err := ioutil.ReadAll(req.Body)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
			return
		}

		process_opts := process.NewDefaultProcessBatchOptions()
		process_opts.APIKey = api_key

		response_set, err := process.ProcessBatch(input, process_opts)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		// please log this somewhere...

		log.Println("TIMING", response_set.RequestKey.String(), response_set.Timing)

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
