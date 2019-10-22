package sequoia

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tungyao/sequoia/caches"
	"github.com/tungyao/tjson"
	"github.com/tungyao/yell"
	"strconv"
)

var log = yell.New(yell.Config{
	Path:     "",
	FileName: "sequoia",
}, "[SEQUOIA]")

const (
	INSERT = iota
	UPDATE
	DELETE
	Select
)

type FUNC interface {
	Use(dbname string, pwd string, name string) *DB
	//TODO------------------------------------------ MAIN FUNC
	Select(string) *DB
	Insert(string) *DB
	Delete(string) *DB
	Update(string) *DB
	//TODO------------------------------------------ OPERATION
	Key(interface{}) *DB
	Where(key map[string]string) *DB
	FindOne(column ...string) map[string]interface{}
	All(column ...string) []map[string]interface{}
	Done() int64
	//TODO-------------------------------------- AUXILIARY FUNC
	Count() *DB
	Limit(str ...int) *DB
	IsExits() *DB
	IsCache(bool) *DB
	Sort(string, string) *DB
	//TODO---------------------------------------- NATIVE FUNC
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
	MaxOpen   int
	MaxIde    int
	Cache     *caches.Conn
	iscache   bool
}
type Config struct {
	MaxOpen, MaxIde int
	Cache           bool
}

func (d *DB) Use(dbname string, pwd string, name string) *DB {
	d.db = dbname
	db, err := sql.Open("mysql", name+":"+pwd+"@tcp(localhost)/"+dbname+"?charset=utf8")
	db.SetMaxOpenConns(d.MaxOpen)
	db.SetMaxIdleConns(d.MaxIde)
	db.Ping()
	if err != nil {
		panic(err)
	}
	d.kel = db
	d.formatSql = make(map[string]string)
	return d
}

//TODO------------------------------------------ MAIN FUNC S
func (d *DB) Insert(table string) *DB {
	format := make(map[string]string)
	format["type"] = "insert into `" + table + "`"
	kel := d.kel
	return &DB{
		op:        INSERT,
		db:        d.db,
		table:     table,
		sql:       "",
		kel:       kel,
		formatSql: format,
		MaxOpen:   d.MaxOpen,
		MaxIde:    d.MaxIde,
		Cache:     d.Cache,
		iscache:   d.iscache,
	}
	//return d
}
func (d *DB) Update(table string) *DB {
	//d.op = UPDATE
	//d.table = table
	//d.formatSql["type"] = "update `" + table + "`"
	//return d
	format := make(map[string]string)
	format["type"] = "update `" + table + "`"
	kel := d.kel
	return &DB{
		op:        UPDATE,
		db:        d.db,
		table:     table,
		sql:       "",
		kel:       kel,
		formatSql: format,
		MaxOpen:   d.MaxOpen,
		MaxIde:    d.MaxIde,
		Cache:     d.Cache,
		iscache:   d.iscache,
	}
}
func (d *DB) Delete(table string) *DB {
	//d.op = DELETE
	//d.formatSql["type"] = "delete from `" + s + "`"
	//return d
	format := make(map[string]string)
	format["type"] = "delete from `" + table + "`"
	kel := d.kel
	return &DB{
		op:        Select,
		db:        d.db,
		table:     table,
		sql:       "",
		kel:       kel,
		formatSql: format,
		MaxOpen:   d.MaxOpen,
		MaxIde:    d.MaxIde,
		Cache:     d.Cache,
		iscache:   d.iscache,
	}
}
func (d *DB) Select(table string) *DB {
	format := make(map[string]string)
	format["select"] = "select "
	format["from"] = " from `" + table + "`"
	kel := d.kel
	return &DB{
		op:        Select,
		db:        d.db,
		table:     table,
		sql:       "",
		kel:       kel,
		formatSql: format,
		MaxOpen:   d.MaxOpen,
		MaxIde:    d.MaxIde,
		Cache:     d.Cache,
		iscache:   d.iscache,
	}
	//d.op = Select
	//d.formatSql = make(map[string]string)
	//d.formatSql["select"] = "select "
	//d.formatSql["from"] = " from `" + tablename + "`"
	//return d
}

