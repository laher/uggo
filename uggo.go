package uggo

import (
	"strings"
)


//note: allow '-help' to be used as single-hyphen (to assist the unitiated)
func Gnuify(call []string) []string {
	return GnuifyWithExceptions(call, []string{ "-help" })
}

func contains(slice []string, subject string) bool {
	for _, item := range slice {
		if item == subject {
			return true
		}
	}
	return false
}

func GnuifyWithExceptions(call, exceptions []string) []string {
	splut := []string{}
	for _, item := range call {
		if strings.HasPrefix(item, "-") && !strings.HasPrefix(item, "--") && !contains(exceptions, item) {
			for _, letter := range item[1:] {
				splut = append(splut, "-"+string(letter))
			}
		} else {
			splut = append(splut, item)
		}
	}
	return splut
}
