package nbtest

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/fatih/color"
)

type JSError struct {
	Message string `json:"message"`
	Stack string `json:"stack"`
	LineNumber int `json:"lineNumber"`
	ColumnNumber int `json:"columnNumber"`
	FileName string `json:"fileName"`
}

type TestResult struct {
	Pass bool `json:"pass"`
	Error JSError `json:"error"`
}

type TestRunner struct {
	timeout time.Duration
	testPath string
	ctx context.Context 

	cancelCtx context.CancelFunc
}

func NewTestRunner(testPath string, headless bool, timeout time.Duration) *TestRunner {
	opts := defaultNBTestExecAllocatorOptions[:]

	// if (headless) {
	// 	opts = append(opts, chromedp.Headless)
	// }

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)

	taskCtx, _ := chromedp.NewContext(allocCtx)

	// ensure the first tab is created (this way the browser doesn't keep getting closed)
	if err := chromedp.Run(taskCtx); err != nil {
		panic(err)
	}

	return &TestRunner {
		timeout: timeout,
		testPath: testPath,
		ctx: taskCtx,
		cancelCtx: cancel,
	}
}

func (r *TestRunner) printError(name string, message string, err error) error {
	if errors.Is(err, context.DeadlineExceeded) {
		fmt.Fprintf(
			color.Output,
			"%s %s - %s %s\n",
			color.HiRedString("FAIL"),
			name,
			color.YellowString((fmt.Sprintf("Timeout exceeded (%s)", r.timeout))),
			message,
		)
		return fmt.Errorf("%s: timeout exceeded %w", message, err)
	}

	return fmt.Errorf("%s: %w", message, err)
}

func (r *TestRunner) testNotebook(targetURL string, name string) error {
	taskCtx, cancel := chromedp.NewContext(r.ctx)
	defer cancel()

	ctx, cancel := context.WithTimeout(taskCtx, r.timeout)
	defer cancel()

	err := chromedp.Run(ctx, chromedp.Navigate(targetURL))
	if err != nil {
		log.Fatal(err)
	}

	if err := chromedp.Run(ctx, chromedp.WaitReady("starboard-notebook")); err != nil {
		return r.printError(name, "waiting for starboard-notebook element", err)
	}

	if err := chromedp.Run(ctx, chromedp.WaitReady(".nbtest-started")); err != nil {
		return r.printError(name, "waiting to start running all cells", err)
	}

	if err := chromedp.Run(ctx, chromedp.WaitReady(".nbtest-done")); err != nil {
		return r.printError(name, "waiting for notebook to be run completely", err)
	}

	// // get project link text
	var testResult TestResult
	if err := chromedp.Run(ctx, chromedp.Evaluate("window.__nbTestResult", &testResult)); err != nil {
		return r.printError(name, "retrieving test result from browser", err)
	}

	if (testResult.Pass) {
		fmt.Fprintf(color.Output, "%s %s\n", color.GreenString("PASS"), name)
	} else {
		fmt.Fprintf(color.Output, "%s %s %s\n", color.HiRedString("FAIL"), name, color.RedString(testResult.Error.Message))
		return fmt.Errorf("nbtest fail")
	}

	return nil
}