//TODO------------------------------------------ OPERATION S
func (d *DB) Key(k interface{}) *DB {
	sl := keyForInsertOrUpdate(k, d.op)
	//switch d.op {
	//case INSERT:
	//	d.sql += sl
	//case UPDATE:
	//	d.sql += sl
	//}
	d.formatSql["key"] = sl
	return d
}
func (d *DB) Where(key map[string]string) *DB {
	d.formatSql["where"] = " where " + ConvertMapStringAnd(key)
	return d
}
func (d *DB) All(column ...string) []map[string]interface{} {
	d.formatSql["column"] = setColumn(column)
	d.sql = d.formatSql["select"] + d.formatSql["column"] + d.formatSql["from"] + d.formatSql["where"] + d.formatSql["sort"] + d.formatSql["limit"]
	if d.Cache != nil && d.iscache {
		hash := d.Cache.HGet(d.sql)
		if hash != nil {
			log.Println("get caches", hash)
			return ConvertStringToArray(hash.Value.(string))
		}
	}
	log.Println("***SQL***", d.sql)
	sq := d.sql
	rows, err := d.kel.Query(sq)
	defer func() {
		d.formatSql = make(map[string]string)
		d.sql = ""
		d.iscache = false
		_ = rows.Close()
	}()
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
		columno := make(map[string]interface{}, length)
		for i := 0; i < length; i++ {
			columnName := columns[i]
			columnValue := columnPointers[i].(*interface{})
			//data[n][columnName] = *columnValue
			columno[columnName] = string((*columnValue).([]uint8))
		}
		data = append(data, columno)
		n++

	}
	if d.Cache != nil && d.iscache && len(data) != 0 {
		d.Cache.HSet(caches.Cache{
			Key:   d.sql,
			Value: tjson.Encode(data),
			Time:  0,
		})
	}
	return data
}
func (d *DB) Done() int64 {
	defer func() {
		d.formatSql = make(map[string]string)
		d.iscache = false
	}()
	switch d.op {
	case INSERT:
		d.sql = d.formatSql["type"] + d.formatSql["key"]
	case UPDATE:
		d.sql = d.formatSql["type"] + d.formatSql["key"] + d.formatSql["where"]
	case DELETE:
		d.sql = d.formatSql["type"] + d.formatSql["where"]
	}
	log.Println("***SQL***", d.sql)
	sq := d.sql
	stmt, err := d.kel.Prepare(sq)
	toError(err)
	res, err := stmt.Exec()
	toError(err)
	id, err := res.LastInsertId()
	toError(err)
	return id
}
func (d *DB) FindOne(column ...string) map[string]interface{} {
	d.formatSql["column"] = setColumn(column)
	d.formatSql["limit"] = " limit 1"
	d.sql = d.formatSql["select"] + d.formatSql["column"] + d.formatSql["from"] + d.formatSql["where"] + d.formatSql["sort"] + d.formatSql["limit"]
	if d.Cache != nil && d.iscache {
		hash := d.Cache.HGet(d.sql)
		if hash != nil {
			log.Println("get caches", d.sql)
			return tjson.Decode([]byte(hash.Value.(string)))
		}
	}
	sq := d.sql
	log.Println("***SQL***", sq)
	rows, err := d.kel.Query(sq)
	defer func() {
		d.formatSql = make(map[string]string)
		d.sql = ""
		d.iscache = false
		_ = rows.Close()
	}()
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
	da := B2S(data).(map[string]interface{})
	if d.Cache != nil && d.iscache && len(da) != 0 {
		d.Cache.HSet(caches.Cache{
			Key:   d.sql,
			Value: tjson.Encode(da),
			Time:  0,
		})
	}
	return da
}

//TODO-------------------------------------- AUXILIARY FUNC S
func (d *DB) Limit(str ...int) *DB {
	d.formatSql["limit"] = " limit " + strconv.FormatInt(int64(str[0]), 10) + "," + strconv.FormatInt(int64(str[1]), 10)
	return d
}
func (d *DB) IsCache(b bool) *DB {
	d.iscache = b
	return d
}
func (d *DB) Count() *DB {
	return d
}
func (d *DB) IsExits() *DB {
	return d
}
func (d *DB) Sort(field, string2 string) *DB {
	d.formatSql["sort"] = " order by " + field + " " + string2
	return d
}

//TODO---------------------------------------- NATIVE FUNC S
func (d *DB) Command(sql string) *DB {
	d.sql = sql
	return d
}
func (d *DB) Execute() (int64, int64) {
	log.Println("***SQL***", d.sql)
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
	if d.Cache != nil && d.iscache {
		hash := d.Cache.HGet(d.sql)
		if hash != nil {
			log.Println("get caches", hash)
			return ConvertStringToArray(hash.Value.(string))
		}
	}
	log.Println("***SQL***", d.sql)
	rows, err := d.kel.Query(d.sql)
	defer func() {
		d.formatSql = make(map[string]string)
		d.sql = ""
		d.iscache = false
		_ = rows.Close()
	}()
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
		columno := make(map[string]interface{}, length)
		for i := 0; i < length; i++ {
			columnName := columns[i]
			columnValue := columnPointers[i].(*interface{})
			columno[columnName] = string((*columnValue).([]uint8))
		}
		data = append(data, columno)
		n++
	}
	if d.Cache != nil && d.iscache && len(data) != 0 {
		d.Cache.HSet(caches.Cache{
			Key:   d.sql,
			Value: tjson.Encode(data),
			Time:  0,
		})
	}
	return data
}

//TODO-------------------------------------- OTHER FUNC
func toError(err error) {
	if err != nil {
		log.Panic(err)
	}
	return
}
func NewDB(c Config) *DB {
	if c.Cache {
		cc := caches.New()
		return &DB{
			MaxOpen: c.MaxOpen,
			MaxIde:  c.MaxIde,
			Cache:   cc,
		}
	}
	return &DB{
		MaxOpen: c.MaxOpen,
		MaxIde:  c.MaxIde,
	}
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
