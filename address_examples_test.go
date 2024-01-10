package emails_test

import "github.com/ThisGuyCodes/emails"

var validAddresses = []string{
	"simple@example.com",
	"very.common@example.com",
	"x@example.com",
	"long.email-address-with-hyphens@and.subdomains.example.com",
	"user.name+tag+sorting@example.com",
	"name/surname@example.com",
	"admin@example",
	"example@s.example",
	`" "@example.org`,
	`"john..doe"@example.org`,
	`".john.doe"@example.org`,
	`mailhost!username@example.org`,
	`"very.(),:;<>[]\".VERY.\"very@\\ \"very\".unusual"@strange.example.com`,
	"user%example.com@example.org",
	"user-@example.org",
	"postmaster@[123.123.123.123]",
	"postmaster@[IPv6:2001:0db8:85a3:0000:0000:8a2e:0370:7334]",
	"_test@[IPv6:2001:0db8:85a3:0000:0000:8a2e:0370:7334]",
	"_test@[IPv6:2001:0db8:85a3::8a2e:0370:7334]",
}

type expectedErr struct {
	email string
	err   error
}

var invalidAddresses = []expectedErr{
	{"abc.example.com", emails.ErrNoAt},
	{"a@b@c@example.com", emails.ErrInvalidLocal},
	{`a"b(c)d,e:f;g<h>i[j\k]l@example.com`, emails.ErrInvalidLocal},
	{`just"not"right@example.com`, emails.ErrInvalidLocal},
	{`this is"not\allowed@example.com`, emails.ErrInvalidLocal},
	{`this\ still\"not\\allowed@example.com`, emails.ErrInvalidLocal},
	{"1234567890123456789012345678901234567890123456789012345678901234+x@example.com", emails.ErrInvalidLocal},

	{"i.like.underscores@but_they_are_not_allowed_in_this_part", emails.ErrInvalidDomain},
	{"example@-dashes-at-start-are-invalid", emails.ErrInvalidDomain},
	{"example@dashes-at-end-are-invalid-", emails.ErrInvalidDomain},
	{"all-numbers-in-domain-are-wrong@123.com", emails.ErrInvalidDomain},
	{"domains@must-be-less-than-64-characters-12345678901234567890123456789012", emails.ErrInvalidDomain},
	{"domain-labels@are.non.zero..length", emails.ErrInvalidDomain},
	{"invalid.IPv4@[abc]", emails.ErrInvalidDomain},
	{"too.many.octets@[123.123.123.123.123]", emails.ErrInvalidDomain},
	{"too.few.octets@[123.123.123]", emails.ErrInvalidDomain},
	{"invalid.IPv6@[IPv6:123456::]", emails.ErrInvalidDomain},
	{"invalid.IPv6@[IPv6:defg::]", emails.ErrInvalidDomain},
	{"invalid.IPv6@[IPv6:abcd::efg::a]", emails.ErrInvalidDomain},
}
