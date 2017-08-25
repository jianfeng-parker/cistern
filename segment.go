package cistern

import (
	"fmt"
	"time"
	"errors"
	"sync"
)

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
