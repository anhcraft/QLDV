package storage

import (
	"github.com/cheshir/ttlcache"
	"time"
)

var cache *ttlcache.Cache

func init() {
	cache = ttlcache.New(time.Minute)
	// TODO Save cache to disk?
}

func GetCache(namespace string, key string) (interface{}, bool) {
	return cache.Get(ttlcache.StringKey(namespace + key))
}

func SetCache(namespace string, key string, val interface{}, duration time.Duration) {
	cache.Set(ttlcache.StringKey(namespace+key), val, duration)
}

func PurgeCache(namespace string, key string) {
	cache.Delete(ttlcache.StringKey(namespace + key))
}
