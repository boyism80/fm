package actor

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type junitSuite struct {
	XMLName  xml.Name    `xml:"testsuite"`
	Name     string      `xml:"name,attr"`
	Tests    int         `xml:"tests,attr"`
	Failures int         `xml:"failures,attr"`
	Time     string      `xml:"time,attr"`
	Cases    []junitCase `xml:"testcase"`
}

type junitCase struct {
	Name    string        `xml:"name,attr"`
	Time    string        `xml:"time,attr"`
	Failure *junitFailure `xml:"failure,omitempty"`
}

type junitFailure struct {
	Message string `xml:"message,attr"`
	Text    string `xml:",chardata"`
}

func (r *RunnerActor) report() {
	var sb strings.Builder
	passed, failed := 0, 0
	sb.WriteString("\n=== INTEGRATION TEST RESULTS ===\n")
	for _, res := range r.results {
		status := "PASS"
		if res.Passed {
			passed++
		} else {
			status = "FAIL"
			failed++
		}
		fmt.Fprintf(&sb, "[%s] %s  %.1fs\n", status, res.Name, res.Elapsed.Seconds())
		for _, f := range res.Failures {
			fmt.Fprintf(&sb, "       - %s\n", f)
		}
	}
	elapsed := time.Since(r.started)
	sb.WriteString("=== SUMMARY ===\n")
	fmt.Fprintf(&sb, "Total %d · Passed %d · Failed %d · Elapsed %.1fs\n", len(r.results), passed, failed, elapsed.Seconds())
	fmt.Print(sb.String())

	if r.junit == "" {
		return
	}
	doc := junitSuite{Name: "fm-bot", Tests: len(r.results), Failures: failed, Time: fmt.Sprintf("%.3f", elapsed.Seconds())}
	for _, res := range r.results {
		c := junitCase{Name: res.Name, Time: fmt.Sprintf("%.3f", res.Elapsed.Seconds())}
		if res.Passed == false {
			c.Failure = &junitFailure{Message: strings.Join(res.Failures, "; "), Text: strings.Join(res.Failures, "\n")}
		}
		doc.Cases = append(doc.Cases, c)
	}
	data, err := xml.MarshalIndent(doc, "", "  ")
	if err != nil {
		log.Printf("junit: %v", err)
		return
	}
	if err := os.MkdirAll(filepath.Dir(r.junit), 0o755); err != nil {
		log.Printf("junit: %v", err)
		return
	}
	if err := os.WriteFile(r.junit, append([]byte(xml.Header), data...), 0o644); err != nil {
		log.Printf("junit: %v", err)
	}
}
