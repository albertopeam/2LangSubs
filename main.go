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
	// error message to be displayed when a subtitle hasn't been found
	errMsg string
	// max percentage of errors allowed to consider a successfull mix of subtitles. inclusive
	maxPercentageErrAllowed int
	// shows a custom divider between the subtitles
	divider string
}

func main() {
	args := parseArgs()

	s1 := readFromFile(args.f1)
	s2 := readFromFile(args.f2)

	numErr, percErr := mix(s1, s2, args.tolerance, args.searchOffset, args.errMsg, args.divider)

	if percErr <= float64(args.maxPercentageErrAllowed) {
		writeToFile(args.output, s1)
		fmt.Println("succesfully write subtitles to", args.output)
		fmt.Printf("total Errors %v, percentage of errors %v\n", numErr, percErr)
	} else {
		fmt.Printf("percentage of errors %v exceeds the limit percentage %v\n", numErr, args.maxPercentageErrAllowed)
		os.Exit(1)
	}
}

func parseArgs() args {
	var help = "parameter is mandatory. Run help(-h) command for more info"
	var f1, f2, fO, errMsg, divider string
	var t, sO, mpea int
	flag.StringVar(&f1, "s1", "", "main subtitle file. mandatory")
	flag.StringVar(&f2, "s2", "", "secondary subtitle file. mandatory")
	flag.StringVar(&fO, "sO", "", "output subtitle file. mandatory")
	flag.IntVar(&t, "t", 500, "max tolerance for subtitle start time")
	flag.IntVar(&sO, "so", 10, "search max offset, number of adjacent subtitles to be analyzed, forwards and backwards")
	flag.IntVar(&mpea, "mpe", 5, "max percentage of errors allowed to consider a successfull mix of subtitles")
	flag.StringVar(&errMsg, "em", "%- not found subtitle -%", "error message to be displayed when a subtitle hasn't been found")
	flag.StringVar(&divider, "d", "", "shows a custom divider between the subtitles")
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
	return args{
		f1:                      f1,
		f2:                      f2,
		output:                  fO,
		tolerance:               tolerance,
		searchOffset:            sO,
		errMsg:                  errMsg,
		maxPercentageErrAllowed: mpea,
		divider:                 divider,
	}
}

func readFromFile(fn string) *astisub.Subtitles {
	s1, err := astisub.OpenFile(fn)
	if err != nil {
		fmt.Println(fn, "could not be read", err)
		os.Exit(1)
	}
	return s1
}

// mutates s1 mixing the contents of s2
// input: s1 and s2 represnting subtitles for two languages,
// tolerance that the subs can have in StartAt,
// searchOffset the number of subtitles to search for both forward and backward
// errMsg the error to be displayed if subtitle can't be found
// output: the number of errors over the total length and the percent that it represents
func mix(s1 *astisub.Subtitles, s2 *astisub.Subtitles, tolerance time.Duration, searchOffset int, errMsg string, divider string) (int, float64) {
	length := Max(len(s1.Items), len(s2.Items))
	numErr := 0
	for _, item := range s1.Items {
		i2Item, err := search(item, s2.Items, searchOffset, tolerance)
		if err != nil {
			numErr += 1
			item.Lines = append(item.Lines, createLine(errMsg))
		} else {
			if len(divider) > 0 {
				item.Lines = append(item.Lines, createLine(divider))
			}
			item.Lines = append(item.Lines, i2Item.Lines...)
		}
	}
	percErr := float64(numErr) / float64(length) * 100
	return numErr, percErr
}

func createLine(t string) astisub.Line {
	li := astisub.LineItem{Text: t}
	s := []astisub.LineItem{li}
	return astisub.Line{Items: s}
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
