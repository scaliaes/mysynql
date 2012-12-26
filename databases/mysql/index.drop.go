package mysql

import (
	"fmt"
)

func (index *Index) Drop() string {
	sql := ""
	if "PRIMARY" == index.Name {
		sql += "DROP PRIMARY KEY"
	} else {
		sql += fmt.Sprintf("DROP KEY `%s`", index.Name)
	}

	return sql
}
