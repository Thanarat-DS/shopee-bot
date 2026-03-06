package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"shopee-bot/internal/browser"
	"shopee-bot/internal/config"
	"shopee-bot/internal/cookies"
	"shopee-bot/internal/worker"

	"github.com/chromedp/chromedp"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.Println("Starting Shopee Bot...")

	// 1. Load config
	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatalf("Failed to load config.json: %v", err)
	}

	// 2. Load cookies
	cookieList, err := cookies.Load("cookies.json")
	if err != nil {
		log.Fatalf("Failed to load cookies.json: %v", err)
	}
	log.Printf("Loaded %d cookies", len(cookieList))

	// 3. Initialize browser allocator
	browserCtx, cancelBrowser := browser.New(cfg.Headless, cfg.Debug)
	defer cancelBrowser()

	// Handle Graceful Shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("\nReceived shutdown signal. Closing browser and exiting...")
		cancelBrowser()
		os.Exit(0)
	}()

	// Ensure the browser is started and inject cookies
	if err := chromedp.Run(browserCtx, chromedp.ActionFunc(func(ctx context.Context) error {
		// Wait just to trigger the browser launch natively so we can inject cookies globally
		return nil
	})); err != nil {
		log.Fatalf("Failed to initialize browser: %v", err)
	}

	log.Println("Injecting cookies into browser session...")
	if err := cookies.Inject(browserCtx, cookieList); err != nil {
		log.Fatalf("Failed to inject cookies: %v", err)
	}

	// 4. Start workers
	log.Printf("Initializing %d workers...", cfg.WorkerCount)
	var wg sync.WaitGroup
	startSignal := make(chan struct{})

	for i := 1; i <= cfg.WorkerCount; i++ {
		w := &worker.Worker{
			ID:      i,
			Config:  cfg,
			Browser: browserCtx,
		}
		wg.Add(1)
		go w.Run(startSignal, &wg)
	}

	// 5. Worker synchronization - send start signal to all at once
	log.Println("All workers initialized. Starting simultaneous purchase attempt...")
	close(startSignal) // Broadcasts to all workers waiting on <-startSignal

	// 6. Wait for all workers to finish
	wg.Wait()
	log.Println("All workers have completed their tasks. Shutting down gracefully.")
}
