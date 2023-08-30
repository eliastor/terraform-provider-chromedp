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
		ctx, _ := chromedp.NewExecAllocator(parentCtx, append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Env("POWEREDBY=eliastor"))...)
		return chromedp.NewContext(ctx)
	}
}

func chromedpCtxWithRemoteChrome(remote string) ctxCreatorFunc {
	return func(parentCtx context.Context) (context.Context, context.CancelFunc) {
		ctx, _ := chromedp.NewRemoteAllocator(parentCtx, remote, chromedp.NoModifyURL)
		return chromedp.NewContext(ctx)
	}
}
