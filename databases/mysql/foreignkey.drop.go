package mysql

import (
  "fmt"
  "github.com/scalia/mysynql/log"
)

func (fk *ForeignKey) Drop(conn *Connection, table string) {
  sql := fmt.Sprintf("ALTER TABLE `%s` DROP FOREIGN KEY `%s`", table, fk.Name)

  log.Debug(sql)

  _, _, err := conn.Query(sql)
  if nil != err {
    panic(err)
  }
}
