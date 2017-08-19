package cistern

import (
	"sync"
	"time"
	"fmt"
	"os"
	"io"
	"encoding/gob"
	"errors"
	"github.com/spaolacci/murmur3"
)

const DefaultExpiration int64 = 30

type Cache2 struct {
	segments [256]segment
}

// 清理过期数据项
func (c *Cache2) ClearExpired() {
	for _, s := range c.segments {
		s.clearExpired()
	}
}

func (c *Cache2) Set(k string, v []byte, expiration int64) error {
	return c.segments[hashID(k)].set(k, v, expiration)
}

func (c *Cache2) Get(k string) ([]byte, bool) {
	return c.segments[hashID(k)].get(k)
}

func (c *Cache2) Add(k string, v []byte, expiration int64) error {
	return c.segments[hashID(k)].add(k, v, expiration)

}

func (c *Cache2) Delete(k string) {
	c.segments[hashID(k)].del(k)
}

// 清空缓存
func (c *Cache2) Clean() {
	for _, s := range c.segments {
		s.clean()
	}
}
func (c *Cache2) Count() int {
	var count int
	for _, s := range c.segments {
		count += s.count()
	}
	return count
}

func (c *Cache2) Expired(k string) bool {
	return c.segments[hashID(k)].expired(k)
}

func (c *Cache2) WriteFile(file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	if err = c.Write(f); err != nil {
		return err
	}
	return f.Close()
}

func (c *Cache2) ReadFile(file string) error {
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

func (c *Cache2) Read(r io.Reader) error {
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

func (c *Cache2) setNode(k string, n node) {
	c.segments[hashID(k)].setNode(k, n)
}

func (c *Cache2) Write(w io.Writer) (err error) {
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

func NewCache2(size int) (cache *Cache2) {
	if size < 256*1024 {
		size = 256 * 1024
	}
	cache = new(Cache2)
	for i := 0; i < 256; i++ {
		cache.segments[i] = newSegment(size/256, i)
	}
	return
}

type segment struct {
	lock  sync.Mutex
	nodes map[string]node
	id    int
}

func newSegment(size, id int) (seg segment) {
	seg.id = id
	seg.nodes = make(map[string]node, size)
	return
}

func (s *segment) clearExpired() {
	now := time.Now().UnixNano()
	s.lock.Lock()
	defer s.lock.Unlock()
	for k, v := range s.nodes {
		if v.expireAt > 0 && now > v.expireAt {
			s.del(k)
		}
	}
}

func (s *segment) del(k string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.nodes, k)
}

func (s *segment) set(k string, v []byte, expiration int64) (err error) {
	var e int64
	if expiration < 0 {
		err = errors.New("invalid expiration for seconds")
	}
	if expiration == 0 {
		expiration = DefaultExpiration // 缓存默认过期时间:30s
	}
	if expiration > 0 {
		e = time.Now().Unix() + expiration
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	s.nodes[k] = node{
		value:    v,
		expireAt: e,
	}
	return
}

func (s *segment) get(k string) ([]byte, bool) {
	node, found := s.nodes[k]
	if !found {
		return nil, false
	}
	if node.expired() {
		return nil, false
	}
	return node.value, true
}

func (s *segment) add(k string, v []byte, expiration int64) error {
	if _, found := s.get(k); found {
		return fmt.Errorf("k:%s already exist", k)
	}
	return s.set(k, v, expiration)
}

func (s *segment) clean() {
	s.lock.Lock()
	defer s.lock.Unlock()
	for k := range s.nodes {
		delete(s.nodes, k)
	}
}

func (s *segment) count() int {
	s.lock.Lock()
	defer s.lock.Unlock()
	return len(s.nodes)
}

func (s *segment) expired(k string) bool {
	if n, found := s.nodes[k]; found {
		return n.expired()
	}
	return true
}

func (s *segment) setNode(k string, n node) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.nodes[k] = n
}

type node struct {
	value    []byte
	expireAt int64
}

func (n *node) expired() bool {
	return time.Now().Unix() > n.expireAt
}
