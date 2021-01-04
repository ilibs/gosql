# gosql
The package based on [sqlx](https://github.com/jmoiron/sqlx), It's simple and keep simple

<a href="https://github.com/ilibs/gosql/actions"><img src="https://github.com/ilibs/gosql/workflows/gosql/badge.svg" alt="Build Status"></a>
<a href="https://codecov.io/gh/ilibs/gosql"><img src="https://codecov.io/gh/ilibs/gosql/branch/master/graph/badge.svg" alt="codecov"></a>
<a href="https://goreportcard.com/report/github.com/ilibs/gosql"><img src="https://goreportcard.com/badge/github.com/ilibs/gosql" alt="Go Report Card
"></a>
<a href="https://godoc.org/github.com/ilibs/gosql"><img src="https://godoc.org/github.com/ilibs/gosql?status.svg" alt="GoDoc"></a>
<a href="https://opensource.org/licenses/mit-license.php" rel="nofollow"><img src="https://badges.frapsoft.com/os/mit/mit.svg?v=103"></a>

⚠️ Because of some disruptive changes, The current major version is upgraded to V2，If you continue with V1, you can check out the v1 branches [https://github.com/ilibs/gosql/tree/v1](https://github.com/ilibs/gosql/tree/v1)

## V2 ChangeLog
- Remove the second argument to the Model() and Table() functions and replace it with WithTx(tx)
- Remove Model interface DbName() function,use the Use() function 
- Uniform API design specification, see [APIDESIGN](APIDESIGN.md)
- Relation add `connection:"db2"` struct tag, Solve the cross-library connection problem caused by deleting DbName()
- Discard the WithTx function

## Usage

Connection database and use sqlx original function,See the https://github.com/jmoiron/sqlx

```go
import (
    _ "github.com/go-sql-driver/mysql" //mysql driver
    "github.com/ilibs/gosql/v2"
)

func main(){
    configs := make(map[string]*gosql.Config)

    configs["default"] = &gosql.Config{
        Enable:  true,
        Driver:  "mysql",
        Dsn:     "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Asia%2FShanghai",
        ShowSql: true,
    }

    //connection database
    gosql.Connect(configs)
    gosql.QueryRowx("select * from users where id = 1")
}

```

Use `default` database, So you can use wrapper function

```go
//Exec
gosql.Exec("insert into users(name,email,created_at,updated_at) value(?,?,?,?)","test","test@gmail.com",time.Now(),time.Now())

//Queryx
rows,err := gosql.Queryx("select * from users")
for rows.Next() {
    user := &Users{}
    err = rows.StructScan(user)
}
rows.Close()

//QueryRowx
user := &Users{}
err := gosql.QueryRowx("select * from users where id = ?",1).StructScan(user)

//Get
user := &Users{}
err := gosql.Get(user,"select * from users where id = ?",1)

//Select
users := make([]Users)
err := gosql.Select(&users,"select * from users")

//Change database
db := gosql.Use("test")
db.Queryx("select * from tests")
```

You can also set the default database connection name

```go
gosql.SetDefaultLink("log")
gosql.Connect(configs)
```

> `gosql.Get` etc., will use the configuration with the connection name `log` 

## Using struct

```go
type Users struct {
	Id        int       `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Status    int       `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (u *Users) TableName() string {
	return "users"
}

func (u *Users) PK() string {
	return "id"
}

//Get
user := &Users{}
gosql.Model(user).Where("id=?",1).Get()

//All
user := make([]Users,0)
gosql.Model(&user).All()

//Create and auto set CreatedAt
gosql.Model(&User{Name:"test",Email:"test@gmail.com"}).Create()

//Update
gosql.Model(&User{Name:"test2",Email:"test@gmail.com"}).Where("id=?",1).Update()
//If you need to update the zero value, you can do so
gosql.Model(&User{Status:0}).Where("id=?",1).Update("status")

//Delete
gosql.Model(&User{}).Where("id=?",1).Delete()

```

If you use struct to generate where conditions

```go
//Get where id = 1 and name = "test1"
user := &Users{Id:1,Name:"test1"}
gosql.Model(&user).Get()

//Update default use primary key as the condition
gosql.Model(&User{Id:1,Name:"test2"}).Update()
//Use custom conditions
//Builder => UPDATE users SET `id`=?,`name`=?,`updated_at`=? WHERE (status = ?)
gosql.Model(&User{Id:1,Name:"test2"}).Where("status = ?",1).Update()

//Delete
gosql.Model(&User{Id:1}).Delete()
```

But the zero value is filtered by default, you can specify fields that are not filtered. For example

```go
user := &Users{Id:1,Status:0}
gosql.Model(&user).Get("status")
```

> You can use the [genstruct](https://github.com/fifsky/genstruct) tool to quickly generate database structs

## Transaction
The `Tx` function has a callback function, if an error is returned, the transaction rollback

```go
gosql.Tx(func(tx *gosql.DB) error {
    for id := 1; id < 10; id++ {
        user := &Users{
            Id:    id,
            Name:  "test" + strconv.Itoa(id),
            Email: "test" + strconv.Itoa(id) + "@test.com",
        }
		
		//v2 support, do some database operations in the transaction (use 'tx' from this point, not 'gosql')
        tx.Model(user).Create()

        if id == 8 {
            return errors.New("interrupt the transaction")
        }
    }

    //query with transaction
    var num int
    err := tx.QueryRowx("select count(*) from user_id = 1").Scan(&num)

    if err != nil {
        return err
    }

    return nil
})
```

> If you need to invoke context, you can use `gosql.Txx`

Now support gosql.Begin() or gosql.Use("other").Begin() for example:
```go
tx, err := gosql.Begin()
if err != nil {
    return err
}

for id := 1; id < 10; id++ {
    _, err := tx.Exec("INSERT INTO users(id,name,status,created_at,updated_at) VALUES(?,?,?,?,?)", id, "test"+strconv.Itoa(id), 1, time.Now(), time.Now())
    if err != nil {
        return tx.Rollback()
    }
}

return tx.Commit()
```

## Automatic time
If your fields contain the following field names, they will be updated automatically

```
AUTO_CREATE_TIME_FIELDS = []string{
    "create_time",
    "create_at",
    "created_at",
    "update_time",
    "update_at",
    "updated_at",
}
AUTO_UPDATE_TIME_FIELDS = []string{
    "update_time",
    "update_at",
    "updated_at",
}
```


## Using Map
`Create` `Update` `Delete` `Count` support `map[string]interface`,For example:

```go
//Create
gosql.Table("users").Create(map[string]interface{}{
    "id":         1,
    "name":       "test",
    "email":      "test@test.com",
    "created_at": "2018-07-11 11:58:21",
    "updated_at": "2018-07-11 11:58:21",
})

//Update
gosql.Table("users").Where("id = ?", 1).Update(map[string]interface{}{
    "name":  "fifsky",
    "email": "fifsky@test.com",
})

//Delete
gosql.Table("users").Where("id = ?", 1).Delete()

//Count
gosql.Table("users").Where("id = ?", 1).Count()

//Change database
gosql.Use("db2").Table("users").Where("id = ?", 1).Count()

//Transaction `tx`
tx.Table("users").Where("id = ?", 1}).Count()
```


## sql.Null*
Now Model support sql.Null* field's, Note, however, that if sql.Null* is also filtered by zero values,For example

```go
type Users struct {
	Id          int            `db:"id"`
	Name        string         `db:"name"`
	Email       string         `db:"email"`
	Status      int            `db:"status"`
	SuccessTime sql.NullString `db:"success_time" json:"success_time"`
	CreatedAt   time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at" json:"updated_at"`
}

user := &Users{
    Id: 1,
    SuccessTime: sql.NullString{
        String: "2018-09-03 00:00:00",
        Valid:  false,
    }
}

err := gosql.Model(user).Get()
```

Builder SQL:

```
Query: SELECT * FROM users WHERE (id=?);
Args:  []interface {}{1}
Time:  0.00082s
```

If `sql.NullString` of `Valid` attribute is false, SQL builder will ignore this zero value


## gosql.Expr
Reference GORM Expr, Resolve update field self-update problem
```go
gosql.Table("users").Update(map[string]interface{}{
    "id":2,
    "count":gosql.Expr("count+?",1)
})
//Builder SQL
//UPDATE `users` SET `count`=count + ?,`id`=?; [1 2]
```


## "In" Queries

Because database/sql does not inspect your query and it passes your arguments directly to the driver, it makes dealing with queries with IN clauses difficult:

```go
SELECT * FROM users WHERE level IN (?);
```

`sqlx.In` is encapsulated In `gosql` and can be queried using the following schema

```go
var levels = []int{4, 6, 7}
rows, err := gosql.Queryx("SELECT * FROM users WHERE level IN (?);", levels)

//or

user := make([]Users, 0)
err := gosql.Select(&user, "select * from users where id in(?)",[]int{1,2,3})
```

## Relation
gosql used the golang structure to express the relationships between tables,You only need to use the `relation` Tag to specify the associated field, see example

⚠️ Since version v2, the relation query across library connections needs to be specified using `connection` tag


```go
type MomentList struct {
	models.Moments
	User   *models.Users    `json:"user" db:"-" relation:"user_id,id"`         //one-to-one
	Photos []models.Photos `json:"photos" db:"-" relation:"id,moment_id" connection:"db2"`     //one-to-many
}
```

Get single result

```go
moment := &MomentList{}
err := gosql.Model(moment).Where("status = 1 and id = ?",14).Get()
//output User and Photos and you get the result
```

SQL:

```sql
2018/12/06 13:27:54
	Query: SELECT * FROM `moments` WHERE (status = 1 and id = ?);
	Args:  []interface {}{14}
	Time:  0.00300s

2018/12/06 13:27:54
	Query: SELECT * FROM `moment_users` WHERE (id=?);
	Args:  []interface {}{5}
	Time:  0.00081s

2018/12/06 13:27:54
	Query: SELECT * FROM `photos` WHERE (moment_id=?);
	Args:  []interface {}{14}
	Time:  0.00093s
```

Get list result, many-to-many

```go
var moments = make([]MomentList, 0)
err := gosql.Model(&moments).Where("status = 1").Limit(10).All()
//You get the total result  for *UserMoment slice
```

SQL:

```sql
2018/12/06 13:50:59
	Query: SELECT * FROM `moments` WHERE (status = 1) LIMIT 10;
	Time:  0.00319s

2018/12/06 13:50:59
	Query: SELECT * FROM `moment_users` WHERE (id in(?));
	Args:  []interface {}{[]interface {}{5}}
	Time:  0.00094s

2018/12/06 13:50:59
	Query: SELECT * FROM `photos` WHERE (moment_id in(?, ?, ?, ?, ?, ?, ?, ?, ?, ?));
	Args:  []interface {}{[]interface {}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}}
	Time:  0.00087s
```


Relation Where:

```go
moment := &MomentList{}
err := gosql.Relation("User" , func(b *gosql.Builder) {
    //this is builder instance,
    b.Where("gender = 0")
}).Get(moment , "select * from moments")
```

## Hooks
Hooks are functions that are called before or after creation/querying/updating/deletion.

If you have defiend specified methods for a model, it will be called automatically when creating, updating, querying, deleting, and if any callback returns an error, `gosql` will stop future operations and rollback current transaction.

```
// begin transaction
BeforeChange
BeforeCreate
// update timestamp `CreatedAt`, `UpdatedAt`
// save
AfterCreate
AfterChange
// commit or rollback transaction
```
Example:

```go
func (u *Users) BeforeCreate(ctx context.Context) (err error) {
  if u.IsValid() {
    err = errors.New("can't save invalid data")
  }
  return
}

func (u *Users) AfterCreate(ctx context.Context, tx *gosql.DB) (err error) {
  if u.Id == 1 {
    u.Email = ctx.Value("email")
    tx.Model(u).Update()
  }
  return
}
```

> BeforeChange and AfterChange only  used in create/update/delete

All Hooks:

```
BeforeChange
AfterChange
BeforeCreate
AfterCreate
BeforeUpdate
AfterUpdate
BeforeDelete
AfterDelete
BeforeFind
AfterFind
```

Hook func type supports multiple ways:

```
func (u *Users) BeforeCreate()
func (u *Users) BeforeCreate() (err error)
func (u *Users) BeforeCreate(tx *gosql.DB)
func (u *Users) BeforeCreate(tx *gosql.DB) (err error)
func (u *Users) BeforeCreate(ctx context.Context)
func (u *Users) BeforeCreate(ctx context.Context) (err error)
func (u *Users) BeforeCreate(ctx context.Context, tx *rsql.DB)
func (u *Users) BeforeCreate(ctx context.Context, tx *rsql.DB) (err error)

```

 If you want to use `context` feature, you need to use below function while start a sql, or the context in callback will be nil:

1. ` gosql.WithContext(ctx).Model(...)`
1. ` gosql.Use("xxx").WithContext(ctx).Model(...)`


## Thanks

sqlx https://github.com/jmoiron/sqlx
