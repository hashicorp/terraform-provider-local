package blocktoattr

import (
	"github.com/hashicorp/hcl2/ext/dynblock"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcldec"
	"github.com/hashicorp/terraform/configs/configschema"
)

// ExpandedVariables finds all of the global variables referenced in the
// given body with the given schema while taking into account the possibilities
// both of "dynamic" blocks being expanded and the possibility of certain
// attributes being written instead as nested blocks as allowed by the
// FixUpBlockAttrs function.
//
// This function exists to allow variables to be analyzed prior to dynamic
// block expansion while also dealing with the fact that dynamic block expansion
// might in turn produce nested blocks that are subject to FixUpBlockAttrs.
//
// This is intended as a drop-in replacement for dynblock.VariablesHCLDec,
// which is itself a drop-in replacement for hcldec.Variables.
func ExpandedVariables(body hcl.Body, schema *configschema.Block) []hcl.Traversal {
	rootNode := dynblock.WalkVariables(body)
	return walkVariables(rootNode, body, schema)
}

func walkVariables(node dynblock.WalkVariablesNode, body hcl.Body, schema *configschema.Block) []hcl.Traversal {
	givenRawSchema := hcldec.ImpliedSchema(schema.DecoderSpec())
	ambiguousNames := ambiguousNames(schema)
	effectiveRawSchema := effectiveSchema(givenRawSchema, body, ambiguousNames, false)
	vars, children := node.Visit(effectiveRawSchema)

	for _, child := range children {
		if blockS, exists := schema.BlockTypes[child.BlockTypeName]; exists {
			vars = append(vars, walkVariables(child.Node, child.Body(), &blockS.Block)...)
		} else if attrS, exists := schema.Attributes[child.BlockTypeName]; exists {
			synthSchema := schemaForCtyType(attrS.Type.ElementType())
			vars = append(vars, walkVariables(child.Node, child.Body(), synthSchema)...)
		}
	}

	return vars
}
