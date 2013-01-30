package mysql

import (
	"github.com/scalia/mysynql/log"
	"github.com/scalia/mysynql/options"
	"encoding/xml"
)

func (database *Database) Apply(conn *Connection, noData bool, conflictStrategy string, deleteTables bool) bool {
	opts := & options.ProgramOptions

	current := ReadConnection(opts.Host, opts.User, opts.Pass, opts.SchemaName, options.StringList{}, options.StringList{}, false, false)

	// Disable checks.
	_, _, err := conn.Query("SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0")
	if nil != err {
		panic(err)
	}
	_, _, err = conn.Query("SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0")
	if nil != err {
		panic(err)
	}

	if opts.Debug {
		xml, err := xml.MarshalIndent(current, "", "\t")
		if nil != err {
			panic(err)
		}

		log.Debug(string(xml))
	}

	channel, count := make(chan bool), 0
	// Process tables in XML.
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

		// Truncate tables if specified on restore.
		if deleteTables {
			database.Tables[index].Truncate(conn)
		}

		if present {	// Alter table.
			go database.Tables[index].Alter(channel, conn, database.Name, &current.Tables[position], noData, conflictStrategy)
		} else {
			go database.Tables[index].Create(channel, conn, database.Name)
		}
		count++
	}

	if ! noData {
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
				go current.Tables[index].Drop(conn, channel)
				count++
			}
		}
	}

	result := true
	for i:= 0; i<count; i++ {
		result = result && <- channel
	}
	close(channel)

	_, _, err = conn.Query("SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS")
	if nil != err {
		panic(err)
	}
	_, _, err = conn.Query("SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS")
	if nil != err {
		panic(err)
	}

	return result
}
