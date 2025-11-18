package database

import (
	"fmt"
	"log"

	"github.com/xtokio/cisco"
)

// VersionInfo defines the structure for the parsed "show version" output.
// It's used as an intermediate struct within the parsing function.
type VersionInfo struct {
	Hardware      string
	Version       string
	Release       string
	SoftwareImage string
	SerialNumber  string
	Uptime        string
	Restarted     string
	ReloadReason  string
	Rommon        string
}

// Show_version connects to a switch, runs "show version", and returns the parsed data as a map.
func Show_version(switch_id int64, switch_hostname string) error {
	show_version_data, err := cisco.Show_version(switch_hostname)
	if err != nil {
		return err
	}

	process_show_version(show_version_data, switch_id, switch_hostname)

	return nil
}

func process_show_version(versionData map[string]string, switch_id int64, switch_hostname string) {
	// Establish the database connection.
	db, err := DB_connect()
	if err != nil {
		log.Print(err)
	}
	defer db.Close()
	updateQuery := fmt.Sprintf("UPDATE switches SET `hardware` = '%s', `version` = '%s', `release` = '%s', `software_image` = '%s', `serial` = '%s', `uptime` = '%s', `restarted` = '%s', `reload_reason` = '%s', `rommon` = '%s' WHERE id = %d",
		versionData["Hardware"],
		versionData["Version"],
		versionData["Release"],
		versionData["SoftwareImage"],
		versionData["SerialNumber"],
		versionData["Uptime"],
		versionData["Restarted"],
		versionData["ReloadReason"],
		versionData["Rommon"],
		switch_id,
	)

	affectedRows, err := Execute_query(db, updateQuery)
	if err != nil {
		log.Printf("Update query failed: %v", err)
	} else {
		log.Printf("%d :: %s :: Show Version :: UPDATE successful. Rows affected: %d\n", switch_id, switch_hostname, affectedRows)
	}
}
