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
	for _, column := range index.Columns {
		sql += "`" + column + "`"
	}
	sql += ")"

	return sql
}
