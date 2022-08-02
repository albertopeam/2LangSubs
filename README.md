# 2LangSubs

* Script to create a single subtitle file from two subtitles in two different languages

## Specs

* Combine subtitles from diff languagues
* Pick the order of the main(top) and secondary(bottom) subtitle
* Select a divider to quickly distinguish where both start
* Extend subtitles time to let user time to read. It is a ratio of the window time to avoiding collisions
* ? Configurable error rate: handles situations where are diffs in timing with subs or some parts have not been translated

## Usage

* help: go run main.go math.go slice.go -h
* usage: go run main.go -s1 file1 -s2 file2
  * ie: go run main.go math.go slice.go -s1 testAssets/Thr-Spanish.srt -s2 testAssets/Thr.2011.720p.BrRip.264.English.srt -sO testAssets/output.srt -t 500 -so 10 -em "-not found subtitle-" -mpe 5 -d "*-*"

## Tasks

* Read from files and structure data sets. OK
* Algorythm to join both data sets. OK
* Investigate subs with diffs in subtitles and find an algorythm capable of handle it. OK
* Investigate if divider occupies to much space. Implement if not
* Algorythm to extend subtitles for a ratio or crop til collision

## Improvements

* Split into diff files
* Test
* Use alias for astisub.Item: `type Item astisub.Item`

## Assets

* [YIFY](https://yifysubtitles.org/movie-imdb/tt0800369)
