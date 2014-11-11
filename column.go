package gsd

import (
// "fmt"
)

/********** column **********/

type column interface {
	Alias() string
}

/********** Columns **********/

func C(distinct bool) *Columns {
	return &Columns{distinct: distinct}
}

type Columns struct {
	distinct bool
	columns  []column
}

func (this *Columns) Distinct(d bool) *Columns {
	this.distinct = d
	return this
}

// add table columns
func (this *Columns) Add(t Table, cols ...string) *Columns {
	for _, col := range cols {
		this.columns = append(this.columns, &normalColumn{table: t, column: col})
	}
	return this
}

// add table column with alias
func (this *Columns) AddA(t Table, col, alias string) *Columns {
	this.columns = append(this.columns, &normalColumn{table: t, column: col, alias: alias})
	return this
}

// add expression column, like 'COUNT(*)'
func (this *Columns) AddE(expr, alias string) *Columns {
	this.columns = append(this.columns, &exprColumn{expr: expr, alias: alias})
	return this
}

/********** NormalColumn **********/

type normalColumn struct {
	table  Table
	column string
	alias  string
}

func (this *normalColumn) Alias() string {
	return this.alias
}

/********** ExprColumn **********/

type exprColumn struct {
	expr  string
	alias string
}

func (this *exprColumn) Alias() string {
	return this.alias
}
