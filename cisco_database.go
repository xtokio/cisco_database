package cisco_database

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/xtokio/cisco_database/database"
)

func check_device(host string) int {
	port := "22"
	address := net.JoinHostPort(host, port)

	timeout := 1 * time.Second
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return 0
	}

	conn.Close()

	return 1
}

func truncate_tables() {
	log.Printf("Trunkating tables...")
	database.Truncate_table("interfaces_status")
	database.Truncate_table("power_modules")
	database.Truncate_table("power_interfaces")
	database.Truncate_table("vlans")
	database.Truncate_table("cdp_neighbors")
	database.Truncate_table("lldp_neighbors")
	database.Truncate_table("mac_address_table")
	database.Truncate_table("show_running_config")
	database.Truncate_table("akips_interface_usage")
	database.Truncate_table("arp_table")
	database.Truncate_table("vendors")
}

func update_interfaces() {
	log.Printf("Updating interfaces mac_address")
	database.Update_interfaces_mac_address()
	log.Println("Waiting 3 seconds before next Update...")
	time.Sleep(3 * time.Second)

	log.Printf("Updating interfaces ip_address")
	database.Update_interfaces_ip_address()
	log.Println("Waiting 3 seconds before next Update...")
	time.Sleep(3 * time.Second)

	log.Printf("Updating interfaces vlan_id")
	database.Update_interfaces_vlan_id()
	log.Println("Waiting 3 seconds before next Update...")
	time.Sleep(3 * time.Second)

	log.Printf("Updating interfaces vlan_name")
	database.Update_interfaces_vlan_name()
}

func Update_interfaces_by_switch_id(switch_id int64) {
	log.Printf("Updating interfaces mac_address")
	database.Update_interfaces_mac_address_by_switch_id(switch_id)
	time.Sleep(1 * time.Second)

	log.Printf("Updating interfaces ip_address")
	database.Update_interfaces_ip_address_by_switch_id(switch_id)
	time.Sleep(1 * time.Second)

	log.Printf("Updating interfaces vlan_id")
	database.Update_interfaces_vlan_id_by_switch_id(switch_id)
	time.Sleep(1 * time.Second)

	log.Printf("Updating interfaces vlan_name")
	database.Update_interfaces_vlan_name_by_switch_id(switch_id)
}

func process_switch(switch_id int64, fqdn string) {
	err := database.Show_running_config(switch_id, fqdn)
	if err != nil {
		log.Printf("ERROR [Show_running_config] %s: %v", fqdn, err)
		return
	}

	err = database.Show_version(switch_id, fqdn)
	if err != nil {
		log.Printf("ERROR [Show_version] %s: %v", fqdn, err)
		return
	}
	err = database.Show_interfaces(switch_id, fqdn)
	if err != nil {
		log.Printf("ERROR [Show_interfaces] %s: %v", fqdn, err)
		return
	}

	err = database.Show_interfaces_status(switch_id, fqdn)
	if err != nil {
		log.Printf("ERROR [Show_interfaces_status] %s: %v", fqdn, err)
		return
	}

	err = database.Show_cdp_neighbors(switch_id, fqdn)
	if err != nil {
		log.Printf("ERROR [Show_cdp_neighbors] %s: %v", fqdn, err)
		return
	}

	err = database.Show_lldp_neighbors(switch_id, fqdn)
	if err != nil {
		log.Printf("ERROR [Show_lldp_neighbors] %s: %v", fqdn, err)
		return
	}

	err = database.Show_vlan(switch_id, fqdn)
	if err != nil {
		log.Printf("ERROR [Show_vlan] %s: %v", fqdn, err)
		return
	}

	err = database.Show_power_inline(switch_id, fqdn)
	if err != nil {
		log.Printf("ERROR [Show_power_inline] %s: %v", fqdn, err)
		return
	}

	err = database.Show_mac_address_table(switch_id, fqdn)
	if err != nil {
		log.Printf("ERROR [Show_mac_address_table] %s: %v", fqdn, err)
		return
	}

	// Akips
	err = database.Akips_get_interface_usage(switch_id, fqdn)
	if err != nil {
		log.Printf("ERROR AKIPS [Interface_usage] %s: %v", fqdn, err)
		return
	}

}

