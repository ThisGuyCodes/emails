package emails

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

func ValidShape(email string) (bool, error) {
	atLoc := strings.LastIndex(email, "@")
	if atLoc == -1 {
		return false, ErrNoAt
	}
	local, domain := email[:atLoc], email[atLoc+1:]
	errs := []error{}
	if valid, err := validLocal(local); !valid {
		errs = append(errs, err)
	}
	if valid, err := validDomain(domain); !valid {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return false, errors.Join(errs...)
	}
	return true, nil
}

var localUnquotedAcceptableChars = map[rune]bool{}

func init() {
	for _, r := range "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!#$%&'*+-/=?^_`{|}~" {
		localUnquotedAcceptableChars[r] = true
	}
}

var localQuotedAcceptableChars = map[rune]bool{}

func init() {
	for _, r := range "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!#$%&'*+-/=?^_`{|}~. (),:;<>@[]\\" {
		localQuotedAcceptableChars[r] = true
	}
}

func validLocal(local string) (bool, error) {
	if len(local) == 0 || len(local) > 64 {
		return false, ErrInvalidLocal
	}
	segments := parseLocalDotSegments(local)
	for _, segment := range segments {
		if !validLocalDotSegment(segment) {
			return false, ErrInvalidLocal
		}
	}
	return true, nil
}

type localDotSegment struct {
	content string
	quoted  bool
	valid   bool
}

func validLocalDotSegment(segment localDotSegment) bool {
	if !segment.valid {
		return false
	}

	acceptableChars := localUnquotedAcceptableChars
	if segment.quoted {
		acceptableChars = localQuotedAcceptableChars
	}

	for _, r := range segment.content {
		if !acceptableChars[r] {
			return false
		}
	}
	return true
}

func parseLocalDotSegments(local string) []localDotSegment {
	segments := []localDotSegment{}
	thisSegment := []rune{}
	quoted := false
	inEscape := false
	valid := true
	firstChar := true
	for _, r := range local {
		if inEscape {
			thisSegment = append(thisSegment, r)
			inEscape = false
			continue
		}

		if firstChar {
			firstChar = false
			valid = true
			quoted = false

			if len(segments) > 0 && segments[len(segments)-1].quoted { // previous segment was quoted
				if r == '.' {
					continue
				}
			}

			switch r {
			case '"':
				quoted = true
				continue
			case '.':
				valid = false
			}
		}

		if r == '\\' {
			inEscape = true
			continue
		}

		if quoted {
			switch r {
			case '"':
				segments = append(segments, localDotSegment{string(thisSegment), quoted, valid})
				firstChar = true
				thisSegment = thisSegment[:0]
				continue
			}
		} else {
			switch r {
			case '"':
				valid = false
			case '.':
				segments = append(segments, localDotSegment{string(thisSegment), quoted, valid})
				firstChar = true
				thisSegment = thisSegment[:0]
				continue
			}
		}

		thisSegment = append(thisSegment, r)
	}

	if len(thisSegment) > 0 || !firstChar {
		segments = append(segments, localDotSegment{string(thisSegment), quoted, valid})
	}
	return segments
}

func validDomain(domain string) (bool, error) {
	switch classifyDomain(domain) {
	case dnsDomainType:
		labels := strings.Split(domain, ".")
		for _, label := range labels {
			if !validDNSLabel(label) {
				return false, ErrInvalidDomain
			}
		}
		return true, nil
	case ipv4DomainType:
		octets := strings.Split(domain[1:len(domain)-2], ".")
		if len(octets) != 4 {
			return false, ErrInvalidDomain
		}
		for _, octet := range octets {
			if octetInt, err := strconv.Atoi(octet); err != nil {
				return false, ErrInvalidDomain
			} else {
				if octetInt < 0 || octetInt > 255 {
					return false, ErrInvalidDomain
				}
			}
		}
		return true, nil
	case ipv6DomainType:
		octets := strings.Split(domain[6:len(domain)-1], ":")
		emptyOctetSeen := false
		for _, octet := range octets {
			if len(octet) == 0 {
				if emptyOctetSeen {
					return false, ErrInvalidDomain
				}
				emptyOctetSeen = true
				continue
			} else if len(octet) > 4 {
				return false, ErrInvalidDomain
			}
			if !isHexDigit.MatchString(octet) {
				return false, ErrInvalidDomain
			}
		}
		if len(octets) > 8 {
			return false, ErrInvalidDomain
		} else if len(octets) < 8 && !emptyOctetSeen {
			return false, ErrInvalidDomain
		}
		return true, nil
	default:
		return false, ErrInvalidDomain
	}
}

var isHexDigit = regexp.MustCompile(`^[0-9a-fA-F]+$`)

type domainType byte

const (
	dnsDomainType domainType = iota
	ipv4DomainType
	ipv6DomainType
)

func classifyDomain(domain string) domainType {
	if domain[0] == '[' && domain[len(domain)-1] == ']' {
		if strings.HasPrefix(domain, "[IPv6:") {
			return ipv6DomainType
		}
		return ipv4DomainType
	}
	return dnsDomainType
}

var validLabelRegexp = regexp.MustCompile(`^[a-zA-Z0-9-]+$`)
var allNumbersRegexp = regexp.MustCompile(`^[0-9]+$`)

func validDNSLabel(label string) bool {
	if len(label) == 0 || len(label) > 63 {
		return false
	}
	if label[0] == '-' || label[len(label)-1] == '-' {
		return false
	}
	if allNumbersRegexp.MatchString(label) {
		return false
	}
	return validLabelRegexp.MatchString(label)
}
