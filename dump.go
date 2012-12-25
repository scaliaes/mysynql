package main

import (
	"github.com/scalia/mysynql/log"
	"github.com/scalia/mysynql/options"
	"github.com/scalia/mysynql/databases/mysql"
	"encoding/xml"
	"fmt"
	"time"
)

func dump() {
	startTime := time.Now()

	log.Verbose("Dumping database")

	opts := & options.ProgramOptions

	database := mysql.ReadConnection(opts.Host, opts.User, opts.Pass, opts.SchemaName, true)

	if opts.Debug {
		xml, err := xml.MarshalIndent(database, "", "\t")
		if nil != err {
			panic(err)
		}

		log.Debug(string(xml))
	}

	endTime := time.Now()
	log.Verbose(fmt.Sprintf("Completed in %v.", endTime.Sub(startTime)))
}
