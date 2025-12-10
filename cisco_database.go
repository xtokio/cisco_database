package cisco_database

import (
	"log"
	"time"
)

func Truncate_tables() {
	log.Printf("Trunkating tables...")
	Truncate_table("interfaces_status")
	Truncate_table("power_modules")
	Truncate_table("power_interfaces")
	Truncate_table("vlans")
	Truncate_table("cdp_neighbors")
	Truncate_table("lldp_neighbors")
	Truncate_table("mac_address_table")
	Truncate_table("show_running_config")
	Truncate_table("akips_interface_usage")
	Truncate_table("arp_table")
	Truncate_table("vendors")
}

func Update_interfaces() {
	log.Printf("Updating interfaces mac_address")
	Update_interfaces_mac_address()
	log.Println("Waiting 3 seconds before next Update...")
	time.Sleep(3 * time.Second)

	log.Printf("Updating interfaces ip_address")
	Update_interfaces_ip_address()
	log.Println("Waiting 3 seconds before next Update...")
	time.Sleep(3 * time.Second)

	log.Printf("Updating interfaces fqdn")
	Update_interfaces_fqdn()
	log.Println("Waiting 3 seconds before next Update...")
	time.Sleep(3 * time.Second)

	log.Printf("Updating interfaces vlan_id")
	Update_interfaces_vlan_id()
	log.Println("Waiting 3 seconds before next Update...")
	time.Sleep(3 * time.Second)

	log.Printf("Updating interfaces vlan_name")
	Update_interfaces_vlan_name()
}

func Update_interfaces_by_switch_id(switch_id int64) {
	log.Printf("Updating interfaces mac_address")
	Update_interfaces_mac_address_by_switch_id(switch_id)
	time.Sleep(1 * time.Second)

	log.Printf("Updating interfaces ip_address")
	Update_interfaces_ip_address_by_switch_id(switch_id)
	time.Sleep(1 * time.Second)

	log.Printf("Updating interfaces fqdn")
	Update_interfaces_fqdn_by_switch_id(switch_id)
	time.Sleep(1 * time.Second)

	log.Printf("Updating interfaces vlan_id")
	Update_interfaces_vlan_id_by_switch_id(switch_id)
	time.Sleep(1 * time.Second)

	log.Printf("Updating interfaces vlan_name")
	Update_interfaces_vlan_name_by_switch_id(switch_id)
}

func Process_switch(switch_id int64, fqdn string) {
	err := Show_running_config(switch_id, fqdn)
	if err != nil {
		log.Printf("ERROR [Show_running_config] %s: %v", fqdn, err)
		return
	}

	err = Show_version(switch_id, fqdn)
	if err != nil {
		log.Printf("ERROR [Show_version] %s: %v", fqdn, err)
		return
	}
	err = Show_interfaces(switch_id, fqdn)
	if err != nil {
		log.Printf("ERROR [Show_interfaces] %s: %v", fqdn, err)
		return
	}

	err = Show_interfaces_status(switch_id, fqdn)
	if err != nil {
		log.Printf("ERROR [Show_interfaces_status] %s: %v", fqdn, err)
		return
	}

	err = Show_cdp_neighbors(switch_id, fqdn)
	if err != nil {
		log.Printf("ERROR [Show_cdp_neighbors] %s: %v", fqdn, err)
		return
	}

	err = Show_lldp_neighbors(switch_id, fqdn)
	if err != nil {
		log.Printf("ERROR [Show_lldp_neighbors] %s: %v", fqdn, err)
		return
	}

	err = Show_vlan(switch_id, fqdn)
	if err != nil {
		log.Printf("ERROR [Show_vlan] %s: %v", fqdn, err)
		return
	}

	err = Show_power_inline(switch_id, fqdn)
	if err != nil {
		log.Printf("ERROR [Show_power_inline] %s: %v", fqdn, err)
		return
	}

	err = Show_mac_address_table(switch_id, fqdn)
	if err != nil {
		log.Printf("ERROR [Show_mac_address_table] %s: %v", fqdn, err)
		return
	}

	// Akips
	err = Akips_get_interface_usage(switch_id, fqdn)
	if err != nil {
		log.Printf("ERROR AKIPS [Interface_usage] %s: %v", fqdn, err)
		return
	}

}
