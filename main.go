package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/asticode/go-astisub"
)

type args struct {
	// input file for main lang subtitles
	f1 string
	// input file for secondary lang subtitles
	f2 string
	// output file where mixed subs must be written
	output string
	// tolerance in ms for a valid subtitle
	tolerance time.Duration
	// max search offset to be applied when searching subs
	searchOffset int
}

func main() {
	args := parseArgs()

	s1 := readFromFile(args.f1)
	s2 := readFromFile(args.f2)

	sO := mix(s1, s2, args.tolerance, args.searchOffset)

	writeToFile(args.output, sO)
}

func parseArgs() args {
	var help = "parameter is mandatory. Run help(-h) command for more info"
	var f1, f2, fO string
	var t, sO int
	flag.StringVar(&f1, "s1", "", "main subtitle file. mandatory")
	flag.StringVar(&f2, "s2", "", "secondary subtitle file. mandatory")
	flag.StringVar(&fO, "sO", "", "output subtitle file. mandatory")
	flag.IntVar(&t, "t", 500, "max tolerance for subtitle start time")
	flag.IntVar(&sO, "so", 10, "search max offset, number of adjacent subtitles to be analyzed")
	flag.Parse()
	if len(f1) == 0 {
		fmt.Println("s1", help)
		os.Exit(1)
	}
	if len(f2) == 0 {
		fmt.Println("s2", help)
		os.Exit(1)
	}
	if len(fO) == 0 {
		fmt.Println("sO", help)
		os.Exit(1)
	}
	tolerance := time.Duration(time.Duration(t) * time.Millisecond)
	return args{f1: f1, f2: f2, output: fO, tolerance: tolerance, searchOffset: sO}
}

func readFromFile(fn string) *astisub.Subtitles {
	s1, err := astisub.OpenFile(fn)
	if err != nil {
		fmt.Println(fn, "could not be read", err)
		os.Exit(1)
	}
	return s1
}

func mix(s1 *astisub.Subtitles, s2 *astisub.Subtitles, tolerance time.Duration, searchOffset int) *astisub.Subtitles {
	length := Max(len(s1.Items), len(s2.Items))
	errors := 0
	for _, item := range s1.Items {
		i2Item, err := search(item, s2.Items, searchOffset, tolerance)
		if err != nil {
			errors += 1
			notFoundItem := astisub.LineItem{Text: "%- not found -%"} //TODO: review /inject this message
			notFoundItems := []astisub.LineItem{notFoundItem}
			notFoundLine := astisub.Line{Items: notFoundItems}
			item.Lines = append(item.Lines, notFoundLine)
		} else {
			item.Lines = append(item.Lines, i2Item.Lines...)
		}
	}

	fmt.Println("total Errors", errors)
	fmt.Println("total errors %", float64(errors)/float64(length)*100)

	return s1
}

func search(i1 *astisub.Item, i2 []*astisub.Item, searchOffset int, tolerance time.Duration) (*astisub.Item, error) {
	index := i1.Index
	length := len(i2)
	var target *astisub.Item
	for i := 0; i < searchOffset; i++ {
		nextIndex := index - 1 + i
		if nextIndex < length {
			if startAtIsInRange(i1, i2[nextIndex], tolerance) {
				target = i2[nextIndex]
				break
			}
		}
		previousIndex := index - 1 - i
		if previousIndex >= 0 && previousIndex < length {
			if startAtIsInRange(i1, i2[previousIndex], tolerance) {
				target = i2[previousIndex]
				break
			}
		}
	}
	if target != nil {
		return target, nil
	} else {
		return i1, fmt.Errorf("search failed for index %v %v", i1.Index, i1.Lines)
	}
}

func startAtIsInRange(i *astisub.Item, t *astisub.Item, tolerance time.Duration) bool {
	startAtDiff := math.Abs(i.StartAt.Seconds() - t.StartAt.Seconds())
	if startAtDiff < tolerance.Seconds() {
		return true
	} else {
		return false
	}
}

func writeToFile(fn string, sO *astisub.Subtitles) {
	err := sO.Write(fn)
	if err != nil {
		fmt.Println(fn, "could not be written", err)
		os.Exit(1)
	}
}
