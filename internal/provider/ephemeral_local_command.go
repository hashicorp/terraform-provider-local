// Copyright IBM Corp. 2017, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
)

var _ ephemeral.EphemeralResource = (*localCommandEphemeral)(nil)

func NewLocalCommandEphemeral() ephemeral.EphemeralResource {
	return &localCommandEphemeral{}
}

type localCommandEphemeral struct{}

func (e *localCommandEphemeral) Metadata(ctx context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_command"
}

func (e *localCommandEphemeral) Schema(ctx context.Context, req ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
}

func (e *localCommandEphemeral) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
}
