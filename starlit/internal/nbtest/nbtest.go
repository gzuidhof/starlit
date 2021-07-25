package nbtest

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/gzuidhof/starlit/starlit/internal/server"
)

func TestPath(testPath string, serveMode bool, headless bool) {
	filesToTest := GatherFilesToTest(testPath)

	if len(filesToTest) == 0 {
		log.Printf("No notebook files found under %s", testPath)
		os.Exit(0)
	}

	serverUrl := server.StartNBTestServer(testPath)
	runner := NewTestRunner(testPath, headless, time.Second * 15)

	hadError := false

	for _, filepath := range filesToTest {
		targetURL := serverUrl + "/nbtest/" + filepath
		err := runner.testNotebook(targetURL, filepath)

		if (err != nil) {
			if (serveMode) {
				fmt.Fprintf(color.Output, "%s %s\n", color.HiBlackString("^^^^"), color.HiBlackString(targetURL))

			}
			hadError = true
		}
	}

	// fmt.Println("Done")
	if (serveMode) {
		done := make(chan bool)
		log.Printf("Serving nbtest on %s/nbtest/", serverUrl)
		<- done
	}

	runner.cancelCtx()
	if (hadError) {
		os.Exit(1)
	}
}