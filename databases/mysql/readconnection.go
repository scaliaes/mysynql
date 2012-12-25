package mysql

import (
	"github.com/scalia/mysynql/options"
)
	
func ReadConnection(host, user, pass, dbname string, dataTables options.StringList) *Database {
	database := new(Database)

	conn := NewConnection(host, user, pass, dbname)
	database.ReadConnection(&conn, dataTables)

	return database
}