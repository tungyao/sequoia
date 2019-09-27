package tndb

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

const (
	INSERT = iota
	UPDATE
	DELETE
	Select
)

type FUNC interface {
	Use(dbname string, pwd string, name string) *DB

	Select(string) *DB
	All(column ...string) []map[string]interface{}
	FindOne(column ...string) map[string]interface{}
	Count() *DB
	IsExits() *DB

	Insert(string) *DB
	Update(string) *DB
	Key(interface{}) *DB
	Where(key map[string]string) *DB
	Done() int64

	Command(string) *DB
	Execute() (int64, int64)
	Query() []map[string]interface{}
}
type DB struct {
	op        int
	db        string
	table     string
	sql       string
	kel       *sql.DB
	formatSql map[string]string
}

func (d *DB) Select(tablename string) *DB {
	d.op = Select
	d.formatSql = make(map[string]string)
	d.formatSql["select"] = "select "
	d.formatSql["from"] = " from " + tablename
	return d
}
func (d *DB) All(column ...string) []map[string]interface{} {
	d.formatSql["column"] = setColumn(column)
	d.sql = d.formatSql["select"] + d.formatSql["column"] + d.formatSql["from"] + d.formatSql["where"]
	fmt.Println(d.sql)
	rows, err := d.kel.Query(d.sql)
	toError(err)
	columns, _ := rows.Columns()
	length := len(columns)
	data := make([]map[string]interface{}, 0)
	n := 0
	for rows.Next() {
		value := make([]interface{}, length)
		columnPointers := make([]interface{}, length)
		for i := 0; i < length; i++ {
			columnPointers[i] = &value[i]
		}
		rows.Scan(columnPointers...)
		//data[n] = make(map[string]interface{})
		for i := 0; i < length; i++ {
			columnName := columns[i]
			columnValue := columnPointers[i].(*interface{})
			//data[n][columnName] = *columnValue
			data = append(data, map[string]interface{}{columnName: *columnValue})
		}
		n++

	}
	return data
}
func setColumn(column ...[]string) string {
	var tos string = ""
	if len(column) != 0 {
		for _, v := range column[0] {
			tos += string(v) + ","
		}
		return tos[:len(tos)-1]
	} else {
		tos = " * "
		return tos
	}

}
func (d *DB) FindOne(column ...string) map[string]interface{} {

	d.formatSql["column"] = setColumn(column)
	d.formatSql["limit"] = "limit 1"

	d.sql = d.formatSql["select"] + d.formatSql["column"] + d.formatSql["from"] + d.formatSql["where"] + d.formatSql["limit"]

	rows, err := d.kel.Query(d.sql)
	toError(err)
	columns, _ := rows.Columns()
	length := len(columns)
	data := make(map[string]interface{})
	for rows.Next() {
		value := make([]interface{}, length)
		columnPointers := make([]interface{}, length)
		for i := 0; i < length; i++ {
			columnPointers[i] = &value[i]
		}
		rows.Scan(columnPointers...)
		for i := 0; i < length; i++ {
			columnName := columns[i]
			columnValue := columnPointers[i].(*interface{})
			data[columnName] = *columnValue
		}
	}
	return B2S(data).(map[string]interface{})
}

func (d *DB) Count() *DB {
	return d
}
func (d *DB) IsExits() *DB {
	return d
}

//TODO 使用数据库 Use
func (d *DB) Use(dbname string, pwd string, name string) *DB {
	d.db = dbname
	db, err := sql.Open("mysql", "root:"+pwd+"@tcp(localhost)/"+dbname+"?charset=utf8")
	fmt.Println(err)
	d.kel = db
	d.formatSql = make(map[string]string)
	return d
}

//TODO 插入数据 HEAD
func (d *DB) Insert(table string) *DB {
	d.op = INSERT
	d.table = table
	d.formatSql["type"] = "insert into " + table
	return d
}

//TODO 升级数据
func (d *DB) Update(table string) *DB {
	d.op = UPDATE
	d.table = table
	d.formatSql["type"] = "update " + table
	return d
}

//TODO 插入数据 / 升级 KEY
func (d *DB) Key(k interface{}) *DB {
	sl := keyForInsertOrUpdate(k)
	//switch d.op {
	//case INSERT:
	//	d.sql += sl
	//case UPDATE:
	//	d.sql += sl
	//}
	d.formatSql["key"] = sl
	return d
}

//TODO
func (d *DB) Where(key map[string]string) *DB {
	d.formatSql["where"] = " where " + ConvertMapString(key)
	return d
}

//TODO 数据
func (d *DB) Done() int64 {
	switch d.op {
	case INSERT:
		d.sql = d.formatSql["type"] + d.formatSql["key"]
		break
	case UPDATE:
		d.sql = d.formatSql["type"] + d.formatSql["key"] + d.formatSql["where"]

	}
	stmt, _ := d.kel.Prepare(d.sql)
	res, err := stmt.Exec()
	toError(err)
	id, err := res.LastInsertId()
	toError(err)
	fmt.Println(d.sql)
	return id
}
func (d *DB) Command(sql string) *DB {
	d.sql = sql
	return d
}
func (d *DB) Execute() (int64, int64) {
	stmt, err := d.kel.Prepare(d.sql)
	toError(err)
	res, err := stmt.Exec()
	toError(err)
	id, err := res.LastInsertId()
	toError(err)
	rw, err := res.RowsAffected()
	return rw, id
}
func (d *DB) Query() []map[string]interface{} {
	rows, err := d.kel.Query(d.sql)
	toError(err)
	columns, err := rows.Columns()
	toError(err)
	length := len(columns)
	data := make([]map[string]interface{}, 0)
	n := 0
	for rows.Next() {
		value := make([]interface{}, length)
		columnPointers := make([]interface{}, length)
		for i := 0; i < length; i++ {
			columnPointers[i] = &value[i]
		}
		err = rows.Scan(columnPointers...)
		toError(err)
		//data[n] = make(map[string]interface{})
		for i := 0; i < length; i++ {
			columnName := columns[i]
			columnValue := columnPointers[i].(*interface{})
			//data[n][columnName] = *columnValue
			data = append(data, map[string]interface{}{columnName: *columnValue})
		}
		n++

	}
	return data
}
func toError(err error) {
	if err != nil {
		log.Panic(err)
	}
	return
}
