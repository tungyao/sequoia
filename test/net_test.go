package test

import (
	"../caches"
	"testing"
)

func TestRedis(t *testing.T) {
	//caches.HSet(caches.Cache{
	//	Key:"test",
	//	Value:"hello",
	//})
	//caches.HSet(caches.Cache{
	//	Key:   "test",
	//	Value: "hellos",
	//	Time:  0,
	//})
	Value := cache.HGet("test").Value
	t.Log(Value)
	//log.Println(fmt.Sprint(map[string]string{}))
}
