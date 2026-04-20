// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var (
	_ provider.Provider              = &UtilsProvider{}
	_ provider.ProviderWithFunctions = &UtilsProvider{}
)

// UtilsProvider defines the provider implementation.
type UtilsProvider struct {
	version string
}

func (p *UtilsProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "utils"
	resp.Version = p.version
}

func (p *UtilsProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

func (p *UtilsProvider) Configure(_ context.Context, _ provider.ConfigureRequest, _ *provider.ConfigureResponse) {
}

func (p *UtilsProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}

func (p *UtilsProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

func (p *UtilsProvider) Functions(_ context.Context) []func() function.Function {
	return []func() function.Function{
		NewCIDRContainsFunction,
		NewCIDROverlapsFunction,
		NewCIDRNoOverlapFunction,
		NewParseTree,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &UtilsProvider{
			version: version,
		}
	}
}
