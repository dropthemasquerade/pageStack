// Command screenshot is a chromedp example demonstrating how to take a
// screenshot of a specific element and of the entire browser viewport.
package main

import (
	"context"
	"io/ioutil"
	"log"
	"math"
  "fmt"
  "sync"
  "gopkg.in/yaml.v2"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// elementScreenshot takes a screenshot of a specific element.
func elementScreenshot(urlstr, sel string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.WaitVisible(sel, chromedp.ByID),
		chromedp.Screenshot(sel, res, chromedp.NodeVisible, chromedp.ByID),
	}
}

// fullScreenshot takes a screenshot of the entire browser viewport.
//
// Liberally copied from puppeteer's source.
//
// Note: this will override the viewport emulation settings.
func fullScreenshot(urlstr string, quality int64, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
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

func dispatchAction(steps Steps, wg *sync.WaitGroup){
  defer wg.Done()
  fmt.Printf("Worker %d starting\n")
  for i, v := range steps.Steps {
    fmt.Println("each step>>> ",i,  v)
    selectAction(v)
  }
  fmt.Printf("Worker %d done\n")
}

func selectAction(step Step){
  switch step.Cmd {
  case "click":
      println("i is click")
  case "screen":
      println("i is screen")
  case "getText":
      println("i is getText")
  default:
      println("type not found")
  }
}

func doAction(){
  // create context
  ctx, cancel := chromedp.NewContext(context.Background())
  defer cancel()

  // capture screenshot of an element
  var buf []byte
  if err := chromedp.Run(ctx, elementScreenshot(`https://www.google.com/`, `#main`, &buf)); err != nil {
    log.Fatal(err)
  }
  if err := ioutil.WriteFile("elementScreenshot.png", buf, 0644); err != nil {
    log.Fatal(err)
  }

  // capture entire browser viewport, returning png with quality=90
  if err := chromedp.Run(ctx, fullScreenshot(`https://brank.as/`, 90, &buf)); err != nil {
    log.Fatal(err)
  }
  if err := ioutil.WriteFile("fullScreenshot.png", buf, 0644); err != nil {
    log.Fatal(err)
  }
}

type Step struct {
  Cmd string `yaml:"cmd"`
  Location string `yaml:"location"`
  Value string `yaml:"value,omitempty"`
  Name string `yaml:"name,omitempty"`
  Desc string `yaml:"desc,omitempty"`
}

type Steps struct {
  Version string `yaml:"version"`
  GroupName string `yaml:"groupName"`
  Entrance string `yaml:"entrance,omitempty"`
  Steps []Step `yaml:"steps"`
}

func getSteps(d string) {
  // filter files by configure the case here
  // case about your product, which you want to test.
  files, err := ioutil.ReadDir(d)
  if err != nil {
      log.Fatal(err)
  }
  var wg sync.WaitGroup
  for _, f := range files {
    // fmt.Println("lll", f)
    var c Steps
    p := "./outputs/" + f.Name()
    yamlFile, err := ioutil.ReadFile(p)
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }
    // fmt.Println("~~~", yamlFile)
    err = yaml.Unmarshal(yamlFile, &c)
    if err != nil {
        log.Fatalf("Unmarshal: %v", err)
    }
    wg.Add(1)
    go dispatchAction(c, &wg)
  }
  wg.Wait()
}

func main() {
  // doAction()
  getSteps("./outputs/")
  // dispatchAction(steps)
  // selectAction()

}