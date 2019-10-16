package gosql

import (
	"fmt"
	"testing"
)

func TestSQLBuilder_queryString(t *testing.T) {

	b := &SQLBuilder{
		dialect: mustGetDialect("mysql"),
		table:   "users",
		order:   "id desc",
		limit:   "0",
		offset:  "10",
	}

	b.Where("id = ?", 1)

	if b.queryString() != "SELECT * FROM `users` WHERE (id = ?) ORDER BY id desc LIMIT 0 OFFSET 10;" {
		t.Error("sql builder query error", b.queryString())
	}
	fmt.Println(b.queryString())
}

func TestSQLBuilder_queryForceIndexString(t *testing.T) {

	b := &SQLBuilder{
		dialect:    mustGetDialect("mysql"),
		table:      "users",
		order:      "id desc",
		forceIndex: "idx_user",
		limit:      "0",
		offset:     "10",
	}

	b.Where("id = ?", 1)

	if b.queryString() != "SELECT * FROM `users` force index(idx_user) WHERE (id = ?) ORDER BY id desc LIMIT 0 OFFSET 10;" {
		t.Error("sql builder query error", b.queryString())
	}
	fmt.Println(b.queryString())
}

func TestSQLBuilder_insertString(t *testing.T) {

	b := &SQLBuilder{
		dialect: mustGetDialect("mysql"),
		table:   "users",
	}

	query := b.insertString(map[string]interface{}{
		"id":         1,
		"name":       "test",
		"email":      "test@test.com",
		"created_at": "2018-07-11 11:58:21",
		"updated_at": "2018-07-11 11:58:21",
	})

	if query != "INSERT INTO `users` (`created_at`,`email`,`id`,`name`,`updated_at`) VALUES(?,?,?,?,?);" {
		t.Error("sql builder insert error", query)
	}
}

func TestSQLBuilder_updateString(t *testing.T) {

	b := &SQLBuilder{
		dialect: mustGetDialect("mysql"),
		table:   "users",
	}

	b.Where("id = ?", 1)
	query := b.updateString(map[string]interface{}{
		"name":  "test",
		"email": "test@test.com",
	})

	if query != "UPDATE `users` SET `email`=?,`name`=? WHERE (id = ?);" {
		t.Error("sql builder update error", query)
	}

	fmt.Println(query, b.args)
}

func TestSQLBuilder_deleteString(t *testing.T) {

	b := &SQLBuilder{
		dialect: mustGetDialect("mysql"),
		table:   "users",
	}

	b.Where("id = ?", 1)

	query := b.deleteString()

	if query != "DELETE FROM `users` WHERE (id = ?);" {
		t.Error("sql builder delete error", query)
	}
}

func TestSQLBuilder_countString(t *testing.T) {

	b := &SQLBuilder{
		dialect: mustGetDialect("mysql"),
		table:   "users",
	}
	b.Where("id = ?", 1)

	query := b.countString()

	if query != "SELECT count(*) FROM `users` WHERE (id = ?);" {
		t.Error("sql builder count error", query)
	}
}

func TestSQLBuilder_Dialect(t *testing.T) {
	testData := map[string]string{
		"mysql":    "INSERT INTO `users` (`created_at`,`email`,`id`,`name`,`updated_at`) VALUES(?,?,?,?,?);",
		"postgres": `INSERT INTO "users" ("created_at","email","id","name","updated_at") VALUES(?,?,?,?,?);`,
		"sqlite3":  `INSERT INTO "users" ("created_at","email","id","name","updated_at") VALUES(?,?,?,?,?);`,
	}

	for k, v := range testData {
		b := &SQLBuilder{
			dialect: mustGetDialect(k),
			table:   "users",
		}

		query := b.insertString(map[string]interface{}{
			"id":         1,
			"name":       "test",
			"email":      "test@test.com",
			"created_at": "2018-07-11 11:58:21",
			"updated_at": "2018-07-11 11:58:21",
		})

		if query != v {
			t.Error(fmt.Sprintf("sql builder %s dialect insert error", k), query)
		}
	}
}
