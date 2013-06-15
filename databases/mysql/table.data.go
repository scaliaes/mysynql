package mysql

import (
  "fmt"
  "github.com/scalia/mysynql/log"
)

func (table *Table) Data(conn *Connection, conflictStrategy string) {
  log.Log(fmt.Sprintf("Inserting data into table `%s`", table.Name))

  if 0 == len(table.Rows) {
    return
  }

  row := table.Rows[0]
  sqlSelect, sqlInsert := "", ""

  sqlSelect = fmt.Sprintf("SELECT COUNT(*) FROM `%s` WHERE", table.Name)

  switch conflictStrategy {
  case "fail":
    sqlInsert = fmt.Sprintf("INSERT INTO `%s` (", table.Name)
  case "skip":
    sqlInsert = fmt.Sprintf("INSERT IGNORE INTO `%s` (", table.Name)
  case "replace":
    sqlInsert = fmt.Sprintf("REPLACE INTO `%s` (", table.Name)
  }

  first := true
  for _, field := range row.Fields {
    if !first {
      sqlSelect += " AND"
      sqlInsert += ", "
    }

    sqlSelect += fmt.Sprintf(" (`%s`=? OR (`%s` IS NULL AND ? IS NULL))", field.Name, field.Name)
    sqlInsert += fmt.Sprintf("`%s`", field.Name)

    first = false
  }
  sqlInsert += ")"

  sqlInsert += " VALUES ("
  first = true
  for i := len(row.Fields); i > 0; i-- {
    if !first {
      sqlInsert += ", "
    }
    sqlInsert += "?"
    first = false
  }
  sqlInsert += ")"

  log.Debug(sqlSelect)
  stmtSelect, err := conn.Prepare(sqlSelect)
  if nil != err { // Unknown error happened.
    panic(err)
  }

  log.Debug(sqlInsert)
  stmtInsert, err := conn.Prepare(sqlInsert)
  if nil != err { // Unknown error happened.
    panic(err)
  }

  for _, row := range table.Rows {
    paramsSelect := make([]interface{}, 0)
    paramsInsert := make([]interface{}, 0)

    for _, field := range row.Fields {
      if field.IsNull {
        paramsSelect = append(paramsSelect, nil, nil)
        paramsInsert = append(paramsInsert, nil)
      } else {
        paramsSelect = append(paramsSelect, field.Value, field.Value)
        paramsInsert = append(paramsInsert, field.Value)
      }
    }

    rows, res, err := stmtSelect.Exec(paramsSelect...)
    if nil != err {
      panic(err)
    }

    // Insert row
    if 0 == rows[0].Int(res.Map("COUNT(*)")) {
      _, _, err = stmtInsert.Exec(paramsInsert...)
      if nil != err {
        panic(err)
      }
    }
  }
}
