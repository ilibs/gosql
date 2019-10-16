package gosql

type Mapper struct {
	db *DB
	SQLBuilder
}

// Table select table name
func Table(t string) *Mapper {
	db := &DB{database: Sqlx(defaultLink)}
	return &Mapper{db: db, SQLBuilder: SQLBuilder{table: t, dialect: newDialect(db.DriverName())}}
}

func (m *Mapper) ShowSQL() *Mapper {
	m.db.logging = true
	return m
}

//Where
func (m *Mapper) Where(str string, args ...interface{}) *Mapper {
	m.SQLBuilder.Where(str, args...)
	return m
}

//Update data from to map[string]interface
func (m *Mapper) Update(data map[string]interface{}) (affected int64, err error) {
	result, err := m.db.Exec(m.updateString(data), m.args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

//Create data from to map[string]interface
func (m *Mapper) Create(data map[string]interface{}) (lastInsertId int64, err error) {
	result, err := m.db.Exec(m.insertString(data), m.args...)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

//Delete data from to map[string]interface
func (m *Mapper) Delete() (affected int64, err error) {
	result, err := m.db.Exec(m.deleteString(), m.args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

//Count data from to map[string]interface
func (m *Mapper) Count() (num int64, err error) {
	err = m.db.Get(&num, m.countString(), m.args...)
	return num, err
}
