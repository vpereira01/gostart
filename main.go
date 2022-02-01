package main

import (
	"flag"
	"fmt"
	"math/big"
	"vscode/gostart/gen"
	"vscode/gostart/gonn"
	"vscode/gostart/norm"
)

func main() {
	bitSize := uint(150)
	numberOfRecords := uint(64000)
	rawRecordsFileName := fmt.Sprintf("%v_%v_%v.csv", "./data/rawRecords", bitSize, numberOfRecords)
	recordsFileName := fmt.Sprintf("%v_%v_%v.csv", "./data/records", bitSize, numberOfRecords)

	fGenerateRawRecords := flag.Bool("gen", false, "Generate raw records")
	fNormalizeRecords := flag.Bool("norm", false, "Normalize raw records into records")
	fTrain := flag.Bool("train", false, "Train based on normalized record")
	fCalc := flag.Bool("calc", false, "Do some calculations")
	flag.Parse()

	if flag.NFlag() != 1 {
		flag.Usage()
	} else if *fGenerateRawRecords {
		gen.GenerateRawRecords(rawRecordsFileName, bitSize, numberOfRecords)
	} else if *fNormalizeRecords {
		norm.NormalizeRecords(rawRecordsFileName, bitSize, recordsFileName)
	} else if *fTrain {
		gonn.Train(recordsFileName)
	} else if *fCalc {
		Calc()
	}

	// easy rsa number 10263280077814176196883978050069
}

func Calc() {
	// numberNMin, _ := new(big.Int).SetString("804783843082352595001719217734972675795934649", 10)
	// numberNMax, _ := new(big.Int).SetString("1425777428768627105133862492622722299351364001", 10)
	smallerPrimeMin, _ := new(big.Int).SetString("28334200780378472749589", 10)
	// smallerPrimeMax, _ := new(big.Int).SetString("37746249243703610892493", 10)
	// biggerPrimeMin, _ := new(big.Int).SetString("28381622674296876536389", 10)
	biggerPrimeMax, _ := new(big.Int).SetString("37778776000131138816929", 10)

	// convert max/mins to big float to avoid multiple conversions
	// numberNMinAsBigFloat := new(big.Float).SetInt(numberNMin)
	// numberNMaxAsBigFloat := new(big.Float).SetInt(numberNMax)
	smallerPrimeMinAsBigFloat := new(big.Float).SetInt(smallerPrimeMin)
	// smallerPrimeMaxAsBigFloat := new(big.Float).SetInt(smallerPrimeMax)
	// biggerPrimeMinAsBigFloat := new(big.Float).SetInt(biggerPrimeMin)
	biggerPrimeMaxAsBigFloat := new(big.Float).SetInt(biggerPrimeMax)

	// baseValue, _ := new(big.Int).SetString("1231853026813173440488109514492006479965925417", 10)
	// convertedValue, _ := new(big.Float).SetString("0.6877191545527104")
	estimatedSmallerPrime, _ := new(big.Float).SetString("0.59046845")
	// estimatedBiggerPrime, _ := new(big.Float).SetString("0.84594564")

	diff := new(big.Float)
	muldiff := new(big.Float)
	result := new(big.Float)

	diff.Sub(biggerPrimeMaxAsBigFloat, smallerPrimeMinAsBigFloat)
	muldiff.Mul(estimatedSmallerPrime, diff)
	result.Add(muldiff, smallerPrimeMinAsBigFloat)
	resultNum := new(big.Int)
	result.Int(resultNum)
	fmt.Printf("resultNum: %v\n", resultNum)
	realNum, _ := new(big.Int).SetString("33110995368503928176591", 10)
	fmt.Printf("resultNum: %v\n", realNum.Sub(realNum, resultNum))
}
