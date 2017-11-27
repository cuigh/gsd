# gsd

> [DEPRECATED] This package now is a part of [auxo](https://github.com/cuigh/auxo), see [auxo/gsd](https://github.com/cuigh/auxo/tree/master/db/gsd).

gsd is a Simple, fluent SQL data access framework. It supports various types of database, like mysql/mssql/sqlite etc.

## Install

To install gsd, just use `go get` command. 

```
$ go get github.com/cuigh/gsd
```

## Usage

gsd's API is very simular to native SQL sytax.

## Configure

For now, gsd only supports initializing databases from config file. There is a sample config file in the package(database.sql.conf):

```
<databases>
	<database name="Test" provider="mysql">
		<setting name="ConnString" value="user:password@tcp(localhost:3306)/Test?parseTime=true"/>
		<setting name="MaxIdleConns" value="1"/>
		<setting name="MaxOpenConns" value="100"/>
	</database>
</databases>
```
You must set [ConfigPath] before you go to next step:

```
gsd.ConfigPath = "./database.sql.conf"
```
Now you can open a database:

```
db, err := Open("Test")
......
```
### INSERT

```
v := gsd.InsertValues{
	"ID":         10,
	"NAME":    	  "Clothes",
	"COUNT": 	  0,
	"ENTER_USER": 1,
	"ENTER_TIME": time.Now(),
}
r, err := db.Insert("Category").Values(v).Result()
```
### DELETE

```
f := gsd.F().Add("ID", 10)
r, err := db.Delete("Category").Where(f).Result()
```
### UPDATE

```
f := gsd.F().Add("ID", 2)
v := gsd.UpdateValues{
	"ENTER_TIME": gsd.UV(time.Now()),
	"COUNT": gsd.UVT(gsd.UPDATE_INC, 1),
}
r, err := db.Update("Category").Set(v).Where(f).Result()
```
### SELECT

```
type Category struct {
	ID    	  int32
	Name  	  string
	Count 	  int32
	EnterUser int32		`gsd:"ENTER_USER"`
	EnterTime time.Time `gsd:"ENTER_TIME"`
}

t := gsd.T("Category")
f := gsd.F().AddT("ID", gsd.FILTER_GT, 0)
r := db.Select(t1.C("ID", "NAME", "ENTER_USER", "ENTER_TIME")).From(t).Where(f).Limit(0, 3).Rows()
objs := []*Category{}
if err := r.All(&objs); err != nil {
	log.Fatal(err)
}
```
### TRANSACTION

```
err := db.Transact(func(tx Transaction) error {
	obj := Category{}	
	err := tx.Execute("SELECT ID, NAME, ENTER_TIME FROM Category WHERE ID=?", 2).Row().ScanObj(&obj)
	if err != nil{
		return err
	}

	f := gsd.F().Add("ID", 3)
	v := gsd.UpdateValues{
		"Name":        gsd.UV(obj.Name),
	}
	_, err = tx.Update("Category").Set(v).Where(f).Result()
	return err
})
if err != nil {
	log.Fatal(err)
}
```
