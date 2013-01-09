package mysql

import (
	"github.com/scalia/mysynql/log"
	"fmt"
)

func (table *Table) Alter(channel chan bool, conn *Connection, dbName string, current *Table, noData bool, conflictStrategy string) {
	log.Log(fmt.Sprintf("Processing diff of table `%s`", table.Name))

	defer func() {
		r := recover()
		if nil != r {
			log.Error(fmt.Sprintf("%s", r))
		}
		channel <- nil == r
	}()

	// Truncate, if needed.
	if ! noData && table.TruncateTable {
		table.Truncate(conn)
	}

	sql := ""

	// Drop unknown columns.
	first := true
	for _, currentColumn := range current.Columns {
		found := false
		for _, column := range table.Columns {
			if currentColumn.Name == column.Name {
				found = true
				break
			}
		}

		if ! found {
			if ! first {
				sql += ","
			}
			sql += fmt.Sprintf("\n\tDROP COLUMN `%s`", currentColumn.Name)
			first = false
		}
	}

	// Add and modify columns.
	for position, column := range table.Columns {
		found := false
		for _, currentColumn := range current.Columns {
			if column.Name == currentColumn.Name {
				switch false {
				case column.FullType == currentColumn.FullType: fallthrough
				case column.Null == currentColumn.Null: fallthrough
				case column.Default == currentColumn.Default: fallthrough
				case column.Charset == currentColumn.Charset: fallthrough
				case column.Collation == currentColumn.Collation: fallthrough
				case column.Extra == currentColumn.Extra:
					if ! first {
						sql += ","
					}
					sql += fmt.Sprintf("\n\tCHANGE COLUMN `%s` `%s` %s", column.Name, column.Name, column.Definition())
					first = false
				}

				found = true
				break;
			}
		}

		// Add column
		if ! found {
			if ! first {
				sql += ","
			}
			sql += fmt.Sprintf("\n\tADD COLUMN `%s` %s", column.Name, column.Definition())
			if 0 == position {
				sql += " FIRST"
			} else {
				sql += fmt.Sprintf(" AFTER `%s`", table.Columns[position-1].Name)
			}
			first = false
		}
	}

	// Drop unnecessary indexes.
	for _, currentIndex := range current.Indexes {
		found := false
		for _, index := range table.Indexes {
			if index.Name == currentIndex.Name {
				found = true
				break
			}
		}
		if ! found {
			if ! first {
				sql += ","
			}
			sql += "\n\t" + currentIndex.Drop()
			first = false
		}
	}

	// Process indexes.
	for _, index := range table.Indexes {
		found := false
		for _, currentIndex := range current.Indexes {
			if index.Name == currentIndex.Name {
				found = true

				equals := true
				if len(index.Columns) != len(currentIndex.Columns) {
					equals = false
				} else {
					for position := range index.Columns {
						if index.Columns[position] != currentIndex.Columns[position] {
							equals = false
							break
						}
					}
				}

				switch false {
				case equals: fallthrough
				case index.Unique == currentIndex.Unique: fallthrough
				case index.Collation == currentIndex.Collation: fallthrough
				case index.Null == currentIndex.Null: fallthrough
				case index.Type == currentIndex.Type:
					if ! first {
						sql += ","
					}
					sql += "\n\t" + index.Drop()
					first, found = false, false
				}
				break
			}
		}

		if ! found {
			if ! first {
				sql += ","
			}
			sql += fmt.Sprintf("\n\tADD %s", index.Definition())
			first = false
		}
	}

	// Drop unnecessary foreign keys.
	for _, currentFk := range current.ForeignKeys {
		found := false
		for _, fk := range table.ForeignKeys {
			if currentFk.Name == fk.Name {
				found = true
				break
			}
		}
		if ! found {
			currentFk.Drop(conn, table.Name)
		}
	}

	// Process foreign keys.
	for _, fk := range table.ForeignKeys {
		found := false
		for _, currentFk := range current.ForeignKeys {
			if fk.Name == currentFk.Name {
				found = true

				equals := true
				if len(fk.Columns) != len(currentFk.Columns) {
					equals = false
				} else {
					for position := range fk.Columns {
						if fk.Columns[position] != currentFk.Columns[position] {
							equals = false
							break
						}
					}
				}

				db := ""
				if dbName != fk.Schema {
					db = fk.Schema
				} else {
					db = conn.DbName
				}
				switch false {
				case equals: fallthrough
				case db == currentFk.Schema: fallthrough
				case fk.Table == currentFk.Table: fallthrough
				case fk.OnUpdate == currentFk.OnUpdate: fallthrough
				case fk.OnDelete == currentFk.OnDelete:
					fk.Drop(conn, table.Name)
					found = false
				}
				break
			}
		}

		if ! found {
			if ! first {
				sql += ","
			}
			sql += fmt.Sprintf("\n\tADD %s", fk.Definition(dbName))
			first = false
		}
	}

	// Exec query.
	if "" != sql {
		sql = fmt.Sprintf("ALTER TABLE `%s`", table.Name) + sql
		log.Debug(sql)
		_, _, err := conn.Query(sql)
		if nil != err {
			panic(err)
		}
	}

	if ! noData {
		// Add data.
		table.Data(conn, conflictStrategy)
	}
}
