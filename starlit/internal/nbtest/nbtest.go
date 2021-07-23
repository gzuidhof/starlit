package nbtest

import (
	"context"
	"log"

	"github.com/chromedp/chromedp"
	"github.com/gzuidhof/starlit/starlit/internal/server"
)


func TestPath(testPath string, stayAlive bool) {
	serverUrl := server.StartNBTestServer(testPath)
	target := serverUrl + "/nbtest/example_content/yaml.sb"

	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// run task list
	var res []string
	err := chromedp.Run(ctx,
		chromedp.Navigate(target),
		chromedp.Evaluate(`Object.keys(window);`, &res),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("window object keys: %v", res)

	if (stayAlive) {
		done := make(chan bool)
		log.Printf("Serving nbtest on %s/nbtest/", serverUrl)
		<- done
	}
}