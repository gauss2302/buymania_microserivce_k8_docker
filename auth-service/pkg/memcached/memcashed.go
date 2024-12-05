package memcached

import "github.com/bradfitz/gomemcache/memcache"

type Client struct {
	client *memcache.Client
}

func NewClient(port, host string) *Client {
	return &Client{
		client: memcache.New(host + ":" + port),
	}
}
