package norm

import (
	"encoding/csv"
	"fmt"
	"io"
	"math/big"
	"os"
	"reflect"
	"vscode/gostart/gen"
)

func NormalizeRecords(rawRecordsFileName string, bitSize uint, recordsFileName string) {
	rawRecordsFile, osOpenErr := os.Open(rawRecordsFileName)
	if osOpenErr != nil {
		panic(fmt.Sprintf("Failed to open csv file, %v\n", osOpenErr))
	}
	// defer file close so the file is always closed
	defer rawRecordsFile.Close()

	csvReader := csv.NewReader(rawRecordsFile)
	headers, err := csvReader.Read()
	if err != nil {
		panic(fmt.Sprintf("Failed to read %v first record, err %v", rawRecordsFileName, err))
	}
	if !reflect.DeepEqual(gen.RawRecordFieldNames(), headers) {
		panic(fmt.Sprintf("Failed to read %v headers, expected %v headers, got headers %v", rawRecordsFileName, gen.RawRecordFieldNames(), headers))
	}

	// initialize min/max
	numberNMin := new(big.Int)
	numberNMax := new(big.Int)
	smallerPrimeMin := new(big.Int)
	smallerPrimeMax := new(big.Int)
	biggerPrimeMin := new(big.Int)
	biggerPrimeMax := new(big.Int)
	primesSumMin := new(big.Int)
	primesSumMax := new(big.Int)

	numberNMax.SetInt64(0)
	smallerPrimeMax.SetInt64(0)
	biggerPrimeMax.SetInt64(0)
	primesSumMax.SetInt64(0)

	// calculate max possible value according to bitSize
	bitSizeAsBigInt := new(big.Int)
	bitSizeAsBigInt.SetInt64(int64(bitSize))
	numberNMin.SetInt64(2)
	numberNMin.Exp(numberNMin, bitSizeAsBigInt, nil)

	smallerPrimeMin.Set(numberNMin)
	biggerPrimeMin.Set(numberNMin)
	primesSumMin.Set(numberNMin)

	// read first record since for "increment" is executed after
	record, err := csvReader.Read()
	for ; record != nil && err == nil; record, err = csvReader.Read() {
		readRawRecord := convertToRawRecord(record)

		if numberNMin.Cmp(readRawRecord.NumberN) > 0 {
			numberNMin.Set(readRawRecord.NumberN)
		}
		if smallerPrimeMin.Cmp(readRawRecord.SmallerPrime) > 0 {
			smallerPrimeMin.Set(readRawRecord.SmallerPrime)
		}
		if biggerPrimeMin.Cmp(readRawRecord.BiggerPrime) > 0 {
			biggerPrimeMin.Set(readRawRecord.BiggerPrime)
		}
		if primesSumMin.Cmp(readRawRecord.PrimesSum) > 0 {
			primesSumMin.Set(readRawRecord.PrimesSum)
		}

		if numberNMax.Cmp(readRawRecord.NumberN) < 0 {
			numberNMax.Set(readRawRecord.NumberN)
		}
		if smallerPrimeMax.Cmp(readRawRecord.SmallerPrime) < 0 {
			smallerPrimeMax.Set(readRawRecord.SmallerPrime)
		}
		if biggerPrimeMax.Cmp(readRawRecord.BiggerPrime) < 0 {
			biggerPrimeMax.Set(readRawRecord.BiggerPrime)
		}
		if primesSumMax.Cmp(readRawRecord.PrimesSum) < 0 {
			primesSumMax.Set(readRawRecord.PrimesSum)
		}
	}
	if err != nil && err != io.EOF {
		panic(fmt.Sprintf("Failed to read %v a record, err %v", rawRecordsFileName, err))
	}

	// reset file handler to begin of the file
	rawRecordsFile.Seek(0, 0)
	headers, err = csvReader.Read()
	if err != nil {
		panic(fmt.Sprintf("Failed to re-read %v first record, err %v", rawRecordsFileName, err))
	}

	// create summary CSV file
	{
		// create records file
		recordsFile, osCreateErr := os.Create(recordsFileName + "_summar.csv")
		if osCreateErr != nil {
			panic(fmt.Sprintf("Failed to create csv file, %v\n", osCreateErr))
		}
		// defer file close so the file is always closed
		defer recordsFile.Close()

		// write header
		csvWriter := csv.NewWriter(recordsFile)
		csvWriteErr := csvWriter.Write([]string{
			"numberNMin",
			"numberNMax",
			"smallerPrimeMin",
			"smallerPrimeMax",
			"biggerPrimeMin",
			"biggerPrimeMax",
			"primesSumMin",
			"primesSumMax"})
		if csvWriteErr != nil {
			panic(fmt.Sprintf("Failed to write header to csv file, %v\n", csvWriteErr))
		}
		csvWriteErr = csvWriter.Write([]string{
			fmt.Sprintf("%v", numberNMin),
			fmt.Sprintf("%v", numberNMax),
			fmt.Sprintf("%v", smallerPrimeMin),
			fmt.Sprintf("%v", smallerPrimeMax),
			fmt.Sprintf("%v", biggerPrimeMin),
			fmt.Sprintf("%v", biggerPrimeMax),
			fmt.Sprintf("%v", primesSumMin),
			fmt.Sprintf("%v", primesSumMax),
		},
		)
		if csvWriteErr != nil {
			panic(fmt.Sprintf("Failed to write header to csv file, %v\n", csvWriteErr))
		}
		// defer flush because write does not immediatly write records
		// and the defer close will not flush before closing the file
		defer csvWriter.Flush()
	}

	// convert max/mins to big float to avoid multiple conversions
	numberNMinAsBigFloat := new(big.Float).SetInt(numberNMin)
	numberNMaxAsBigFloat := new(big.Float).SetInt(numberNMax)
	smallerPrimeMinAsBigFloat := new(big.Float).SetInt(smallerPrimeMin)
	smallerPrimeMaxAsBigFloat := new(big.Float).SetInt(smallerPrimeMax)
	biggerPrimeMinAsBigFloat := new(big.Float).SetInt(biggerPrimeMin)
	biggerPrimeMaxAsBigFloat := new(big.Float).SetInt(biggerPrimeMax)
	primesSumMinAsBigFloat := new(big.Float).SetInt(primesSumMin)
	primesSumMaxAsBigFloat := new(big.Float).SetInt(primesSumMax)

	// Shaddy normalization, using the same scale for both smaller and bigger prime
	biggerPrimeMinAsBigFloat.Set(smallerPrimeMinAsBigFloat)
	smallerPrimeMaxAsBigFloat.Set(biggerPrimeMaxAsBigFloat)

	// create records file
	recordsFile, osCreateErr := os.Create(recordsFileName)
	if osCreateErr != nil {
		panic(fmt.Sprintf("Failed to create csv file, %v\n", osCreateErr))
	}
	// defer file close so the file is always closed
	defer recordsFile.Close()

	// write header
	csvWriter := csv.NewWriter(recordsFile)
	csvWriteErr := csvWriter.Write(RecordFieldNames())
	if csvWriteErr != nil {
		panic(fmt.Sprintf("Failed to write header to csv file, %v\n", csvWriteErr))
	}
	// defer flush because write does not immediatly write records
	// and the defer close will not flush before closing the file
	defer csvWriter.Flush()

	// read first raw record since for "increment" is executed after
	record, err = csvReader.Read()
	for ; record != nil && err == nil; record, err = csvReader.Read() {
		readRawRecord := convertToRawRecord(record)
		normalizedRecord := normalizeRecord(
			readRawRecord,
			numberNMinAsBigFloat,
			numberNMaxAsBigFloat,
			smallerPrimeMinAsBigFloat,
			smallerPrimeMaxAsBigFloat,
			biggerPrimeMinAsBigFloat,
			biggerPrimeMaxAsBigFloat,
			primesSumMinAsBigFloat,
			primesSumMaxAsBigFloat)

		csvWriteErr := csvWriter.Write(normalizedRecord.AsStringSlice())
		if csvWriteErr != nil {
			panic(fmt.Sprintf("Failed to write csv file, %v\n", csvWriteErr))
		}
	}
}

