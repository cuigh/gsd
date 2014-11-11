package gsd

type Table interface {
	Name() string
	Alias() string
	Prefix() string
	C(cols ...string) *Columns
	G(cols ...string) *Groupers
	S(st sortType, cols ...string) *Sorters
}

func T(name string) Table {
	return &basicTable{
		name: name,
	}
}

func TA(name, alias string) Table {
	return &basicTable{
		name:  name,
		alias: alias,
	}
}

type basicTable struct {
	name  string
	alias string
}

func (this *basicTable) Name() string {
	return this.name
}

func (this *basicTable) Alias() string {
	return this.alias
}

func (this *basicTable) Prefix() string {
	if this.alias == "" {
		return this.name
	} else {
		return this.alias
	}
}

func (this *basicTable) C(cols ...string) *Columns {
	return new(Columns).Add(this, cols...)
}

func (this *basicTable) G(cols ...string) *Groupers {
	return new(Groupers).AddT(this, cols...)
}

func (this *basicTable) S(st sortType, cols ...string) *Sorters {
	return new(Sorters).AddT(st, this, cols...)
}
