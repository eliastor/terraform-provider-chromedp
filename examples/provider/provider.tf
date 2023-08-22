terraform {
  required_providers {
    chromedp = {
      source = "eliastor/chromedp"
    }
  }
}

provider "chromedp" {
  endpoint = "ws://localhost:3000"
}
