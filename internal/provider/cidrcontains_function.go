// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"net/netip"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = CIDRContainsFunction{}

func NewCIDRContainsFunction() function.Function {
	return CIDRContainsFunction{}
}

type CIDRContainsFunction struct{}

func (f CIDRContainsFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "cidrcontains"
}

func (f CIDRContainsFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Checks whether a CIDR prefix contains another CIDR prefix",
		MarkdownDescription: "Returns true if the containing CIDR prefix fully encompasses the " +
			"contained CIDR prefix. Both arguments must be in CIDR notation. " +
			"To check a single IP, use a `/32` (IPv4) or `/128` (IPv6) prefix. " +
			"Supports both IPv4 and IPv6, but both arguments must be the same address family.",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "containing_prefix",
				MarkdownDescription: "CIDR prefix to check against (e.g. `10.0.0.0/8` or `fd00::/16`)",
			},
			function.StringParameter{
				Name:                "contained_prefix",
				MarkdownDescription: "CIDR prefix to check (e.g. `10.1.0.0/16` or `10.1.2.3/32`)",
			},
		},
		Return: function.BoolReturn{},
	}
}

func (f CIDRContainsFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var containing, contained string

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &containing, &contained))
	if resp.Error != nil {
		return
	}

	outer, err := netip.ParsePrefix(containing)
	if err != nil {
		resp.Error = function.NewArgumentFuncError(0, fmt.Sprintf("invalid CIDR prefix %q: %s", containing, err))
		return
	}

	inner, err := netip.ParsePrefix(contained)
	if err != nil {
		resp.Error = function.NewArgumentFuncError(1, fmt.Sprintf("invalid CIDR prefix %q: %s", contained, err))
		return
	}

	if outer.Addr().Is4() != inner.Addr().Is4() {
		resp.Error = function.NewFuncError("address family mismatch: both prefixes must be the same address family (IPv4 or IPv6)")
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, prefixContainsPrefix(outer, inner)))
}

// prefixContainsPrefix returns true if outer fully contains inner.
// The inner prefix must be at least as specific (longer mask) as the outer,
// and the inner network address must fall within the outer prefix.
func prefixContainsPrefix(outer, inner netip.Prefix) bool {
	if inner.Bits() < outer.Bits() {
		return false
	}
	return outer.Contains(inner.Addr())
}
