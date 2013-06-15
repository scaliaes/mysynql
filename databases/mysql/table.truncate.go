package mysql

import (
  "fmt"
  "github.com/scalia/mysynql/log"
)

func (table *Table) Truncate(conn *Connection) {
  log.Log(fmt.Sprintf("Truncating table `%s`", table.Name))

  sql := fmt.Sprintf("TRUNCATE TABLE `%s`", table.Name)
  log.Debug(sql)

  _, _, err := conn.Query(sql)
  if nil != err {
    panic(err)
  }
}
