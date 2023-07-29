package provider

import (
	"context"
	"os"

	"github.com/chromedp/chromedp"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure ChromedpProvider satisfies various provider interfaces.
var _ provider.Provider = &ChromedpProvider{}

// ChromedpProvider defines the provider implementation.
type ChromedpProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type providerData struct {
	ctxCreator ctxCreatorFunc
}

// ChromedpProviderModel describes the provider data model.
type ChromedpProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
}

func (p *ChromedpProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "chromedp"
	resp.Version = p.version
}

func (p *ChromedpProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Description: `URL to chromedp websocket. Must be like "ws://hostname" or "ws://hostname:port".
Can be set through CHROMEDP_ENDPOINT environment variable.
If no endpoint is defined, chromedp launches existing installation of chrome (google-chrome) from $PATH.`,
				Optional: true,
			},
		},
	}
}

func pingChrome(ctx context.Context) (bool, error) {
	err := chromedp.Run(ctx, chromedp.Navigate("about://blank"))
	if err != nil {
		return false, err
	}
	return true, nil
}

func (p *ChromedpProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	endpoint := os.Getenv("CHROMEDP_ENDPOINT")

	var data ChromedpProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Endpoint.ValueString() != "" {
		endpoint = data.Endpoint.ValueString()
	}

	var ctxCreator ctxCreatorFunc
	if endpoint != "" {
		ctxCreator = chromedpCtxWithRemoteChrome(endpoint)
	} else {
		ctxCreator = chromedpCtxWithLocalChrome()
	}
	dpCtx, cancel := ctxCreator(ctx)
	defer cancel()
	_, err := pingChrome(dpCtx)
	if err != nil {
		resp.Diagnostics.AddError("Cannot start to chromedp", err.Error())
		return
	}

	resourcesData := &providerData{
		ctxCreator: ctxCreator,
	}

	resp.DataSourceData = resourcesData
	resp.ResourceData = resourcesData
}

func (p *ChromedpProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *ChromedpProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewRecipeDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ChromedpProvider{
			version: version,
		}
	}
}
