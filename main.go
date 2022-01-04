package main

import (
	"flag"
	"vscode/gostart/gendata"
)

func main() {
	fGenerateRawRecords := flag.Bool("gen", false, "Generate raw records")
	flag.Parse()

	if *fGenerateRawRecords {
		gendata.GenerateRawRecords("./data/gendataCenas", 150, 9)
	} else {
		flag.Usage()
	}

	// easy rsa number 10263280077814176196883978050069
}
