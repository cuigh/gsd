package gsd

/********** insertInfo **********/

type insertInfo struct {
	table  string
	values map[string]interface{}
	// todo: add subquery supporting, like: insert into X(a,b,c) select a1, b1, c1 from Y
}

/********** insertContext **********/

type insertContext struct {
	exe  executor
	b    builder
	info *insertInfo
}

func newInsertContext(exe executor, b builder, info *insertInfo) *insertContext {
	return &insertContext{
		exe:  exe,
		b:    b,
		info: info,
	}
}

func (this *insertContext) Values(values InsertValues) InsertResultClause {
	this.info.values = values
	return this
}

func (this *insertContext) Result() (InsertResult, error) {
	ctx := newBuildContext()
	err := this.b.BuildInsert(ctx, this.info)
	if err != nil {
		return nil, err
	}
	return this.exe.Exec(ctx.GetSql(), ctx.GetParams()...)
}

/********** InsertValues **********/

type InsertValues map[string]interface{}

func (this InsertValues) Get(name string) (value interface{}, ok bool) {
	value, ok = this[name]
	return
}

func (this InsertValues) Set(name string, value interface{}) InsertValues {
	this[name] = value
	return this
}

func (this InsertValues) Remove(name string) InsertValues {
	delete(this, name)
	return this
}
