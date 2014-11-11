package gsd

/********** Common **********/

type ResultClause interface {
	Result() (Result, error)
}

type RowClause interface {
	Row() Row
	Rows() Rows
}

/********** Select Clauses **********/

type SelectClause interface {
	From(t Table) FromClause
}

type FromClause interface {
	LimitClause
	RowClause
	Join(t Table, on Filters) FromClause
	LeftJoin(t Table, on Filters) FromClause
	RightJoin(t Table, on Filters) FromClause
	FullJoin(t Table, on Filters) FromClause
	Where(f Filters) WhereClause
}

type WhereClause interface {
	RowClause
	LimitClause
	OrderBy(s *Sorters) OrderByClause
	GroupBy(g *Groupers) GroupByClause
}

type GroupByClause interface {
	LimitClause
	OrderBy(s *Sorters) OrderByClause
	Having(f Filters) HavingClause
}

type HavingClause interface {
	RowClause
	LimitClause
	OrderBy(s *Sorters) OrderByClause
}

type OrderByClause interface {
	RowClause
	LimitClause
}

type LimitClause interface {
	Limit(skip, take int32) RowClause
}

/********** Update Clauses **********/

type UpdateClause interface {
	Set(values UpdateValues) SetClause
}

type SetClause interface {
	ResultClause
	Where(f Filters) ResultClause
}

/********** Delete Clauses **********/

type DeleteClause interface {
	Where(f Filters) ResultClause
}

/********** Insert Clauses **********/

type InsertClause interface {
	Values(values InsertValues) InsertResultClause
}

type InsertResultClause interface {
	Result() (InsertResult, error)
}

/********** Execute Clauses **********/

type ExecuteClause interface {
	ExecuteResultClause
	RowClause
}

type ExecuteResultClause interface {
	Result() (ExecuteResult, error)
}
