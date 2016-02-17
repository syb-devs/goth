package email

import (
	"regexp"

	"bitbucket.org/syb-devs/goth/validate"
	regexpval "bitbucket.org/syb-devs/goth/validate/regexp"
)

// Regexp is the regular expression used to validate emails
var Regexp = regexp.MustCompile("(^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+$)")

func init() {
	validate.RegisterRule("email", &regexpval.Rule{Regexp: Regexp})
}
