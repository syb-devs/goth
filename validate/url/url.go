package url

import (
	"regexp"

	"github.com/syb-devs/goth/validate"
	regexpval "github.com/syb-devs/goth/validate/regexp"
)

// Regexp is the regular expression used to validate URL's
var Regexp = regexp.MustCompile(`^(https?:\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?$`)

func init() {
	validate.RegisterRule("url", &regexpval.Rule{Regexp: Regexp})
}
