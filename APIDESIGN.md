# API Design

## Model interface
```go
type IModel interface {
	TableName() string
	PK() string
}
```

> Remove the V1 version DbName(),Use the Use() function instead

## Use sqlx
```go
gosql.DB() //return native sqlx
```
## Change database
```go
gosql.Use(name string) 
gosql.Use(db).Table("xxxx").Where("id = ?",1).Update(map[string]interface{}{"name":"test"})
gosql.Use(db).Model(&Users{}}).Get()
```

## Transaction context switching
```go
gosql.WithTx(tx *sqlx.Tx)
gosql.WithTx(tx).Table("xxxx").Where("id = ?",1).Get(&user)
gosql.WithTx(tx).Model(&Users{}).Get()
```
