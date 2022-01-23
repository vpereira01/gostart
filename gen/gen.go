package gen

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/csv"
	"fmt"
	"log"
	"math/big"
	"os"
	"runtime"
	"sync"
)

func GenerateRawRecords(rawRecordsFileName string, bitSize uint, numberOfRecords uint) {
	if runtime.NumCPU() <= 1 {
		panic("Requires num of cpus greater than 1")
	}
	numOfRoutines := uint(runtime.NumCPU() - 1)
	quotientRecordsPerRoutine := numberOfRecords / numOfRoutines
	remainderRecordsPerRoutine := numberOfRecords - (quotientRecordsPerRoutine * numOfRoutines)

	generatedRawRecords := make(chan RawRecord, 8)
	finishedWaitGroup := new(sync.WaitGroup)
	finishedWaitGroup.Add(int(numOfRoutines))

	// set trigger to close results channel
	go func() {
		finishedWaitGroup.Wait()
		close(generatedRawRecords)
	}()

	rawRecordsFile, osCreateErr := os.Create(rawRecordsFileName)
	if osCreateErr != nil {
		panic(fmt.Sprintf("Failed to create csv file, %v\n", osCreateErr))
	}
	// defer file close so the file is always closed
	defer rawRecordsFile.Close()

	csvWriter := csv.NewWriter(rawRecordsFile)
	csvWriteErr := csvWriter.Write(RawRecordFieldNames())
	if csvWriteErr != nil {
		panic(fmt.Sprintf("Failed to write header to csv file, %v\n", csvWriteErr))
	}
	// defer flush because write does not immediatly write records
	// and the defer close will not flush before closing the file
	defer csvWriter.Flush()

	// launch go routines
	var workerNum uint
	leftRemainderRecordsPerRoutine := remainderRecordsPerRoutine
	for workerNum = 1; workerNum <= numOfRoutines; workerNum++ {
		var nextRecordsPerRoutine uint
		if leftRemainderRecordsPerRoutine > 0 {
			nextRecordsPerRoutine = quotientRecordsPerRoutine + 1
			leftRemainderRecordsPerRoutine--
		} else {
			nextRecordsPerRoutine = quotientRecordsPerRoutine
		}
		go rawRecordGenerator(generatedRawRecords, finishedWaitGroup, bitSize, nextRecordsPerRoutine)
	}

	count := 0
	// read generated raw records and store them on csv
	for rawRecord := range generatedRawRecords {
		count++
		csvWriteErr := csvWriter.Write(rawRecord.AsStringSlice())
		if csvWriteErr != nil {
			panic(fmt.Sprintf("Failed to write csv file, %v\n", csvWriteErr))
		}
		if count%5000 == 0 {
			fmt.Printf("wrote %v records\n", count)
		}
	}
}

type RawRecord struct {
	NumberN      *big.Int
	SmallerPrime *big.Int
	BiggerPrime  *big.Int
	PrimesSum    *big.Int
}

func RawRecordFieldNames() []string {
	return []string{"NumberN", "SmallerPrime", "BiggerPrime", "PrimesSum"}
}

func (source RawRecord) AsStringSlice() []string {
	return []string{source.NumberN.String(), source.SmallerPrime.String(), source.BiggerPrime.String(), source.PrimesSum.String()}
}

func rawRecordGenerator(resultsChannel chan RawRecord, finishedWaitGroup *sync.WaitGroup, bitSize uint, numberOfRecords uint) {
	var count uint
	for count = 1; count <= numberOfRecords; count++ {
		//Generate RSA keys
		privateKey, err := rsa.GenerateKey(rand.Reader, int(bitSize))
		if err != nil {
			log.Printf("Failed to generate key, %v\n", err)
			// decrease count so this is repeated
			count--
			continue
		}
		if len(privateKey.Primes) != 2 {
			log.Printf("More than 2 primes generated, %v\n", privateKey.Primes)
			// decrease count so this is repeated
			count--
			continue
		}

		numberN := new(big.Int)
		smallerPrime := new(big.Int)
		biggerPrime := new(big.Int)
		primesSum := new(big.Int)

		numberN.Set(privateKey.N)
		if privateKey.Primes[0].Cmp(privateKey.Primes[1]) < 0 {
			smallerPrime.Set(privateKey.Primes[0])
			biggerPrime.Set(privateKey.Primes[1])
		} else {
			smallerPrime.Set(privateKey.Primes[1])
			biggerPrime.Set(privateKey.Primes[0])
		}
		primesSum.Add(smallerPrime, biggerPrime)

		resultsChannel <- RawRecord{
			NumberN:      numberN,
			SmallerPrime: smallerPrime,
			BiggerPrime:  biggerPrime,
			PrimesSum:    primesSum,
		}
	}

	finishedWaitGroup.Done()
}
