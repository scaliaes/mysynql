package mysql

import (
	"github.com/scalia/mysynql/log"
	"fmt"
)

func (table *Table) Alter(channel chan bool, conn *Connection, current *Table) {
	log.Log(fmt.Sprintf("Processing diff of table `%s`", table.Name))

	defer func() {
		r := recover()
		if nil != r {
			log.Error(fmt.Sprintf("%s", r))
		}
		channel <- nil == r
	}()

	sql := fmt.Sprintf("ALTER TABLE `%s`", table.Name)
	log.Log(sql)

	if table.TruncateTable {
		table.Truncate(conn)
	}

	table.Data(conn)
}
