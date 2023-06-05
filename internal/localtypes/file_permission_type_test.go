// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package localtypes

import (
	"context"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestFilePermissionTypeValidator(t *testing.T) {
	t.Parallel()

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
			expectedErr: regexp.MustCompile(`bad mode permission: string must be expressed in octal numeric notation: 9999`),
		},
		{
			val:         "7",
			expectedErr: regexp.MustCompile(`bad mode permission: string length should be 3 or 4 digits: 7`),
		},
		{
			val:         "00700",
			expectedErr: regexp.MustCompile(`bad mode permission: string length should be 3 or 4 digits: 00700`),
		},
		{
			val:         "-1",
			expectedErr: regexp.MustCompile(`bad mode permission: string length should be 3 or 4 digits: -1`),
		},
	}

	matchErr := func(diags diag.Diagnostics, r *regexp.Regexp) bool {
		// err must match one provided
		for _, err := range diags {
			if r.MatchString(err.Detail()) {
				return true
			}
		}

		return false
	}

	for i, tc := range testCases {
		diags := NewFilePermissionType().Validate(context.Background(), tftypes.NewValue(tftypes.String, tc.val), path.Empty())

		if !diags.HasError() && tc.expectedErr == nil {
			continue
		}

		if diags.HasError() && tc.expectedErr == nil {
			t.Fatalf("expected test case %d to produce no errors, got %v", i, diags.Errors())
		}

		if !matchErr(diags, tc.expectedErr) {
			t.Fatalf("expected test case %d to produce error matching \"%s\", got %v", i, tc.expectedErr, diags.Errors())
		}
	}
}
