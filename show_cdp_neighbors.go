package cisco_database

import (
	"fmt"
	"log"
	"strings"

	"github.com/xtokio/cisco"
)

// CdpNeighbor defines the structure for a single CDP neighbor entry.
type CdpNeighbor struct {
	Neighbor          string
	Interface         string
	HoldTime          string
	Capability        string
	Platform          string
	NeighborInterface string
}

func Show_cdp_neighbors(switch_id int64, switch_hostname string) error {
	show_cdp_neighbors_data, err := cisco.Show_cdp_neighbors(switch_hostname)
	if err != nil {
		return err
	}

	// Check the length of the slice, not the map.
	if len(show_cdp_neighbors_data) == 0 {
		log.Printf("Warning: Parsing completed for %s, but no cdp_neighbors were found.", switch_hostname)
		return nil
	}

	// Establish the database connection.
	db, err := DB_connect()
	if err != nil {
		log.Print(err)
	}
	defer db.Close()

	// Delete records
	deleteQuery := fmt.Sprintf("DELETE FROM cdp_neighbors WHERE switch_id = %d AND DATE(created_at) = CURDATE()", switch_id)
	Execute_query(db, deleteQuery)

	sqlStr := "INSERT INTO `cdp_neighbors` (`switch_id`, `interface`, `neighbor_name`, `neighbor_interface`, `capabilities`, `platform`) VALUES "
	// 2. Create slices for the placeholder strings and the actual values
	var valueStrings []string
	var valueArgs []any // Use []any (or []interface{})

	// This is the placeholder group for a SINGLE row (25 placeholders)
	placeholderRow := "(?, ?, ?, ?, ?, ?)"

	// Iterate over the slice, which maintains the correct order.
	for _, details := range show_cdp_neighbors_data {
		// Add the placeholder group for this row
		valueStrings = append(valueStrings, placeholderRow)
		// Add the values for this row, in the *exact* same order as the columns above
		valueArgs = append(valueArgs,
			switch_id, // The switch_id from the function argument
			details.Interface,
			details.Neighbor,
			details.NeighborInterface,
			details.Capability,
			details.Platform,
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

	log.Printf("%d :: %s :: Show CDP Neighbors :: %d records inserted.\n", switch_id, switch_hostname, len(show_cdp_neighbors_data))

	return nil
}
