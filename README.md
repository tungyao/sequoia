## Go for MySQL simple controller
### Example
---
Initialize connection,you can use two ways:
```go
//`test` is database
var db sequoia.FUNC = sequoia.NewDB(sequoia.Config{
	MaxOpen: 2000,
	MaxIde:  1000,
	Cache:   true,
}).Use("test", "123456", "root")
```
Select one column
```go
//`test` is table
//Where(map[string]string)
//FindOne(field name)

db.Select("test").Where(map[string]string{"id":"2"}).FindOne("id", "name")
==> map[string]interface{}
```
Select one column
```go
//`test` is table
//FindOne(field name)

db.Select("test").FindOne("id", "name")
==> map[string]interface{}
```
Select All
```go
db.Select("test").All("id","name")
==>[map[string]interface{}]
```
Insert One
```go
//Return the last ID

db.Insert("user").Key([]string{`NULL`, "123123123", `NULL`}).Done()
==> int64
```
Update One
```go
//Return the last ID

db.Update("t_user").Key(map[string]string{"email":"asdasdas"}).Where(map[string]string{"id":"10"}).Done()
==> int64
```
#Last , you can use Native SQL
```go
db.Command(sql).Query()
==> [map[string]interface{}]

db.Command(sql).Execute()
==> (int64,int64)

```
## Join Us