package match

import "strings"

// FirstPrefixRule returns the first matching rule’s profile (prefix match, order matters).
func FirstPrefixRule(rawURL string, rules []Rule) (profile string, ok bool) {
	for _, r := range rules {
		if r.Match == "" {
			continue
		}
		if strings.HasPrefix(rawURL, r.Match) {
			return r.Profile, true
		}
	}
	return "", false
}

type Rule struct {
	Match   string
	Profile string
}
