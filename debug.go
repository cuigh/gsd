package gsd

import (
	"errors"
)

// debugger is the interface that wraps the Debug method.
type debugger interface {
	// Debug method is used to print sql and args for debugging
	Debug() (sql string, args []interface{}, err error)
}

// Debug prints query clause and args of query context for debugging
func Debug(ctx interface{}) (sql string, args []interface{}, err error) {
	if d, ok := ctx.(debugger); ok {
		return d.Debug()
	}

	err = errors.New("对象不支持Debug方法")
	return
}

func (this *deleteContext) Debug() (sql string, args []interface{}, err error) {
	ctx := newBuildContext()
	err = this.b.BuildDelete(ctx, this.info)
	if err == nil {
		sql, args = ctx.GetSql(), ctx.GetParams()
	}
	return
}

func (this *insertContext) Debug() (sql string, args []interface{}, err error) {
	ctx := newBuildContext()
	err = this.b.BuildInsert(ctx, this.info)
	if err == nil {
		sql, args = ctx.GetSql(), ctx.GetParams()
	}
	return
}

func (this *updateContext) Debug() (sql string, args []interface{}, err error) {
	ctx := newBuildContext()
	err = this.b.BuildUpdate(ctx, this.info)
	if err == nil {
		sql, args = ctx.GetSql(), ctx.GetParams()
	}
	return
}

func (this *selectContext) Debug() (sql string, args []interface{}, err error) {
	ctx := newBuildContext()
	err = this.b.BuildSelect(ctx, this.info)
	if err == nil {
		sql, args = ctx.GetSql(), ctx.GetParams()
	}
	return
}
