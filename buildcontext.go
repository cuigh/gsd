package gsd

import (
	"bytes"
	"fmt"
)

type buildContext struct {
	sql    *bytes.Buffer
	params []interface{}
}

func newBuildContext() *buildContext {
	return &buildContext{
		sql: new(bytes.Buffer),
	}
}

func (this *buildContext) AppendSql(strs ...string) *buildContext {
	for _, s := range strs {
		this.sql.WriteString(s)
	}
	return this
}

func (this *buildContext) AppendSqlF(format string, args ...interface{}) *buildContext {
	this.sql.WriteString(fmt.Sprintf(format, args...))
	return this
}

func (this *buildContext) AddParam(params ...interface{}) *buildContext {
	this.params = append(this.params, params...)
	return this
}

func (this *buildContext) GetSql() string {
	return this.sql.String()
}

func (this *buildContext) GetParams() []interface{} {
	return this.params
}
