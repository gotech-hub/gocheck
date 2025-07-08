package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var path string
	var htmlOutput bool
	var jsonOutput bool
	flag.StringVar(&path, "path", ".", "Path to scan")
	flag.BoolVar(&htmlOutput, "html", true, "Generate HTML report")
	flag.BoolVar(&jsonOutput, "json", true, "Generate JSON report")
	flag.Parse()

	err := Scan(path, htmlOutput, jsonOutput)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
