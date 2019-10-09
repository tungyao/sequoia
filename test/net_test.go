package test

import (
	"../cache"
	"testing"
)

func TestRedis(t *testing.T) {
	//cache.HSet(cache.Cache{
	//	Key:"test",
	//	Value:"hello",
	//})
	//cache.HSet(cache.Cache{
	//	Key:   "test",
	//	Value: "hellos",
	//	Time:  0,
	//})
	Value := cache.HGet("test").Value
	t.Log(Value)
	//log.Println(fmt.Sprint(map[string]string{}))
}
