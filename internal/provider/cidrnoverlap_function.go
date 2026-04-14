// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"net/netip"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ function.Function = CIDRNoOverlapFunction{}

func NewCIDRNoOverlapFunction() function.Function {
	return CIDRNoOverlapFunction{}
}

type CIDRNoOverlapFunction struct{}

func (f CIDRNoOverlapFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "cidrnoverlap"
}

func (f CIDRNoOverlapFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Checks that a list of CIDR prefixes do not overlap with each other",
		MarkdownDescription: "Returns true if no pair of CIDR prefixes in the list overlaps. " +
			"All elements must be in CIDR notation and the same address family (IPv4 or IPv6). " +
			"An empty list or single-element list returns true.",
		Parameters: []function.Parameter{
			function.ListParameter{
				Name:                "prefixes",
				ElementType:         types.StringType,
				MarkdownDescription: "List of CIDR prefix strings to check for mutual non-overlap",
			},
		},
		Return: function.BoolReturn{},
	}
}

func (f CIDRNoOverlapFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var prefixes []string

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &prefixes))
	if resp.Error != nil {
		return
	}

	parsed := make([]netip.Prefix, len(prefixes))
	for i, s := range prefixes {
		var err error
		parsed[i], err = netip.ParsePrefix(s)
		if err != nil {
			resp.Error = function.NewArgumentFuncError(0, fmt.Sprintf("invalid CIDR prefix at index %d %q: %s", i, s, err))
			return
		}
	}

	for i := range len(parsed) {
		for j := i + 1; j < len(parsed); j++ {
			if parsed[i].Addr().Is4() != parsed[j].Addr().Is4() {
				resp.Error = function.NewFuncError(fmt.Sprintf(
					"address family mismatch at index %d (%s) and %d (%s): all prefixes must be the same address family (IPv4 or IPv6)",
					i, prefixes[i], j, prefixes[j],
				))
				return
			}
			if parsed[i].Overlaps(parsed[j]) {
				resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, false))
				return
			}
		}
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, true))
}
