package str

import (
	"net"
	"net/url"
	"regexp"
	"unicode/utf8"

	iso6391 "github.com/emvi/iso-639-1"
	country "github.com/mikekonan/go-countries"
	"golang.org/x/text/currency"

	"github.com/egsam98/ecto"
)

// Min restricts string length with a lower inclusive bound
func Min(length uint) ecto.Test[string] {
	return ecto.Test[string]{
		Error: ecto.Errorf("must be at least %d characters long", length),
		Func:  func(v *string) bool { return utf8.RuneCountInString(*v) >= int(length) },
	}
}

// Max restricts string length with an upper inclusive bound
func Max(length uint) ecto.Test[string] {
	return ecto.Test[string]{
		Error: ecto.Errorf("must be at most %d characters long", length),
		Func:  func(v *string) bool { return utf8.RuneCountInString(*v) <= int(length) },
	}
}

// Regex validates string against regular expression
func Regex(regex *regexp.Regexp) ecto.Test[string] {
	return ecto.Test[string]{
		Error: ecto.Error("must match regex " + regex.String()),
		Func:  func(v *string) bool { return regex.MatchString(*v) },
	}
}

func URL() ecto.Test[string] {
	return ecto.Test[string]{
		Error: "invalid URL",
		Func: func(v *string) bool {
			uri, err := url.Parse(*v)
			return err == nil && uri.Scheme != "" && uri.Host != ""
		},
	}
}

// Currency ISO-4217 standard for currencies
func Currency() ecto.Test[string] {
	return ecto.Test[string]{
		Error: "invalid ISO-4217 currency",
		Func: func(v *string) bool {
			_, err := currency.ParseISO(*v)
			return err == nil
		},
	}
}

// IP address (see net.IP)
func IP() ecto.Test[string] {
	return ecto.Test[string]{
		Error: "invalid IP address",
		Func:  func(v *string) bool { return net.ParseIP(*v) != nil },
	}
}

// Lang ISO-639-1 standard for languages
func Lang() ecto.Test[string] {
	return ecto.Test[string]{
		Error: "invalid ISO-639-1 language code",
		Func:  func(v *string) bool { return iso6391.ValidCode(*v) },
	}
}

// Country ISO-3166 standard for countries
func Country() ecto.Test[string] {
	return ecto.Test[string]{
		Error: "invalid ISO-3166 country code",
		Func: func(v *string) bool {
			_, ok := country.ByAlpha2CodeStr(*v)
			return ok
		},
	}
}
