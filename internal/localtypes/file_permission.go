package localtypes

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type FilePermission struct {
	types.String
}

func (f FilePermission) FromTerraform5Value(value tftypes.Value) error {
	return nil
}
