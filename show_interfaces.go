package cisco_database

import (
	"fmt"
	"log"
	"strings"

	"github.com/xtokio/cisco"
)

// Show_interfaces connects to a switch, gets interface data, and returns it as a map.
func Show_interfaces(switch_id int64, switch_hostname string) error {
	show_interface_data, err := cisco.Show_interfaces(switch_hostname)
	if err != nil {
		return err
	}

	// Check the length of the slice, not the map.
	if len(show_interface_data) == 0 {
		log.Printf("Show Interfaces ::Warning: Parsing completed for %s, but no interfaces were found.", switch_hostname)
		return nil
	}

	process_show_interfaces(show_interface_data, switch_id, switch_hostname)

	// Return the data
	return nil
}

func process_show_interfaces(interfacesSlice []cisco.InterfaceDetails, switch_id int64, switch_hostname string) {
	// Establish the database connection.
	db, err := DB_connect()
	if err != nil {
		log.Print(err)
	}
	defer db.Close()

	// Delete records from today
	deleteQuery := fmt.Sprintf("DELETE FROM interfaces WHERE switch_id = %d AND DATE(created_at) = CURDATE()", switch_id)
	Execute_query(db, deleteQuery)

	sqlStr := "INSERT INTO `interfaces` (" +
		"`switch_id`, `interface`, `description`, `ip_address`, `link_status`, `protocol_status`, `hardware_type`, `reliability`, `txload`, `rxload`, `mtu`, `duplex`, `speed`, `media_type`, `bandwidth`, `delay`, `encapsulation`, `last_input`, `last_output`, `last_output_hang`, `queue_strategy`, `input_rate`, `output_rate`, `input_packets`, `output_packets`, `runts`, `giants`, `throttles`, `input_errors`, `output_errors`, `crc_errors`, `collisions`) VALUES "
	// 2. Create slices for the placeholder strings and the actual values
	var valueStrings []string
	var valueArgs []any // Use []any (or []interface{})

	// This is the placeholder group for a SINGLE row (25 placeholders)
	placeholderRow := "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	// Iterate over the slice, which maintains the correct order.
	for _, details := range interfacesSlice {
		// Add the placeholder group for this row
		valueStrings = append(valueStrings, placeholderRow)
		// Add the values for this row, in the *exact* same order as the columns above
		valueArgs = append(valueArgs,
			switch_id, // The switch_id from the function argument
			details.Interface,
			details.Description,
			details.IPAddress,
			details.LinkStatus,
			details.ProtocolStatus,
			details.Hardware,
			details.Reliability,
			details.TxLoad,
			details.RxLoad,
			details.Mtu,
			details.Duplex,
			details.Speed,
			details.MediaType,
			details.Bandwidth,
			details.Delay,
			details.Encapsulation,
			details.LastInput,
			details.LastOutput,
			details.OutputHang,
			details.QueueStrategy,
			details.InputRateBps,
			details.OutputRateBps,
			details.PacketsInput,
			details.PacketsOutput,
			details.Runts,
			details.Giants,
			details.Throttles,
			details.InputErrors,
			details.OutputErrors,
			details.CrcErrors,
			details.Collisions,
		)
	}
	finalQuery := sqlStr + strings.Join(valueStrings, ",")

	tx, err := db.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction for %s: %v", switch_hostname, err)
		return
	}

	// Use tx.Exec() with the final query and the flat slice of all values
	_, err = tx.Exec(finalQuery, valueArgs...)
	if err != nil {
		// Roll back the transaction if there's an error
		tx.Rollback()
		log.Printf("Failed to execute bulk insert for %s: %v", switch_hostname, err)
		log.Printf("Failed query: %s", finalQuery) // Log the query for debugging
		return
	}

	// If successful, commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Printf("Failed to commit bulk insert transaction for %s: %v", switch_hostname, err)
		return
	}

	log.Printf("%d :: %s :: Show Interfaces :: %d records inserted.\n", switch_id, switch_hostname, len(interfacesSlice))

}
