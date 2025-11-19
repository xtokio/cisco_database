package cisco_database

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/xtokio/cisco"
)

// MacAddressEntry defines the structure for a single entry in the MAC address table.
type MacAddressEntry struct {
	Interface  string
	MacAddress string
	VlanID     string
	Type       string // e.g., DYNAMIC, STATIC, SECURE
}

// Show_mac_address_table constructs the command, runs it, and processes the output.
func Show_mac_address_table(switch_id int64, switch_hostname string) error {
	interfaces := Mac_address_table_interfaces(switch_id)
	if len(interfaces) == 0 {
		log.Printf("%s :: mac address-table records not found in Database", switch_hostname)
		return nil
	}
	interfacesFilter, ok := interfaces[0]["interfaces"].(string)
	if !ok || interfacesFilter == "" {
		// Log and return if the key is missing, not a string, or is an empty string
		log.Printf("%s :: interfaces[0][\"interfaces\"] is missing, nil, or not a string. Cannot proceed.", switch_hostname)
		return nil
	}

	// 1. Construct the Cisco command
	command := fmt.Sprintf("show mac address-table | exclude %s", interfacesFilter)

	outputString, err := cisco.RunCommand(switch_hostname, command)
	if err != nil {
		return err
	}

	// 2. Parse the output
	mac_table_data, err := parseMacAddressTable(outputString)
	if err != nil {
		log.Printf("Error during parsing 'show mac address-table' output for %s: %v", switch_hostname, err)
		return fmt.Errorf("error during parsing 'show mac address-table' output for %s: %v", switch_hostname, err)
	}

	if len(mac_table_data) == 0 {
		log.Printf("Show MAC Address Table :: Warning: Parsing completed for %s, but no MAC entries were found for filter '%s'.", switch_hostname, interfacesFilter)
		return nil
	}

	db, err := DB_connect()
	if err != nil {
		log.Print(err)
	}
	defer db.Close()

	// Start the database transaction
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction for %s: %v", switch_hostname, err)
		return err
	}
	// Defer a rollback. If tx.Commit() is called, this becomes a no-op.
	defer tx.Rollback()

	// Define a safe batch size
	const batchSize = 1000

	// These are the "template" parts of your query
	sqlStr := "INSERT INTO `mac_address_table` (`switch_id`, `interface`, `mac_address`, `vlan_id`, `type`) VALUES "
	placeholderRow := "(?, ?, ?, ?, ?)"

	// Iterate over the mac_table_data in chunks of 'batchSize'
	for i := 0; i < len(mac_table_data); i += batchSize {
		// Determine the end of the current batch
		end := i + batchSize
		if end > len(mac_table_data) {
			end = len(mac_table_data)
		}

		// Slice the full dataset to get the current batch
		batch := mac_table_data[i:end]

		// --- Build the query for THIS BATCH ---
		var valueStrings []string
		var valueArgs []any

		// Iterate over the *batch* (not the full set)
		for _, details := range batch {
			// Add the placeholder group for this row
			valueStrings = append(valueStrings, placeholderRow)

			// Add the values for this row
			valueArgs = append(valueArgs,
				switch_id,
				details.Interface,
				details.MacAddress,
				details.VlanID,
				details.Type,
			)
		}

		// Construct the final query *for this batch*
		finalQuery := sqlStr + strings.Join(valueStrings, ",")

		// Use tx.Exec() with the batch's query and values
		_, err = tx.Exec(finalQuery, valueArgs...)
		if err != nil {
			// If ANY batch fails, roll back the entire transaction
			// This is the ONLY place Rollback should be called.
			// tx.Rollback()
			log.Printf("Failed to execute bulk insert batch for %s: %v", switch_hostname, err)
			log.Printf("Failed query: %s", finalQuery) // Log the query for debugging
			return err
		}
	} // --- End of batch loop ---

	// If all batches were successful, commit the transaction
	// This is the ONLY place Commit should be called.
	err = tx.Commit()
	if err != nil {
		log.Printf("Failed to commit bulk insert transaction for %s: %v", switch_hostname, err)
		return err
	}

	// Log the success message
	log.Printf("%d :: %s :: Show MAC Address Table :: %d records found.\n", switch_id, switch_hostname, len(mac_table_data))

	return nil
}

// parseMacAddressTable takes the raw output and extracts MacAddressEntry structs.
func parseMacAddressTable(rawOutput string) ([]MacAddressEntry, error) {
	var macEntries []MacAddressEntry
	reEntry := regexp.MustCompile(`^\s*\*?\s*(\d+)\s+([\w\.]+)\s+([\w]+)(?:\s+[\w\-])*\s+(\S+)`)

	lines := strings.Split(rawOutput, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip header, separator lines, and summary lines
		if len(line) == 0 ||
			strings.Contains(line, "Mac Address Table") ||
			strings.Contains(line, "Vlan") ||
			strings.Contains(line, "----") ||
			strings.Contains(line, "Total Mac Addresses") ||
			strings.Contains(line, "CPU") { // Often the 'CPU' entries are less relevant for port checks
			continue
		}

		if matches := reEntry.FindStringSubmatch(line); len(matches) == 5 {
			entry := MacAddressEntry{
				// Clean up the VLAN ID in case the '*' was captured with it
				VlanID:     strings.TrimSpace(matches[1]),
				MacAddress: matches[2],
				Type:       matches[3],
				Interface:  matches[4],
			}
			macEntries = append(macEntries, entry)
		}
	}

	return macEntries, nil
}
