package batch

import (
       "fmt"
	"github.com/whosonfirst/go-whosonfirst-api"
	"time"
)

type BatchRequest map[string]string

type BatchRequestKey struct {
     APIKey string
     InputHash string
}

func (k BatchRequestKey) String() string {
     return fmt.Sprintf("%s#%s", k.APIKey, k.InputHash)
}

type BatchResponse struct {
	Index       int
	APIResponse api.APIResponse
	Timing      time.Duration
}

type BatchRequestSet struct {
	APIKey   string
	Requests []BatchRequest
	RequestKey BatchRequestKey
}

type BatchResponseSet struct {
     Responses []interface{}
     RequestKey	BatchRequestKey
     Timing    time.Duration
}     
