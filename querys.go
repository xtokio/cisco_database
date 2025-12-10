package cisco_database

import (
	"log"
	"strconv"
)

func Device_all() []map[string]interface{} {
	// Establish the database connection.
	db, err := DB_connect()
	if err != nil {
		log.Print(err)
	}
	defer db.Close()

	rows, err := Return_query(db, "SELECT * from switches") // 'malloy-spine-1.hub.nd.edu','drni-n7010-itc.hub.nd.edu','lgomezreswitch.hub.nd.edu'
	if err != nil {
		log.Printf("Error reading data: %v", err)
	}

	return rows
}

func Device_reachable() []map[string]interface{} {
	// Establish the database connection.
	db, err := DB_connect()
	if err != nil {
		log.Print(err)
	}
	defer db.Close()

	rows, err := Return_query(db, "SELECT count(*) as count from switches WHERE reachable = 1")
	if err != nil {
		log.Printf("Error reading data: %v", err)
	}

	return rows
}

func Device_unreachable() []map[string]interface{} {
	// Establish the database connection.
	db, err := DB_connect()
	if err != nil {
		log.Print(err)
	}
	defer db.Close()

	rows, err := Return_query(db, "SELECT count(*) as count from switches WHERE reachable = 0")
	if err != nil {
		log.Printf("Error reading data: %v", err)
	}

	return rows
}

func Device_by_id(switch_id string) []map[string]interface{} {
	// Establish the database connection.
	db, err := DB_connect()
	if err != nil {
		log.Print(err)
	}
	defer db.Close()

	rows, err := Return_query(db, "SELECT * from switches WHERE id = "+switch_id)
	if err != nil {
		log.Printf("Error reading data: %v", err)
	}

	return rows
}

func Interfaces_by_switch_id(switch_id string) []map[string]interface{} {
	// Establish the database connection.
	db, err := DB_connect()
	if err != nil {
		log.Print(err)
	}
	defer db.Close()

	rows, err := Return_query(db, "SELECT * from view_interfaces WHERE switch_id = "+switch_id)
	if err != nil {
		log.Printf("Error reading data: %v", err)
	}

	return rows
}

func Vlan_names(switch_id int64) []map[string]interface{} {
	// Establish the database connection.
	db, err := DB_connect()
	if err != nil {
		log.Print(err)
	}
	defer db.Close()

	rows, err := Return_query(db, "SELECT i.switch_id,i.interface,v.vlan_id,v.vlan_name FROM interfaces AS i JOIN vlans AS v ON i.switch_id = v.switch_id AND i.vlan_id = v.vlan_id WHERE i.switch_id = "+strconv.FormatInt(switch_id, 10)+" AND FIND_IN_SET(i.interface, v.interfaces) > 0")
	if err != nil {
		log.Printf("Error reading data: %v", err)
	}

	return rows
}

func Mac_address_table_interfaces(switch_id int64) []map[string]interface{} {
	// Establish the database connection.
	db, err := DB_connect()
	if err != nil {
		log.Print(err)
	}
	defer db.Close()

	rows, err := Return_query(db, "SELECT CONCAT(GROUP_CONCAT(CONCAT('_', interface, '_') SEPARATOR '|'),'|_CPU_') as interfaces from interfaces_status where switch_id = "+strconv.FormatInt(switch_id, 10)+" and status = 'connected' and vlan_id like '%trunk%' and date(created_at) = CURDATE()")
	if err != nil {
		log.Printf("Error reading data: %v", err)
	}

	return rows
}

func Truncate_table(table_name string) {
	// Establish the database connection.
	db, err := DB_connect()
	if err != nil {
		log.Print(err)
	}
	defer db.Close()

	Execute_query(db, "TRUNCATE TABLE "+table_name)
	if err != nil {
		log.Printf("%s :: Error truncating table: %v", table_name, err)
	}

}

func Update_interfaces_mac_address() {
	// Establish the database connection.
	db, err := DB_connect()
	if err != nil {
		log.Print(err)
	}
	defer db.Close()

	Execute_query(db, "UPDATE interfaces i JOIN (SELECT switch_id,interface,DATE(created_at) AS date_created,MIN(mac_address) AS mac_address FROM mac_address_table GROUP BY switch_id,interface,DATE(created_at)) m ON i.switch_id = m.switch_id AND i.interface = m.interface AND DATE(i.created_at) = m.date_created SET i.mac_address = m.mac_address WHERE DATE(i.created_at) = CURDATE()")
	if err != nil {
		log.Printf("%s :: Error updating interfaces mac_address: %v", "Interfaces mac_address", err)
	}
}
func Update_interfaces_mac_address_by_switch_id(switch_id int64) {
	// Establish the database connection.
	db, err := DB_connect()
	if err != nil {
		log.Print(err)
	}
	defer db.Close()

	Execute_query(db, "UPDATE interfaces i JOIN (SELECT switch_id,interface,DATE(created_at) AS date_created,MIN(mac_address) AS mac_address FROM mac_address_table GROUP BY switch_id,interface,DATE(created_at)) m ON i.switch_id = m.switch_id AND i.interface = m.interface AND DATE(i.created_at) = m.date_created SET i.mac_address = m.mac_address WHERE i.switch_id = "+strconv.FormatInt(switch_id, 10)+" AND DATE(i.created_at) = CURDATE()")
	if err != nil {
		log.Printf("%s :: Error updating interfaces mac_address: %v", "Interfaces mac_address", err)
	}
}

