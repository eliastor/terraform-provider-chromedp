package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	providerConfig = `
provider "chromedp" {
	endpoint = "ws://localhost:3000"
}
`
)

var testProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"chromedp": providerserver.NewProtocol6WithError(New("test")()),
}

func testPreCheck(t *testing.T) {

}
