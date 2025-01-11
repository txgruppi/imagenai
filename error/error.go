package error

import (
	"maps"
	"slices"
	"sort"
	"strings"
)

type ErrorMap map[string]error

func (t ErrorMap) Error() string {
	var s strings.Builder
	keys := make([]string, 0, len(t))
	keys = slices.AppendSeq(keys, maps.Keys(t))
	sort.Strings(keys)
	for i, k := range keys {
		s.WriteString(k)
		s.WriteString(": ")
		s.WriteString(t[k].Error())
		s.WriteString(";")
		if i < len(keys)-1 {
			s.WriteString(" ")
		}
	}
	return s.String()
}

func (t ErrorMap) Unwrap() []error {
	errs := make([]error, 0, len(t))
	errs = slices.AppendSeq(errs, maps.Values(t))
	return errs
}
