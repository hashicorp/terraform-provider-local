package stringmodifier

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// stringDefaultModifier is a plan modifier that sets a default value for a
// types.StringType attribute when it is not configured. The attribute must be
// marked as Optional and Computed. When setting the state during the resource
// Create, Read, or Update methods, this default value must also be included or
// the Terraform CLI will generate an error.
type stringDefaultModifier struct {
	Default string
}

func (m stringDefaultModifier) Description(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to %s", m.Default)
}

func (m stringDefaultModifier) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to `%s`", m.Default)
}

func (m stringDefaultModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// If the value is known, do not set default value.
	if !req.PlanValue.IsUnknown() {
		return
	}

	resp.PlanValue = types.StringValue(m.Default)
}

func StringDefault(defaultValue string) planmodifier.String {
	return stringDefaultModifier{
		Default: defaultValue,
	}
}
