package validate

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	reZIP   = regexp.MustCompile(`^[A-Za-z0-9\- ]{3,10}$`)
	reEmail = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
)

func Region(s string) (string, bool) {
	s = strings.TrimSpace(s)
	return s, s != "" && reZIP.MatchString(s)
}

func Email(s string) (string, bool) {
	s = strings.TrimSpace(s)
	return s, s != "" && reEmail.MatchString(s)
}
func Qty(s string) int {
	n, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil || n < 1 {
		return 1
	}
	if n > 50 {
		return 50
	} // clamp to avoid abuse
	return n
}