func Update_interfaces_ip_address() {
	// Establish the database connection.
	db, err := DB_connect()
	if err != nil {
		log.Print(err)
	}
	defer db.Close()

	Execute_query(db, "UPDATE interfaces JOIN arp_table ON interfaces.mac_address = arp_table.mac_address SET interfaces.ip_address = arp_table.ip_address WHERE DATE(interfaces.created_at) = CURDATE()")
	if err != nil {
		log.Printf("%s :: Error updating interfaces ip_address: %v", "Interfaces ip_address", err)
	}
}

func Update_interfaces_ip_address_by_switch_id(switch_id int64) {
	// Establish the database connection.
	db, err := DB_connect()
	if err != nil {
		log.Print(err)
	}
	defer db.Close()

	Execute_query(db, "UPDATE interfaces JOIN arp_table ON interfaces.mac_address = arp_table.mac_address SET interfaces.ip_address = arp_table.ip_address WHERE interfaces.switch_id = "+strconv.FormatInt(switch_id, 10)+" AND DATE(interfaces.created_at) = CURDATE()")
	if err != nil {
		log.Printf("%s :: Error updating interfaces ip_address: %v", "Interfaces ip_address", err)
	}
}

func Update_interfaces_fqdn() {
	// Establish the database connection.
	db, err := DB_connect()
	if err != nil {
		log.Print(err)
	}
	defer db.Close()

	Execute_query(db, "UPDATE interfaces JOIN fqdn_table ON interfaces.mac_address = fqdn_table.mac_address AND interfaces.ip_address = fqdn_table.ip_address SET interfaces.fqdn = fqdn_table.fqdn WHERE DATE(interfaces.created_at) = CURDATE()")
	if err != nil {
		log.Printf("%s :: Error updating interfaces fqdn: %v", "Interfaces fqdn", err)
	}
}

func Update_interfaces_fqdn_by_switch_id(switch_id int64) {
	// Establish the database connection.
	db, err := DB_connect()
	if err != nil {
		log.Print(err)
	}
	defer db.Close()

	Execute_query(db, "UPDATE interfaces JOIN fqdn_table ON interfaces.mac_address = fqdn_table.mac_address AND interfaces.ip_address = fqdn_table.ip_address SET interfaces.fqdn = fqdn_table.fqdn WHERE interfaces.switch_id = "+strconv.FormatInt(switch_id, 10)+" AND DATE(interfaces.created_at) = CURDATE()")
	if err != nil {
		log.Printf("%s :: Error updating interfaces fqdn: %v", "Interfaces fqdn", err)
	}
}

func Update_interfaces_vlan_id() {
	// Establish the database connection.
	db, err := DB_connect()
	if err != nil {
		log.Print(err)
	}
	defer db.Close()

	Execute_query(db, "update interfaces join interfaces_status on interfaces.switch_id = interfaces_status.switch_id and interfaces.interface = interfaces_status.interface and date(interfaces.created_at) = date(interfaces_status.created_at) set interfaces.vlan_id = interfaces_status.vlan_id where date(interfaces.created_at) = curdate()")
	if err != nil {
		log.Printf("%s :: Error updating interfaces vlan_id: %v", "Interfaces vlan_id", err)
	}
}
func Update_interfaces_vlan_id_by_switch_id(switch_id int64) {
	// Establish the database connection.
	db, err := DB_connect()
	if err != nil {
		log.Print(err)
	}
	defer db.Close()

	Execute_query(db, "update interfaces join interfaces_status on interfaces.switch_id = interfaces_status.switch_id and interfaces.interface = interfaces_status.interface and date(interfaces.created_at) = date(interfaces_status.created_at) set interfaces.vlan_id = interfaces_status.vlan_id where interfaces.switch_id = "+strconv.FormatInt(switch_id, 10)+" AND date(interfaces.created_at) = curdate()")
	if err != nil {
		log.Printf("%s :: Error updating interfaces vlan_id: %v", "Interfaces vlan_id", err)
	}
}

func Update_interfaces_vlan_name() {
	// Establish the database connection.
	db, err := DB_connect()
	if err != nil {
		log.Print(err)
	}
	defer db.Close()

	Execute_query(db, "update interfaces join vlans on interfaces.switch_id = vlans.switch_id and interfaces.vlan_id = vlans.vlan_id and date(interfaces.created_at) = date(vlans.created_at) set interfaces.vlan_name = vlans.vlan_name where date(interfaces.created_at) = curdate()")
	if err != nil {
		log.Printf("%s :: Error updating interfaces vlan_id: %v", "Interfaces vlan_id", err)
	}
}
func Update_interfaces_vlan_name_by_switch_id(switch_id int64) {
	// Establish the database connection.
	db, err := DB_connect()
	if err != nil {
		log.Print(err)
	}
	defer db.Close()

	Execute_query(db, "update interfaces join vlans on interfaces.switch_id = vlans.switch_id and interfaces.vlan_id = vlans.vlan_id and date(interfaces.created_at) = date(vlans.created_at) set interfaces.vlan_name = vlans.vlan_name where interfaces.switch_id = "+strconv.FormatInt(switch_id, 10)+" AND date(interfaces.created_at) = curdate()")
	if err != nil {
		log.Printf("%s :: Error updating interfaces vlan_id: %v", "Interfaces vlan_id", err)
	}
}