type Record struct {
	NumberN      float64
	SmallerPrime float64
	BiggerPrime  float64
	PrimesSum    float64
}

func RecordFieldNames() []string {
	return []string{"NumberN", "SmallerPrime", "BiggerPrime", "PrimesSum"}
}

func (source Record) AsStringSlice() []string {
	return []string{fmt.Sprintf("%v", source.NumberN), fmt.Sprintf("%v", source.SmallerPrime), fmt.Sprintf("%v", source.BiggerPrime), fmt.Sprintf("%v", source.PrimesSum)}
}

// normalize record using min-max strategy
// reference https://docs.microsoft.com/en-us/azure/machine-learning/studio-module-reference/normalize-data#-how-to-configure-normalize-data
func normalizeRecord(readRawRecord gen.RawRecord,
	numberNMin *big.Float,
	numberNMax *big.Float,
	smallerPrimeMin *big.Float,
	smallerPrimeMax *big.Float,
	biggerPrimeMin *big.Float,
	biggerPrimeMax *big.Float,
	primesSumMin *big.Float,
	primesSumMax *big.Float) Record {

	var normalizedRecord Record
	dividend := new(big.Float)
	divisor := new(big.Float)
	quotient := new(big.Float)

	// normalize NumberN
	dividend.SetInt(readRawRecord.NumberN)
	dividend.Sub(dividend, numberNMin)
	divisor.Sub(numberNMax, numberNMin)
	quotient.Quo(dividend, divisor)
	normalizedRecord.NumberN, _ = quotient.Float64()

	// normalize SmallerPrime
	dividend.SetInt(readRawRecord.SmallerPrime)
	dividend.Sub(dividend, smallerPrimeMin)
	divisor.Sub(smallerPrimeMax, smallerPrimeMin)
	quotient.Quo(dividend, divisor)
	normalizedRecord.SmallerPrime, _ = quotient.Float64()

	// normalize BiggerPrime
	dividend.SetInt(readRawRecord.BiggerPrime)
	dividend.Sub(dividend, biggerPrimeMin)
	divisor.Sub(biggerPrimeMax, biggerPrimeMin)
	quotient.Quo(dividend, divisor)
	normalizedRecord.BiggerPrime, _ = quotient.Float64()

	// normalize PrimesSum
	dividend.SetInt(readRawRecord.PrimesSum)
	dividend.Sub(dividend, primesSumMin)
	divisor.Sub(primesSumMax, primesSumMin)
	quotient.Quo(dividend, divisor)
	normalizedRecord.PrimesSum, _ = quotient.Float64()

	return normalizedRecord
}

