package mysql

import (
	"github.com/scalia/mysynql/log"
	"fmt"
)

func (table *Table) Data(conn *Connection) {
	log.Log(fmt.Sprintf("Inserting data into table `%s`", table.Name))

	for _, row := range table.Rows {
		sqlSelect, sqlInsert := "", ""
		
		sqlSelect = fmt.Sprintf("SELECT COUNT(*) FROM `%s` WHERE", table.Name)
		sqlInsert = fmt.Sprintf("INSERT INTO `%s` (", table.Name)

		first := true
		params := make([]interface{}, 0)
		for _, field := range row.Fields {
			if ! first {
				sqlSelect += " AND"
				sqlInsert += ", "
			}

			sqlSelect += fmt.Sprintf(" `%s`=?", field.Name)
			sqlInsert += fmt.Sprintf("`%s`", field.Name)
			params = append(params, field.Value)

			first = false
		}
		sqlInsert += ")"

		stmt, err := conn.Prepare(sqlSelect)
		if nil != err { // Unknown error happened.
				panic(err)
		}

		// Bind parameters
		stmt.Bind(params...)

		rows, res, err := stmt.Exec()
		if nil != err {
			panic(err)
		}

		// Insert row
		if 0 == rows[0].Int(res.Map("COUNT(*)")) {
			sqlInsert += " VALUES ("
			first = true
			for i:= len(params); i>0; i-- {
				if ! first {
					sqlInsert +=", "
				}
				sqlInsert += "?"
				first = false
			}
			sqlInsert += ")"

			log.Debug(sqlInsert)

			stmt, err = conn.Prepare(sqlInsert)
			if nil != err { // Unknown error happened.
					panic(err)
			}
			// Bind parameters
			stmt.Bind(params...)

			_, _, err = stmt.Exec()
			if nil != err {
				panic(err)
			}
		}
	}
}
