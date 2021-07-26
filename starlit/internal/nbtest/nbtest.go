/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */
package nbtest

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/gzuidhof/starlit/starlit/internal/fs/assetfs"
	"github.com/gzuidhof/starlit/starlit/internal/server"
	"github.com/spf13/viper"
	"github.com/xxjwxc/gowp/workpool"
)

func getConcurrency() int {	
	con := viper.GetInt("nbtest.concurrency")

	if (con < 0) {
		log.Fatalf("Concurrency must be positive, or zero for number of cores")
	}

	if con == 0 {
		con = runtime.NumCPU()
	}

	return con
}

func getTimeout() time.Duration {
	timeoutOpt := viper.GetFloat64("nbtest.timeout")

	if (timeoutOpt < 0) {
		log.Fatalf("Timeout must be positive, or zero for no timeout")
	}

	if (timeoutOpt == 0) {
		return time.Hour * 24 * 10 // Not actually infinite.. But 10 days is a pretty long time.
	}

	return time.Duration(float64(time.Second) * timeoutOpt)
}

// getStarboardFolder returns the folder to serve artifacts from, or ""
// in case no such folder was passed
func getArtifactsFolder(viperConfigKey string) string {
	f := viper.GetString(viperConfigKey)
	if (f == "") {
		return "";
	}

	if strings.HasPrefix(f, "http://") || strings.HasPrefix(f, "https://") {
		return ""
	}

	stat, err := os.Stat(f)

	if err != nil {
		log.Fatalf("Invalid folder for %s: %v", viperConfigKey, err)
	}

	if (!stat.IsDir()) {
		log.Fatalf("Invalid folder for %s: not a directory", viperConfigKey)
	}

	return f
}

func TestPath(testPath string) {
	filesToTest := GatherFilesToTest(testPath)

	if len(filesToTest) == 0 {
		log.Printf("No notebook files found under %s", testPath)
		os.Exit(0)
	}

	port := viper.GetString("server.port")
	viper.Set("serve_filepath", testPath)
	serveFS := assetfs.GetAssetFileSystems()
	serverURL := fmt.Sprintf("http://localhost:%s", port)

	go func() {
		log.Fatal(server.StartNBTestServer(testPath, serveFS, port, getArtifactsFolder("nbtest.starboard_artifacts"), getArtifactsFolder("nbtest.pyodide_artifacts")))
	}()
	

	runner := NewTestRunner(testPath, viper.GetBool("nbtest.headless"), viper.GetBool("nbtest.serve"), getTimeout())
	serveMode := viper.GetBool("nbtest.serve")

	if serveMode {
		fmt.Fprintf(color.Output, "%s %s\n", color.HiBlackString("Serving nbtest on"), color.BlueString(serverURL + "/nbtest/"))
	}

	hadError := false
	start := time.Now()

	wp := workpool.New(getConcurrency())
	for _, p := range filesToTest {
		filepath := p
		wp.Do(func() error {
			targetURL := serverURL + "/nbtest/" + filepath
			result := runner.testNotebook(targetURL, filepath)
			if (result.Status == "FAIL") {
				hadError = true
			}

			runner.PrintTestResult(result)
			return nil
		})
	}
	wp.Wait()
	runner.cancelCtx()

	timing := color.HiBlackString(fmt.Sprintf("(%s)", time.Since(start)))
	if (hadError) {
		fmt.Fprintf(color.Output, "\n%s %s\n", color.RedString("Done testing, one or more tests failed"), timing )
	} else {
		fmt.Fprintf(color.Output, "\n%s %s\n", color.CyanString("Done testing"), timing )
	}

	if (serveMode) {
		done := make(chan bool)
		
		<- done
	}
	
	if (hadError) {
		os.Exit(1)
	}
}