package lib

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"text/template"

	"github.com/moraesvic/shenme/types"
	gopinyin "github.com/mozillazg/go-pinyin"
	"github.com/siongui/gojianfan"
)

var (
	chineseHeader  *regexp.Regexp
	otherHeader    *regexp.Regexp
	definition     *regexp.Regexp
	unwantedMarkup *regexp.Regexp
	pinyinArgs     gopinyin.Args

	definitionsHTMLTemplate *template.Template
)

func init() {
	log.SetFlags(log.Lmicroseconds)

	chineseHeader = regexp.MustCompile("^==Chinese==$")
	otherHeader = regexp.MustCompile("^==[^=].*[^=]==$")
	definition = regexp.MustCompile("^# ")
	unwantedMarkup = regexp.MustCompile(`(^# |\[\[|\]\])`)

	pinyinArgs = gopinyin.NewArgs()
	pinyinArgs.Style = gopinyin.Tone
	pinyinArgs.Heteronym = true

	_definitionsHTMLTemplate, err := template.New("definitionsHTML").
		Parse(`<ol>{{range .}}{{printf "<li>%s</li>" .}}{{end}}</ol>`)

	if err != nil {
		log.Fatalf("Error compiling template.")
	}

	definitionsHTMLTemplate = _definitionsHTMLTemplate
}

type Definitions []string

func (definitions Definitions) Length() int {
	return len(definitions)
}

func (definitions Definitions) String() string {
	var sb strings.Builder

	for index, line := range definitions {
		sb.WriteString(fmt.Sprintf("%d. %s\n", index+1, line))
	}

	return sb.String()
}

func (definitions Definitions) HTML() string {
	buf := new(bytes.Buffer)
	definitionsHTMLTemplate.Execute(buf, definitions)
	return buf.String()
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

func Traditional(simplified string) string {
	traditional := gojianfan.S2T(simplified)
	log.Printf("%s maps to Traditional Chinese %s", simplified, traditional)
	return traditional
}

func GetWikiURL(traditional string) string {
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

func _Pinyin(char string) string {
	result := gopinyin.Pinyin(char, pinyinArgs)[0]
	return fmt.Sprintf("(%s)", strings.Join(result, ", "))
}

func Pinyin(chars string) string {
	runes := []rune(chars)
	if len(runes) == 1 {
		return _Pinyin(chars)
	}

	results := []string{}
	for _, r := range runes {
		results = append(results, _Pinyin(string(r)))
	}

	return fmt.Sprintf("[%s]", strings.Join(results, ", "))
}

type Definer struct{}

func (Definer) _Define(traditional string) types.IDefinitionBoth {
	wikiURL := GetWikiURL(traditional)
	log.Printf("Obtaining definitions for %s at %s\n", traditional, wikiURL)
	return GetDefinitions(wikiURL)
}

func (d Definer) Define(traditional string) types.IDefinitionString {
	specialized, ok := d._Define(traditional).(types.IDefinitionString)
	if !ok {
		log.Fatalf("Type assertion failed.")
	}
	return specialized
}

func (d Definer) DefineHTML(traditional string) types.IDefinitionHTML {
	specialized, ok := d._Define(traditional).(types.IDefinitionHTML)
	if !ok {
		log.Fatalf("Type assertion failed.")
	}
	return specialized
}
