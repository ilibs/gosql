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
gosql.Sqlx() //return native sqlx
```
## Change database
```go
gosql.Use(name string) 
gosql.Use(db).Table("xxxx").Where("id = ?",1).Update(map[string]interface{}{"name":"test"})
gosql.Use(db).Model(&Users{}}).Get()
```

## Transaction context switching
```go
gosql.Use(db).Tx(func(tx *gosql.DB){
    tx.Table("xxxx").Where("id = ?",1).Get(&user)
    tx.Model(&Users{}).Get()	
})
```
