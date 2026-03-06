# Work in progress
There are still many things to do.

# SHIT Bot – Shopee Hyper Instant Trigger Bot

**Tagline:** "Deals drop. SHIT hits. Items gone." 🦶💥

## Overview
SHIT is a hyper-fast bot for grabbing Shopee deals instantly.  
Never miss a flash sale, hot deal, or limited stock again.

## Features

- **Headless Chrome Automation:** Uses chromedp to control a Chrome instance via CDP.
- **Preloaded Cookies:** Bypass manual login by providing a preloaded session.
- **Concurrent Workers:** Uses Goroutines to create multiple tabs simultaneously attempting to buy the target item.
- **Performance Optimized:** Optionally disables images, fonts, and media downloads to speed up page loads.
- **Browser Re-use:** A single Chrome instance is launched and orchestrates multiple isolated browser contexts (tabs).

## Getting Started

### Prerequisites

- Go 1.20+
- A Chromium-based browser installed on your machine.

### Usage

Run the compiled executable:

```bash
go run ./cmd/main.go
```
OR build .exe file
``` build .exe file
go build -o shopee-bot.exe ./cmd/main.go
```

### Bot detection
https://abrahamjuliot.github.io/creepjs/

### Configuration (`config.json`)

- `product_url`: The Shopee URL of the item you want to buy.
- `worker_count`: Number of simultaneous buyer threads.
- `headless`: Whether to run Chrome without a GUI.
- `debug`: Enable this to see Chrome visibly (forces `headless: false`).
- `timeout_seconds`: Global timeout for navigation and buying actions.
- `selectors`: CSS rules used for finding the buy out buttons locally.

Enjoy fast checking out!
