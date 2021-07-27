/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */
package nbtest

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/spf13/viper"
)

type JSError struct {
	Message string `json:"message"`
	Stack string `json:"stack"`
	LineNumber int `json:"lineNumber"`
	ColumnNumber int `json:"columnNumber"`
	FileName string `json:"fileName"`
}

type TestResult struct {
	Name string
	URL string

	// "PASS"" | "FAIL" | "SKIP"
	Status string
	Message string

	Timing time.Duration
	InternalError error
}

type browserTestResult struct {
	// "PASS"" | "FAIL" | "SKIP"
	Status string `json:"status"`
	Error JSError `json:"error"`
}

type TestRunner struct {
	timeout time.Duration
	serveMode bool

	testPath string
	ctx context.Context 

	cancelCtx context.CancelFunc
}

func NewTestRunner(testPath string, headless bool, serveMode bool, timeout time.Duration) *TestRunner {
	opts := defaultNBTestExecAllocatorOptions[:]

	if (headless) {
		opts = append(opts, chromedp.Headless)
	}
	execPath := viper.GetString("nbtest.browser_exec_path")
	if execPath != "" {
		opts = append(opts, chromedp.ExecPath(execPath))
	}
	
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	taskCtx, _ := chromedp.NewContext(allocCtx)

	// ensure the first tab is created (this way the browser doesn't keep getting closed)
	if err := chromedp.Run(taskCtx); err != nil {
		panic(err)
	}

	return &TestRunner {
		timeout: timeout,
		serveMode: serveMode,
		testPath: testPath,
		ctx: taskCtx,
		cancelCtx: cancel,
	}
}

func (r *TestRunner) testNotebook(targetURL string, name string) *TestResult {
	taskCtx, cancel := chromedp.NewContext(r.ctx)
	defer cancel()

	ctx, cancel := context.WithTimeout(taskCtx, r.timeout)
	defer cancel()

	tr := &TestResult{
		URL: targetURL,
		Name: name,
		Status: "FAIL", // We overwrite it in the other cases
	}
	defer func (t time.Time) {
		tr.Timing = time.Since(t)
	}(time.Now())

	err := chromedp.Run(ctx, chromedp.Navigate(targetURL))
	if err != nil {
		tr.InternalError = err;
		tr.Message = "waiting for browser to open page"
		return tr
	}

	if err := chromedp.Run(ctx, chromedp.WaitReady(".nbtest-started")); err != nil {
		tr.InternalError = err;
		tr.Message = "waiting to start running all cells"
		return tr
	}

	if err := chromedp.Run(ctx, chromedp.WaitReady(".nbtest-done")); err != nil {
		tr.InternalError = err;
		tr.Message = "waiting for notebook to be run completely"
		return tr
	}

	var testResult browserTestResult
	if err := chromedp.Run(ctx, chromedp.Evaluate("window.__nbTestResult", &testResult)); err != nil {
		tr.InternalError = err;
		tr.Message = "retrieving test result from browser"
		return tr
	}

	if (testResult.Status == "PASS") {
		tr.Status = "PASS"
		return tr
	} else if (testResult.Status == "FAIL"){
		tr.Status = "FAIL"
		tr.Message = testResult.Error.Message
		return tr
	} else if (testResult.Status == "SKIP") {
		tr.Status = "SKIP"
		return tr
	} else { // Should never happen
		tr.InternalError = fmt.Errorf("unknown nbtest result status: %s", testResult.Status)
		tr.Message = "Unkown nbtest result status from browser message"
		return tr
	}
}
