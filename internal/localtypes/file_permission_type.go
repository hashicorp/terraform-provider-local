// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package localtypes

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ basetypes.StringTypable     = FilePermissionType{}
	_ basetypes.StringValuable    = FilePermissionValue{}
	_ xattr.ValidateableAttribute = FilePermissionValue{}
)

type FilePermissionType struct {
	basetypes.StringType
}

func (t FilePermissionType) Equal(o attr.Type) bool {
	other, ok := o.(FilePermissionType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t FilePermissionType) String() string {
	return "FilePermissionType"
}

func (t FilePermissionType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := FilePermissionValue{
		StringValue: in,
	}

	return value, nil
}

func (t FilePermissionType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.StringType.ValueFromTerraform(ctx, in)

	if err != nil {
		return nil, err
	}

	stringValue, ok := attrValue.(basetypes.StringValue)

	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	stringValuable, diags := t.ValueFromString(ctx, stringValue)

	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting StringValue to StringValuable: %v", diags)
	}

	return stringValuable, nil
}

func (t FilePermissionType) ValueType(ctx context.Context) attr.Value {
	return FilePermissionValue{}
}

func NewFilePermissionType() FilePermissionType {
	return FilePermissionType{StringType: types.StringType}
}

type FilePermissionValue struct {
	basetypes.StringValue
}

func (v FilePermissionValue) Equal(o attr.Value) bool {
	other, ok := o.(FilePermissionValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v FilePermissionValue) Type(ctx context.Context) attr.Type {
	return FilePermissionType{}
}

// ValidateAttribute checks that the given input string is a valid file permission string,
// expressed in numeric notation.
// See: https://en.wikipedia.org/wiki/File-system_permissions#Numeric_notation
func (v FilePermissionValue) ValidateAttribute(ctx context.Context, req xattr.ValidateAttributeRequest, resp *xattr.ValidateAttributeResponse) {
	if v.IsNull() {
		return
	}

	if v.IsUnknown() {
		return
	}

	fp := v.ValueString()

	if len(fp) < 3 || len(fp) > 4 {
		resp.Diagnostics.Append(diag.NewAttributeErrorDiagnostic(req.Path,
			"Invalid File Permission String Value",
			"bad mode permission: string length should be 3 or 4 digits: "+fp))
		return
	}

	fileMode, err := strconv.ParseInt(fp, 8, 64)
	if err != nil || fileMode > 0777 || fileMode < 0 {
		resp.Diagnostics.Append(diag.NewAttributeErrorDiagnostic(req.Path,
			"Invalid File Permission String Value",
			"bad mode permission: string must be expressed in octal numeric notation: "+fp))
		return
	}
}
