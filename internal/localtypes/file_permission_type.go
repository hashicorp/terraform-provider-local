package localtypes

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ types.StringTypable    = FilePermissionType{}
	_ xattr.TypeWithValidate = FilePermissionType{}
)

type FilePermissionType struct {
	types.StringTypable
}

// Validate checks that the given input string is a valid file permission string,
// expressed in numeric notation.
// See: https://en.wikipedia.org/wiki/File-system_permissions#Numeric_notation
func (f FilePermissionType) Validate(ctx context.Context, value tftypes.Value, path path.Path) diag.Diagnostics {
	if value.IsNull() {
		//return diag.Diagnostics{
		//		//	diag.NewAttributeErrorDiagnostic(path,
		//		//		"Invalid File Permission String Value",
		//		//		"The File Permission value cannot be null"),
		//		//}
		return nil
	}

	if !value.IsKnown() {
		//return diag.Diagnostics{
		//	diag.NewAttributeErrorDiagnostic(path,
		//		"Invalid File Permission String Value",
		//		"The File Permission value cannot be unknown"),
		//}
		return nil
	}
	var fp string
	err := value.As(&fp)
	if err != nil {
		return diag.Diagnostics{
			diag.NewAttributeErrorDiagnostic(path,
				"Invalid File Permission String Value",
				"An unexpected error occurred while converting the file permission to a string value"+
					"Error: "+err.Error()),
		}
	}

	if len(fp) < 3 || len(fp) > 4 {
		return diag.Diagnostics{
			diag.NewAttributeErrorDiagnostic(path,
				"Invalid File Permission String Value",
				"bad mode permission: string length should be 3 or 4 digits: "+fp),
		}
	}

	fileMode, err := strconv.ParseInt(fp, 8, 64)
	if err != nil || fileMode > 0777 || fileMode < 0 {
		return diag.Diagnostics{
			diag.NewAttributeErrorDiagnostic(path,
				"Invalid File Permission String Value",
				"string must be expressed in octal numeric notation: "+fp),
		}
	}
	return diag.Diagnostics{}
}
