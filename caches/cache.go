package caches

import (
	"crypto/sha1"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

//var (
//	Conn net.Conn
//	err error
//)
type Cache struct {
	Key   string
	Value interface{}
	Time  int64
}
type Conn struct {
	Con net.Conn
	sw  sync.RWMutex
}

func New() *Conn {
	conn, err := net.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		panic(err)
		return nil
	}
	return &Conn{
		Con: conn,
	}
}
func (con Conn) Expire(key string, i int64) {
	go func() {
		s := "EXPIRE" + " " + key + " " + strconv.FormatInt(i, 10)
		_, err := con.Con.Write(format(s))
		if err != nil {
			log.Panic(err)
		}
	}()
}
func (con Conn) Pexpire(key string, i int64) {
	go func() {
		s := "PEXPIRE" + " " + key + " " + strconv.FormatInt(i, 10)
		_, err := con.Con.Write(format(s))
		if err != nil {
			log.Panic(err)
		}
	}()
}
func (con Conn) HPexpire(key string, i int64) {
	go func() {
		s := "PEXPIRE" + " " + md(key) + " " + strconv.FormatInt(i, 10)
		_, err := con.Con.Write(format(s))
		if err != nil {
			log.Panic(err)
		}
	}()
}
func (con Conn) HExpire(key string, i int64) {
	go func() {
		s := "EXPIRE" + " " + md(key) + " " + strconv.FormatInt(i, 10)
		_, err := con.Con.Write(format(s))
		if err != nil {
			log.Panic(err)
		}
	}()
}
func (con Conn) Set(cache Cache) {
	go func() {
		s := "set" + " " + cache.Key + " \"" + fmt.Sprint(cache.Value) + "\""
		_, err := con.Con.Write(format(s))
		if err != nil {
			log.Panic(err)
		}
	}()
}
func (con Conn) HSet(cache Cache) {
	go func() {
		s := "set" + " " + md(cache.Key) + " " + fmt.Sprint(cache.Value)
		log.Print(s)
		n, err := con.Con.Write(format(s))
		log.Println(n)
		if err != nil {
			log.Panic(err)
		}
	}()
}
func (con Conn) Get(cacheName string) *Cache {
	data := make(chan *Cache)
	go func() {
		con.sw.Lock()
		defer con.sw.Unlock()
		s := "get" + " " + cacheName
		_, err := con.Con.Write(format(s))
		if err != nil {
			log.Panic(err)
		}
		var msg = make([]byte, 4096)
		n, _ := con.Con.Read(msg)
		if n == len(cacheName) || string(msg[:n]) == "$-1\r\n" {
			data <- nil
		}
		data <- &Cache{
			Key:   cacheName,
			Value: strings.Split(string(msg[:n]), "\r\n")[1],
			Time:  0,
		}
	}()
	return <-data
}
func (con Conn) HGet(cacheName string) *Cache {
	data := make(chan *Cache)
	go func() {
		con.sw.Lock()
		defer con.sw.Unlock()
		s := "get" + " " + md(cacheName)
		_, err := con.Con.Write(format(s))
		if err != nil {
			log.Panic(err)
		}
		var msg = make([]byte, 4096)
		n, _ := con.Con.Read(msg)
		if n == len(cacheName) || string(msg[:n]) == "$-1\r\n" {
			data <- nil
		}
		data <- &Cache{
			Key:   cacheName,
			Value: strings.Split(string(msg[:n]), "\r\n")[1],
			Time:  0,
		}
	}()
	return <-data
}
func md(s string) string {
	Sha1Inst := sha1.New()
	Sha1Inst.Write([]byte(s))
	Result := Sha1Inst.Sum([]byte(""))
	return fmt.Sprintf("%x", Result)
}
func format(s string) []byte {
	var pro string
	ret := strings.Split(s, " ")
	for k, v := range ret {
		if k == 0 {
			pro = fmt.Sprintf("*%d\r\n", len(ret))
		}
		pro += fmt.Sprintf("$%d\r\n%s\r\n", len(v), v)
	}
	return []byte(pro)
}
func TestRedis() {
	conn, err := net.Dial("tcp", "127.0.0.1:6379")
	defer conn.Close()
	if err != nil {
		os.Exit(0)
	}
	go func() {
		_, err = conn.Write([]byte("PING\r\n"))
		if err != nil {
			log.Panic(err)
		}
	}()

	var msg = make([]byte, 4096)

	_, err = conn.Read(msg)

	log.Println(string(msg))
	_ = conn.Close()
}
