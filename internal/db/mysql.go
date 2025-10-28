package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// MySQL implements the Driver interface for MySQL
type MySQL struct {
	conn *sql.DB
}

// NewMySQL creates a new MySQL driver
func NewMySQL() *MySQL {
	return &MySQL{}
}

// Connect connects to MySQL database
func (m *MySQL) Connect(dsn string) error {
	var err error
	m.conn, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open MySQL connection: %w", err)
	}

	err = m.conn.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping MySQL: %w", err)
	}

	return nil
}

// Close closes the connection
func (m *MySQL) Close() error {
	if m.conn != nil {
		return m.conn.Close()
	}
	return nil
}

// GetDatabases returns list of databases
func (m *MySQL) GetDatabases() ([]string, error) {
	if m.conn == nil {
		return nil, fmt.Errorf("not connected")
	}

	rows, err := m.conn.Query("SHOW DATABASES")
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
		if db != "information_schema" && db != "mysql" && db != "performance_schema" && db != "sys" {
			databases = append(databases, db)
		}
	}

	return databases, rows.Err()
}

// GetTables returns list of tables in a database
func (m *MySQL) GetTables(database string) ([]string, error) {
	if m.conn == nil {
		return nil, fmt.Errorf("not connected")
	}

	query := fmt.Sprintf("SHOW TABLES FROM `%s`", database)
	rows, err := m.conn.Query(query)
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
func (m *MySQL) ExecuteQuery(query string) (*QueryResult, error) {
	if m.conn == nil {
		return nil, fmt.Errorf("not connected")
	}

	result := &QueryResult{}

	// Try as SELECT query first
	rows, err := m.conn.Query(query)
	if err != nil {
		// Try as DML/DDL statement
		res, execErr := m.conn.Exec(query)
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
func (m *MySQL) GetDriverName() string {
	return "MySQL"
}
