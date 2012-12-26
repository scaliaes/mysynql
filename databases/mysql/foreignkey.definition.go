package mysql

import (
	"fmt"
)

func (fk *ForeignKey) Definition() string {
	sql := fmt.Sprintf("CONSTRAINT `%s` FOREIGN KEY (", fk.Name)
	first := true
	for _, column := range fk.Columns {
		if ! first {
			sql += ","
		}
		sql += "`" + column.Referencer + "`"
		first = false
	}

	sql += fmt.Sprintf(") REFERENCES `%s`.`%s` (", fk.Schema, fk.Table)
	first = true
	for _, column := range fk.Columns {
		if ! first {
			sql += ","
		}
		sql += "`" + column.Referenced + "`"
		first = false
	}

	sql += fmt.Sprintf(") ON DELETE %s ON UPDATE %s", fk.OnDelete, fk.OnUpdate)

	return sql
}
