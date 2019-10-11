package test

import (
	"../caches"
	"testing"
)

func TestRedis(t *testing.T) {
	caches.TestRedis()
	//caches.HSet(caches.Cache{
	//	Key:"test",
	//	Value:"hello",
	//})
	//r := caches.New()
	//r.HSet(caches.Cache{
	//	Key:   "test",
	//	Value: "testtest",
	//	Time:  0,
	//})
	//for i := 0; i < 100; i++ {
	//	t.Log(r.HGet("test"))
	//
	//}
	//time.Sleep(time.Second)
	//log.Println(fmt.Sprint(map[string]string{}))
}
