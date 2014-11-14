package gsd

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"
)

// Scanner valid types: int64, float64, bool, []byte, string, time.Time, nil - for NULL values

/********** NullTime **********/

type NullInt64 sql.NullInt64

// Scan implements the Scanner interface.
func (this *NullInt64) Scan(value interface{}) (err error) {
	return (*sql.NullInt64)(this).Scan(value)
}

/********** NullTime **********/

type NullFloat64 sql.NullFloat64

// Scan implements the Scanner interface.
func (this *NullFloat64) Scan(value interface{}) (err error) {
	return (*sql.NullFloat64)(this).Scan(value)
}

/********** NullTime **********/

type NullBool sql.NullBool

// Scan implements the Scanner interface.
func (this *NullBool) Scan(value interface{}) (err error) {
	return (*sql.NullBool)(this).Scan(value)
}

/********** NullTime **********/

type NullString sql.NullString

// Scan implements the Scanner interface.
func (this *NullString) Scan(value interface{}) (err error) {
	return (*sql.NullString)(this).Scan(value)
}

/********** NullTime **********/

type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface.
func (this *NullTime) Scan(value interface{}) (err error) {
	if value == nil {
		this.Time, this.Valid = time.Time{}, false
		return
	}

	if this.Time, this.Valid = value.(time.Time); this.Valid {
		return
	}

	return fmt.Errorf("Can't convert %T to time.Time", value)
}

// Value implements the driver Valuer interface.
func (this NullTime) Value() (driver.Value, error) {
	if this.Valid {
		return this.Time, nil
	}

	return nil, nil
}

/********** NullInt32 **********/

type NullInt32 struct {
	Int32 int32
	Valid bool // Valid is true if Int32 is not NULL
}

// Scan implements the Scanner interface.
func (this *NullInt32) Scan(value interface{}) (err error) {
	if value == nil {
		this.Int32, this.Valid = 0, false
		return
	}

	switch v := value.(type) {
	case int64:
		this.Int32, this.Valid = int32(v), true
		return
	}

	this.Valid = false
	return fmt.Errorf("Can't convert %T to int32", value)
}

// Value implements the driver Valuer interface.
func (this NullInt32) Value() (driver.Value, error) {
	if this.Valid {
		return this.Int32, nil
	}

	return nil, nil
}

/********** NullFloat32 **********/

type NullFloat32 struct {
	Float32 float32
	Valid   bool // Valid is true if Float32 is not NULL
}

// Scan implements the Scanner interface.
func (this *NullFloat32) Scan(value interface{}) (err error) {
	if value == nil {
		this.Float32, this.Valid = 0, false
		return
	}

	switch v := value.(type) {
	case float64:
		this.Float32, this.Valid = float32(v), true
		return
	}

	this.Valid = false
	return fmt.Errorf("Can't convert %T to float32", value)
}

// Value implements the driver Valuer interface.
func (this NullFloat32) Value() (driver.Value, error) {
	if this.Valid {
		return this.Float32, nil
	}

	return nil, nil
}
