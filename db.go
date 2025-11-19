package cisco_database

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// db_connect handles the database connection setup.
// It reads credentials from environment variables and returns a database handle.
func DB_connect() (*sql.DB, error) {
	// --- 1. Get Credentials from Environment Variables ---
	username := os.Getenv("MYSQL_DATABASE_USERNAME")
	password := os.Getenv("MYSQL_DATABASE_PASSWORD")

	if username == "" || password == "" {
		return nil, fmt.Errorf("error: MYSQL_DATABASE_USERNAME and MYSQL_DATABASE_PASSWORD environment variables must be set")
	}

	// --- 2. Construct the Data Source Name (DSN) ---
	dbName := "devices"
	// Use ?parseTime=true to handle DATE/DATETIME types correctly.
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?parseTime=true", username, password, dbName)

	// --- 3. Open and Verify the Connection ---
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error preparing database connection: %w", err)
	}

	err = db.Ping()
	if err != nil {
		db.Close() // Close the connection if ping fails
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	return db, nil
}

// Execute_query runs a SQL statement that does not return rows (e.g., INSERT, UPDATE, DELETE).
// For an INSERT, it returns the last inserted ID.
// For UPDATE or DELETE, it returns the number of rows affected.
func Execute_query(db *sql.DB, query string) (int64, error) {
	// Execute the query
	result, err := db.Exec(query)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %w", err)
	}

	// Check if the query was an INSERT
	if strings.HasPrefix(strings.ToUpper(strings.TrimSpace(query)), "INSERT") {
		// Get the last inserted ID
		id, err := result.LastInsertId()
		if err != nil {
			return 0, fmt.Errorf("failed to get last insert ID: %w", err)
		}
		return id, nil
	}

	// For UPDATE, DELETE, or other commands, get the number of rows affected.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}
	return rowsAffected, nil
}

// Return_query executes a SELECT statement and returns the results dynamically.
// It returns a slice of maps, where each map represents a row (column_name -> value).
func Return_query(db *sql.DB, query string) ([]map[string]interface{}, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	// Get column names from the result set.
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}
	columnCount := len(columns)

	// The final results slice.
	var results []map[string]interface{}

	// Iterate over each row.
	for rows.Next() {
		// Create a slice of empty interfaces to hold the values for scanning.
		values := make([]interface{}, columnCount)
		valuePtrs := make([]interface{}, columnCount)
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		// Scan the row's values into our value pointers.
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Create a map to hold the row data (column_name -> value).
		rowData := make(map[string]interface{})
		for i, colName := range columns {
			val := values[i]

			// Convert []byte to string for better readability, if applicable.
			if b, ok := val.([]byte); ok {
				rowData[colName] = string(b)
			} else {
				rowData[colName] = val
			}
		}
		results = append(results, rowData)
	}

	// Check for any errors that occurred during the iteration.
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %w", err)
	}

	return results, nil
}
