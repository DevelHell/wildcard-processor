package wildcard_processor

import (
	"testing"
	"strings"
)

func TestAllowsRcpt(t *testing.T) {
	allowedHosts := []string{
		"spam4.me",
		"grr.la",
		"newhost.com",
		"example.*",
		"*.test",
		"wild*.card",
		"multiple*wild*cards.*",
	}

	c1 := WildcardConfig{WildcardHosts: strings.Join(allowedHosts, ",")}
	w1 := newWildcardProcessor(&c1)

	testTable := map[string]bool{
		"spam4.me":                true,
		"dont.match":              false,
		"example.com":             true,
		"another.example.com":     false,
		"anything.test":           true,
		"wild.card":               true,
		"wild.card.com":           false,
		"multipleXwildXcards.com": true,
	}

	for host, allows := range testTable {
		if res := w1.allowsRcpt(host); res != allows {
			t.Error(host, ": expected", allows, "but got", res)
		}
	}

	// only wildcard - should match anything
	c2 := WildcardConfig{WildcardHosts: "*"}
	w2 := newWildcardProcessor(&c2)

	if !w2.allowsRcpt("match.me") {
		t.Error("match.me: expected true but got false")
	}
}
