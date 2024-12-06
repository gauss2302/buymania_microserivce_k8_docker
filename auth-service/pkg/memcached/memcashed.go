package memcached

import "github.com/bradfitz/gomemcache/memcache"

type MemcachedWrapper struct {
	Client *memcache.Client // Make the internal client public
}

func NewClient(host, port string) *MemcachedWrapper {
	return &MemcachedWrapper{
		Client: memcache.New(host + ":" + port),
	}
}
