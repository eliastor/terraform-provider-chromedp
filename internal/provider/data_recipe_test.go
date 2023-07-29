package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRecipeDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testPreCheck(t)
		},
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + testRecipeDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.chromedp_recipe.test", "id", "placeholder"),

					resource.TestCheckResourceAttr("data.chromedp_recipe.test", "actions.0.0", "navigate"),
					resource.TestCheckResourceAttr("data.chromedp_recipe.test", "values.text", "package main\n\nimport (\n\t\"fmt\"\n\t\"time\"\n)\n\nvar c chan int\n\nfunc handle(int) {}\n\nfunc main() {\n\tselect {\n\tcase m := <-c:\n\t\thandle(m)\n\tcase <-time.After(10 * time.Second):\n\t\tfmt.Println(\"timed out\")\n\t}\n}\n"),
				),
			},
		},
	})
}

const testRecipeDataSourceConfig = `
data "chromedp_recipe" "test" {
	screenshot_filename = "test.png"
	actions = [
	  ["navigate", "https://pkg.go.dev/time"],
	  ["wait_visible", "body footer"],
	  ["click", "#example-After", "visible"],
	  # ["value", "#example-After textarea", "text"],
	  ["value", "#example-After div.Documentation-exampleDetailsBody textarea", "text"],
	  ["text", "div.Documentation-function:has(#After) p", "description"],
	  ["click", "#example-After div.Documentation-exampleButtonsContainer button.Documentation-exampleRunButton"],
	  ["sleep", "3s"],
	  ["text", "#example-After div.Documentation-exampleDetailsBody pre span.Documentation-exampleOutput", "runtext"],
	]
  }
  output "test" {
	value = data.chromedp_recipe.test.values
  }
`

//TODO: add cookie test
//Remove dependency on internet in tests
