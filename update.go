package gsd

/********** updateInfo **********/

type updateInfo struct {
	table  string
	values map[string]*updateValue
	where  Filters
}

/********** updateContext **********/

type updateContext struct {
	exe  executor
	b    builder
	info *updateInfo
}

func newUpdateContext(exe executor, b builder, info *updateInfo) *updateContext {
	return &updateContext{
		exe:  exe,
		b:    b,
		info: info,
	}
}

func (this *updateContext) Set(values UpdateValues) SetClause {
	this.info.values = values
	return this
}

func (this *updateContext) Where(f Filters) ResultClause {
	this.info.where = f
	return this
}

func (this *updateContext) Result() (Result, error) {
	ctx := newBuildContext()
	err := this.b.BuildUpdate(ctx, this.info)
	if err != nil {
		return nil, err
	}
	return this.exe.Exec(ctx.GetSql(), ctx.GetParams()...)
}

/********** UpdateType **********/

type updateType int8

const (
	UPDATE_EQ updateType = iota
	UPDATE_INC
	UPDATE_XP
)

/********** UpdateValue **********/

type updateValue struct {
	ut  updateType
	val interface{}
}

func UV(value interface{}) *updateValue {
	return &updateValue{ut: UPDATE_EQ, val: value}
}

func UVT(ut updateType, value interface{}) *updateValue {
	return &updateValue{ut: ut, val: value}
}

/********** UpdateValues **********/

type UpdateValues map[string]*updateValue

func (this UpdateValues) Get(name string) (value *updateValue, ok bool) {
	value, ok = this[name]
	return
}

func (this UpdateValues) Set(name string, value *updateValue) UpdateValues {
	this[name] = value
	return this
}

func (this UpdateValues) Remove(name string) UpdateValues {
	delete(this, name)
	return this
}
