package main

import "flag"

func main() {
	path := flag.String("file", "urls.txt", "path to URL file")
	flag.Parse()
}
