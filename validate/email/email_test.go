package email_test

import (
	"testing"

	"bitbucket.org/syb-devs/goth/validate/email"
)

func TestEmailRegex(t *testing.T) {
	tests := []struct {
		Text  string
		Match bool
	}{
		{"cucu", false},
		{"me@gmail.com", true},
		{"john.doe@localhost", false},
		{"john.doe-himself@localhost.com", true},
	}

	for _, test := range tests {
		actual := email.Regexp.MatchString(test.Text)
		if actual != test.Match {
			t.Errorf("expecting match for email %s to be %v but was %v", test.Text, test.Match, actual)
		}
	}
}
