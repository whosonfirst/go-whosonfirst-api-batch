package request

import (
	"github.com/whosonfirst/go-whosonfirst-api-batch"
	"github.com/whosonfirst/go-whosonfirst-hash"
)

func NewRequestKey(api_key string, input []byte) (*batch.BatchRequestKey, error) {

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

	return &request_key, nil
}
