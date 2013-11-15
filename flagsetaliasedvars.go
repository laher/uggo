//extend FlagSet with support for flag aliasMap
package uggo

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"
)

type FlagSetWithAliases struct {
	*flag.FlagSet
	name           string
	argUsage       string
	out            io.Writer
	aliasMap       map[string][]string
	isPrintUsage   *bool //optional usage (you can use your own, or none, instead)
	isPrintVersion *bool //optional (you can use your own, or none, instead)
	version        string
}

func NewFlagSet(desc string, errorHandling flag.ErrorHandling) FlagSetWithAliases {
	fs := flag.NewFlagSet(desc, errorHandling)
	fs.SetOutput(ioutil.Discard)
	return FlagSetWithAliases{fs, desc, "", os.Stderr, map[string][]string{}, nil, nil, "unknown"}
}

func NewFlagSetDefault(name, argUsage, version string) FlagSetWithAliases {
	fs := flag.NewFlagSet(name+" "+argUsage, flag.ContinueOnError)
	fs.SetOutput(ioutil.Discard)
	tmp := false
	tmp2 := false
	flagSet := FlagSetWithAliases{fs, name, argUsage, os.Stderr, map[string][]string{}, &tmp, &tmp2, version}
	flagSet.BoolVar(flagSet.isPrintUsage, "help", false, "Show this help")
	flagSet.BoolVar(flagSet.isPrintVersion, "version", false, "Show version")
	flagSet.version = version
	return flagSet
}

func (flagSet FlagSetWithAliases) ProcessHelpOrVersion() bool {
	if flagSet.IsHelp() {
		flagSet.Usage()
		return true
	} else if flagSet.IsVersion() {
		flagSet.PrintVersion()
		return true
	}
	return false
}

//convenience method for storing 'usage' behaviour
func (flagSet FlagSetWithAliases) IsHelp() bool {
	return *flagSet.isPrintUsage
}

//convenience method for storing 'get version' behaviour
func (flagSet FlagSetWithAliases) IsVersion() bool {
	return *flagSet.isPrintVersion
}

func (flagSet FlagSetWithAliases) PrintVersion() {
	fmt.Fprintf(flagSet.out, "`%s` version: '%s'\n", flagSet.name, flagSet.version)
}

func (flagSet FlagSetWithAliases) Usage() {
	fmt.Fprintf(flagSet.out, "`%s %s`:\n", flagSet.name, flagSet.argUsage)
	flagSet.PrintDefaults()
}

func (flagSet FlagSetWithAliases) SetOutput(out io.Writer) {
	flagSet.out = out
}

func (flagSet FlagSetWithAliases) AliasedBoolVar(p *bool, items []string, def bool, description string) {
	flagSet.RecordAliases(items, "bool")
	for _, item := range items {
		flagSet.BoolVar(p, item, def, description)
	}
}

func (flagSet FlagSetWithAliases) AliasedDurationVar(p *time.Duration, items []string, def time.Duration, description string) {
	flagSet.RecordAliases(items, "duration")
	for _, item := range items {
		flagSet.DurationVar(p, item, def, description)
	}
}

func (flagSet FlagSetWithAliases) AliasedFloat64Var(p *float64, items []string, def float64, description string) {
	flagSet.RecordAliases(items, "float64")
	for _, item := range items {
		flagSet.Float64Var(p, item, def, description)
	}
}

func (flagSet FlagSetWithAliases) AliasedIntVar(p *int, items []string, def int, description string) {
	flagSet.RecordAliases(items, "int")
	for _, item := range items {
		flagSet.IntVar(p, item, def, description)
	}
}

func (flagSet FlagSetWithAliases) AliasedInt64Var(p *int64, items []string, def int64, description string) {
	flagSet.RecordAliases(items, "int64")
	for _, item := range items {
		flagSet.Int64Var(p, item, def, description)
	}
}

func (flagSet FlagSetWithAliases) AliasedStringVar(p *string, items []string, def string, description string) {
	flagSet.RecordAliases(items, "string")
	for _, item := range items {
		flagSet.StringVar(p, item, def, description)
	}
}

func (flagSet FlagSetWithAliases) isAlternative(name string) bool {
	for _, altSlice := range flagSet.aliasMap {
		for _, alt := range altSlice {
			if alt == name {
				return true
			}
		}
	}
	return false
}

func (flagSet FlagSetWithAliases) RecordAliases(items []string, typ string) {
	var key string
	for i, item := range items {
		if i == 0 {
			key = item
			if _, ok := flagSet.aliasMap[key]; !ok {
				flagSet.aliasMap[key] = []string{}
			}
		} else {
			//key is same as before
			flagSet.aliasMap[key] = append(flagSet.aliasMap[key], item)
		}
	}
}

func (flagSet FlagSetWithAliases) Parse(call []string) error {
	gnuified := Gnuify(call)
	return flagSet.FlagSet.Parse(gnuified)
}

func (flagSet FlagSetWithAliases) PrintDefaults() {
	flagSet.PrintDefaultsTo(flagSet.out)
}

func (flagSet FlagSetWithAliases) PrintDefaultsTo(out io.Writer) {
	flagSet.FlagSet.VisitAll(func(fl *flag.Flag) {
		format := "-%s=%s"
		l := 0
		alts, isAliased := flagSet.aliasMap[fl.Name]
		if isAliased {
			li, _ := fmt.Fprintf(out, "  ")
			l += li
			if len(fl.Name) > 1 {
				li, _ := fmt.Fprint(out, "-")
				l += li
			}
			//no known straightforward way to test for boolean types
			if fl.DefValue == "false" {
				li, _ = fmt.Fprintf(out, "-%s", fl.Name)
				l += li
			} else {
				li, _ = fmt.Fprintf(out, format, fl.Name, fl.DefValue)
				l += li
			}
			fmt.Fprint(out, " ")
			l += 1
			for _, alt := range alts {
				if len(alt) > 1 {
					li, _ := fmt.Fprint(out, "-")
					l += li
				}
				li, _ := fmt.Fprintf(out, "-%s ", alt)
				l += li
			}
		} else if !flagSet.isAlternative(fl.Name) {
			li, _ := fmt.Fprint(out, "  ")
			l += li
			if len(fl.Name) > 1 {
				li, _ := fmt.Fprint(out, "-")
				l += li
			}
			if fl.DefValue == "false" {
				li, _ = fmt.Fprintf(out, "-%s", fl.Name)
				l += li
			} else {
				li, _ = fmt.Fprintf(out, format, fl.Name, fl.DefValue)
				l += li
			}
		} else {
			//fmt.Fprintf(out, "alias %s\n", fl.Name)
		}
		if !flagSet.isAlternative(fl.Name) {
			for l < 25 {
				l += 1
				fmt.Fprintf(out, " ")
			}
			fmt.Fprintf(out, ": %s\n", fl.Usage)
		}

	})
	fmt.Fprintln(out, "")
}