func convertToRawRecord(record []string) gen.RawRecord {
	rawRecordNumberOfColumns := len(gen.RawRecordFieldNames())
	var parsedRawRecord gen.RawRecord

	if len(record) != rawRecordNumberOfColumns {
		panic(fmt.Sprintf("Failed to parse a raw record, unexpected number of columns, expected %v columns, got this %v", rawRecordNumberOfColumns, record))
	}

	numberN := new(big.Int)
	smallerPrime := new(big.Int)
	biggerPrime := new(big.Int)
	primesSum := new(big.Int)

	_, numberWasSet := numberN.SetString(record[0], 10)
	if !numberWasSet {
		panic(fmt.Sprintf("Failed to parse a raw record number, read value as string %v", record[0]))
	}

	_, numberWasSet = smallerPrime.SetString(record[1], 10)
	if !numberWasSet {
		panic(fmt.Sprintf("Failed to parse a raw record number, read value as string %v", record[1]))
	}

	_, numberWasSet = biggerPrime.SetString(record[2], 10)
	if !numberWasSet {
		panic(fmt.Sprintf("Failed to parse a raw record number, read value as string %v", record[2]))
	}

	_, numberWasSet = primesSum.SetString(record[3], 10)
	if !numberWasSet {
		panic(fmt.Sprintf("Failed to parse a raw record number, read value as string %v", record[3]))
	}

	parsedRawRecord.NumberN = numberN
	parsedRawRecord.SmallerPrime = smallerPrime
	parsedRawRecord.BiggerPrime = biggerPrime
	parsedRawRecord.PrimesSum = primesSum

	return parsedRawRecord
}
