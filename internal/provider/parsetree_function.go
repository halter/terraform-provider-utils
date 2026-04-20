// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ function.Function = ParseTree{}

func NewParseTree() function.Function {
	return ParseTree{}
}

var pathSegmentAttrTypes = map[string]attr.Type{
	"base_path":   types.StringType,
	"path":        types.StringType,
	"parent_path": types.StringType,
	"value":       types.StringType,
}

type pathSegment struct {
	basePath   string
	path       string
	parentPath string
	value      string
}

type ParseTree struct{}

func (f ParseTree) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "parsetree"
}

func (f ParseTree) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Parses a map of delimited paths into a flat list of tree segments.",
		MarkdownDescription: "Takes a map keyed by delimited paths (e.g. `a/b/c`) and returns a list of " +
			"objects describing each node in the tree. Intermediate path prefixes that are not present as " +
			"explicit keys are synthesized with an empty value. Each object contains: `base_path` (leaf " +
			"segment name), `path` (full delimited path, unique), `parent_path` (full path of parent, or " +
			"`\"\"` for roots), and `value` (the string value from the input map, or `\"\"` for synthesized " +
			"intermediates). The output list is sorted lexicographically by `path`.",
		Parameters: []function.Parameter{
			function.MapParameter{
				Name:                "paths",
				ElementType:         types.StringType,
				MarkdownDescription: "Map keyed by delimited path with string values.",
			},
			function.StringParameter{
				Name:                "delimiter",
				MarkdownDescription: "Delimiter used to split path keys into segments. Must be non-empty.",
			},
		},
		Return: function.ListReturn{
			ElementType: types.ObjectType{AttrTypes: pathSegmentAttrTypes},
		},
	}
}

func (f ParseTree) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	paths := make(map[string]string)
	var delimiter string

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &paths, &delimiter))
	if resp.Error != nil {
		return
	}

	if len(delimiter) == 0 {
		resp.Error = function.NewArgumentFuncError(1, "Invalid delimiter: delimiter length must be greater than 0")
		return
	}

	segments := buildPathSegments(paths, delimiter)

	elements := make([]attr.Value, 0, len(segments))
	for _, seg := range segments {
		obj, diags := types.ObjectValue(pathSegmentAttrTypes, map[string]attr.Value{
			"base_path":   types.StringValue(seg.basePath),
			"path":        types.StringValue(seg.path),
			"parent_path": types.StringValue(seg.parentPath),
			"value":       types.StringValue(seg.value),
		})
		if diags.HasError() {
			resp.Error = function.ConcatFuncErrors(resp.Error, function.FuncErrorFromDiags(ctx, diags))
			return
		}
		elements = append(elements, obj)
	}

	result, diags := types.ListValue(types.ObjectType{AttrTypes: pathSegmentAttrTypes}, elements)
	if diags.HasError() {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.FuncErrorFromDiags(ctx, diags))
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, result))
}

type treeNode struct {
	parentPath string
	value      string
}

func buildPathSegments(paths map[string]string, delimiter string) []pathSegment {
	nodes := make(map[string]*treeNode)

	for key, val := range paths {
		parts := strings.Split(key, delimiter)
		for i := range parts {
			full := strings.Join(parts[:i+1], delimiter)
			if _, ok := nodes[full]; !ok {
				parent := ""
				if i > 0 {
					parent = strings.Join(parts[:i], delimiter)
				}
				nodes[full] = &treeNode{parentPath: parent}
			}
		}
		nodes[key].value = val
	}

	fullPaths := make([]string, 0, len(nodes))
	for fp := range nodes {
		fullPaths = append(fullPaths, fp)
	}
	sort.Strings(fullPaths)

	out := make([]pathSegment, 0, len(fullPaths))
	for _, fp := range fullPaths {
		n := nodes[fp]
		parts := strings.Split(fp, delimiter)
		out = append(out, pathSegment{
			basePath:   parts[len(parts)-1],
			path:       fp,
			parentPath: n.parentPath,
			value:      n.value,
		})
	}
	return out
}
