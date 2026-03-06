package worker

import (
	"context"
	"log"
	"sync"
	"time"

	"shopee-bot/internal/browser"
	"shopee-bot/internal/config"

	"github.com/chromedp/chromedp"
)

type Worker struct {
	ID      int
	Config  *config.Config
	Browser context.Context
}

// Run executes the worker's logic to purchase the product.
func (w *Worker) Run(startSignal <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	// Wait for the synchronous start signal
	<-startSignal

	log.Printf("Worker %d: Starting attempt", w.ID)

	// Create a new tab (target) context attached to the main browser allocator
	tabCtx, cancelTab := chromedp.NewContext(w.Browser)
	defer cancelTab()

	// Optionally apply resource blocking for this tab
	if err := browser.EnableResourceBlocking(tabCtx); err != nil {
		log.Printf("Worker %d: Failed to enable resource blocking: %v", w.ID, err)
	}

	// Apply stealth injection for this tab
	if err := browser.EnableStealth(tabCtx); err != nil {
		log.Printf("Worker %d: Failed to inject stealth js: %v", w.ID, err)
	}

	timeoutCtx, cancelTimeout := context.WithTimeout(tabCtx, time.Duration(w.Config.TimeoutSeconds)*time.Second)
	defer cancelTimeout()

	// Record the time to measure performance metrics
	start := time.Now()

	log.Printf("Worker %d: Navigating to %s", w.ID, w.Config.ProductURL)

	err := chromedp.Run(timeoutCtx,
		chromedp.EmulateViewport(1920, 1080),
		chromedp.Navigate(w.Config.ProductURL),
		chromedp.WaitVisible(w.Config.Selectors.BuyButton, chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			loadToVisible := time.Since(start)
			log.Printf("Worker %d: Page loaded & Buy button visible in %v", w.ID, loadToVisible)
			return nil
		}),
		chromedp.Click(w.Config.Selectors.BuyButton, chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			loadToClick := time.Since(start)
			log.Printf("Worker %d: Successfully clicked Buy button in %v", w.ID, loadToClick)
			return nil
		}),
		// Depending on Shopee logic, you might wait for checkout button
		// chromedp.WaitVisible(w.Config.Selectors.CheckoutButton, chromedp.ByQuery),
	)

	if err != nil {
		log.Printf("Worker %d: Failed on main attempt: %v", w.ID, err)
		w.retry(timeoutCtx, 1) // Retry once
	} else {
		log.Printf("Worker %d: Completed checkout initiation", w.ID)
	}
}

func (w *Worker) retry(ctx context.Context, attempts int) {
	for i := 0; i < attempts; i++ {
		log.Printf("Worker %d: Retrying attempt %d", w.ID, i+1)

		err := chromedp.Run(ctx,
			chromedp.EmulateViewport(1920, 1080),
			chromedp.Navigate(w.Config.ProductURL),
			chromedp.WaitVisible(w.Config.Selectors.BuyButton, chromedp.ByQuery),
			chromedp.Click(w.Config.Selectors.BuyButton, chromedp.ByQuery),
		)

		if err == nil {
			log.Printf("Worker %d: Retry successful", w.ID)
			return
		}
		log.Printf("Worker %d: Retry failed - %v", w.ID, err)
	}
}