func Cisco_database_all() {
	// Record the start time.
	start := time.Now()

	// Check if there are new switches to add
	database.Akips_get_switch_list()

	// Truncate tables
	truncate_tables()

	// Establish the database connection.
	db, db_err := database.DB_connect()
	if db_err != nil {
		log.Print(db_err)
	}
	defer db.Close()

	switch_records := database.Device_all()
	total_records := len(switch_records)

	// --- Concurrency and Batching Configuration ---
	const maxConcurrency = 50
	const batchSize = 100 // Your desired batch size

	// Create a buffered channel to act as a semaphore, limiting concurrency.
	// This pool is shared across all batches.
	pool := make(chan struct{}, maxConcurrency)

	log.Printf("Starting to process %d total devices in batches of %d...", total_records, batchSize)

	// Loop over the switch_records in batches
	for i := 0; i < total_records; i += batchSize {
		// Calculate the end index for the current batch
		end := i + batchSize
		if end > total_records {
			end = total_records
		}

		// Get the slice for the current batch
		batch := switch_records[i:end]

		log.Printf("Processing batch %d/%d (records %d to %d)...", (i/batchSize)+1, (total_records+batchSize-1)/batchSize, i, end-1)

		// Use a WaitGroup to wait for *this batch* to finish.
		var wg sync.WaitGroup

		// Process all items in the current batch
		for _, row := range batch {
			// Extract variables inside the loop
			switch_id := row["id"].(int64)
			fqdn := row["fqdn"].(string)

			reachable := check_device(fqdn)

			// Update `reachable` from switches
			updateQuery := fmt.Sprintf("UPDATE switches SET reachable = %d WHERE id = %d", reachable, switch_id)
			database.Execute_query(db, updateQuery)

			if reachable == 1 {
				wg.Add(1)          // Increment the batch WaitGroup counter.
				pool <- struct{}{} // Acquire a spot in the pool.

				// Launch a goroutine to process the switch.
				// **IMPORTANT:** Pass switch_id and fqdn as arguments
				// to avoid the loop variable capture bug.
				go func(id int64, host string) {
					defer wg.Done()           // Decrement the counter when the goroutine completes.
					defer func() { <-pool }() // Release the spot in the pool.
					process_switch(id, host)
				}(switch_id, fqdn) // Pass the current values to the goroutine
			}
		}

		// Wait for all goroutines *in this batch* to finish.
		wg.Wait()
		log.Printf("Batch %d/%d finished.", (i/batchSize)+1, (total_records+batchSize-1)/batchSize)

		// If this is not the last batch, wait 3 seconds
		if end < total_records {
			log.Println("Waiting 3 seconds before next batch...")
			time.Sleep(3 * time.Second)
		}
	}

	// Update Interfaces
	update_interfaces()

	// All batches are done.
	duration := time.Since(start)
	message := "Devices: " + fmt.Sprint(total_records) + " | Process duration: " + duration.String() + "\n"
	fmt.Print("\n" + message + "\n")
}

func Cisco_database_by_id(switch_id int64) {
	// Record the start time.
	start := time.Now()

	switch_records := database.Device_by_id(strconv.FormatInt(switch_id, 10))
	total_records := len(switch_records)

	// --- Concurrency and Batching Configuration ---
	const maxConcurrency = 1

	// Create a buffered channel to act as a semaphore, limiting concurrency.
	// This pool is shared across all batches.
	pool := make(chan struct{}, maxConcurrency)

	// Use a WaitGroup to wait for *this batch* to finish.
	var wg sync.WaitGroup

	// Process all items in the current batch
	for _, row := range switch_records {
		// Extract variables inside the loop
		switch_id := row["id"].(int64)
		fqdn := row["fqdn"].(string)

		reachable := check_device(fqdn)

		if reachable == 1 {
			wg.Add(1)          // Increment the batch WaitGroup counter.
			pool <- struct{}{} // Acquire a spot in the pool.

			go func(id int64, host string) {
				defer wg.Done()           // Decrement the counter when the goroutine completes.
				defer func() { <-pool }() // Release the spot in the pool.
				process_switch(id, host)
			}(switch_id, fqdn) // Pass the current values to the goroutine
		}
	}

	// Wait for all goroutines *in this batch* to finish.
	wg.Wait()

	// Update Interfaces
	Update_interfaces_by_switch_id(switch_id)

	// All batches are done.
	duration := time.Since(start)
	message := "Devices: " + fmt.Sprint(total_records) + " | Process duration: " + duration.String() + "\n"
	fmt.Print("\n" + message + "\n")

}
