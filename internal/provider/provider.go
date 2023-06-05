// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var (
	_ provider.Provider = (*localProvider)(nil)
)

func New() provider.Provider {
	return &localProvider{}
}

type localProvider struct{}

func (p *localProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "local"
}

func (p *localProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

}

func (p *localProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewLocalFileDataSource,
		NewLocalSensitiveFileDataSource,
	}
}

func (p *localProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewLocalFileResource,
		NewLocalSensitiveFileResource,
	}
}

func (p *localProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

type fileChecksums struct {
	md5Hex       string
	sha1Hex      string
	sha256Hex    string
	sha256Base64 string
	sha512Hex    string
	sha512Base64 string
}

func genFileChecksums(data []byte) fileChecksums {
	var checksums fileChecksums

	md5Sum := md5.Sum(data)
	checksums.md5Hex = hex.EncodeToString(md5Sum[:])

	sha1Sum := sha1.Sum(data)
	checksums.sha1Hex = hex.EncodeToString(sha1Sum[:])

	sha256Sum := sha256.Sum256(data)
	checksums.sha256Hex = hex.EncodeToString(sha256Sum[:])
	checksums.sha256Base64 = base64.StdEncoding.EncodeToString(sha256Sum[:])

	sha512Sum := sha512.Sum512(data)
	checksums.sha512Hex = hex.EncodeToString(sha512Sum[:])
	checksums.sha512Base64 = base64.StdEncoding.EncodeToString(sha512Sum[:])

	return checksums
}
