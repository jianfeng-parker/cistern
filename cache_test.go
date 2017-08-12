package cistern

import (
	"testing"
	"time"
	"github.com/jianfeng-parker/cistern"
)

// 从缓存中获取对象
func TestCache_Get(t *testing.T) {
	defaultExpired, _ := time.ParseDuration("1m")
	cleanInterval, _ := time.ParseDuration("3s")
	c := cistern.NewCache(defaultExpired, cleanInterval)
	k1 := "k1"
	v1 := "v1"
	expiration, _ := time.ParseDuration("5s")
	c.Set(k1, v1, expiration)
	if _, found := c.Get(k1); !found {
		t.Errorf("not found k1")
	}
}


func TestCache_Get_With_Expiration(t *testing.T) {
	defaultExpired, _ := time.ParseDuration("1m")
	cleanInterval, _ := time.ParseDuration("10s")
	c := cistern.NewCache(defaultExpired, cleanInterval)
	e, _ := time.ParseDuration("1s")
	c.Set("name", "吴建峰", e)
	sleep, _ := time.ParseDuration("2s")
	time.Sleep(sleep)
	if _, found := c.Get("name"); found {
	   t.Errorf("data should be expired")
	}
}
