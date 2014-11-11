package gsd

import (
	"database/sql"
)

type Transaction interface {
	Insert(table string) InsertClause
	Delete(table string) DeleteClause
	Update(table string) UpdateClause
	Select(columns *Columns) SelectClause
	Execute(query string, args ...interface{}) ExecuteClause
}

type transaction struct {
	tx *sql.Tx
	b  builder
}

func newTransaction(tx *sql.Tx, b builder) *transaction {
	return &transaction{
		tx: tx,
		b:  b,
	}
}

func (this *transaction) Insert(table string) InsertClause {
	return newInsertContext(this.tx, this.b, &insertInfo{table: table})
}

func (this *transaction) Delete(table string) DeleteClause {
	return newDeleteContext(this.tx, this.b, &deleteInfo{table: table})
}

func (this *transaction) Update(table string) UpdateClause {
	return newUpdateContext(this.tx, this.b, &updateInfo{table: table})
}

func (this *transaction) Select(columns *Columns) SelectClause {
	return newSelectContext(this.tx, this.b, &selectInfo{columns: columns.columns})
}

func (this *transaction) Execute(query string, args ...interface{}) ExecuteClause {
	return newExecuteContext(this.tx, query, args)
}

func (this *transaction) Commit() error {
	return this.tx.Commit()
}

func (this *transaction) Rollback() error {
	return this.tx.Rollback()
}
