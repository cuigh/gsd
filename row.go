package gsd

import (
	"database/sql"
	// "fmt"
	"reflect"
	"strings"
)

var ErrNoRows = sql.ErrNoRows

/********** Row **********/

type Row interface {
	Scan(dest ...interface{}) error
	ScanObj(obj interface{}) error
}

type row struct {
	exe  executor
	sql  string
	args []interface{}
	rows *sql.Rows
	err  error
}

func (this *row) Scan(values ...interface{}) error {
	err := this.prepareRows()
	if err != nil {
		return err
	}
	defer this.rows.Close()

	if this.rows.Next() {
		return this.rows.Scan(values...)
	} else {
		if err := this.rows.Err(); err != nil {
			return err
		} else {
			return ErrNoRows
		}
	}

}

func (this *row) ScanObj(obj interface{}) error {
	err := this.prepareRows()
	if err != nil {
		return err
	}
	defer this.rows.Close()

	columns, err := this.rows.Columns()
	if err != nil {
		return err
	}

	values := make([]interface{}, len(columns))
	for i := 0; i < len(values); i++ {
		var val interface{}
		values[i] = &val
	}

	if this.rows.Next() {
		if err := this.rows.Scan(values...); err != nil {
			return err
		}
	} else {
		if err := this.rows.Err(); err != nil {
			return err
		} else {
			return ErrNoRows
		}
	}

	ti := getTypeInfo(reflect.TypeOf(obj).Elem())
	for i, col := range columns {
		if fi := ti.GetFieldInfo(strings.ToLower(col)); fi != nil {
			v := values[i].(*interface{})
			fillField(obj, fi.name, *v)
		}
	}
	return nil
}

func (this *row) prepareRows() error {
	if this.err == nil && this.rows == nil {
		this.rows, this.err = this.exe.Query(this.sql, this.args...)
	}
	return this.err
}

/********** Rows **********/

type Rows interface {
	All(objs interface{}) (err error)
	For(f func(r Row) error) error
}

type rows struct {
	exe     executor
	sql     string
	args    []interface{}
	err     error
	rows    *sql.Rows
	columns []string
	values  []interface{}
}

func (this *rows) Scan(values ...interface{}) error {
	return this.rows.Scan(values...)
}

func (this *rows) ScanObj(obj interface{}) error {
	ti := getTypeInfo(reflect.TypeOf(obj).Elem())
	return this.scanObj(ti, obj)
}

// All reads all records and push them to objs, objs must be a pointer to struct array, like: objs := &[]*Object{}
func (this *rows) All(objs interface{}) (err error) {
	// if kind := reflect.ValueOf(objs).Kind(); kind != reflect.Ptr {
	// 	return nil
	// }
	err = this.prepareRows()
	if err != nil {
		return
	}
	defer this.rows.Close()

	err = this.prepareColumns()
	if err != nil {
		return
	}

	reflect.ValueOf(objs).Elem().SetLen(0)
	t := reflect.TypeOf(objs).Elem().Elem().Elem()
	ti := getTypeInfo(t)
	for this.rows.Next() {
		obj := reflect.New(t).Interface()
		err = this.scanObj(ti, obj)
		if err != nil {
			return err
		}

		v := reflect.Append(reflect.ValueOf(objs).Elem(), reflect.ValueOf(obj))
		reflect.ValueOf(objs).Elem().Set(v)
	}

	return this.rows.Err()
}

func (this *rows) For(f func(r Row) error) error {
	err := this.prepareRows()
	if err != nil {
		return err
	}
	defer this.rows.Close()

	err = this.prepareColumns()
	if err != nil {
		return err
	}

	for this.rows.Next() {
		if err = f(this); err != nil {
			return err
		}
	}

	return this.rows.Err()
}

func (this *rows) scanObj(ti *typeInfo, obj interface{}) error {
	// for i := 0; i < len(this.values); i++ {
	// 	var val interface{}
	// 	this.values[i] = &val
	// }
	if err := this.rows.Scan(this.values...); err != nil {
		return err
	}

	for i, col := range this.columns {
		if fi := ti.GetFieldInfo(strings.ToLower(col)); fi != nil {
			v := this.values[i].(*interface{})
			fillField(obj, fi.name, *v)
		}
	}
	return nil
}

func (this *rows) prepareRows() error {
	if this.err == nil && this.rows == nil {
		this.rows, this.err = this.exe.Query(this.sql, this.args...)
	}
	return this.err
}

func (this *rows) prepareColumns() (err error) {
	if this.columns == nil {
		this.columns, err = this.rows.Columns()
		if err == nil {
			this.values = make([]interface{}, len(this.columns))
			for i := 0; i < len(this.values); i++ {
				var val interface{}
				this.values[i] = &val
			}
		}
	}
	return
}
