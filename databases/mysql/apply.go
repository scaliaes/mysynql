package mysql

import (
	"github.com/scalia/mysynql/log"
	"fmt"
)

func Apply(database *Database, host, user, pass, dbname string, noData bool, conflictStrategy string) {
	defer func() {
		if r := recover(); nil != r {
			log.Error(fmt.Sprintf("%s", r))
		}
	}()

	conn := NewConnection(host, user, pass, dbname)

	database.Apply(&conn, noData, conflictStrategy)
}
