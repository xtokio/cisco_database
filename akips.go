package cisco_database

import (
	"fmt"
	"log"
	"strings"

	"github.com/xtokio/akips"
)

type SwitchList struct {
	Switch string
	IP     string
}

func Akips_get_switch_list() error {
	var missing_switch_list []SwitchList

	switch_list := akips.Switch_list()

	switch_records := Device_all()
	// Iterate through each line from the URL list
	for _, current_switch := range switch_list {
		found := false // Flag to track if we find a match

		// Now, check this single line against all your database records
		for _, row := range switch_records {
			if strings.Contains(current_switch.Switch, row["fqdn"].(string)) {
				found = true // We found a match
				break        // Stop checking other records for this line
			}
		}

		// If, after checking all database records, the line was still not found...
		if !found {
			single_switch := SwitchList{}
			single_switch.Switch = current_switch.Switch
			single_switch.IP = current_switch.IP
			missing_switch_list = append(missing_switch_list, single_switch)
		}
	}

	if len(missing_switch_list) > 0 {
		sqlStr := "INSERT INTO `switches` (`fqdn`,`ip_address`) VALUES "
		// 2. Create slices for the placeholder strings and the actual values
		var valueStrings []string
		var valueArgs []any // Use []any (or []interface{})

		// This is the placeholder group for a SINGLE row (25 placeholders)
		placeholderRow := "(?,?)"

		for _, current_switch := range missing_switch_list {
			// Add the placeholder group for this row
			valueStrings = append(valueStrings, placeholderRow)
			// Add the values for this row, in the *exact* same order as the columns above
			valueArgs = append(valueArgs,
				current_switch.Switch,
				current_switch.IP,
			)
		}

		finalQuery := sqlStr + strings.Join(valueStrings, ",")

		// Establish the database connection.
		db, err := DB_connect()
		if err != nil {
			log.Print(err)
		}
		defer db.Close()

		tx, err := db.Begin()
		if err != nil {
			log.Printf("Failed to begin transaction for %s: %v", "Akips switch list", err)
		}

		// Use tx.Exec() with the final query and the flat slice of all values
		_, err = tx.Exec(finalQuery, valueArgs...)
		if err != nil {
			// Roll back the transaction if there's an error
			tx.Rollback()
			log.Printf("Failed to execute bulk insert for %s: %v", "Akips switch list", err)
			log.Printf("Failed query: %s", finalQuery) // Log the query for debugging
		}

		// If successful, commit the transaction
		err = tx.Commit()
		if err != nil {
			log.Printf("Failed to commit bulk insert transaction for %s: %v", "Akips switch list", err)
		}
	}

	fmt.Printf("AKIPS Swicht List :: %d records inserted.\n", len(missing_switch_list))

	return nil

}

func Akips_get_interface_usage(switch_id int64, switch_hostname string) error {
	interface_usage := akips.Interface_usage(switch_hostname)

	// Check the length of the slice, not the map.
	if len(interface_usage) == 0 {
		log.Printf("AKIPS Interface Usage ::Warning: Parsing completed for %s, but no interfaces were found.", switch_hostname)
		return nil
	}

	// Establish the database connection.
	db, err := DB_connect()
	if err != nil {
		log.Print(err)
	}
	defer db.Close()

	// Delete records from today
	deleteQuery := fmt.Sprintf("DELETE FROM akips_interface_usage WHERE switch_id = %d", switch_id)
	Execute_query(db, deleteQuery)

	sqlStr := "INSERT INTO `akips_interface_usage` (" +
		"`switch_id`, `interface`, `status`, `last_change`, `days`, `hours`, `minutes`) VALUES "
	// 2. Create slices for the placeholder strings and the actual values
	var valueStrings []string
	var valueArgs []any // Use []any (or []interface{})

	// This is the placeholder group for a SINGLE row (25 placeholders)
	placeholderRow := "(?, ?, ?, ?, ?, ?, ?)"

	for _, current_interface := range interface_usage {
		// Add the placeholder group for this row
		valueStrings = append(valueStrings, placeholderRow)
		// Add the values for this row, in the *exact* same order as the columns above
		valueArgs = append(valueArgs,
			switch_id, // The switch_id from the function argument
			current_interface.Interface,
			current_interface.Status,
			current_interface.Last_change,
			current_interface.Days,
			current_interface.Hours,
			current_interface.Minutes,
		)
	}

	finalQuery := sqlStr + strings.Join(valueStrings, ",")

	tx, err := db.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction for %s: %v", switch_hostname, err)
		return err
	}

	// Use tx.Exec() with the final query and the flat slice of all values
	_, err = tx.Exec(finalQuery, valueArgs...)
	if err != nil {
		// Roll back the transaction if there's an error
		tx.Rollback()
		log.Printf("Failed to execute bulk insert for %s: %v", switch_hostname, err)
		log.Printf("Failed query: %s", finalQuery) // Log the query for debugging
		return err
	}

	// If successful, commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Printf("Failed to commit bulk insert transaction for %s: %v", switch_hostname, err)
		return err
	}

	log.Printf("%d :: %s :: AKIPS Interface Usage :: %d records inserted.\n", switch_id, switch_hostname, len(interface_usage))

	return nil
}
