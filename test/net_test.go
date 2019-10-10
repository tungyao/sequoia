package test

import (
	"../caches"
	"testing"
	"time"
)

func TestRedis(t *testing.T) {
	//caches.HSet(caches.Cache{
	//	Key:"test",
	//	Value:"hello",
	//})
	r := caches.New()
	//r.HSet(caches.Cache{
	//	Key:   "test",
	//	Value: "testtest",
	//	Time:  0,
	//})
	time.Sleep(time.Second)
	for i := 0; i < 100; i++ {
		r.HGet("test")
	}
	//log.Println(fmt.Sprint(map[string]string{}))
}
