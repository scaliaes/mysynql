package mysql

import (
	"fmt"
	mymysql "github.com/ziutek/mymysql/mysql"
	"github.com/scalia/mysynql/log"
//	"strings"
)

func (table *Table) Apply(conn *Connection, channel chan bool) {
	log.Log(fmt.Sprintf("Processing table \"%s\"", table.Name))
	defer func() {
		r := recover()
		if nil != r {
			log.Error(fmt.Sprintf("%s", r))
		}
		channel <- nil == r
	}()

return
	sql := fmt.Sprintf("SHOW FULL COLUMNS FROM `%s`", table.Name)
	log.Debug(sql)
	stmt, err := conn.Prepare(sql)
	if nil != err { // Unknown error happened.
		panic(err)
	}

	rows, res, err := stmt.Exec()
	if nil == err {	// Alter table.
//		table.alter()
	} else if mymysql.ER_NO_SUCH_TABLE == err.(*mymysql.Error).Code {	// Create table.
//		sql := table.create()
//		log.Debug(sql)
		return
	} else {
		panic(err)
	}

	for _, field := range res.Fields() {
		fmt.Print(field.Name, "\t\t")
	}
	fmt.Println()
	for _, row := range rows {
		for _, field := range res.Fields() {
			index := res.Map(field.Name)
			if nil == row[index] {
				fmt.Print("NULL\t\t")
			} else {
				fmt.Print(row.Str(index), "\t\t")
			}
		}
		fmt.Println()
	}
}
