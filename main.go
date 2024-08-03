package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/siongui/gojianfan"
)

var (
	chineseHeader  *regexp.Regexp
	otherHeader    *regexp.Regexp
	definition     *regexp.Regexp
	unwantedMarkup *regexp.Regexp
)

func init() {
	chineseHeader = regexp.MustCompile("^==Chinese==$")
	otherHeader = regexp.MustCompile("^==[^=].*[^=]==$")
	definition = regexp.MustCompile("^# ")
	unwantedMarkup = regexp.MustCompile(`(^# |\[\[|\]\])`)
}

type Definitions []string

func (definitions Definitions) String() string {
	var sb strings.Builder

	for index, line := range definitions {
		sb.WriteString(fmt.Sprintf("%d. %s\n", index+1, line))
	}

	return sb.String()
}

func RawWikiTextToDefinitions(text string) Definitions {
	var definitions []string = []string{}
	var reading bool = false

	lines := strings.Split(text, "\n")
	log.Printf("Processing %d lines of raw wiki text.", len(lines))

	for _, line := range lines {
		if chineseHeader.MatchString(line) {
			reading = true
			continue
		}

		if reading && otherHeader.MatchString(line) {
			reading = false
			break
		}

		if reading && definition.MatchString(line) {
			definition := unwantedMarkup.ReplaceAllString(line, "")
			definitions = append(definitions, definition)
		}
	}

	log.Printf("Returning %d definitions.", len(definitions))
	return definitions
}

func GetWikiURL(input string) string {
	traditional := gojianfan.S2T(input)
	log.Printf("%s maps to Traditional Chinese %s", input, traditional)

	encoded := url.QueryEscape(traditional)
	log.Printf("URL-encoded: %s", encoded)

	wikiURL := fmt.Sprintf("https://en.wiktionary.org/wiki/%s?action=raw", encoded)
	log.Printf("Wiki URL: %s", wikiURL)

	return wikiURL
}

func GetDefinitions(wikiURL string) Definitions {
	res, err := http.Get(wikiURL)
	if err != nil {
		log.Printf("Cannot download wiki page %s", wikiURL)
		return []string{}
	}

	if res.StatusCode != 200 {
		log.Printf("Bad status code %d for wiki page %s", res.StatusCode, wikiURL)
		return []string{}
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Cannot read the body for wiki page %s", wikiURL)
		return []string{}
	}

	return RawWikiTextToDefinitions(string(body))
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <chinese-word>\n", os.Args[0])
	}

	input := os.Args[1]

	wikiURL := GetWikiURL(input)
	fmt.Printf("Obtaining definitions for %s at %s\n", input, wikiURL)

	definitions := GetDefinitions(wikiURL)
	fmt.Print(definitions.String())
}
