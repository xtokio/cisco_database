CREATE TABLE IF NOT EXISTS `switches` (
  `id` INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
  `fqdn` TEXT NULL,
  `hostname` TEXT NULL,
  `ip_address` TEXT NULL,
  `mac_address` TEXT NULL,
  `location` TEXT NULL,
  `software_image` TEXT NULL,
  `version` TEXT NULL,
  `release` TEXT NULL,
  `rommon` TEXT NULL,
  `uptime` TEXT NULL,
  `reload_reason` TEXT NULL,
  `hardware` TEXT NULL,
  `serial` TEXT NULL,
  `restarted` TEXT NULL,
  `reachable` INT NULL,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `interfaces` (
  `id` INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
  `switch_id` INT NOT NULL,
  `interface` TEXT NULL,
  `description` TEXT NULL,
  `mac_address` TEXT NULL,
  `ip_address` TEXT NULL,
  `fqdn` TEXT NULL,
  `link_status` TEXT NULL,
  `protocol_status` TEXT NULL,
  `ise_port` INT DEFAULT 0,
  `vendor` TEXT NULL,
  `hardware_type` TEXT NULL,
  `capabilities` TEXT NULL,
  `power_module` INT NULL,
  `power` TEXT NULL,
  `power_device` TEXT NULL,
  `mtu` TEXT NULL,
  `reliability` TEXT NULL,
  `txload` TEXT NULL,
  `rxload` TEXT NULL,
  `duplex` TEXT NULL,
  `speed` TEXT NULL,
  `media_type` TEXT NULL,
  `bandwidth` TEXT NULL,
  `delay` TEXT NULL,
  `encapsulation` TEXT NULL,
  `last_input` TEXT NULL,
  `last_output` TEXT NULL,
  `last_output_hang` TEXT NULL,
  `queue_strategy` TEXT NULL,
  `input_rate` TEXT NULL,
  `output_rate` TEXT NULL,
  `input_packets` TEXT NULL,
  `output_packets` TEXT NULL,
  `runts` TEXT NULL,
  `giants` TEXT NULL,
  `throttles` TEXT NULL,
  `input_errors` TEXT NULL,
  `output_errors` TEXT NULL,
  `crc_errors` TEXT NULL,
  `collisions` TEXT NULL,
  `vlan_id` TEXT NULL,
  `vlan_name` TEXT NULL,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `interfaces_status` (
  `id` INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
  `switch_id` INT NOT NULL,
  `interface` TEXT NULL,
  `description` TEXT NULL,
  `status` TEXT NULL,
  `vlan_id` TEXT NULL,
  `duplex` TEXT NULL,
  `speed` TEXT NULL,
  `type` TEXT NULL,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `power_modules` (
  `id` INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
  `switch_id` INT NOT NULL,
  `module` INT NULL,
  `available` TEXT NULL,
  `used` TEXT NULL,
  `remaining` TEXT NULL,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `power_interfaces` (
  `id` INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
  `switch_id` INT NOT NULL,
  `interface` TEXT NULL,
  `admin` TEXT NULL,
  `oper` TEXT NULL,
  `power` TEXT NULL,
  `device` TEXT NULL,
  `class` TEXT NULL,
  `max` TEXT NULL,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `vlans` (
  `id` INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
  `switch_id` INT NOT NULL,
  `vlan_id` TEXT NULL,
  `vlan_name` TEXT NULL,
  `status` TEXT NULL,
  `interfaces` TEXT NULL,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `cdp_neighbors` (
  `id` INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
  `switch_id` INT NOT NULL,
  `interface` TEXT NULL,
  `neighbor_name` TEXT NULL,
  `neighbor_interface` TEXT NULL,
  `capabilities` TEXT NULL,
  `platform` TEXT NULL,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `lldp_neighbors` (
  `id` INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
  `switch_id` INT NOT NULL,
  `interface` TEXT NULL,
  `neighbor_name` TEXT NULL,
  `neighbor_interface` TEXT NULL,
  `capabilities` TEXT NULL,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `mac_address_table` (
  `id` INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
  `switch_id` INT NOT NULL,
  `interface` TEXT NULL,
  `mac_address` TEXT NULL,
  `vlan_id` TEXT NULL,
  `type` TEXT NULL,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `show_running_config` (
  `id` INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
  `switch_id` INT NOT NULL,
  `interface` TEXT NULL,
  `configuration` TEXT NULL,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `akips_interface_usage` (
  `id` INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
  `switch_id` INT NOT NULL,
  `interface` TEXT NULL,
  `status` TEXT NULL,
  `last_change` TEXT NULL,
  `days` TEXT NULL,
  `hours` TEXT NULL,
  `minutes` TEXT NULL,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `vendors` (
  `id` INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
  `mac_address` TEXT NOT NULL,
  `vendor` TEXT NULL,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `arp_table` (
  `id` INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
  `source` TEXT NOT NULL,
  `mac_address` TEXT NOT NULL,
  `ip_address` TEXT NOT NULL,
  `protocol` TEXT NULL,
  `age` TEXT NULL,
  `type` TEXT NULL,
  `interface` TEXT NULL,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `fqdn_table` (
  `id` INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
  `mac_address` TEXT NOT NULL,
  `ip_address` TEXT NOT NULL,
  `fqdn` TEXT NULL,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `ise_ip_phones` (
  `id` INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
  `mac_address` TEXT NOT NULL,
  `profile` TEXT NOT NULL,
  `switch_ip_address` TEXT NULL,
  `interface` TEXT NULL,
  `last_change` TEXT NULL,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

ALTER TABLE `mac_address_table` ADD INDEX `idx_mac_date` (mac_address(20), created_at);
ALTER TABLE `interfaces` ADD INDEX `idx_sw_if_date` (switch_id, interface(10), created_at);
ALTER TABLE `interfaces_status` ADD INDEX `idx_sw_if_date` (switch_id, interface(10), created_at);
ALTER TABLE `lldp_neighbors` ADD INDEX `idx_sw_if_date` (switch_id, interface(10), created_at);
ALTER TABLE `cdp_neighbors` ADD INDEX `idx_sw_if_date` (switch_id, interface(10), created_at);
ALTER TABLE `power_interfaces` ADD INDEX `idx_sw_if_date` (switch_id, interface(10), created_at);
ALTER TABLE `power_modules` ADD INDEX `idx_sw_date` (switch_id, created_at);
ALTER TABLE `vlans` ADD INDEX `idx_sw_date` (switch_id, created_at);
ALTER TABLE `show_running_config` ADD INDEX `idx_sw_if_date` (switch_id, interface(10), created_at);
ALTER TABLE `akips_interface_usage` ADD INDEX `idx_sw_if_date` (switch_id, interface(10), created_at);
ALTER TABLE `arp_table` ADD INDEX `idx_mac_date` (mac_address(20), created_at);
ALTER TABLE `fqdn_table` ADD INDEX `idx_ip_date` (ip_address(20), created_at);
ALTER TABLE `ise_ip_phones` ADD INDEX `idx_mac_date` (mac_address(20), created_at);

CREATE OR REPLACE VIEW `view_interfaces` AS
SELECT
  switches.id as switch_id,
  switches.ip_address as switch_ip_address,
	switches.fqdn,
	interfaces.id as interface_id,
	interfaces.interface,
	interfaces.mac_address,
	interfaces.ip_address,
	vendors.vendor,
	interfaces.description,
	interfaces.link_status as status,
	CASE
	  WHEN ise_check.interface IS NOT NULL THEN 'ISE' ELSE NULL
	END AS ise,
	interfaces.vlan_id,
	interfaces.vlan_name,
	akips_interface_usage.last_change,
	interfaces.created_at
FROM interfaces
JOIN switches ON switches.id = interfaces.switch_id
LEFT JOIN (
    SELECT
      DISTINCT switch_id,
      interface
    FROM
      show_running_config
    WHERE
      `configuration` LIKE '%authentication priority dot1x mab%'
      AND date(created_at) = CURDATE()
  ) AS ise_check ON interfaces.switch_id = ise_check.switch_id
  AND interfaces.interface = ise_check.interface
LEFT JOIN vendors ON SUBSTRING(REPLACE(interfaces.mac_address, '.', ''), 1, 6) = vendors.mac_address
JOIN akips_interface_usage ON akips_interface_usage.switch_id = interfaces.switch_id AND akips_interface_usage.interface = interfaces.interface
WHERE
	DATE(interfaces.created_at) = CURDATE() ORDER BY interfaces.id
