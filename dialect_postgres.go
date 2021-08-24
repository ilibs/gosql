package gosql

import "strconv"

type postgresDialect struct {
	commonDialect
	count int
}

func init() {
	RegisterDialect("postgres", &postgresDialect{})
}

func (postgresDialect) GetName() string {
	return "postgres"
}

func (p *postgresDialect) Placeholder() string {
	p.count++
	return "$" + strconv.Itoa(p.count)
}
