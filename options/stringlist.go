package options

import (
  "fmt"
  "strings"
)

type StringList []string

func (l *StringList) String() string {
  return fmt.Sprint(*l)
}

// Set is the method to set the flag value, part of the flag.Value interface.
// Set's argument is a string to be parsed to set the flag.
// It's a comma-separated list, so we split it.
func (l *StringList) Set(value string) error {
  for _, elem := range strings.Split(value, ",") {
    *l = append(*l, elem)
  }
  return nil
}

func (l *StringList) Contains(search string) bool {
  for _, elem := range *l {
    if search == elem {
      return true
    }
  }
  return false
}
