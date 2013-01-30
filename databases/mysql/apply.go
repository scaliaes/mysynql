package mysql

import (
	"github.com/scalia/mysynql/log"
	"fmt"
)

func Apply(database *Database, host, user, pass, dbname string, noData bool, conflictStrategy string, deleteTables bool) (result bool) {
	defer func() {
		if r := recover(); nil != r {
			log.Error(fmt.Sprintf("%s", r))
			result = false
		}
	}()

	conn := NewConnection(host, user, pass, dbname)

	result = database.Apply(&conn, noData, conflictStrategy, deleteTables)
	return
}
