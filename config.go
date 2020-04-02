package gosql

// Config is database connection configuration
type Config struct {
	Enable       bool   `yml:"enable" toml:"enable" json:"enable"`
	Driver       string `yml:"driver" toml:"driver" json:"driver"`
	Dsn          string `yml:"dsn" toml:"dsn" json:"dsn"`
	MaxOpenConns int    `yml:"max_open_conns" toml:"max_open_conns" json:"max_open_conns"`
	MaxIdleConns int    `yml:"max_idle_conns" toml:"max_idle_conns" json:"max_idle_conns"`
	MaxLifetime  int    `yml:"max_lefttime" toml:"max_lefttime" json:"max_lefttime"`
	ShowSql      bool   `yml:"show_sql" toml:"show_sql" json:"show_sql"`
}
