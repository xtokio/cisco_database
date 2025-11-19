package cisco_database

import (
	"fmt"
	"log"
	"strings"

	"github.com/xtokio/cisco"
)

// InterfaceStatus defines the structure for a single network interface entry.
type InterfaceStatus struct {
	Interface   string
	Description string
	Status      string
	VlanID      string
	Duplex      string
	Speed       string
	Type        string
}

func Show_interfaces_status(switch_id int64, switch_hostname string) error {
	show_interface_status_data, err := cisco.Show_interfaces_status(switch_hostname)
	if err != nil {
		return err
	}

	// Check the length of the slice, not the map.
	if len(show_interface_status_data) == 0 {
		log.Printf("Show Interface Status :: Warning: Parsing completed for %s, but no interfaces were found.", switch_hostname)
		return nil
	}

	// Establish the database connection.
	db, err := DB_connect()
	if err != nil {
		log.Print(err)
	}
	defer db.Close()

	// Delete records
	deleteQuery := fmt.Sprintf("DELETE FROM interfaces_status WHERE switch_id = %d", switch_id)
	Execute_query(db, deleteQuery)

	// 2. Define the update query.
	// NOTE: We MUST include 'interface_name' in the WHERE clause to update
	// each interface with its specific Vlan.
	insertQuery := "INSERT INTO `interfaces_status` (`switch_id`, `interface`, `description`, `status`, `vlan_id`, `duplex`, `speed`, `type`) VALUES "
	// 2. Create slices for the placeholder strings and the actual values
	var valueStrings []string
	var valueArgs []any // Use []any (or []interface{})

	// This is the placeholder group for a SINGLE row (25 placeholders)
	placeholderRow := "(?, ?, ?, ?, ?, ?, ?, ?)"

	// Iterate over the slice, which maintains the correct order.
	for _, details := range show_interface_status_data {
		// Add the placeholder group for this row
		valueStrings = append(valueStrings, placeholderRow)
		// Add the values for this row, in the *exact* same order as the columns above
		valueArgs = append(valueArgs,
			switch_id, // The switch_id from the function argument
			details.Interface,
			details.Description,
			details.Status,
			details.VlanID,
			details.Duplex,
			details.Speed,
			details.Type,
		)
	}

	finalQuery := insertQuery + strings.Join(valueStrings, ",")

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

	log.Printf("%d :: %s :: Show Interface Status :: %d records updated.\n", switch_id, switch_hostname, len(show_interface_status_data))

	return nil
}
