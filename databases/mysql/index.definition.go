package mysql

import (
	"fmt"
)

func (index *Index) Definition() string {
	sql := ""

	if "PRIMARY" == index.Name {	// It's a primary index.
		sql += "PRIMARY KEY ("
	} else {	// It's a regular index.
		if index.Unique {
			sql += "UNIQUE "
		}
		sql += fmt.Sprintf("KEY `%s` (", index.Name)
	}
	first := true
	for _, column := range index.Columns {
		if ! first {
			sql += ", "
		}
		sql += "`" + column + "`"
		first = false
	}
	sql += ")"

	return sql
}
