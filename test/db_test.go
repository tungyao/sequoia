package test

import (
	"../../sequoia"
	"testing"
)

var db sequoia.FUNC = sequoia.NewDB(2000, 1000).Use("test", "123456", "root")

func TestDb(t *testing.T) {
	//db.Insert("test").Key(map[string]string{"name":"asdzxc"}).Done()
	//db.Update("test").Key(map[string]string{"name": "asdasdas"}).Where(map[string]string{"name": "asdzxc"}).Done()
	//data := db.Select("test").Where(map[string]string{"id":"8"}).FindOne("id", "name")
	//fmt.Println(len(data))
	//data2:=db.Select("test").All("id","name")
	//fmt.Println(data2)
	//sql := "insert into test set name='hello'"
	//la, rw := db.Command(sql).Execute()
	//fmt.Println(la, rw)
	db.Delete("test").Where(map[string]string{"id": "2"}).Done()
}
