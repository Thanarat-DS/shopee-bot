package browser

import (
	"context"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func EnableStealth(ctx context.Context) error {

	js := `
/******** webdriver ********/
Object.defineProperty(navigator, 'webdriver', {
	get: () => undefined
});

/******** languages ********/
Object.defineProperty(navigator, 'languages', {
	get: () => ['en-US','en']
});

/******** plugins ********/
Object.defineProperty(navigator, 'plugins', {
	get: () => [1,2,3,4,5]
});

/******** platform ********/
Object.defineProperty(navigator, 'platform', {
	get: () => 'Win32'
});

/******** hardware ********/
Object.defineProperty(navigator, 'hardwareConcurrency', {
	get: () => 8
});

/******** permissions ********/
const originalQuery = window.navigator.permissions.query;
window.navigator.permissions.query = (parameters) => (
	parameters.name === 'notifications'
		? Promise.resolve({ state: Notification.permission })
		: originalQuery(parameters)
);

/******** userAgentData ********/
Object.defineProperty(navigator, 'userAgentData', {
	get: () => ({
		brands: [
			{ brand: "Google Chrome", version: "145" },
			{ brand: "Chromium", version: "145" }
		],
		mobile: false,
		platform: "Windows"
	})
});

/******** connection ********/
Object.defineProperty(navigator, 'connection', {
    get: () => ({
        downlink: 10,
        effectiveType: '4g',
        rtt: 50,
        saveData: false,
        downlinkMax: 100
    })
});

/******** contacts ********/
Object.defineProperty(navigator, 'contacts', {
    get: () => ({
        select: () => Promise.resolve([])
    })
});

/******** contentIndex ********/
Object.defineProperty(navigator, 'contentIndex', {
    get: () => ({
        add: () => Promise.resolve(),
        delete: () => Promise.resolve(),
        getAll: () => Promise.resolve([])
    })
});

/******** WebGL spoof ********/
const getParameter = WebGLRenderingContext.prototype.getParameter;
WebGLRenderingContext.prototype.getParameter = function(parameter) {

	if (parameter === 37445) {
		return "Intel Inc.";
	}

	if (parameter === 37446) {
		return "Intel Iris OpenGL Engine";
	}

	return getParameter(parameter);
};

/******** screen ********/
Object.defineProperty(screen, 'width', {get: () => 1920});
Object.defineProperty(screen, 'height', {get: () => 1080});
Object.defineProperty(screen, 'availWidth', {get: () => 1920});
Object.defineProperty(screen, 'availHeight', {get: () => 1040});

/******** DOMRect patch ********/
const originalGetBoundingClientRect = Element.prototype.getBoundingClientRect;

Element.prototype.getBoundingClientRect = function() {
	const rect = originalGetBoundingClientRect.apply(this, arguments);

	return {
		x: rect.x,
		y: rect.y,
		width: rect.width,
		height: rect.height,
		top: rect.top,
		right: rect.right,
		bottom: rect.bottom,
		left: rect.left
	};
};

/******** chrome runtime ********/
window.chrome = {
	runtime: {}
};
`

	return chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, err := page.AddScriptToEvaluateOnNewDocument(js).Do(ctx)
			return err
		}),
	)
}
