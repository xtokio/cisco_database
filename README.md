# Cisco database

Save cisco switch commands output to a MySQL database

## Installation

Cisco package
```bash
go get github.com/xtokio/cisco_database
```

ENV Variables `you will need a Cisco privilege level 15 password`
```bash
export CISCO_USERNAME="your_cisco_username"
export CISCO_PASSWORD="your_cisco_password"

export MYSQL_DATABASE_USERNAME="your_mysql_username"
export MYSQL_DATABASE_PASSWORD="your_mysql_password"
```

Database schema inside `database/schema.sql`

Add the dependency to your `main.go` file:

  ```go
 	import (
    "github.com/xtokio/cisco_database"
	)
  ```

## Contributing

1. Fork it (<https://github.com/xtokio/cisco_database/fork>)
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create a new Pull Request

## Contributors

- [Luis Gomez](https://github.com/xtokio) - creator and maintainer