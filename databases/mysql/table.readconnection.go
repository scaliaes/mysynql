package mysql

import (
	"github.com/scalia/mysynql/log"
	"fmt"
)

func (table *Table) ReadConnection(conn *Connection, channel chan bool, readData, truncateData bool) {
	log.Log(fmt.Sprintf("Reading table `%s`.`%s`", conn.DbName, table.Name))

	defer func() {
		r := recover()
		if nil != r {
			log.Error(fmt.Sprintf("%s", r))
		}
		channel <- nil == r
	}()

	// Read columns.
	sql := "SELECT COLUMN_NAME, DATA_TYPE, CHARACTER_MAXIMUM_LENGTH, COLUMN_TYPE, IS_NULLABLE, COLUMN_DEFAULT, CHARACTER_SET_NAME, COLLATION_NAME, EXTRA" +
		" FROM INFORMATION_SCHEMA.COLUMNS" +
		" WHERE TABLE_SCHEMA=? AND TABLE_NAME=?" +
		" ORDER BY ORDINAL_POSITION ASC"
	stmt, err := conn.Prepare(sql)
	if nil != err { // Unknown error happened.
		panic(err)
	}
	stmt.Bind(conn.DbName, table.Name)

	rows, res, err := stmt.Exec()
	if nil != err {
		panic(err)
	}

	table.Columns = make([]Column, len(rows))
	for index, row := range rows {
		table.Columns[index].Name = row.Str(res.Map("COLUMN_NAME"))
		table.Columns[index].Type = row.Str(res.Map("DATA_TYPE"))
		table.Columns[index].Length = row.Int64(res.Map("CHARACTER_MAXIMUM_LENGTH"))
		table.Columns[index].FullType = row.Str(res.Map("COLUMN_TYPE"))
		table.Columns[index].Null = "YES" == row.Str(res.Map("IS_NULLABLE"))
		table.Columns[index].Default = row.Str(res.Map("COLUMN_DEFAULT"))
		table.Columns[index].Charset = row.Str(res.Map("CHARACTER_SET_NAME"))
		table.Columns[index].Collation = row.Str(res.Map("COLLATION_NAME"))
		table.Columns[index].Extra = row.Str(res.Map("EXTRA"))
	}

	// Read indexes.
	sql = "SELECT INDEX_NAME, COLUMN_NAME, NON_UNIQUE, COLLATION, NULLABLE, INDEX_TYPE" +
		" FROM INFORMATION_SCHEMA.STATISTICS" +
		" WHERE TABLE_SCHEMA=? AND TABLE_NAME=?" +
		" ORDER BY INDEX_NAME, SEQ_IN_INDEX ASC"

	stmt, err = conn.Prepare(sql)
	if nil != err { // Unknown error happened.
		panic(err)
	}
	stmt.Bind(conn.DbName, table.Name)

	rows, res, err = stmt.Exec()
	if nil != err {
		panic(err)
	}

	table.Indexes = make([]Index, 0)
	index := Index{}
	some := false
	for _, row := range rows {
		newName := row.Str(res.Map("INDEX_NAME"))

		if ("" != index.Name) && (newName != index.Name) {
			table.Indexes = append(table.Indexes, index)
			index.Columns = []string{}
		}

		index.Name = newName
		index.Columns = append(index.Columns, row.Str(res.Map("COLUMN_NAME")))
		index.Unique = "0" == row.Str(res.Map("NON_UNIQUE"))
		index.Collation = row.Str(res.Map("COLLATION"))
		index.Null = "YES" == row.Str(res.Map("NULLABLE"))
		index.Type = row.Str(res.Map("INDEX_TYPE"))
		some = true
	}
	if some {
		table.Indexes = append(table.Indexes, index)
	}

	// Read foreign keys. Only MySQL >= 5.1.16.
	sql = "SELECT K.CONSTRAINT_NAME, COLUMN_NAME, UNIQUE_CONSTRAINT_SCHEMA, K.REFERENCED_TABLE_NAME, REFERENCED_COLUMN_NAME, UPDATE_RULE, DELETE_RULE" +
		" FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE K" +
			" JOIN INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS R ON K.CONSTRAINT_SCHEMA=R.CONSTRAINT_SCHEMA" +
				" AND K.CONSTRAINT_NAME=R.CONSTRAINT_NAME" +
		" WHERE K.CONSTRAINT_SCHEMA=?" +
			" AND K.TABLE_NAME=?" +
		" ORDER BY K.CONSTRAINT_NAME, ORDINAL_POSITION ASC"

	stmt, err = conn.Prepare(sql)
	if nil != err { // Unknown error happened.
		panic(err)
	}
	stmt.Bind(conn.DbName, table.Name)

	rows, res, err = stmt.Exec()
	if nil != err {
		panic(err)
	}

	table.ForeignKeys = make([]ForeignKey, 0)
	fk := ForeignKey{}
	some = false
	for _, row := range rows {
		newName := row.Str(res.Map("CONSTRAINT_NAME"))

		if ("" != fk.Name) && (newName != fk.Name) {
			table.ForeignKeys = append(table.ForeignKeys, fk)
			fk.Columns = []ColumnReference{}
		}

		fk.Name = newName
		referencer := row.Str(res.Map("COLUMN_NAME"))
		referenced := row.Str(res.Map("REFERENCED_COLUMN_NAME"))
		fk.Columns = append(fk.Columns, ColumnReference{referencer, referenced} )

		fk.Schema = row.Str(res.Map("UNIQUE_CONSTRAINT_SCHEMA"))
		fk.Table = row.Str(res.Map("REFERENCED_TABLE_NAME"))
		fk.OnUpdate = row.Str(res.Map("UPDATE_RULE"))
		fk.OnDelete = row.Str(res.Map("DELETE_RULE"))
		some = true
	}
	if some {
		table.ForeignKeys = append(table.ForeignKeys, fk)
	}

	table.TruncateTable = truncateData
	if readData {
		log.Log(fmt.Sprintf("Dumping data from table `%s`.`%s`", conn.DbName, table.Name))
		// Read data
		sql = fmt.Sprintf("SELECT * FROM `%s`", table.Name)
		stmt, err = conn.Prepare(sql)
		if nil != err { // Unknown error happened.
			panic(err)
		}

		rows, res, err = stmt.Exec()
		if nil != err {
			panic(err)
		}

		table.Rows = make([]Row, 0)
		fields := res.Fields()
		for _, row := range rows {
			r := Row{}
			for _, field := range fields {
				index := res.Map(field.Name)
				isnull := nil == row[index]
				r.Fields = append(r.Fields, Field{field.Name, isnull, row.Str(index)})
			}
			table.Rows = append(table.Rows, r)
		}
	}

	log.Log(fmt.Sprintf("Done reading table `%s`.`%s`", conn.DbName, table.Name))
}

