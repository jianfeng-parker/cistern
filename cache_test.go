package cistern

import (
	"testing"
	"time"
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

// 过期数据不该被GET到
func TestCache2_Get_With_Expiration(t *testing.T) {
	c := NewCache2(100)
	c.Set("name", []byte("吴建峰"), 1)
	sleep, _ := time.ParseDuration("2s")
	time.Sleep(sleep)
	if _, found := c.Get("name"); found {
		t.Errorf("data should be expired")
	}
}

//// 缓存数据写、读 文件，UT运行完将文件删除
//func TestCache2_WriteFile(t *testing.T) {
//	c := NewCache2(100)
//	c.Set("name", []byte("吴建峰"), 20)
//	if err := c.WriteFile("/workspace/tmp/log"); err != nil {
//		t.Errorf("write cache to file failure:\n")
//	} else {
//		c2 := NewCache2(100)
//		c2.ReadFile("/workspace/tmp/log")
//		if v, found := c2.Get("name"); !found {
//			t.Errorf("could not read from file")
//		} else {
//			t.Logf("read data:%s from file", v)
//		}
//		os.Remove("/workspace/tmp/log")
//	}
//}
