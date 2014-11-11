package gsd

import (
	// "fmt"
	"reflect"
	"strings"
	"sync"
)

var (
	_TypeLocker sync.Mutex
	_TypeCaches map[string]*typeInfo = make(map[string]*typeInfo)
)

type typeInfo struct {
	fields map[string]*fieldInfo
}

func (this *typeInfo) GetFieldInfo(field string) *fieldInfo {
	if fi, ok := this.fields[field]; ok {
		return fi
	}
	return nil
}

type fieldInfo struct {
	name string
	t    reflect.Type
}

func getTypeInfo(t reflect.Type) *typeInfo {
	key := t.PkgPath() + "," + t.Name()

	_TypeLocker.Lock()
	ti, ok := _TypeCaches[key]
	if !ok {
		// t := reflect.TypeOf(obj)
		// if t.Kind() != reflect.Ptr {
		// 	return fmt.Errorf("only pointer to a struct is valid")
		// }

		ti = &typeInfo{
			fields: make(map[string]*fieldInfo),
		}
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			col := f.Tag.Get("gsd")
			if col == "" {
				col = f.Name
			}

			ti.fields[strings.ToLower(col)] = &fieldInfo{
				name: f.Name,
				t:    f.Type,
			}
		}
		_TypeCaches[key] = ti

	}
	_TypeLocker.Unlock()
	return ti
}

func fillField(obj interface{}, field string, val interface{}) {
	//================================================
	// values that drivers must be able to handle:
	//
	// nil
	// int64
	// float64
	// bool
	// []byte
	// string   [*] everywhere except from Rows.Next.
	// time.Time
	//================================================
	if val == nil {
		return
	}

	f := reflect.ValueOf(obj).Elem().FieldByName(field)
	switch f.Kind() {
	case reflect.Int64:
		f.SetInt(val.(int64))
	case reflect.Float64:
		f.SetFloat(val.(float64))
	case reflect.Bool:
		f.SetBool(val.(bool))
	// case reflect.Slice:
	// todo: handle []byte
	// case reflect.Struct:
	// todo: handle time.Time
	default:
		f.Set(reflect.ValueOf(val).Convert(f.Type()))
	}
}
