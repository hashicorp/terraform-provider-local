package provider

import (
	"regexp"
	"testing"
)

func TestValidateNoTrailingSlash(t *testing.T) {
	testCases := []struct {
		val         string
		expectedErr *regexp.Regexp
	}{
		{
			val: "0777",
		},
		{
			val: "0644",
		},
		{
			val:         "9999",
			expectedErr: regexp.MustCompile(`bad mode for file - must be three octal digits: 9999`),
		},
		{
			val:         "7",
			expectedErr: regexp.MustCompile(`bad mode for file - string length should be 3 or 4 digits: 7`),
		},
		{
			val:         "00700",
			expectedErr: regexp.MustCompile(`bad mode for file - string length should be 3 or 4 digits: 00700`),
		},
		{
			val:         "-1",
			expectedErr: regexp.MustCompile(`bad mode for file - must be three octal digits: -1`),
		},
	}

	matchErr := func(errs []error, r *regexp.Regexp) bool {
		// err must match one provided
		for _, err := range errs {
			if r.MatchString(err.Error()) {
				return true
			}
		}

		return false
	}

	for i, tc := range testCases {
		_, errs := validateMode(tc.val, "test_property")

		if len(errs) == 0 && tc.expectedErr == nil {
			continue
		}

		if len(errs) != 0 && tc.expectedErr == nil {
			t.Fatalf("expected test case %d to produce no errors, got %v", i, errs)
		}

		if !matchErr(errs, tc.expectedErr) {
			t.Fatalf("expected test case %d to produce error matching \"%s\", got %v", i, tc.expectedErr, errs)
		}
	}
}
