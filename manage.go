package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"golang.org/x/net/html"
)

func load(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	defer file.Close()

	node, err := html.Parse(file)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}

	panic("Not implemented.")
}

func add(tag string) {
	index, err := load("index.html")
	if err != nil {
		log.Printf("Error: load index: %v", err)
	}

	html.Render(os.Stdout, index)
}

func main() {
	log.SetFlags(0)

	mode := flag.String("mode", "add", "mode to run in (supported: add)")
	tag := flag.String("tag", "", "name of tag to edit")
	flag.Parse()

	switch *mode {
	case "add":
		add(*tag)
	default:
		log.Printf("Error: unknown mode: %q", *mode)
		flag.Usage()
		os.Exit(2)
	}
}
