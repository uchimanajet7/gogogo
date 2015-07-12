// gogogo project main.go
package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
)

const appName = "gogogo"
const appVersion = "v0.0.1"

func showHelp() func() {
	return func() {
		helpText := fmt.Sprintf("\nUsage: %s [options]\n\n", appName)
		helpText += "options:\n"

		options := make(map[string]string)

		// list options
		flag.VisitAll(func(f *flag.Flag) {
			if value, ok := options[f.Usage]; ok {
				// key exists
				if len(f.Name) > len(value) {
					options[f.Usage] = fmt.Sprintf("-%s, --%s", value, f.Name)
				} else {
					options[f.Usage] = fmt.Sprintf("-%s, --%s", f.Name, value)
				}
			} else {
				// key does not exist
				options[f.Usage] = f.Name
			}
		})

		// to store the keys in slice in sorted order
		var keys []string
		textLength := 0
		for k, v := range options {
			keys = append(keys, k)

			if len(v) > textLength {
				textLength = len(v)
			}
		}
		sort.Strings(keys)

		// to prepare the output format
		for _, k := range keys {
			optionText := fmt.Sprintf("%s%s", options[k], strings.Repeat(" ", textLength-len(options[k])))
			helpText += fmt.Sprintf("  %s  %s\n", optionText, k)
		}

		// show a help string
		fmt.Println(helpText)
	}
}

var (
	helpFlag    bool
	versionFlag bool
	loopFlag    int
)

func resolveArgs() int {
	// register flag name
	flag.BoolVar(&helpFlag, "help", false, "show this help message and exit")
	flag.BoolVar(&helpFlag, "h", false, "show this help message and exit")
	flag.BoolVar(&versionFlag, "version", false, "show version message and exit")
	flag.BoolVar(&versionFlag, "v", false, "show version message and exit")
	flag.IntVar(&loopFlag, "loop", 555, "designate number of runs of a loop (default: 555)")
	flag.IntVar(&loopFlag, "l", 555, "designate number of runs of a loop (default: 555)")

	// set help func
	flag.Usage = showHelp()
	flag.Parse()

	// show help
	if helpFlag {
		flag.Usage()
		return 1
	}
	// show version
	if versionFlag {
		fmt.Printf("%s %s\n", appName, appVersion)
		return 1
	}

	return 0
}

func loopProc() {
	var wg sync.WaitGroup

	for i := 0; i < loopFlag; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			// fmt.Printf("  loop count [%d] finished.\n", index)
		}(i)
	}
	wg.Wait()
}

func main() {
	var exCode int
	defer func() { os.Exit(exCode) }()

	// environment variable is set up in order to correspond to multi-core CPU
	if envvar := os.Getenv("GOMAXPROCS"); envvar == "" {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	// resolve command line
	exCode = resolveArgs()
	if exCode != 0 {
		return
	}

	// run loop and output
	stime := time.Now()
	fmt.Printf("[%s] --- start loop --- \n", stime.Format("2006/01/02 15:04:05"))
	fmt.Printf("\n  loop count: %d\n\n", loopFlag)

	loopProc()

	etime := time.Now()
	fmt.Printf("[%s] --- end loop --- \n", etime.Format("2006/01/02 15:04:05"))

	fmt.Printf("\n  total runtime: %v\n", etime.Sub(stime))
}
