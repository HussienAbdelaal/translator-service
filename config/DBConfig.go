package config

type DBConfig struct {
	Host     string
	Username string
	Password string
	Database string
}

func (c *DBConfig) GetDSN() string {
	// Construct the DSN (Data Source Name) string for PostgreSQL
	// Example: "postgres://username:password@host/database"
	dsn := "postgres://" + c.Username + ":" + c.Password + "@" + c.Host + ":" + "/" + c.Database
	return dsn
}
