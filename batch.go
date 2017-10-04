package batch

import (
	"github.com/whosonfirst/go-whosonfirst-api"
	"time"
)

type BatchRequest map[string]string

type BatchResponse struct {
	Index       int
	APIResponse api.APIResponse
	Timing      time.Duration
}

type BatchRequestSet struct {
	APIKey   string
	Requests []BatchRequest
}

type BatchResponseSet struct {
     Responses []interface{}
     Timing    time.Duration
}     
