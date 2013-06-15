package mysql

import (
  "fmt"
)

func (column *Column) Definition() string {
  sql := column.FullType

  // These types have a charset and a collation.
  switch column.Type {
  case "char":
    fallthrough
  case "varchar":
    fallthrough
  case "text":
    fallthrough
  case "enum":
    fallthrough
  case "set":
    sql += fmt.Sprintf(" CHARACTER SET %s COLLATE %s", column.Charset, column.Collation)
  }

  if !column.Null {
    sql += " NOT NULL"
  }
  if "" != column.Extra {
    sql += " " + column.Extra
  }
  if "" != column.Default {
    sql += fmt.Sprintf(" DEFAULT '%s'", column.Default)
  }

  return sql
}
