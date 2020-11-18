// Command screenshot is a chromedp example demonstrating how to take a
// screenshot of a specific element and of the entire browser viewport.
package main

import (
	"context"
	"io/ioutil"
	"log"
	"math"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)


func elementScreenshot(urlstr, sel string, res *[]byte) chromedp.Tasks {
  q := "hello Frank"
  // loc := "//input[contains(@title,'搜索')]"
  return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.WaitVisible(sel, chromedp.ByID),
    // if I do some click here ..
    chromedp.SetValue(sel, q, chromedp.ByID),
		// chromedp.Screenshot(sel, res, chromedp.NodeVisible, chromedp.ByID),

    chromedp.ActionFunc(func(ctx context.Context) error {
			// get layout metrics
			_, _, contentSize, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}

			width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

			// force viewport emulation
			err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).
				Do(ctx)
			if err != nil {
				return err
			}
      var quality int64 = 90
			// capture screenshot
			*res, err = page.CaptureScreenshot().
				WithQuality(quality).
				WithClip(&page.Viewport{
					X:      contentSize.X,
					Y:      contentSize.Y,
					Width:  contentSize.Width,
					Height: contentSize.Height,
					Scale:  1,
				}).Do(ctx)
			if err != nil {
				return err
			}
			return nil
		}),

  }
}


func main() {
	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
  // chromedp.WithLog(log.Printf)
  // ctx, cancel := chromedp.NewContext(context.TODO())
	defer cancel()

	// capture screenshot of an element
	var buf []byte
	if err := chromedp.Run(ctx, elementScreenshot(`https://www.baidu.com/`, `#kw`, &buf)); err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile("baidu.png", buf, 0644); err != nil {
		log.Fatal(err)
	}
}
