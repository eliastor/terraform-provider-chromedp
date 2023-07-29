data "chromedp_recipe" "example" {
  screenshot_filename = "example.png"
  actions = [
    ["navigate", "https://pkg.go.dev/time"],
    ["wait_visible", "body footer"],
    ["click", "#example-After", "visible"],
    ["value", "#example-After div.Documentation-exampleDetailsBody textarea", "text"],
    ["text", "div.Documentation-function:has(#After) p", "description"],
    ["click", "#example-After div.Documentation-exampleButtonsContainer button.Documentation-exampleRunButton"],
    ["sleep", "3s"],
    ["text", "#example-After div.Documentation-exampleDetailsBody pre span.Documentation-exampleOutput", "runtext"],
  ]
}
