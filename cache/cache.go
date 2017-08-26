package cache

import (
	"time"
	"fmt"
	"os"
	"io"
	"encoding/gob"
	"github.com/spaolacci/murmur3"
)

const DefaultExpiration int64 = 30

type Cache struct {
	segments [256]segment
}

// 清理过期数据项
func (c *Cache) ClearExpired() {
	for _, s := range c.segments {
		s.clearExpired()
	}
}

func (c *Cache) Set(k string, v []byte, expiration int64) error {
	return c.segments[hashID(k)].set(k, v, expiration)
}

func (c *Cache) Get(k string) ([]byte, bool) {
	return c.segments[hashID(k)].get(k)
}

func (c *Cache) Add(k string, v []byte, expiration int64) error {
	return c.segments[hashID(k)].add(k, v, expiration)

}

func (c *Cache) Delete(k string) {
	c.segments[hashID(k)].del(k)
}

// 清空缓存
func (c *Cache) Clean() {
	for _, s := range c.segments {
		s.clean()
	}
}
func (c *Cache) Count() int {
	var count int
	for _, s := range c.segments {
		count += s.count()
	}
	return count
}

func (c *Cache) Expired(k string) bool {
	return c.segments[hashID(k)].expired(k)
}

func (c *Cache) WriteFile(file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	if err = c.Write(f); err != nil {
		return err
	}
	return f.Close()
}

func (c *Cache) ReadFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	if err = c.Read(f); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

func (c *Cache) Read(r io.Reader) error {
	decoder := gob.NewDecoder(r)
	segments := make([]segment, 256)
	err := decoder.Decode(&segments)
	if err == nil {
		for i, s := range segments {
			for k, n := range s.nodes {
				if _, found := c.Get(k); !found || c.Expired(k) {
					c.setNode(k, n)
				}
			}
			c.segments[i] = s
		}
	}
	return err
}

func (c *Cache) setNode(k string, n node) {
	c.segments[hashID(k)].setNode(k, n)
}

func (c *Cache) Write(w io.Writer) (err error) {
	encoder := gob.NewEncoder(w)
	defer func() {
		if x := recover(); x != nil {
			err = fmt.Errorf("error occurred when registering items type with Gob library")
		}
	}()
	for _, s := range c.segments {
		gob.Register(s.nodes)
	}
	err = encoder.Encode(&c.segments)
	return
}
func hashID(k string) uint64 {
	return murmur3.Sum64([]byte(k)) & 255
}

func NewCache(size int) (cache *Cache) {
	if size < 256*1024 {
		size = 256 * 1024
	}
	cache = new(Cache)
	for i := 0; i < 256; i++ {
		cache.segments[i] = newSegment(size/256, i)
	}
	return
}

type node struct {
	value    []byte
	expireAt int64
}

func (n *node) expired() bool {
	return time.Now().Unix() > n.expireAt
}
