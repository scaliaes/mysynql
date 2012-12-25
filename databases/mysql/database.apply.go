package mysql

import (
	"github.com/scalia/mysynql/log"
	"github.com/scalia/mysynql/options"
	"fmt"
	"encoding/xml"
)

func (database *Database) Apply(conn *Connection) {
	opts := & options.ProgramOptions

	current := ReadConnection(opts.Host, opts.User, opts.Pass, opts.SchemaName, false)

	if opts.Debug {
		xml, err := xml.MarshalIndent(current, "", "\t")
		if nil != err {
			panic(err)
		}

		log.Debug(string(xml))
	}

	// Delete tables in database and not in XML.
	for index := range current.Tables {
		present := false
		for index2 := range database.Tables {
			if current.Tables[index].Name == database.Tables[index2].Name {
				present = true
				break
			}
		}

		if ! present {
			log.Log(fmt.Sprintf("Dropping table %s.", current.Tables[index].Name))
			sql := fmt.Sprintf("DROP TABLE `%s`", current.Tables[index].Name)
			log.Debug(sql)
		}
	}

	channel, count := make(chan bool), 0
	// Process remaining tables.
	for index := range database.Tables {
		present := false
		position := 0
		for index2 := range current.Tables {
			if database.Tables[index].Name == current.Tables[index2].Name {
				present = true
				position = index2
				break
			}
		}

		if present {	// Alter table.
			log.Log(fmt.Sprintf("Altering table %s.", database.Tables[index].Name))
			go database.Tables[index].Alter(channel, &current.Tables[position])
		} else {
			log.Log(fmt.Sprintf("Creating table %s.", database.Tables[index].Name))
			go database.Tables[index].Create(channel)
		}
		count++
	}

	result := true
	for i:= 0; i<count; i++ {
		result = result && <- channel
	}
	close(channel)

	log.Verbose("BEGIN")

	if result {
		log.Verbose("COMMIT")
	} else {
		log.Verbose("ROLLBACK")
	}
}
