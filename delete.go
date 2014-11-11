package gsd

/********** deleteInfo **********/

type deleteInfo struct {
	table string
	where Filters
}

/********** deleteContext **********/

type deleteContext struct {
	exe  executor
	b    builder
	info *deleteInfo
}

func newDeleteContext(exe executor, b builder, info *deleteInfo) *deleteContext {
	return &deleteContext{
		exe:  exe,
		b:    b,
		info: info,
	}
}

func (this *deleteContext) Where(f Filters) ResultClause {
	this.info.where = f
	return this
}

func (this *deleteContext) Result() (Result, error) {
	ctx := newBuildContext()
	err := this.b.BuildDelete(ctx, this.info)
	if err != nil {
		return nil, err
	}
	return this.exe.Exec(ctx.GetSql(), ctx.GetParams()...)
}
