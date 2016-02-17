package url_test

import (
	"testing"

	"bitbucket.org/syb-devs/goth/validate/url"
)

func TestURLRegex(t *testing.T) {
	tests := []struct {
		Text  string
		Match bool
	}{
		{"cucu", false},
		{"http://google.com", true},
		{"https://google.com", true},
		{"http://00.sub01.domain-name.barcelona", true},
		{"htt://localhost", false},
		{"http://localhost", false},
	}

	for _, test := range tests {
		actual := url.Regexp.MatchString(test.Text)
		if actual != test.Match {
			t.Errorf("expecting match for URL %s to be %v but was %v", test.Text, test.Match, actual)
		}
	}
}
