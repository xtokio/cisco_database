package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/xtokio/cisco"
)

// PowerModuleInfo defines the structure for a power supply module.
type PowerModuleInfo struct {
	Module    string
	Available string
	Used      string
	Remaining string
}

// PowerInterfaceInfo defines the structure for a single PoE interface.
type PowerInterfaceInfo struct {
	Interface string
	Admin     string
	Oper      string
	Power     string // (Watts)
	Device    string
	Class     string
	Max       string // (Watts)
}

// Show_power_inline fetches and processes "show power inline" output.
func Show_power_inline(switch_id int64, switch_hostname string) error {
	show_power_inline_modules_data, show_power_inline_interfaces_data, err := cisco.Show_power_inline(switch_hostname)
	if err != nil {
		return err
	}

	// --- DATABASE OPERATIONS ---
	db, err := DB_connect()
	if err != nil {
		log.Print(err)
	}
	defer db.Close()

	// Process Modules
	if len(show_power_inline_modules_data) > 0 {
		processPowerModules(db, switch_id, switch_hostname, show_power_inline_modules_data)
	} else {
		log.Printf("Warning: No power modules found for %s.", switch_hostname)
	}

	// Process Interfaces
	if len(show_power_inline_interfaces_data) > 0 {
		processPowerInterfaces(db, switch_id, switch_hostname, show_power_inline_interfaces_data)
	} else {
		log.Printf("Warning: No power interfaces found for %s.", switch_hostname)
	}

	return nil
}

// processPowerModules handles the bulk insert for power modules.
func processPowerModules(db *sql.DB, switch_id int64, switch_hostname string, modules []cisco.PowerModuleInfo) {
	deleteQuery := fmt.Sprintf("DELETE FROM power_modules WHERE switch_id = %d", switch_id)
	Execute_query(db, deleteQuery)

	sqlStr := "INSERT INTO `power_modules` (`switch_id`, `module`, `available`, `used`, `remaining`) VALUES "
	var valueStrings []string
	var valueArgs []any
	placeholderRow := "(?, ?, ?, ?, ?)"

	for _, mod := range modules {
		valueStrings = append(valueStrings, placeholderRow)
		valueArgs = append(valueArgs,
			switch_id,
			mod.Module,
			mod.Available,
			mod.Used,
			mod.Remaining,
		)
	}

	finalQuery := sqlStr + strings.Join(valueStrings, ",")
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction for %s (modules): %v", switch_hostname, err)
		return
	}

	_, err = tx.Exec(finalQuery, valueArgs...)
	if err != nil {
		tx.Rollback()
		log.Printf("Failed to execute bulk insert for %s (modules): %v", switch_hostname, err)
		return
	}
	tx.Commit()

	log.Printf("%d :: %s :: Show Power Modules :: %d records inserted.\n", switch_id, switch_hostname, len(modules))
}

// processPowerInterfaces handles the bulk insert for power interfaces.
func processPowerInterfaces(db *sql.DB, switch_id int64, switch_hostname string, interfaces []cisco.PowerInterfaceInfo) {
	deleteQuery := fmt.Sprintf("DELETE FROM power_interfaces WHERE switch_id = %d", switch_id)
	Execute_query(db, deleteQuery)

	// Note: Column names like 'interface' and 'class' might be reserved keywords; use backticks.
	sqlStr := "INSERT INTO `power_interfaces` (`switch_id`, `interface`, `admin`, `oper`, `power`, `device`, `class`, `max`) VALUES "
	var valueStrings []string
	var valueArgs []any
	placeholderRow := "(?, ?, ?, ?, ?, ?, ?, ?)"

	for _, iface := range interfaces {
		valueStrings = append(valueStrings, placeholderRow)
		valueArgs = append(valueArgs,
			switch_id,
			iface.Interface,
			iface.Admin,
			iface.Oper,
			iface.Power,
			iface.Device,
			iface.Class,
			iface.Max,
		)
	}

	finalQuery := sqlStr + strings.Join(valueStrings, ",")
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction for %s (interfaces): %v", switch_hostname, err)
		return
	}

	_, err = tx.Exec(finalQuery, valueArgs...)
	if err != nil {
		tx.Rollback()
		log.Printf("Failed to execute bulk insert for %s (interfaces): %v", switch_hostname, err)
		return
	}
	tx.Commit()

	log.Printf("%d :: %s :: Show Power Interfaces :: %d records inserted.\n", switch_id, switch_hostname, len(interfaces))
}
