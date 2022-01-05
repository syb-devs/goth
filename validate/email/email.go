package email

import (
	"regexp"

	"github.com/syb-devs/goth/validate"
	regexpval "github.com/syb-devs/goth/validate/regexp"
)

// Regexp is the regular expression used to validate emails
var Regexp = regexp.MustCompile("(^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+$)")

func init() {
	validate.RegisterRule("email", &regexpval.Rule{Regexp: Regexp})
}
