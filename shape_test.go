package emails_test

import (
	"errors"
	"strconv"
	"testing"

	"github.com/ThisGuyCodes/emails"
)

func Test_ValidShape(t *testing.T) {
	t.Parallel()

	t.Run("valid", func(t *testing.T) {

		for n, email := range validAddresses {
			t.Run(strconv.Itoa(n), func(t *testing.T) {
				if valid, err := emails.ValidShape(email); !valid {
					t.Errorf("ValidShape(%q) = %v, %v\nwant true, nil", email, valid, err)
				}
			})
		}
	})
	t.Run("invalid", func(t *testing.T) {
		for n, test := range invalidAddresses {
			t.Run(strconv.Itoa(n), func(t *testing.T) {
				if valid, err := emails.ValidShape(test.email); !valid {
					if !errors.Is(err, test.err) {
						t.Errorf("ValidShape(%q) = %v, %v\nwant false, %v", test.email, valid, err, test.err)
					}
				} else {
					t.Errorf("ValidShape(%q) = %v, %v\nwant false, %v", test.email, valid, err, test.err)
				}
			})
		}
	})
}
