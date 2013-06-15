package main

import (
  "encoding/xml"
  "fmt"
  "github.com/scalia/mysynql/databases/mysql"
  "github.com/scalia/mysynql/log"
  "github.com/scalia/mysynql/options"
  "io/ioutil"
  "strings"
  "time"
)

func dump() {
  startTime := time.Now()

  log.Verbose("Dumping database")

  opts := &options.ProgramOptions

  database := mysql.ReadConnection(opts.Host, opts.User, opts.Pass, opts.SchemaName, opts.DataTables, opts.TruncateTables, opts.DataTablesAll, opts.TruncateTablesAll)

  xml, err := xml.MarshalIndent(database, "", "\t")
  if nil != err {
    panic(err)
  }

  str := string(xml)
  log.Debug(str)
  str = strings.TrimSpace(str) + "\n"

  err = ioutil.WriteFile(opts.DumpFile, []byte(str), 0644)
  if nil != err {
    panic(err)
  }

  endTime := time.Now()
  log.Verbose(fmt.Sprintf("Completed in %v.", endTime.Sub(startTime)))
}
