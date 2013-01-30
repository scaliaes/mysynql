package main

import (
	"github.com/scalia/mysynql/log"
	"github.com/scalia/mysynql/options"
	"github.com/scalia/mysynql/databases/mysql"
	"encoding/xml"
	"fmt"
	"time"
	"os"
)

func restore() {
	startTime := time.Now()

	log.Verbose("Restoring database")

	opts := & options.ProgramOptions

	database := mysql.ReadXML(opts.StructureFile)

	if opts.Debug {
		xml, err := xml.MarshalIndent(database, "", "\t")
		if nil != err {
			panic(err)
		}

		log.Debug(string(xml))
	}

	if ! mysql.Apply(database, opts.Host, opts.User, opts.Pass, opts.SchemaName, opts.NoData, opts.ConflictStrategy, opts.DeleteTables) {
		os.Exit(2)
	}

	endTime := time.Now()
	log.Verbose(fmt.Sprintf("Completed in %v.", endTime.Sub(startTime)))
}
