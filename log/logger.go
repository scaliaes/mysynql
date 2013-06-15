package log

import (
  "fmt"
  "github.com/scalia/mysynql/options"
  "io"
  "os"
  "time"
)

type Level int

const (
  NONE = iota
  ERROR
  STANDARD
  VERBOSE
  DEBUG
)

func currentLevel() Level {
  opts := &options.ProgramOptions

  var level Level
  switch {
  case opts.VeryQuiet:
    level = NONE
  case opts.Quiet:
    level = ERROR
  case opts.Verbose:
    level = VERBOSE
  case opts.VeryVerbose:
    level = DEBUG
  default:
    level = STANDARD
  }

  return level
}

func log(level, message string) {
  now := time.Now()
  writer.Write([]byte(fmt.Sprintf("[%s] %s: %s\n", now.Format(time.RFC1123Z), level, message)))
}

var writer io.Writer = os.Stdout
