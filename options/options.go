package options

type Options struct {
  // Connection parameters.
  User, Pass, Host, SchemaName string

  // Modes.
  StructureFile, DumpFile string
  Version                 bool

  // Dump mode parameters.
  DataTables, TruncateTables       StringList
  DataTablesAll, TruncateTablesAll bool

  // Restore mode parameters.
  DeleteTables     bool
  NoData           bool
  ConflictStrategy string

  // Verbosity level.
  VeryQuiet, Quiet, Verbose, VeryVerbose, Debug bool
}
