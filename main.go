package main

import (
	"flag"
	"fmt"
	"vscode/gostart/gen"
	"vscode/gostart/norm"
)

func main() {
	bitSize := uint(150)
	numberOfRecords := uint(64000)
	rawRecordsFileName := fmt.Sprintf("%v_%v_%v.csv", "./data/rawRecords", bitSize, numberOfRecords)
	recordsFileName := fmt.Sprintf("%v_%v_%v.csv", "./data/records", bitSize, numberOfRecords)

	fGenerateRawRecords := flag.Bool("gen", false, "Generate raw records")
	fNormalizeRecords := flag.Bool("norm", false, "Normalize raw records into records")
	flag.Parse()

	if flag.NFlag() != 1 {
		flag.Usage()
	} else if *fGenerateRawRecords {
		gen.GenerateRawRecords(rawRecordsFileName, bitSize, numberOfRecords)
	} else if *fNormalizeRecords {
		norm.NormalizeRecords(rawRecordsFileName, bitSize, recordsFileName)
	}

	// easy rsa number 10263280077814176196883978050069
}
