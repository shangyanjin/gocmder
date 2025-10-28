package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// Postgres implements the Driver interface for PostgreSQL
type Postgres struct {
	conn *sql.DB
}

// NewPostgres creates a new PostgreSQL driver
func NewPostgres() *Postgres {
	return &Postgres{}
}

// Connect connects to PostgreSQL database
func (p *Postgres) Connect(dsn string) error {
	var err error
	p.conn, err = sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open PostgreSQL connection: %w", err)
	}

	err = p.conn.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}

	return nil
}

// Close closes the connection
func (p *Postgres) Close() error {
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}

// GetDatabases returns list of databases
func (p *Postgres) GetDatabases() ([]string, error) {
	if p.conn == nil {
		return nil, fmt.Errorf("not connected")
	}

	query := "SELECT datname FROM pg_database WHERE datistemplate = false"
	rows, err := p.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var databases []string
	for rows.Next() {
		var db string
		if err := rows.Scan(&db); err != nil {
			return nil, err
		}
		// Filter system databases
		if db != "postgres" {
			databases = append(databases, db)
		}
	}

	return databases, rows.Err()
}

// GetTables returns list of tables in a database
func (p *Postgres) GetTables(database string) ([]string, error) {
	if p.conn == nil {
		return nil, fmt.Errorf("not connected")
	}

	query := "SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname = 'public'"
	rows, err := p.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}

	return tables, rows.Err()
}

// ExecuteQuery executes a SQL query
func (p *Postgres) ExecuteQuery(query string) (*QueryResult, error) {
	if p.conn == nil {
		return nil, fmt.Errorf("not connected")
	}

	result := &QueryResult{}

	// Try as SELECT query first
	rows, err := p.conn.Query(query)
	if err != nil {
		// Try as DML/DDL statement
		res, execErr := p.conn.Exec(query)
		if execErr != nil {
			result.Error = execErr
			return result, execErr
		}

		affected, _ := res.RowsAffected()
		result.RowsAffected = affected
		return result, nil
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		result.Error = err
		return result, err
	}
	result.Columns = columns

	// Get rows
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			result.Error = err
			return result, err
		}

		row := make([]string, len(columns))
		for i, val := range values {
			if val == nil {
				row[i] = "NULL"
			} else {
				row[i] = fmt.Sprintf("%v", val)
			}
		}
		result.Rows = append(result.Rows, row)
	}

	result.RowsAffected = int64(len(result.Rows))
	return result, rows.Err()
}

// GetDriverName returns the driver name
func (p *Postgres) GetDriverName() string {
	return "PostgreSQL"
}
