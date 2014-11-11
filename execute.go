package gsd

/********** executeContext **********/

type executeContext struct {
	exe   executor
	query string
	args  []interface{}
}

func newExecuteContext(exe executor, query string, args []interface{}) *executeContext {
	return &executeContext{
		exe:   exe,
		query: query,
		args:  args,
	}
}

func (this *executeContext) Result() (ExecuteResult, error) {
	return this.exe.Exec(this.query, this.args...)
}

func (this *executeContext) Row() Row {
	return &row{
		exe:  this.exe,
		sql:  this.query,
		args: this.args,
	}
}

func (this *executeContext) Rows() Rows {
	return &rows{
		exe:  this.exe,
		sql:  this.query,
		args: this.args,
	}
}
