package cistern

import (
	"fmt"
	"sync"
	"time"
	"os"
	"io"
	"encoding/gob"
)

type Item struct {
	Object     interface{} // 缓存数据项
	Expiration int64       // 数据项过期时间
}

func (i *Item) Expired() bool {
	if i.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > i.Expiration
}

type Cache struct {
	items             map[string]Item
	mu                sync.RWMutex
	defaultExpiration time.Duration
	gcInterval        time.Duration
	stopGC            chan bool
}

// 创建缓存实例
func NewCache(defaultExpiration, gcInterval time.Duration) *Cache {
	c := &Cache{
		defaultExpiration: defaultExpiration,
		gcInterval:        gcInterval,
		items:             map[string]Item{},
		stopGC:            make(chan bool),
	}
	go c.gcLoop() // 启动一个gorountine用于清理过期数据项
	return c
}

// 清理过期数据项
func (c *Cache) GcExpired() {
	now := time.Now().UnixNano()
	c.mu.Lock()
	defer c.mu.Unlock()

	for k, v := range c.items {
		if v.Expiration > 0 && now > v.Expiration {
			c.del(k)
		}
	}
}

func (c *Cache) Set(k string, v interface{}, d time.Duration) {
	var e int64
	if d == 0 {
		d = c.defaultExpiration
	}
	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[k] = Item{
		Object:     v,
		Expiration: e,
	}
}

func (c *Cache) Get(k string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.get(k)
}

func (c *Cache) Add(k string, v interface{}, d time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, found := c.get(k)
	if found {
		return fmt.Errorf("Item %s already exist", k)
	}
	c.set(k, v, d)
	return nil
}

func (c *Cache) Delete(k string) {
	c.mu.Lock()
	c.del(k)
	c.mu.Unlock()
}

// 清空缓存
func (c *Cache) Clean() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = map[string]Item{}
}

func (c *Cache) Count() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.items)
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
	if err = c.Read(f); err != nil{
		f.Close()
		return err
	}
	return f.Close()
}

func (c *Cache) Read(r io.Reader) error {
	decoder := gob.NewDecoder(r)
	items := map[string]Item{}
	err := decoder.Decode(&items)
	if err == nil {
		c.mu.Lock()
		defer c.mu.Unlock()
		for k, v := range items {
			ov, found := c.items[k]
			// 此处将缓存中已经过期或不存在的k从IO中读入
			// 即缓存中已经存在的k不会被IO中相同的k覆盖
			if !found || ov.Expired() {
				c.items[k] = v
			}
		}
	}
	return err
}

func (c *Cache) Write(w io.Writer) (err error) {
	encoder := gob.NewEncoder(w)
	defer func() {
		if x := recover(); x != nil {
			err = fmt.Errorf("error occurred when registering items type with Gob library")
		}
	}()
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, v := range c.items {
		gob.Register(v.Object)
	}
	err = encoder.Encode(&c.items)
	return
}

func (c *Cache) StopGC() {
	c.stopGC <- true
}

// 删除数据项
func (c *Cache) del(k string) {
	delete(c.items, k)
}

// 不加锁的方法，内部调用
func (c *Cache) get(k string) (interface{}, bool) {
	item, found := c.items[k]
	if !found {
		return nil, false
	}
	if item.Expired() {
		return nil, false
	}
	return item.Object, true
}

func (c *Cache) set(k string, v interface{}, d time.Duration) {
	var e int64
	if d == 0 {
		d = c.defaultExpiration
	}
	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}
	c.items[k] = Item{
		Object:     v,
		Expiration: e,
	}
}

func (c *Cache) gcLoop() {
	ticker := time.NewTicker(c.gcInterval)
	for {
		select {
		case <-ticker.C:
			c.GcExpired()
		case <-c.stopGC:
			ticker.Stop()
			return
		}
	}
}
