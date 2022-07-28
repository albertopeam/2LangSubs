# 2LangSubs

* Script to create a single subtitle file from two subtitles in two different languages

## Specs

* Combine subtitles from diff languagues
* Pick the order of the main(top) and secondary(bottom) subtitle
* Select a divider to quickly distinguish where both start
* Extend subtitles time to let user time to read. It is a ratio of the window time to avoiding collisions
* ? Configurable error rate: handles situations where are diffs in timing with subs or some parts have not been translated

## Usage

* help: go run main.go -h
* usage: go run main.go -s1 file1 -s2 file2

## Tasks

* Read from files and structure data sets
* Algorythm to join both data sets
* Investigate if divider occupies to much space. Implement if not
* Algorythm to extend subtitles for a ratio or crop til collision
* Investigate subs with diffs in subtitles and find an algorythm capable of handle it.

## Assets

* [YIFY](https://yifysubtitles.org/movie-imdb/tt0800369)
