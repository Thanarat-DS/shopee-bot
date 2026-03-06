package browser

import (
	"context"
	"log"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

// New initializes a new browser context that can be reused for workers.
func New(headless bool, debug bool) (context.Context, context.CancelFunc) {
	opts := chromedp.DefaultExecAllocatorOptions[:]

	// Anti-detection options
	opts = append(opts,
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.Flag("disable-infobars", true),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36"),
	)

	if !debug {
		if headless {
			opts = append(opts, chromedp.Flag("headless", "new"))
		} else {
			opts = append(opts, chromedp.Flag("headless", false))
		}
		// Additionally ignore certificate errors and other basic stuff for perf
		opts = append(opts, chromedp.Flag("ignore-certificate-errors", true))
		opts = append(opts, chromedp.Flag("blink-settings", "imagesEnabled=false"))
	} else {
		opts = append(opts, chromedp.Flag("headless", false))
	}

	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)

	var ctx context.Context
	var cancel context.CancelFunc

	if debug {
		ctx, cancel = chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	} else {
		ctx, cancel = chromedp.NewContext(allocCtx)
	}

	return ctx, func() {
		cancel()
		allocCancel()
	}
}

// EnableResourceBlocking sets up CDP network blocked URLs on the given target context
func EnableResourceBlocking(ctx context.Context) error {
	// We run it as an Action to set blocked URLs for the current context (tab)
	return chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		return network.SetBlockedURLs([]string{
			"*.jpg", "*.jpeg", "*.png", "*.gif", "*.webp", "*.svg", "data:image/*",
			"*.woff", "*.woff2", "*.ttf", "*.eot",
			"*.mp4", "*.webm", "*.mp3", "*.wav",
		}).Do(ctx)
	}))
}
