package gosql

type sqlite3Dialect struct {
	commonDialect
}

func init() {
	RegisterDialect("sqlite3", &sqlite3Dialect{})
}

func (sqlite3Dialect) GetName() string {
	return "sqlite3"
}
