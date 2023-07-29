default: test doc-gen

.PHONY: test

test:
	podman run --name terraform-provider-chromedp-test --rm -d -p 3000:3000 docker.io/browserless/chrome:latest
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m || podman stop terraform-provider-chromedp-test
	podman stop terraform-provider-chromedp-test
doc-gen:
	go generate ./...
install:
	go install