package main

import (
	"flag"
	"fmt"
	"math/big"
	"vscode/gostart/findpdiff"
	"vscode/gostart/findpsum"
	"vscode/gostart/gen"
	"vscode/gostart/norm"
)

func main() {
	bitSize := uint(64)
	numberOfRecords := uint(256000)
	rawRecordsFileName := fmt.Sprintf("%v_%v_%v.csv", "./data/rawRecords", bitSize, numberOfRecords)
	recordsFileName := fmt.Sprintf("%v_%v_%v.csv", "./data/records", bitSize, numberOfRecords)

	fGenerateRawRecords := flag.Bool("gen", false, "Generate raw records")
	fNormalizeRecords := flag.Bool("norm", false, "Normalize raw records into records")
	fDeNormalizeValue := flag.Bool("denorm", false, "Denormalize value")
	fFindPrimeFactorsSum := flag.Bool("findpsum", false, "Find prime factors through their sum")
	fFindPrimeFactorsDiff := flag.Bool("findpdiff", false, "Find prime factors through their difference")
	flag.Parse()

	if flag.NFlag() != 1 {
		flag.Usage()
	} else if *fGenerateRawRecords {
		gen.GenerateRawRecords(rawRecordsFileName, bitSize, numberOfRecords)
	} else if *fNormalizeRecords {
		norm.NormalizeRecords(rawRecordsFileName, bitSize, recordsFileName)
	} else if *fDeNormalizeValue {
		norm.DeNormalizeValue()
	} else if *fFindPrimeFactorsSum {
		targetNumberN, _ := new(big.Int).SetString("3135626011863327378113", 10)
		primesSumEstimation, _ := new(big.Int).SetString("112029222574", 10)
		findpsum.TryToFindPrimeFactors(targetNumberN, primesSumEstimation)
	} else if *fFindPrimeFactorsDiff {
		targetNumberN, _ := new(big.Int).SetString("3135626011863327378113", 10)
		findpdiff.TryToFindPrimeFactors(targetNumberN)
	}

	// easy rsa number 10263280077814176196883978050069 ~100bits 2seconds by gnu factor
}
