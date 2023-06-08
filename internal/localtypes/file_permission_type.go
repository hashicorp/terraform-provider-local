// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package localtypes

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ basetypes.StringTypable = FilePermissionType{}
	_ xattr.TypeWithValidate  = FilePermissionType{}
)

type FilePermissionType struct {
	basetypes.StringType
}

func NewFilePermissionType() FilePermissionType {
	return FilePermissionType{StringType: types.StringType}
}

// Validate checks that the given input string is a valid file permission string,
// expressed in numeric notation.
// See: https://en.wikipedia.org/wiki/File-system_permissions#Numeric_notation
func (f FilePermissionType) Validate(ctx context.Context, value tftypes.Value, path path.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	if value.IsNull() {
		return diags
	}

	if !value.IsKnown() {
		return diags
	}

	var fp string
	err := value.As(&fp)
	if err != nil {
		diags.Append(diag.NewAttributeErrorDiagnostic(path,
			"Invalid File Permission String Value",
			"An unexpected error occurred while converting the file permission to a string value"+
				"Error: "+err.Error()))
		return diags
	}

	if len(fp) < 3 || len(fp) > 4 {
		diags.Append(diag.NewAttributeErrorDiagnostic(path,
			"Invalid File Permission String Value",
			"bad mode permission: string length should be 3 or 4 digits: "+fp))
		return diags
	}

	fileMode, err := strconv.ParseInt(fp, 8, 64)
	if err != nil || fileMode > 0777 || fileMode < 0 {
		diags.Append(diag.NewAttributeErrorDiagnostic(path,
			"Invalid File Permission String Value",
			"bad mode permission: string must be expressed in octal numeric notation: "+fp))
		return diags
	}
	return diags
}
