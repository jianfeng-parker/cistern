package main

import (
	"cistern"
	"fmt"
	"time"
)

func main() {
	defaultExpired, _ := time.ParseDuration("1m")
	cleanInterval, _ := time.ParseDuration("3s")
	c := cistern.NewCache(defaultExpired, cleanInterval)
	k1 := "k1"
	v1 := "v1"
	expiration,_:= time.ParseDuration("5s")
	c.Set(k1, v1, expiration)
	if v, found := c.Get(k1); found {
		fmt.Println("Found k1:" + v)
	} else {
		fmt.Println("Not found k1")
	}
     sleep,_ :=time.ParseDuration("10s")
	time.Sleep(sleep)
	if v, found := c.Get(k1); found {
		fmt.Println("Found k1:" + v)
	} else {
		fmt.Println("Not found k1")
	}
}
