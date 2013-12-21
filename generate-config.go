package main

import (
	"fmt"
	"io/ioutil"
)

const CONFIG_SAMPLE = `# edit this file and run the application when you're done

# if a socket is specified we will use that
# if tcp is chosen you can use compression
mysql:
 hostname: 127.0.0.1
 port: 3306
 # socket: /tmp/mysql.sock
 username: someuser
 password: somepass
 database: somedb
 compress: false

# if file is given, output goes to file, if postgres parameters
# are given, output is executed straight on the db, socket is
# prioritized if specified.
destination:
 # file: test.sql
 postgres:
   hostname: localhost
   port: 5432
   socket: /var/run/postgresql
   username:
   password:
   database: somedb

# projections can help you align data between the source and
# destination databases, it's basically like a view (and used to be
# implemented as one). It will create a table that only lasts as long as the
# session.
projections:
    pr_players:
     engine: MEMORY
     pk: [hostname]
     body: |
         SELECT hostname, name, timetable, location,
         FROM Player
         WHERE hostname LIKE '%.new.client'
         AND name IS NOT NULL


# table "a" in the source database has been renamed to table "b"
# in the destination database
table_map:
 pr_players: players

# which tables (projections defined in this file included) should be
# synced? If not defined, all tables are synced
only_tables:
 - pr_players

# which tables should NOT be synced
#exclude_tables:
#- table3
#- table4

# the following fiels are not implemented yet, they are inherited from
# py-mysql2pgsql, leave them as-is.
merge: true

# if supress_data is true, only the schema definition will be exported/migrated, and not the data
supress_data: false

# if supress_ddl is true, only the data will be exported/imported, and not the schema
supress_ddl: false

# if force_truncate is true, forces a table truncate before table loading
force_truncate: false

# if timezone is true, forces to append/convert to UTC tzinfo mysql data
timezone: false
`

type GenerateConfigCommand struct {
	Force bool `short:"f" long:"force" description:"force config file overwrite"`
}

func (x *GenerateConfigCommand) Execute(args []string) error {
	path := "config.yml"
	if FileExists(path) && !x.Force {
		return fmt.Errorf("File %v already exists, use the -f flag if you want to overwrite", path)
	}

	fmt.Printf("Generating config...")
	err := ioutil.WriteFile("config.yml", []byte(CONFIG_SAMPLE), 0644)
	if err != nil {
		fmt.Println("ERROR")
		return err
	}

	fmt.Println("DONE")
	return nil
}

func init() {
	var genconfig GenerateConfigCommand
	parser.AddCommand("generate-config",
		"Generate a sample config file in the current directory",
		"Generate a sample config file in the current directory",
		&genconfig)
}
