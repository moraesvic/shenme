package main

import (
	"fmt"
	"os"

	"github.com/moraesvic/shenme/lib"
	"github.com/moraesvic/shenme/types"
)

func CLI(i types.IDefinerString, input string) {
	traditional := lib.Traditional(input)

	wikiURL := lib.GetWikiURL(traditional)
	fmt.Printf("Obtaining definitions for %s %s at %s\n", input, lib.Pinyin(input), wikiURL)

	definitions := i.Define(traditional)
	fmt.Print(definitions.String())

	if definitions.Length() == 0 {
		fmt.Println("No definitions were found.")
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <chinese-word>\n", os.Args[0])
		panic(0)
	}

	input := os.Args[1]
	CLI(lib.Definer{}, input)
}
