package test

import (
	"../../sequoia"
	"testing"
)

var db sequoia.FUNC = sequoia.NewDB(sequoia.Config{
	MaxOpen: 2000,
	MaxIde:  1000,
	Cache:   false,
}).Use("test", "123456", "root")

func TestDb(t *testing.T) {
	//d:=`[{"id":"Mw=="}]`
	//sequoia.ConvertStringToArray(d)aa
	//db.Insert("test").Key(map[string]string{"name": "asdzxc", "a": "5"}).Done()
	//db.Update("test").Key(map[string]string{"name": "asdasdas"}).Where(map[string]string{"name": "asdzxc"}).Done()
	//data := db.Select("test").All("name","id")
	//data := db.Select("test").Sort("addtime", "desc").All("name", "id")
	//t.Log(data)
	//log.Println(string(data[0]["name"].([]uint8)))
	//data2:=db.Select("test").All("id","name")
	//fmt.Println(data2)
	//sql := "insert into test set name='hello'"
	//la, rw := db.Command(sql).Execute()
	//fmt.Println(la, rw)
	//db.Delete("test").Where(map[string]string{"id": "2"}).Done()
}
