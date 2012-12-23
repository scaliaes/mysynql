package main

import (
	"github.com/scalia/mysynql/log"
	"github.com/scalia/mysynql/options"
	"github.com/scalia/mysynql/databases/mysql"
	"encoding/xml"
	"fmt"
	"time"
)

func restore() {
	startTime := time.Now()

	log.Verbose("Restoring database")

	opts := & options.ProgramOptions

	database := mysql.ReadStructure(opts.StructureFile)

	if opts.Debug {
		xml, err := xml.MarshalIndent(database, "", "\t")
		if nil != err {
			panic(err)
		}

		log.Debug(string(xml))
	}

//	database.Apply()

	endTime := time.Now()
	log.Verbose(fmt.Sprintf("Completed in %v.", endTime.Sub(startTime)))
}
