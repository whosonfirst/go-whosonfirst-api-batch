package parse

import (
	"errors"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-api-batch"
	"strconv"
)

type ParseRequestOptions struct {
	MaxRequests int
}

func NewDefaultParseRequestOptions() *ParseRequestOptions {

	opts := ParseRequestOptions{
		MaxRequests: 10,
	}

	return &opts
}

func ParseRequest(body []byte, opts *ParseRequestOptions) ([]batch.BatchRequest, error) {

	batch := make([]batch.BatchRequest, 0)

	c := gjson.GetBytes(body, "#")
	count := int(c.Int())

	if count == 0 {
		return nil, errors.New("Invalid batch request")
	}

	if count > opts.MaxRequests {
		return nil, errors.New("Too many requests")
	}

	for i := 0; i < count; i++ {

		path := strconv.Itoa(i)
		r := gjson.GetBytes(body, path)

		br := make(map[string]string)

		for k, v := range r.Map() {
			br[k] = v.String()
		}

		batch = append(batch, br)
	}

	if len(batch) == 0 {
		return nil, errors.New("Invalid batch request")
	}

	return batch, nil
}
