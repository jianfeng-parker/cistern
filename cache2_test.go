package cistern

import (
	"testing"
)

// 从缓存中获取对象
func TestCache2_Get(t *testing.T) {
	c := NewCache2(65535)
	k1 := "k1"
	v1 := "v1"
	c.Set(k1, []byte(v1), 30)
	if _, found := c.Get(k1); !found {
		t.Errorf("not found k1")
	}
}
