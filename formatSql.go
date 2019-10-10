package sequoia

import (
	"github.com/tungyao/tjson"
	"strings"
)

type operation interface {
}

type _SQL struct {
}

func ConvertMapStringAnd(m map[string]string) string {
	sl := ""
	for k, v := range m {
		sl += k + "='" + v + "'" + " and "
	}
	return sl[:len(sl)-4]
}
func ConvertMapString(m map[string]string) string {
	sl := ""
	for k, v := range m {
		sl += k + "='" + v + "'" + ","
	}
	return sl[:len(sl)-1]
}
func ConvertArrayString(arr []string) string {
	sl := ""
	for _, v := range arr {
		if v == "NULL" {
			sl += "" + v + ","
		} else {
			sl += "'" + v + "',"
		}

	}
	return "(" + sl[:len(sl)-1] + ")"
}
func ConvertStringToArray(s string) []map[string]interface{} {
	s = strings.Replace(s, "[", "", -1)
	s = strings.Replace(s, "]", "", -1)
	data := make([]map[string]interface{}, 0)
	data = append(data, tjson.Decode(s))
	return data
}
func keyForInsertOrUpdate(k interface{}, op int) string {
	sl := ""
	if op == INSERT {
		switch k.(type) {
		case map[string]string:
			sl = " set " + ConvertMapString(k.(map[string]string))
		case []string:
			sl = " values " + ConvertArrayString(k.([]string))
		default:
			sl = ""
		}
	} else if op == UPDATE {
		sl = " set " + ConvertMapString(k.(map[string]string))
	}
	return sl
}
func B2S(obj interface{}) interface{} { //TODO 将查询数据正确显示 uint8 -> string

	switch obj.(type) {
	case map[string]interface{}:
		data := make(map[string]interface{})
		for k, v := range obj.(map[string]interface{}) {
			data[k] = byteToString(v.([]uint8))
		}
		return data
	case []map[string]interface{}:
		length := len(obj.([]map[string]interface{}))
		data := make([]map[string]interface{}, length)
		for k, v := range obj.([]map[string]interface{}) {
			for j, l := range v {
				data[k][j] = byteToString(l.([]uint8))
			}
		}
		return data

	}
	return nil
}
func byteToString(bs []uint8) string {
	ba := []byte{}
	for _, b := range bs {
		ba = append(ba, byte(b))
	}
	return string(ba)
}
