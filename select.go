package gsd

import (
// "database/sql"
)

/********** selectInfo **********/

type selectInfo struct {
	table    Table
	columns  []column
	distinct bool
	joins    []*joiner
	where    Filters
	groups   []*grouper
	having   Filters
	orders   []*sorter
	skip     int32
	take     int32
}

/********** selectContext **********/

type selectContext struct {
	exe  executor
	b    builder
	info *selectInfo
}

func newSelectContext(exe executor, b builder, info *selectInfo) *selectContext {
	return &selectContext{
		exe:  exe,
		b:    b,
		info: info,
	}
}

func (this *selectContext) From(t Table) FromClause {
	this.info.table = t
	return this
}

func (this *selectContext) Join(t Table, f Filters) FromClause {
	j := &joiner{
		jt: JOIN_INNER,
		t:  t,
		on: f,
	}
	this.info.joins = append(this.info.joins, j)
	return this
}

func (this *selectContext) LeftJoin(t Table, f Filters) FromClause {
	j := &joiner{
		jt: JOIN_LEFT,
		t:  t,
		on: f,
	}
	this.info.joins = append(this.info.joins, j)
	return this
}

func (this *selectContext) RightJoin(t Table, f Filters) FromClause {
	j := &joiner{
		jt: JOIN_RIGHT,
		t:  t,
		on: f,
	}
	this.info.joins = append(this.info.joins, j)
	return this
}

func (this *selectContext) FullJoin(t Table, f Filters) FromClause {
	j := &joiner{
		jt: JOIN_FULL,
		t:  t,
		on: f,
	}
	this.info.joins = append(this.info.joins, j)
	return this
}

func (this *selectContext) Where(f Filters) WhereClause {
	this.info.where = f
	return this
}

func (this *selectContext) Limit(skip, take int32) RowClause {
	this.info.skip = skip
	this.info.take = take
	return this
}

func (this *selectContext) GroupBy(g *Groupers) GroupByClause {
	this.info.groups = g.groupers
	return this
}

func (this *selectContext) Having(f Filters) HavingClause {
	this.info.having = f
	return this
}

func (this *selectContext) OrderBy(s *Sorters) OrderByClause {
	this.info.orders = s.sorters
	return this
}

// func (this *selectContext) Result() (*SelectResult, error) {
// 	ctx := newBuildContext()
// 	err := this.b.BuildSelect(ctx, this.info)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &SelectResult{exe: this.exe, ctx: ctx}, nil
// }

func (this *selectContext) Row() Row {
	ctx := newBuildContext()
	if err := this.b.BuildSelect(ctx, this.info); err != nil {
		return &row{
			exe: this.exe,
			err: err,
		}
	} else {
		return &row{
			exe:  this.exe,
			sql:  ctx.GetSql(),
			args: ctx.GetParams(),
		}
	}
}

func (this *selectContext) Rows() Rows {
	ctx := newBuildContext()
	if err := this.b.BuildSelect(ctx, this.info); err != nil {
		return &rows{
			exe: this.exe,
			err: err,
		}
	} else {
		return &rows{
			exe:  this.exe,
			sql:  ctx.GetSql(),
			args: ctx.GetParams(),
		}
	}
}

/********** joinType **********/

type joinType int8

const (
	JOIN_INNER joinType = iota
	JOIN_LEFT
	JOIN_RIGHT
	JOIN_FULL
)

func (this joinType) String() string {
	switch this {
	case JOIN_LEFT:
		return "LEFT JOIN"
	case JOIN_RIGHT:
		return "RIGHT JOIN"
	case JOIN_FULL:
		return "FULL JOIN"
	default:
		return "JOIN"
	}
}

/********** joiner **********/

type joiner struct {
	jt joinType
	t  Table
	on Filters
}

/********** Sorters **********/

type Sorters struct {
	sorters []*sorter
}

func (this *Sorters) Add(st sortType, cols ...string) *Sorters {
	s := &sorter{st: st, columns: cols}
	this.sorters = append(this.sorters, s)
	return this
}

func (this *Sorters) AddT(st sortType, t Table, cols ...string) *Sorters {
	s := &sorter{st: st, table: t, columns: cols}
	this.sorters = append(this.sorters, s)
	return this
}

/********** sortType **********/

type sortType int8

const (
	SORT_ASC sortType = iota
	SORT_DESC
)

func (this sortType) String() string {
	switch this {
	case SORT_ASC:
		return "ASC"
	case SORT_DESC:
		return "DESC"
	default:
		return ""
	}
}

/********** sorter **********/

type sorter struct {
	st      sortType
	table   Table
	columns []string
}

// create sorter with table columns
// func NewSorter(st sortType, t Table, cols ...string) *Sorter {
// 	return &Sorter{st: st, table: t, columns: cols}
// }

// // create sorter with expression columns
// func NewSorterE(st sortType, cols ...string) *Sorter {
// 	return &Sorter{st: st, columns: cols}
// }

/********** Groupers **********/

type Groupers struct {
	groupers []*grouper
}

func (this *Groupers) Add(cols ...string) *Groupers {
	grouper := &grouper{
		columns: cols,
	}
	this.groupers = append(this.groupers, grouper)
	return this
}

func (this *Groupers) AddT(t Table, cols ...string) *Groupers {
	grouper := &grouper{
		table:   t,
		columns: cols,
	}
	this.groupers = append(this.groupers, grouper)
	return this
}

/********** grouper **********/

type grouper struct {
	table   Table
	columns []string
}
