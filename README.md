# gofind

A script I made to index .txt files in a directory and sub directorys so I can use `cat` and `grep` to search them. I chose Go because 1. I wanted to learn it and 2. indexing multiple files could be parallelized and Go seemed like a good choice.

## How it works

A main thread scans the first level directory and makes a channel for goroutines to dump `word: filename` maps to, then the main thread collects the maps together and writes an index line by line. Words in the index are lowercased and sanitized of any non-alphanumeric start and end characters.
The index file can then be `grep`ed for the word you're looking for.

## Stuff I'd like to do

* [ ] add support for different text document types (pdf, namely, probably docx, odt?)
* [ ] add more terminal options/configuration (thread count, recursion depth, whether the output is sorted)
* [ ] add a way to better search the end result and display it better to the user 
* [ ] reduce memory consumption while indexing, it all dumps into a huge map right now
* [ ] use some naive map-reduce to create the index faster/more distributed
