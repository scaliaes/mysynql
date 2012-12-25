package mysql

import (
	"github.com/scalia/mysynql/log"
	"fmt"
)

func (table *Table) Create(channel chan bool, conn *Connection) {
	log.Log(fmt.Sprintf("Creating table `%s`", table.Name))

	defer func() {
		r := recover()
		if nil != r {
			log.Error(fmt.Sprintf("%s", r))
		}
		channel <- nil == r
	}()

	sql := fmt.Sprintf("CREATE TABLE `%s` (", table.Name)

	first := true

	// Columns
	for _, column := range table.Columns {
		if ! first {
			sql += ","
		}
		sql += fmt.Sprintf("\n\t`%s` %s", column.Name, column.FullType)

		// These types have a charset and a collation.
		switch column.Type {
		case "char": fallthrough
		case "varchar": fallthrough
		case "text": fallthrough
		case "enum": fallthrough
		case "set":
			sql += fmt.Sprintf(" CHARACTER SET %s COLLATE %s", column.Charset, column.Collation)
		}

		if ! column.Null {
			sql += " NOT NULL"
		}
		if "" != column.Extra {
			sql += " " + column.Extra
		}
		if "" != column.Default {
			sql += fmt.Sprintf(" DEFAULT '%s'", column.Default)
		}
		first = false
	}

	// Indexes
	for _, index := range table.Indexes {
		if ! first {
			sql += ","
		}

		sql += "\n\t"

		if "PRIMARY" == index.Name {	// It's a primary index.
			sql += "PRIMARY KEY ("
		} else {	// It's a regular index.
			if index.Unique {
				sql += "UNIQUE "
			}
			sql += fmt.Sprintf("KEY `%s` (", index.Name)
		}
		for _, column := range index.Columns {
			sql += "`" + column + "`"
		}
		sql += ")"
		first = false
	}

	// Foreign keys.
	for _, fk := range table.ForeignKeys {
		if ! first {
			sql += ","
		}

		sql += fmt.Sprintf("\n\tCONSTRAINT `%s` FOREIGN KEY (", fk.Name)
		for _, column := range fk.Columns {
			sql += "`" + column.Referencer + "`"
		}

		sql += fmt.Sprintf(") REFERENCES `%s`.`%s` (", fk.Schema, fk.Table)
		for _, column := range fk.Columns {
			sql += "`" + column.Referenced + "`"
		}

		sql += fmt.Sprintf(") ON DELETE %s ON UPDATE %s", fk.OnDelete, fk.OnUpdate)
		first = false
	}

	sql += fmt.Sprintf("\n) ENGINE=%s DEFAULT COLLATE=%s", table.Engine, table.Collation)
	log.Debug(sql)

	_, _, err := conn.Query(sql)
	if nil != err {
		panic(err)
	}
}
