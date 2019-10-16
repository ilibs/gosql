package gosql

type postgresDialect struct {
	commonDialect
}

func init() {
	RegisterDialect("postgres", &postgresDialect{})
}

func (postgresDialect) GetName() string {
	return "postgres"
}
