package mysql

import (
  "fmt"
  "github.com/scalia/mysynql/log"
)

func (table *Table) Drop(conn *Connection, channel chan bool) {
  log.Log(fmt.Sprintf("Dropping table %s.", table.Name))

  defer func() {
    r := recover()
    if nil != r {
      log.Error(fmt.Sprintf("%s", r))
    }
    channel <- nil == r
  }()

  sql := fmt.Sprintf("DROP TABLE `%s`", table.Name)
  log.Debug(sql)

  _, _, err := conn.Query(sql)
  if nil != err {
    panic(err)
  }
}
