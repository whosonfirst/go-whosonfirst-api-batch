package lock

import (
	gocache "github.com/patrickmn/go-cache"
	"github.com/whosonfirst/go-whosonfirst-api-batch"
	"sync"
	"time"
)

var cache *gocache.Cache
var mu *sync.RWMutex

func init() {
	cache = gocache.New(30*time.Minute, 60*time.Minute)
	mu = new(sync.RWMutex)
}

type GoCacheLock struct {
	batch.BatchRequestLock
}

func NewGoCacheLock() (*GoCacheLock, error) {

	l := GoCacheLock{}
	return &l, nil
}

func (l *GoCacheLock) Get(k *batch.BatchRequestKey) (bool, error) {

	mu.RLock()
	defer mu.RUnlock()

	_, found := cache.Get(k.String())
	return found, nil
}

func (l *GoCacheLock) Set(k *batch.BatchRequestKey) error {

	mu.Lock()
	defer mu.Unlock()

	cache.Set(k.String(), time.Now(), gocache.DefaultExpiration)
	return nil
}

func (l *GoCacheLock) Unset(k *batch.BatchRequestKey) error {

	mu.Lock()
	defer mu.Unlock()

	cache.Delete(k.String())
	return nil
}
