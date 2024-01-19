// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ function.Function = &DirectoryExistsFunction{}

type DirectoryExistsFunction struct{}

func NewDirectoryExistsFunction() function.Function {
	return &DirectoryExistsFunction{}
}

func (f *DirectoryExistsFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "direxists"
}

func (f *DirectoryExistsFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Determines whether a directory exists at a given path.",
		Description: "Given a path string, will return true if the directory exists. " +
			"This function works only with directories. If used with a regular file, FIFO, or other " +
			"special mode, it will return an error.",

		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "path",
				Description: "Relative or absolute path to check for the existence of a directory",
			},
		},
		Return: function.BoolReturn{},
	}
}

func (f *DirectoryExistsFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var directoryPath string

	resp.Diagnostics.Append(req.Arguments.Get(ctx, &directoryPath)...)
	if resp.Diagnostics.HasError() {
		return
	}

	fi, err := os.Stat(directoryPath)
	if err != nil {
		if os.IsNotExist(err) {
			resp.Diagnostics.Append(resp.Result.Set(ctx, types.BoolValue(false))...)
			return
		} else {
			resp.Diagnostics.AddArgumentError(0, "Error checking for directory", err.Error())
			return
		}
	}

	if fi.IsDir() {
		resp.Diagnostics.Append(resp.Result.Set(ctx, types.BoolValue(true))...)
		return
	}

	resp.Diagnostics.AddArgumentError(0, "Invalid file mode detected", fmt.Sprintf("%q was found, but is not a directory", directoryPath))
}
