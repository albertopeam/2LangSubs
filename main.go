package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/asticode/go-astisub"
)

func main() {
	var f1, f2, fO string = parseArgs()

	s1 := readFromFile(f1)
	s2 := readFromFile(f2)

	sO := mix(s1, s2)

	writeToFile(fO, sO)
}

func parseArgs() (string, string, string) {
	var help = "parameter is mandatory. Run help(-h) command for more info"
	var s1, s2, sO string
	flag.StringVar(&s1, "s1", "", "main subtitle file. mandatory")
	flag.StringVar(&s2, "s2", "", "secondary subtitle file. mandatory")
	flag.StringVar(&sO, "sO", "", "output subtitle file. mandatory")
	flag.Parse()
	if len(s1) == 0 {
		fmt.Println("s1", help)
		os.Exit(1)
	}
	if len(s2) == 0 {
		fmt.Println("s2", help)
		os.Exit(1)
	}
	if len(sO) == 0 {
		fmt.Println("sO", help)
		os.Exit(1)
	}
	fmt.Println("Input", s1, "-", s2, "Output", sO)
	return s1, s2, sO
}

func readFromFile(fn string) *astisub.Subtitles {
	s1, err := astisub.OpenFile(fn)
	if err != nil {
		fmt.Println(fn, "could not be read", err)
		os.Exit(1)
	}
	return s1
}

func mix(s1 *astisub.Subtitles, s2 *astisub.Subtitles) *astisub.Subtitles {
	// crear un astisub.Subtitles para poder añadir entradas
	//TODO: los lengths no cuadran...
	// hay que resolver por tiempo, y los ms no cuadran... definir un tiempo max y hacer
	// busqueda a ver cuantos hacen match. resolver con el más cercano? el que tenga resta menor?
	//fmt.Println(s1.Items[0].Index, s1.Items[0].StartAt, s1.Items[0].EndAt, s1.Items[0].Lines)
	// s1.Items[1085].Lines = append(s1.Items[1085].Lines, s2.Items[1085].Lines...)
	// for i, item := range s1.Items[1085].Lines {
	// 	fmt.Println(i, item)
	// }
	// fmt.Println(s1.Items[1085].Lines)

	//TODO: all this stuff must be injected
	//TODO: all this stuff must be injected
	// diff params
	maxDiffD := time.Duration(500 * time.Millisecond) //max diff aceptable
	//search params
	maxOffsetI := 10 // num of subs to search forwards and backwards
	length := Max(len(s1.Items), len(s2.Items))
	li := LastIndex(s1.Items, s2.Items)
	avgLS := li.EndAt.Seconds() / float64(length)
	fmt.Println("avg line seconds", avgLS)
	maxOffsetS := avgLS * float64(maxOffsetI)
	fmt.Println("max offset seconds", maxOffsetS)
	maxOffsetSInt64 := int64(maxOffsetS) * int64(time.Second)
	maxOffsetD := time.Duration(maxOffsetSInt64) // max offset seconds to search a matching subtitle
	fmt.Println("maxOffsetD", maxOffsetD)
	fmt.Println("maxOffsetI", maxOffsetI)
	fmt.Println("maxDiffD", maxDiffD)
	//TODO: all this stuff must be injected
	//TODO: all this stuff must be injected

	errors := 0
	for _, item := range s1.Items {
		i2Item, err := search(item, s2.Items, maxOffsetD, maxOffsetI, maxDiffD)
		fmt.Println("after")
		if err != nil {
			fmt.Println(err)
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

//TODO: not used maxOffsetD
func search(i1 *astisub.Item, i2 []*astisub.Item, maxOffsetD time.Duration, maxOffsetI int, maxDiffD time.Duration) (*astisub.Item, error) {
	index := i1.Index
	minTime := i1.StartAt - maxOffsetD
	maxTime := i1.StartAt + maxOffsetD
	fmt.Println("index", index)
	fmt.Println("time", i1.StartAt)
	fmt.Println("minTime", minTime)
	fmt.Println("maxTime", maxTime)
	length := len(i2)
	var target *astisub.Item
	for i := 1; i <= maxOffsetI; i++ {
		nextIndex := index - 1 + i
		if nextIndex < length {
			fmt.Println("accesing ni", nextIndex)
			if startAtIsInRange(i1, i2[nextIndex], maxDiffD) {
				target = i2[nextIndex]
				break
			}
		}
		previousIndex := index - 1 - i
		if previousIndex >= 0 && previousIndex < length {
			fmt.Println("accesing pi", previousIndex)
			if startAtIsInRange(i1, i2[previousIndex], maxDiffD) {
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

func startAtIsInRange(i *astisub.Item, t *astisub.Item, diff time.Duration) bool {
	startAtDiff := math.Abs(i.StartAt.Seconds() - t.StartAt.Seconds())
	if startAtDiff < diff.Seconds() {
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
