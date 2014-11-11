package gsd

type Result interface {
	RowsAffected() (int64, error)
}

type InsertResult interface {
	Result
	LastInsertId() (int64, error)
}

type ExecuteResult interface {
	Result
	LastInsertId() (int64, error)
}
