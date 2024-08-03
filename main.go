package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/moraesvic/shenme/lib"
	"github.com/moraesvic/shenme/types"
)

var (
	isHTML  *bool
	isDebug *bool
)

func Usage() {
	fmt.Printf("Usage: %s [OPTIONS] <chinese-word>\n", os.Args[0])
	flag.PrintDefaults()
}

func init() {
	flag.Usage = Usage

	isHTML = flag.Bool("html", false, "display as HTML (default is text)")
	isDebug = flag.Bool("debug", false, "enable debug mode")

	flag.Parse()
}

func CLIText(i types.IDefinerString, input string) {
	traditional := lib.Traditional(input)

	wikiURL := lib.GetWikiURL(traditional)
	fmt.Fprintf(os.Stderr, "Obtaining definitions for %s %s at %s\n", input, lib.Pinyin(input), wikiURL)

	definitions := i.Define(traditional)
	fmt.Print(definitions.String())

	if definitions.Length() == 0 {
		fmt.Fprintf(os.Stderr, "No definitions were found.")
	}
}

func CLIHTML(i types.IDefinerHTML, input string) {
	traditional := lib.Traditional(input)

	wikiURL := lib.GetWikiURL(traditional)
	fmt.Fprintf(os.Stderr, "Obtaining definitions for %s %s at %s\n", input, lib.Pinyin(input), wikiURL)

	definitions := i.DefineHTML(traditional)
	fmt.Print(definitions.HTML())
}

func main() {
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	input := flag.Args()[0]
	definer := lib.Definer{}

	if !*isDebug {
		log.SetOutput(io.Discard)
	}

	if *isHTML {
		CLIHTML(definer, input)
	} else {
		CLIText(definer, input)
	}
}
