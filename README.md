# gosql
This is a database operations class library based on [sqlx](https://github.com/jmoiron/sqlx), which is simple and will continue to be simple

<a href="https://travis-ci.org/ilibs/gosql"><img src="https://travis-ci.org/ilibs/gosql.svg" alt="Build Status"></a>
<a href="https://codecov.io/gh/ilibs/gosql"><img src="https://codecov.io/gh/ilibs/gosql/branch/master/graph/badge.svg" alt="codecov"></a>
<a href="https://goreportcard.com/report/github.com/ilibs/gosql"><img src="https://goreportcard.com/badge/github.com/ilibs/gosql" alt="Go Report Card
"></a>
<a href="https://godoc.org/github.com/ilibs/gosql"><img src="https://godoc.org/github.com/ilibs/gosql?status.svg" alt="GoDoc"></a>
<a href="https://opensource.org/licenses/mit-license.php" rel="nofollow"><img src="https://badges.frapsoft.com/os/mit/mit.svg?v=103"></a>

## Usage

Connection database and use sqlx native function,See the https://github.com/jmoiron/sqlx

```go
import "github.com/ilibs/gosql"

func main(){
    configs := make(map[string]*Config)

    configs["default"] = &gosql.Config{
        Enable:  true,
        Driver:  "mysql",
        Dsn:     "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Asia%2FShanghai",
        ShowSql: true,
    }

    //connection database
    gosql.Connect(configs)

    gosql.DB().QueryRowx("select * from users where id = 1")
}

```

Default connection key is `default` config, So you can use a simpler wraper function

```go
//Exec
gosql.Exec("insert into users(name,email,created_at,updated_at) value(?,?,?,?)","test","test@gmail.com",time.Now(),time.Now())

//Queryx
rows,err := gosql.Queryx("select * from users")
for rows.Next() {
    user := &Users{}
    err = rows.StructScan(user)
}

//QueryRowx
user := &Users{}
err := gosql.QueryRowx("select * from users where id = 1").StructScan(user)

//Get
user := &Users{}
err := gosql.Get(user,"select * from users where id = 1")

//Select
users := make([]*Users)
err := gosql.Select(&users,"select * from users")

//Change database
db := gosql.Use("test")
db.Queryx("select * from tests")
```

## CRUD interface

```go
type Users struct {
	Id        int       `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Status    int       `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (u *Users) DbName() string {
	return "default"
}

func (u *Users) TableName() string {
	return "users"
}

func (u *Users) PK() string {
	return "id"
}

//Get
user := &Users{}
gosql.Model(user).Where("id=1").Get()

//All
user := make([]*Users,0)
gosql.Model(&user).All()

//Create and Timestamp Tracking
gosql.Model(&User{Name:"test",Email:"test@gmail.com"}).Create()

//Update
gosql.Model(&User{Name:"test2",Email:"test@gmail.com"}).Where("id=1").Update()
//If you need to update the zero value, you can do so
gosql.Model(&User{Status:0}).Where("id=1").Update("status")

//Delete
gosql.Model(&User{}).Where("id=1").Delete()

```


## Transaction
The `Tx` function has a callback function, if an error is returned, the transaction rollback

```go
gosql.Tx(func(tx *sqlx.Tx) error {
    for id := 1; id < 10; id++ {
        user := &Users{
            Id:    id,
            Name:  "test" + strconv.Itoa(id),
            Email: "test" + strconv.Itoa(id) + "@test.com",
        }

        Model(user, tx).Create()

        if id == 8 {
            return errors.New("interrupt the transaction")
        }
    }

    return nil
})
```

> If you need to invoke context, you can use `gosql.Txx`

## Timestamp Tracking
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
    "id":         id,
    "name":       "test" + strconv.Itoa(int(id)),
    "email":      "test@test.com",
    "created_at": "2018-07-11 11:58:21",
    "updated_at": "2018-07-11 11:58:21",
})

//Update
gosql.Table("users").Where("id = ?", id).Update(map[string]interface{}{
    "name":  "fifsky",
    "email": "fifsky@test.com",
})

//Delete
gosql.Table("users").Where("id = ?", id).Delete()

//Count
gosql.Table("users").Where("id = ?", id).Count()

//Change database
gosql.Use("db2").Table("users").Where("id = ?", id).Count()

//Transaction `tx` is *sqlx.Tx
gosql.Table("users",tx).Where("id = ?", id}).Count()
```

## Thanks

sqlx https://github.com/jmoiron/sqlx