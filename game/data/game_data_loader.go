package data

import (
	"encoding/xml"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type XMLNode struct {
	XMLName  xml.Name
	Name     string    `xml:"name,attr"`
	Value    string    `xml:"value,attr"`
	Children []XMLNode `xml:",any"`
}

func (node *XMLNode) Find(name string) *XMLNode {

	parts := strings.Split(name, ":")
	current := node
	for _, part := range parts {
		found := false
		for _, v := range current.Children {
			if v.Name == part {
				current = &v
				found = true
				break
			}
		}

		if !found {
			return nil
		}
	}

	return current
}

func LoadXmlFiles[T any](root string,
	workerCount int,
	action func(path string) (result *T, err error),
	callback func(percent float32, value *T)) error {
	var allFiles []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err == nil && !d.IsDir() && strings.HasSuffix(d.Name(), ".img.xml") {
			allFiles = append(allFiles, path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	total := len(allFiles)
	if total == 0 {
		return nil
	}

	jobs := make(chan string, total)
	results := make(chan *T, total)
	var wg sync.WaitGroup

	for range workerCount {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range jobs {
				m, err := action(path)

				if err == nil {
					results <- m
				} else {
					log.Println(err)
				}
			}
		}()
	}

	for _, path := range allFiles {
		jobs <- path
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	count := 0
	for m := range results {
		count++
		percent := float32(count) / float32(total) * 100
		callback(percent, m)
	}

	return nil
}
