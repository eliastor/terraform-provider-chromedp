package provider

import (
	"context"

	"github.com/chromedp/chromedp"
)

type ctxCreatorFunc func(ctx context.Context) (context.Context, context.CancelFunc)

var (
	_ ctxCreatorFunc = chromedpCtxWithLocalChrome()
	_ ctxCreatorFunc = chromedpCtxWithRemoteChrome("")
)

func chromedpCtxWithLocalChrome() ctxCreatorFunc {
	return func(parentCtx context.Context) (context.Context, context.CancelFunc) {
		return chromedp.NewContext(parentCtx)
	}
}

func chromedpCtxWithRemoteChrome(remote string) ctxCreatorFunc {
	return func(parentCtx context.Context) (context.Context, context.CancelFunc) {
		ctx, _ := chromedp.NewRemoteAllocator(parentCtx, remote, chromedp.NoModifyURL)
		return chromedp.NewContext(ctx)
	}
}
