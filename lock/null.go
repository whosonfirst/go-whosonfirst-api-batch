package lock

import (
	"github.com/whosonfirst/go-whosonfirst-api-batch"
)

type NullLock struct {
	batch.BatchRequestLock
}

func NewNullLock() (*NullLock, error) {

	l := NullLock{}
	return &l, nil
}

func (l *NullLock) Get(k *batch.BatchRequestKey) (bool, error) {
	return false, nil
}

func (l *NullLock) Set(k *batch.BatchRequestKey) error {
	return nil
}

func (l *NullLock) Unset(k *batch.BatchRequestKey) error {
	return nil
}
