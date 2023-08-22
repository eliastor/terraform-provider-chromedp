package provider

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/chromedp/chromedp"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &RecipeDataSource{}

func NewRecipeDataSource() datasource.DataSource {
	return &RecipeDataSource{}
}

type RecipeDataSource struct {
	data *providerData
}

// type Action []types.String

type RecipeDataSourceModel struct {
	Actions            [][]types.String `tfsdk:"actions"`
	Values             types.Map        `tfsdk:"values"`
	Id                 types.String     `tfsdk:"id"`
	ScreenshotFilename types.String     `tfsdk:"screenshot_filename"`
	ScreenshotSelector types.String     `tfsdk:"screenshot_selector"`
}

func (d *RecipeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_recipe"
}

func (d *RecipeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Recipe runs list of action sequentially.
If "screenshot_filename" is set it makes the screenshot after all actions executed.

All actions are list of string where the first string is action name.

In most cases the second string is selector for the action.

The easiest way to get the selector for the element:

Open Devtools -> select element in DOM (or right click on the element at page and click "inspect element) -> right click on the element in dev tools: copy -> copy selector.

For more information about selectors https://en.wikipedia.org/wiki/CSS#Selector`,
		Attributes: map[string]schema.Attribute{
			"actions": schema.ListAttribute{
				ElementType: types.ListType{
					ElemType: types.StringType,
				},
				Required: true,
				MarkdownDescription: `
List of Actions. Each action is a list of arguments (strings).
Supported actions:
	- **navigate**: navigates the current frame to specific URL.
	
		> ["navigate", "https://github.com/eliastor/terraform-provider-chromedp"]
	
	- **click**: sends a mouse click event to the first element node matching the selector. Last argument "visible" waits for all queried elements are visible. 
	
		> ["click", "#example-After", "visible"]
	
	- **value**: gets value of form, input, textarea, select, or any other element with a ".value" field. Last argument places caught value into "values" attribute under specified key
	
		> ["value", "#example-After textarea", "text"]
	
		in values["text"] one can find caught value.
	
	- **text**: retrieves the visible text of the first element node matching the selector. Last argument places caught value into "values" attribute under specified key
	
		> ["text", "div.Documentation-function:has(#After) p", "description"]
	
		in values["description"] one can find retrieved value.
	
	- **wait_visible**: waits until selector matched element is visible:

		> ["wait_visible", "body footer"]

	- **sleep**: waits specific duration (consisting of sequences of number and unit pairs, like "1.5h" or "1m". Valid time units are "ns", "us", "ms", "s", "m", "h")

		> ["sleep", "3s"]

	- **cookie**: sets the cookie, with arguments: cookie name, value, and optional domain. 
				
		> ["cookie", "key", "value", "example.com"]

				`,
			},

			"id": schema.StringAttribute{
				MarkdownDescription: "id of recipe",
				Computed:            true,
			},
			"values": schema.MapAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: `
Map of output values from **value** and **text** actions.`,
			},
			"screenshot_filename": schema.StringAttribute{
				Optional:    true,
				Description: "If set screenshot at the end of the recipe will be made",
			},
			"screenshot_selector": schema.StringAttribute{
				Optional:    true,
				Validators:  []validator.String{stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("screenshot_filename"))},
				Description: "Requires **screenshot_filename** to be set. Points frame to the selector before making the screenshot",
			},
		},
	}
}

func (d *RecipeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*providerData)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.data = data
}

func (d *RecipeDataSource) run(ctx context.Context, actions []chromedp.Action) error {
	err := chromedp.Run(ctx, actions...)
	return err
}

func (d *RecipeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data RecipeDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Id = types.StringValue("placeholder")

	values := map[string]*string{}

	var actions []chromedp.Action

	tflog.Debug(ctx, "loop over actions")
	for _, actionArgs := range data.Actions {
		tflog.Debug(ctx, "building actions", map[string]interface{}{"args": actionArgs})
		action, err := actionBuilder(actionArgs)
		if err != nil {
			resp.Diagnostics.AddError("wrong action definition", err.Error())
			continue
		}
		actions = append(actions, action.Action(values))
	}

	if resp.Diagnostics.HasError() {
		return
	}

	picbuf := []byte{}
	screenshotPath := data.ScreenshotFilename.ValueString()
	screenshotRequested := screenshotPath != ""
	if screenshotRequested {
		screenshotDir := filepath.Dir(screenshotPath)
		err := os.MkdirAll(screenshotDir, 0600)
		if err != nil {
			resp.Diagnostics.AddError("can't create directory for the screenshot", fmt.Sprintf("%s: %v", screenshotDir, err))
			return
		}

		selector := data.ScreenshotSelector.ValueString()
		if selector != "" {
			actions = append(actions, chromedp.Screenshot(selector, &picbuf, chromedp.NodeVisible))
		} else {
			actions = append(actions, chromedp.CaptureScreenshot(&picbuf))
		}
	}
	dpCtx, cancel := d.data.ctxCreator(ctx)
	defer cancel()

	err := d.run(dpCtx, actions)
	if err != nil {
		resp.Diagnostics.AddError("can't process actions", err.Error())
	}
	var diags diag.Diagnostics

	data.Values, diags = types.MapValueFrom(ctx, types.StringType, values)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if screenshotRequested {
		err = os.WriteFile(screenshotPath, picbuf, 0600)
		if err != nil {
			resp.Diagnostics.AddError("can't save the screenshot:", err.Error())
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
