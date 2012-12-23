package options

type Options struct {
	// Connection parameters.
	User, Pass, Host, SchemaName string

	// Modes.
	StructureFile, DumpFile string
	Version bool

	// Verbosity level.
	VeryQuiet, Quiet, Verbose, VeryVerbose, Debug bool
}
