package mysql

import (
	"fmt"
	"github.com/scalia/mysynql/log"
)

func (table *Table) Alter(channel chan bool, current *Table) {
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
}
