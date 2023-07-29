terraform {
  required_providers {
    chromedp = {
      source = "hashicorp.com/eliastor/chromedp"
    }
  }
}

provider "chromedp" {
  # for an instance: podman run --rm -d -p 3000:3000 docker.io/browserless/chrome:latest
  endpoint = "ws://localhost:3000"
}
