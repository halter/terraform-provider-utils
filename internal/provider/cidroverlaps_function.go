package provider

import (
	"context"
	"fmt"
	"net/netip"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = CIDROverlapsFunction{}

func NewCIDROverlapsFunction() function.Function {
	return CIDROverlapsFunction{}
}

type CIDROverlapsFunction struct{}

func (f CIDROverlapsFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "cidroverlaps"
}

func (f CIDROverlapsFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Checks whether two CIDR prefixes overlap",
		MarkdownDescription: "Returns true if the two CIDR prefixes have any addresses in common. " +
			"For CIDR prefixes, overlap means one fully contains the other — partial overlap is not possible. " +
			"Both arguments must be in CIDR notation and the same address family (IPv4 or IPv6).",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "prefix_a",
				MarkdownDescription: "First CIDR prefix (e.g. `10.0.0.0/8` or `fd00::/16`)",
			},
			function.StringParameter{
				Name:                "prefix_b",
				MarkdownDescription: "Second CIDR prefix (e.g. `10.1.0.0/16` or `fd00:1::/32`)",
			},
		},
		Return: function.BoolReturn{},
	}
}

func (f CIDROverlapsFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var prefixA, prefixB string

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &prefixA, &prefixB))
	if resp.Error != nil {
		return
	}

	a, err := netip.ParsePrefix(prefixA)
	if err != nil {
		resp.Error = function.NewArgumentFuncError(0, fmt.Sprintf("invalid CIDR prefix %q: %s", prefixA, err))
		return
	}

	b, err := netip.ParsePrefix(prefixB)
	if err != nil {
		resp.Error = function.NewArgumentFuncError(1, fmt.Sprintf("invalid CIDR prefix %q: %s", prefixB, err))
		return
	}

	if a.Addr().Is4() != b.Addr().Is4() {
		resp.Error = function.NewFuncError("address family mismatch: both prefixes must be the same address family (IPv4 or IPv6)")
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, a.Overlaps(b)))
}
