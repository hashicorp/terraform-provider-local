// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

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
			"This function works only with directories. If used with a file, the function will return an error.\n\n" +
			"This function behaves similar to the built-in [`fileexists`](https://developer.hashicorp.com/terraform/language/functions/fileexists) function, " +
			"however, `direxists` will not replace filesystem paths including `~` with the current user's home directory path. This functionality can be achieved by using the built-in " +
			"[`pathexpand`](https://developer.hashicorp.com/terraform/language/functions/pathexpand) function with `direxists`, see example below.",

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
	var inputPath string

	resp.Error = req.Arguments.Get(ctx, &inputPath)
	if resp.Error != nil {
		return
	}

	directoryPath := inputPath
	if !filepath.IsAbs(directoryPath) {
		var err error
		directoryPath, err = filepath.Abs(directoryPath)
		if err != nil {
			resp.Error = function.NewArgumentFuncError(0, fmt.Sprintf("Error expanding relative path to absolute path: %s", err))
			return
		}
	}

	directoryPath = filepath.Clean(directoryPath)

	fi, err := os.Stat(directoryPath)
	if err != nil {
		if os.IsNotExist(err) {
			resp.Error = resp.Result.Set(ctx, types.BoolValue(false))
			return
		} else {
			resp.Error = function.NewArgumentFuncError(0, fmt.Sprintf("Error checking for directory: %s", err))
			return
		}
	}

	if fi.IsDir() {
		resp.Error = resp.Result.Set(ctx, types.BoolValue(true))
		return
	}
	resp.Error = function.NewArgumentFuncError(0, fmt.Sprintf("Invalid file mode detected: %q was found, but is not a directory", inputPath))
}
