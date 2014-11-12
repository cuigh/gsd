package gsd

import (
	"database/sql"
	"fmt"
)

var _Databases map[string]*Database = make(map[string]*Database)

/********** executor **********/

type executor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

/********** Database **********/

type Database struct {
	db *sql.DB
	b  builder
}

func (this *Database) Insert(table string) InsertClause {
	return newInsertContext(this.db, this.b, &insertInfo{table: table})
}

func (this *Database) Delete(table string) DeleteClause {
	return newDeleteContext(this.db, this.b, &deleteInfo{table: table})
}

func (this *Database) Update(table string) UpdateClause {
	return newUpdateContext(this.db, this.b, &updateInfo{table: table})
}

func (this *Database) Select(columns *Columns) SelectClause {
	return newSelectContext(this.db, this.b, &selectInfo{columns: columns.columns, distinct: columns.distinct})
}

func (this *Database) Execute(query string, args ...interface{}) ExecuteClause {
	return newExecuteContext(this.db, query, args)
}

// Transact begin a transaction, the transaction will automatic Commit or Rollback according to return value of handler
func (this *Database) Transact(f func(tx Transaction) error) (err error) {
	trans, err := this.db.Begin()
	if err != nil {
		return err
	}

	tx := newTransaction(trans, this.b)

	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
		if err == nil {
			err = tx.Commit()
		} else {
			err = tx.Rollback()
		}
	}()

	err = f(tx)
	return
}

// open database
func Open(name string) (db *Database, err error) {
	var ok bool
	if db, ok = _Databases[name]; ok {
		return
	}

	var cfg *Config
	if cfg, err = GetConfig(name); err != nil {
		return
	}

	// since the config only need be initialized once, we can reuse the locker for initializing databases
	_Locker.Lock()
	defer _Locker.Unlock()

	// double check
	if db, ok = _Databases[name]; ok {
		return
	}

	if db, err = newDatabase(cfg); err == nil {
		_Databases[name] = db
	}
	return
}

func newDatabase(cfg *Config) (*Database, error) {
	var (
		db  *Database
		err error
	)

	db = &Database{}
	switch cfg.Provider {
	case "mssql":
		db.b = &mssqlBuilder{}
	case "mysql":
		db.b = &mysqlBuilder{}
	default:
		return nil, fmt.Errorf("not supported database provider: %s", cfg.Provider)
	}

	db.db, err = newDB(cfg)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func newDB(cfg *Config) (db *sql.DB, err error) {
	connStr, ok := cfg.Settings["ConnString"]
	if !ok {
		return nil, fmt.Errorf("connection string of database [%s] is not configured", cfg.Name)
	}

	db, err = sql.Open(cfg.Driver, connStr)
	if err != nil {
		return
	}

	maxOpenConns := cfg.Settings.Int("MaxOpenConns", 0)
	maxIdleConns := cfg.Settings.Int("MaxIdleConns", 0)
	if maxOpenConns > 0 {
		db.SetMaxOpenConns(maxOpenConns)
	}
	if maxIdleConns > 0 {
		db.SetMaxIdleConns(maxIdleConns)
	}
	return
}
