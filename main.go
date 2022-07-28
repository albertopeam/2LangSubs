package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var s1, s2 string = parseArgs()
	fmt.Println("Input:", s1, "-", s2)
}

func parseArgs() (string, string) {
	var help = "parameter is mandatory. Run help(-h) command for more info"
	var s1, s2 string
	flag.StringVar(&s1, "s1", "", "main subtitle file. mandatory")
	flag.StringVar(&s2, "s2", "", "secondary subtitle file. mandatory")
	flag.Parse()
	if len(s1) == 0 {
		fmt.Println("s1", help)
		os.Exit(1)
	}
	if len(s2) == 0 {
		fmt.Println("s2", help)
		os.Exit(1)
	}
	return s1, s2
}
