package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"slices"
	"strings"
)

var (
	//go:embed tmpl
	tmplFS embed.FS
	tmpl   *template.Template
)

func init() {
	tmpl = template.Must(template.ParseFS(tmplFS, "tmpl/*"))
}

func loadTags() (tags, ignore []string) {
	ignore = []string{"tags.txt"}

	file, err := os.Open("tags.txt")
	if err != nil {
		log.Printf("Error: open tags file: %v", err)
		os.Exit(1)
	}
	defer file.Close()

	s := bufio.NewScanner(file)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "!") {
			ignore = append(ignore, line[1:])
			continue
		}
		tags = append(tags, line)
	}
	if err := s.Err(); err != nil {
		log.Printf("Error: scan tags file: %v", err)
		os.Exit(1)
	}

	return tags, ignore
}

func clean(ignore []string) {
	entries, err := os.ReadDir(".")
	if err != nil {
		log.Printf("Error: read entries from current directory: %v", err)
		os.Exit(1)
	}

	for _, entry := range entries {
		if entry.IsDir() || slices.Contains(ignore, entry.Name()) {
			fmt.Printf("ignore: %v\n", entry.Name())
			continue
		}

		log.Printf("clean: %v\n", entry.Name())
		err := os.Remove(entry.Name())
		if err != nil {
			log.Printf("Error: remove %q: %v", entry.Name(), err)
			os.Exit(1)
		}
	}
}

func generateIndex(tags []string) {
	file, err := os.Create("index.html")
	if err != nil {
		log.Printf("Error: create index: %v", err)
		os.Exit(1)
	}
	defer file.Close()

	fmt.Println("generate: index.html")
	err = tmpl.ExecuteTemplate(file, "index.html", map[string]any{
		"Tags": tags,
	})
	if err != nil {
		log.Printf("Error: execute index template: %v", err)
		os.Exit(1)
	}
}

func generateTagFile(tag string) {
	name := fmt.Sprintf("%v.html", tag)
	file, err := os.Create(name)
	if err != nil {
		log.Printf("Error: create %q tag file: %v", tag, err)
		os.Exit(1)
	}
	defer file.Close()

	fmt.Printf("generate: %v\n", name)
	err = tmpl.ExecuteTemplate(file, "tag.html", map[string]any{
		"Tag": tag,
	})
	if err != nil {
		log.Printf("Error: execute %q tag template: %v", tag, err)
		os.Exit(1)
	}
}

func main() {
	log.SetFlags(0)

	flag.Parse()

	tags, ignore := loadTags()
	clean(ignore)

	generateIndex(tags)
	for _, tag := range tags {
		generateTagFile(tag)
	}
}
