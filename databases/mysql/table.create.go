package mysql

import (
  "fmt"
  "github.com/scalia/mysynql/log"
)

func (table *Table) Create(channel chan bool, conn *Connection, dbName string, noData bool) {
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

  // Columns.
  for _, column := range table.Columns {
    if !first {
      sql += ","
    }
    sql += fmt.Sprintf("\n\t`%s` %s", column.Name, column.Definition())
    first = false
  }

  // Indexes.
  for _, index := range table.Indexes {
    if !first {
      sql += ","
    }

    sql += "\n\t" + index.Definition()
    first = false
  }

  // Foreign keys.
  for _, fk := range table.ForeignKeys {
    if !first {
      sql += ","
    }

    sql += "\n\t" + fk.Definition(dbName)
    first = false
  }

  sql += fmt.Sprintf("\n) ENGINE=%s DEFAULT COLLATE=%s", table.Engine, table.Collation)
  log.Debug(sql)

  _, _, err := conn.Query(sql)
  if nil != err {
    panic(err)
  }

  if !noData {
    // Insert data
    table.Data(conn, "fail")
  }
}
