package gsd

import (
// "bytes"
)

/********** filterType **********/

type filterType int8

const (
	FILTER_EQ filterType = iota
	FILTER_NE
	FILTER_LT
	FILTER_GT
	FILTER_LTE
	FILTER_GTE
	FILTER_IN
	FILTER_LK
)

/********** oneColumnFilter **********/

type oneColumnFilter struct {
	table  Table
	column string
	ft     filterType
	value  interface{}
}

/********** twoColumnFilter **********/

type twoColumnFilter struct {
	table1  Table
	column1 string
	table2  Table
	column2 string
	ft      filterType
}

/********** exprFilter **********/

type exprFilter struct {
	expr string
}

/********** Filters **********/

type Filters interface {
	Not() Filters
	And(f Filters) Filters
	Or(f Filters) Filters
}

type BasicFilters interface {
	Filters
	Add(col string, value interface{}) BasicFilters
	AddT(col string, ft filterType, value interface{}) BasicFilters
	AddF(t Table, col string, ft filterType, value interface{}) BasicFilters
	AddJ(t1 Table, col1 string, ft filterType, t2 Table, col2 string) BasicFilters
	AddE(expr string) BasicFilters
}

/********** basicFilters **********/

type basicFilters struct {
	items []interface{}
}

func F() BasicFilters {
	return &basicFilters{}
}

func (this *basicFilters) Not() Filters {
	return newNotFilters(this)
}

func (this *basicFilters) And(f Filters) Filters {
	return newAndFilters(this, f)
}

func (this *basicFilters) Or(f Filters) Filters {
	return newOrFilters(this, f)
}

// add simple filter
func (this *basicFilters) Add(col string, value interface{}) BasicFilters {
	return this.AddF(nil, col, FILTER_EQ, value)
}

// add filter with specific type
func (this *basicFilters) AddT(col string, ft filterType, value interface{}) BasicFilters {
	return this.AddF(nil, col, ft, value)
}

// add filter with full information
func (this *basicFilters) AddF(t Table, col string, ft filterType, value interface{}) BasicFilters {
	f := &oneColumnFilter{
		table:  t,
		column: col,
		ft:     ft,
		value:  value,
	}
	this.items = append(this.items, f)
	return this
}

// add two columns filter, normally for JOIN clause
func (this *basicFilters) AddJ(t1 Table, col1 string, ft filterType, t2 Table, col2 string) BasicFilters {
	f := &twoColumnFilter{
		table1:  t1,
		column1: col1,
		ft:      ft,
		table2:  t2,
		column2: col2,
	}
	this.items = append(this.items, f)
	return this
}

// add expression filter
func (this *basicFilters) AddE(expr string) BasicFilters {
	f := &exprFilter{
		expr: expr,
	}
	this.items = append(this.items, f)
	return this
}

/********** notFilters **********/

type notFilters struct {
	inner Filters
}

func newNotFilters(f Filters) *notFilters {
	return &notFilters{inner: f}
}

func (this *notFilters) Not() Filters {
	return newNotFilters(this)
}

func (this *notFilters) And(f Filters) Filters {
	return newAndFilters(this, f)
}

func (this *notFilters) Or(f Filters) Filters {
	return newOrFilters(this, f)
}

/********** andFilters **********/

type andFilters struct {
	left  Filters
	right Filters
}

func newAndFilters(l, r Filters) *andFilters {
	return &andFilters{left: l, right: r}
}

func (this *andFilters) Not() Filters {
	return newNotFilters(this)
}

func (this *andFilters) And(f Filters) Filters {
	return newAndFilters(this, f)
}

func (this *andFilters) Or(f Filters) Filters {
	return newOrFilters(this, f)
}

/********** orFilters **********/

type orFilters struct {
	left  Filters
	right Filters
}

func newOrFilters(l, r Filters) *orFilters {
	return &orFilters{left: l, right: r}
}

func (this *orFilters) Not() Filters {
	return newNotFilters(this)
}

func (this *orFilters) And(f Filters) Filters {
	return newAndFilters(this, f)
}

func (this *orFilters) Or(f Filters) Filters {
	return newOrFilters(this, f)
}
