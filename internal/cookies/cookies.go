package cookies

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

type Cookie struct {
	Name     string  `json:"name"`
	Value    string  `json:"value"`
	Domain   string  `json:"domain"`
	Path     string  `json:"path"`
	Expires  float64 `json:"expires"`
	HTTPOnly bool    `json:"httpOnly"`
	Secure   bool    `json:"secure"`
}

func Load(path string) ([]Cookie, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cookies []Cookie
	if err := json.Unmarshal(data, &cookies); err != nil {
		return nil, err
	}
	return cookies, nil
}

func Inject(ctx context.Context, cookies []Cookie) error {
	return chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		for _, c := range cookies {
			expr := cdp.TimeSinceEpoch(time.Unix(int64(c.Expires), 0))
			err := network.SetCookie(c.Name, c.Value).
				WithDomain(c.Domain).
				WithPath(c.Path).
				WithExpires(&expr).
				WithHTTPOnly(c.HTTPOnly).
				WithSecure(c.Secure).
				Do(ctx)
			if err != nil {
				return err
			}
		}
		return nil
	}))
}
