# Terraform Provider ChromeDP

This providers allows you to use [chromedp](https://github.com/chromedp/chromedp) in terraform in form of recipes.


## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.19

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Using the provider

Please refer do examples/.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To use provider in sample terraform code `~/.terraformrc` must be edited:

```hcl
provider_installation {

  dev_overrides {
      "hashicorp.com/eliastor/chromedp" = "<your go bin path>"
  }

  direct {}
}
```


To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

```shell
make testacc
```


## Supported actions

- [x] Navigate
- [x] Wait_visible
- [x] Click
- [x] Value (getting content of forms, inputs, textareas, selects, or any other element with a '.value' field.)
- [x] Text (getting text content of the element)
- [x] Focus
- [x] Cookies (setting cookies)
- [x] Screenshots
## Roadmap 
- [ ] browserless.io support
- [ ] Text
- [ ] Downloads
- [ ] Uploads
- [ ] Send keys
- [ ] Submit
- [ ] Emulate different viewports