package db

// Driver defines the interface for database drivers
type Driver interface {
	Connect(dsn string) error
	Close() error
	GetDatabases() ([]string, error)
	GetTables(database string) ([]string, error)
	ExecuteQuery(query string) (*QueryResult, error)
	GetDriverName() string
}

// QueryResult represents the result of a SQL query
type QueryResult struct {
	Columns      []string
	Rows         [][]string
	RowsAffected int64
	Error        error
}
