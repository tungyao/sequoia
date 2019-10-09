package test

import (
	"../../sequoia"
	"log"
	"testing"
)

var db sequoia.FUNC = sequoia.NewDB(sequoia.Config{
	MaxOpen: 2000,
	MaxIde:  1000,
	Cache:   true,
}).Use("test", "123456", "root")

func TestDb(t *testing.T) {
	//d:=`[{"id":"Mw=="}]`
	//sequoia.ConvertStringToArray(d)
	//db.Insert("test").Key(map[string]string{"name":"asdzxc"}).Done()
	//db.Update("test").Key(map[string]string{"name": "asdasdas"}).Where(map[string]string{"name": "asdzxc"}).Done()
	data := db.Select("test").Where(map[string]string{"id": "5"}).All("name")
	log.Println(data)
	//data2:=db.Select("test").All("id","name")
	//fmt.Println(data2)
	//sql := "insert into test set name='hello'"
	//la, rw := db.Command(sql).Execute()
	//fmt.Println(la, rw)
	//db.Delete("test").Where(map[string]string{"id": "2"}).Done()
}
