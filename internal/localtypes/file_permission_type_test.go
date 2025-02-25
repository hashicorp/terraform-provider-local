// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package localtypes

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestFilePermissionValueValidateAttribute(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		value    FilePermissionValue
		request  xattr.ValidateAttributeRequest
		expected xattr.ValidateAttributeResponse
	}{
		"0777": {
			value: FilePermissionValue{basetypes.NewStringValue("0777")},
			request: xattr.ValidateAttributeRequest{
				Path: path.Root("test"),
			},
			expected: xattr.ValidateAttributeResponse{},
		},
		"0644": {
			value: FilePermissionValue{basetypes.NewStringValue("0644")},
			request: xattr.ValidateAttributeRequest{
				Path: path.Root("test"),
			},
			expected: xattr.ValidateAttributeResponse{},
		},
		"9999": {
			value: FilePermissionValue{basetypes.NewStringValue("9999")},
			request: xattr.ValidateAttributeRequest{
				Path: path.Root("test"),
			},
			expected: xattr.ValidateAttributeResponse{
				Diagnostics: diag.Diagnostics{
					diag.NewAttributeErrorDiagnostic(
						path.Root("test"),
						"Invalid File Permission String Value",
						"bad mode permission: string must be expressed in octal numeric notation: 9999",
					),
				},
			},
		},
		"7": {
			value: FilePermissionValue{basetypes.NewStringValue("7")},
			request: xattr.ValidateAttributeRequest{
				Path: path.Root("test"),
			},
			expected: xattr.ValidateAttributeResponse{
				Diagnostics: diag.Diagnostics{
					diag.NewAttributeErrorDiagnostic(
						path.Root("test"),
						"Invalid File Permission String Value",
						"bad mode permission: string length should be 3 or 4 digits: 7",
					),
				},
			},
		},
		"00700": {
			value: FilePermissionValue{basetypes.NewStringValue("00700")},
			request: xattr.ValidateAttributeRequest{
				Path: path.Root("test"),
			},
			expected: xattr.ValidateAttributeResponse{
				Diagnostics: diag.Diagnostics{
					diag.NewAttributeErrorDiagnostic(
						path.Root("test"),
						"Invalid File Permission String Value",
						"bad mode permission: string length should be 3 or 4 digits: 00700",
					),
				},
			},
		},
		"-1": {
			value: FilePermissionValue{basetypes.NewStringValue("-1")},
			request: xattr.ValidateAttributeRequest{
				Path: path.Root("test"),
			},
			expected: xattr.ValidateAttributeResponse{
				Diagnostics: diag.Diagnostics{
					diag.NewAttributeErrorDiagnostic(
						path.Root("test"),
						"Invalid File Permission String Value",
						"bad mode permission: string length should be 3 or 4 digits: -1",
					),
				},
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			got := xattr.ValidateAttributeResponse{}

			testCase.value.ValidateAttribute(context.Background(), testCase.request, &got)

			if diff := cmp.Diff(testCase.expected, got); diff != "" {
				t.Errorf("unexpected response: %s", diff)
			}
		})
	}
}
