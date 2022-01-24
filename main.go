package main

import (
	"errors"
	"flag"
	"fmt"
	"math"
	"math/big"
	"sync"
	"vscode/gostart/find"
	"vscode/gostart/gen"
	"vscode/gostart/gonn"
	"vscode/gostart/norm"
)

func main() {
	bitSize := uint(64)
	numberOfRecords := uint(256000)
	rawRecordsFileName := fmt.Sprintf("%v_%v_%v.csv", "./data/rawRecords", bitSize, numberOfRecords)
	recordsFileName := fmt.Sprintf("%v_%v_%v.csv", "./data/records", bitSize, numberOfRecords)

	fGenerateRawRecords := flag.Bool("gen", false, "Generate raw records")
	fNormalizeRecords := flag.Bool("norm", false, "Normalize raw records into records")
	fTrain := flag.Bool("train", false, "Train based on normalized record")
	fInvTrans := flag.Bool("invTrans", false, "Do inverse transformation")
	fWIP := flag.Bool("wip", false, "Work in progress")
	flag.Parse()

	if flag.NFlag() != 1 {
		flag.Usage()
	} else if *fGenerateRawRecords {
		gen.GenerateRawRecords(rawRecordsFileName, bitSize, numberOfRecords)
	} else if *fNormalizeRecords {
		norm.NormalizeRecords(rawRecordsFileName, bitSize, recordsFileName)
	} else if *fTrain {
		gonn.Train(recordsFileName)
	} else if *fInvTrans {
		InverseTransform()
	} else if *fWIP {
		WIP()
	}

	// easy rsa number 10263280077814176196883978050069 ~100bits 2seconds
}

func WIP() {
	finishedWaitGroup := new(sync.WaitGroup)

	targetNumberN, _ := new(big.Int).SetString("13438310478517603073", 10)
	realSum, _ := new(big.Int).SetString("7387124034", 10)

	targetNumberN.SetString("13322573880505234223", 10)
	find.FindPrimesSub(targetNumberN)

	return
	primesSumEstimation := new(big.Int).SetInt64(500000)
	primesSumEstimation.Add(realSum, primesSumEstimation)

	// FindPeriods(targetNumberN, primesSumEstimation)
	// FindPeriods2(targetNumberN, primesSumEstimation)
	// FindPeriods3(targetNumberN, primesSumEstimation)
	// return

	finishedWaitGroup.Add(2)
	TryToFindPrimes(finishedWaitGroup, targetNumberN, primesSumEstimation, 1)

	// targetNumberN, _ := new(big.Int).SetString("1137501950159415429602603162643696967245621181", 10)
	// primesSumEstimation, _ := new(big.Int).SetString("67620573609986982832142", 10)

	// finishedWaitGroup := new(sync.WaitGroup)
	// finishedWaitGroup.Add(2)
	// go TryToFindPrimes(finishedWaitGroup, targetNumberN, primesSumEstimation, -1)
	// go TryToFindPrimes(finishedWaitGroup, targetNumberN, primesSumEstimation, 1)
	// finishedWaitGroup.Wait()
}

func FindPeriods(targetNumberN *big.Int, sumEstimation *big.Int) {
	targetNumberNBF := new(big.Float).SetInt(targetNumberN)
	smallerPrimeBF := new(big.Float)
	sumEstimationBF := new(big.Float).SetInt(sumEstimation)
	sumEstimationErrorBF := new(big.Float).SetInt64(0)
	one := new(big.Int).SetInt64(1)
	negOneBF := new(big.Float).SetInt64(-1)
	twoBF := new(big.Float).SetInt64(2)
	tenBF := new(big.Float).SetInt64(10)

	// calculate starting value
	hardCalc, err := WIP2(sumEstimationBF, sumEstimationErrorBF, targetNumberNBF)
	if err != nil {
		panic(fmt.Sprintf("Stopping due to %v", err))
	}
	//    _____________________________________________
	// -╲╱sumEst² - 2 ⋅ c1 ⋅ sumEst - 4 ⋅ numberN + c1² - sumEst + c1
	//  ────────────────────────────────────────────────────────────
	//								2
	smallerPrimeBF.Sub(hardCalc, sumEstimationBF)
	smallerPrimeBF.Add(smallerPrimeBF, sumEstimationErrorBF)
	smallerPrimeBF.Quo(smallerPrimeBF, twoBF)
	smallerPrimeBF.Mul(smallerPrimeBF, negOneBF)

	targetStartingValueSmallerPrime := new(big.Int)
	targetStartingValueSmallerPrimeBF := new(big.Float)
	smallerPrimeBF.Int(targetStartingValueSmallerPrime)
	targetStartingValueSmallerPrime.Add(targetStartingValueSmallerPrime, one)
	fmt.Printf("targetStartingValueSmallerPrime: %v\n", targetStartingValueSmallerPrime)
	targetStartingValueSmallerPrimeBF.SetInt(targetStartingValueSmallerPrime)

	sumEstimationErrorIncrementBF := new(big.Float)

	currentDiff := new(big.Float)
	previousDiff := new(big.Float)

	sumEstimationErrorIncrementBF.SetFloat64(0.1)

	// find base value
	for i := 0; i <= 513; i++ {
		sumEstimationErrorBF.Add(sumEstimationErrorBF, sumEstimationErrorIncrementBF)

		hardCalc, err := WIP2(sumEstimationBF, sumEstimationErrorBF, targetNumberNBF)
		if err != nil {
			fmt.Printf("Stopping due to %v", err)
			break
		}
		//    _____________________________________________
		// -╲╱sumEst² - 2 ⋅ c1 ⋅ sumEst - 4 ⋅ numberN + c1² - sumEst + c1
		//  ────────────────────────────────────────────────────────────
		//								2
		smallerPrimeBF.Sub(hardCalc, sumEstimationBF)
		smallerPrimeBF.Add(smallerPrimeBF, sumEstimationErrorBF)
		smallerPrimeBF.Quo(smallerPrimeBF, twoBF)
		smallerPrimeBF.Mul(smallerPrimeBF, negOneBF)

		currentDiff.Sub(targetStartingValueSmallerPrimeBF, smallerPrimeBF)

		if previousDiff.Cmp(currentDiff) == 0 {
			fmt.Printf("sumEstimationErrorBF: %.12f\n", sumEstimationErrorBF)
			break
		} else if smallerPrimeBF.Cmp(targetStartingValueSmallerPrimeBF) > 0 {
			sumEstimationErrorBF.Sub(sumEstimationErrorBF, sumEstimationErrorIncrementBF)
			sumEstimationErrorIncrementBF.Quo(sumEstimationErrorIncrementBF, tenBF)
			fmt.Printf("sumEstimationErrorIncrementBF: %.12f\n", sumEstimationErrorIncrementBF)
		}
		previousDiff.Set(currentDiff)
		if i == 512 {
			panic("could not find base value")
		}
	}

	sumEstimationErrorIncrementBF.SetFloat64(0.1)

	fmt.Printf("\"run2\": %v\n", "run2")

	targetStartingValueSmallerPrime.Add(targetStartingValueSmallerPrime, one)
	fmt.Printf("targetStartingValueSmallerPrime: %v\n", targetStartingValueSmallerPrime)
	targetStartingValueSmallerPrimeBF.SetInt(targetStartingValueSmallerPrime)

	// find base value
	for i := 0; i <= 513; i++ {
		sumEstimationErrorBF.Add(sumEstimationErrorBF, sumEstimationErrorIncrementBF)

		hardCalc, err := WIP2(sumEstimationBF, sumEstimationErrorBF, targetNumberNBF)
		if err != nil {
			fmt.Printf("Stopping due to %v", err)
			break
		}
		//    _____________________________________________
		// -╲╱sumEst² - 2 ⋅ c1 ⋅ sumEst - 4 ⋅ numberN + c1² - sumEst + c1
		//  ────────────────────────────────────────────────────────────
		//								2
		smallerPrimeBF.Sub(hardCalc, sumEstimationBF)
		smallerPrimeBF.Add(smallerPrimeBF, sumEstimationErrorBF)
		smallerPrimeBF.Quo(smallerPrimeBF, twoBF)
		smallerPrimeBF.Mul(smallerPrimeBF, negOneBF)

		currentDiff.Sub(targetStartingValueSmallerPrimeBF, smallerPrimeBF)

		if previousDiff.Cmp(currentDiff) == 0 {
			fmt.Printf("sumEstimationErrorBF: %.12f\n", sumEstimationErrorBF)
			break
		} else if smallerPrimeBF.Cmp(targetStartingValueSmallerPrimeBF) > 0 {
			sumEstimationErrorBF.Sub(sumEstimationErrorBF, sumEstimationErrorIncrementBF)
			sumEstimationErrorIncrementBF.Quo(sumEstimationErrorIncrementBF, tenBF)
			fmt.Printf("sumEstimationErrorIncrementBF: %.12f\n", sumEstimationErrorIncrementBF)
		}
		previousDiff.Set(currentDiff)
		if i == 512 {
			panic("could not find base value")
		}
	}
}

func FindPeriods2(targetNumberN *big.Int, sumEstimation *big.Int) {
	targetNumberNBF := new(big.Float).SetInt(targetNumberN)
	biggerPrimeBF := new(big.Float)
	sumEstimationBF := new(big.Float).SetInt(sumEstimation)
	sumEstimationErrorBF := new(big.Float).SetInt64(0)
	one := new(big.Int).SetInt64(1)
	twoBF := new(big.Float).SetInt64(2)
	tenBF := new(big.Float).SetInt64(10)

	fmt.Printf("\n#########################################\n\n")

	// calculate starting value
	hardCalc, err := WIP2(sumEstimationBF, sumEstimationErrorBF, targetNumberNBF)
	if err != nil {
		panic(fmt.Sprintf("Stopping due to %v", err))
	}
	//    _____________________________________________
	//  ╲╱sumEst² - 2 ⋅ c1 ⋅ sumEst - 4 ⋅ numberN + c1² + sumEst - c1
	//  ────────────────────────────────────────────────────────────
	//								2
	biggerPrimeBF.Add(hardCalc, sumEstimationBF)
	biggerPrimeBF.Sub(biggerPrimeBF, sumEstimationErrorBF)
	biggerPrimeBF.Quo(biggerPrimeBF, twoBF)

	targetStartingValueBiggerPrime := new(big.Int)
	targetStartingValueBiggerPrimeBF := new(big.Float)

	biggerPrimeBF.Int(targetStartingValueBiggerPrime)
	targetStartingValueBiggerPrime.Sub(targetStartingValueBiggerPrime, one)
	fmt.Printf("targetStartingValueBiggerPrime: %v\n", targetStartingValueBiggerPrime)
	targetStartingValueBiggerPrimeBF.SetInt(targetStartingValueBiggerPrime)

	sumEstimationErrorIncrementBF := new(big.Float)

	currentDiff := new(big.Float)
	previousDiff := new(big.Float)

	sumEstimationErrorIncrementBF.SetFloat64(0.1)

	// find base value
	for i := 0; i <= 513; i++ {
		sumEstimationErrorBF.Add(sumEstimationErrorBF, sumEstimationErrorIncrementBF)

		hardCalc, err := WIP2(sumEstimationBF, sumEstimationErrorBF, targetNumberNBF)
		if err != nil {
			fmt.Printf("Stopping due to %v", err)
			break
		}
		//    _____________________________________________
		//  ╲╱sumEst² - 2 ⋅ c1 ⋅ sumEst - 4 ⋅ numberN + c1² + sumEst - c1
		//  ────────────────────────────────────────────────────────────
		//								2
		biggerPrimeBF.Add(hardCalc, sumEstimationBF)
		biggerPrimeBF.Sub(biggerPrimeBF, sumEstimationErrorBF)
		biggerPrimeBF.Quo(biggerPrimeBF, twoBF)

		currentDiff.Sub(targetStartingValueBiggerPrimeBF, biggerPrimeBF)
		// fmt.Printf("currentDiff: %v\n", currentDiff)

		if previousDiff.Cmp(currentDiff) == 0 {
			fmt.Printf("sumEstimationErrorBF: %.12f\n", sumEstimationErrorBF)
			break
		} else if biggerPrimeBF.Cmp(targetStartingValueBiggerPrimeBF) < 0 {
			sumEstimationErrorBF.Sub(sumEstimationErrorBF, sumEstimationErrorIncrementBF)
			sumEstimationErrorIncrementBF.Quo(sumEstimationErrorIncrementBF, tenBF)
			fmt.Printf("sumEstimationErrorIncrementBF: %.12f\n", sumEstimationErrorIncrementBF)
		}
		previousDiff.Set(currentDiff)
		if i == 512 {
			panic("could not find base value")
		}
	}

	sumEstimationErrorIncrementBF.SetFloat64(0.1)

	fmt.Printf("\"run2\": %v\n", "run2")

	targetStartingValueBiggerPrime.Sub(targetStartingValueBiggerPrime, one)
	fmt.Printf("targetStartingValueSmallerPrime: %v\n", targetStartingValueBiggerPrime)
	targetStartingValueBiggerPrimeBF.SetInt(targetStartingValueBiggerPrime)

	// find base value
	for i := 0; i <= 513; i++ {
		sumEstimationErrorBF.Add(sumEstimationErrorBF, sumEstimationErrorIncrementBF)

		hardCalc, err := WIP2(sumEstimationBF, sumEstimationErrorBF, targetNumberNBF)
		if err != nil {
			fmt.Printf("Stopping due to %v", err)
			break
		}
		//    _____________________________________________
		//  ╲╱sumEst² - 2 ⋅ c1 ⋅ sumEst - 4 ⋅ numberN + c1² + sumEst - c1
		//  ────────────────────────────────────────────────────────────
		//								2
		biggerPrimeBF.Add(hardCalc, sumEstimationBF)
		biggerPrimeBF.Sub(biggerPrimeBF, sumEstimationErrorBF)
		biggerPrimeBF.Quo(biggerPrimeBF, twoBF)

		currentDiff.Sub(targetStartingValueBiggerPrimeBF, biggerPrimeBF)

		if previousDiff.Cmp(currentDiff) == 0 {
			fmt.Printf("sumEstimationErrorBF: %.12f\n", sumEstimationErrorBF)
			break
		} else if biggerPrimeBF.Cmp(targetStartingValueBiggerPrimeBF) < 0 {
			sumEstimationErrorBF.Sub(sumEstimationErrorBF, sumEstimationErrorIncrementBF)
			sumEstimationErrorIncrementBF.Quo(sumEstimationErrorIncrementBF, tenBF)
			fmt.Printf("sumEstimationErrorIncrementBF: %.12f\n", sumEstimationErrorIncrementBF)
		}
		previousDiff.Set(currentDiff)
		if i == 512 {
			panic("could not find base value")
		}
	}
}

func FindPeriods3(targetNumberN *big.Int, sumEstimation *big.Int) {
	targetNumberNBF := new(big.Float).SetInt(targetNumberN)
	sumEstimationBF := new(big.Float).SetInt(sumEstimation)
	sumEstimationErrorBF := new(big.Float).SetInt64(0)
	one := new(big.Int).SetInt64(1)
	tenBF := new(big.Float).SetInt64(10)

	// calculate starting value
	hardCalcBF, err := WIP2(sumEstimationBF, sumEstimationErrorBF, targetNumberNBF)
	if err != nil {
		panic(fmt.Sprintf("Stopping due to %v", err))
	}

	targetStartingValueHardCalc := new(big.Int)
	targetStartingValueHardCalcBF := new(big.Float)
	hardCalcBF.Int(targetStartingValueHardCalc)
	targetStartingValueHardCalc.Sub(targetStartingValueHardCalc, one)
	fmt.Printf("targetStartingValueSmallerPrime: %v\n", targetStartingValueHardCalc)
	targetStartingValueHardCalcBF.SetInt(targetStartingValueHardCalc)

	sumEstimationErrorIncrementBF := new(big.Float)

	currentDiff := new(big.Float)
	previousDiff := new(big.Float)

	sumEstimationErrorIncrementBF.SetFloat64(0.1)

	// find base value
	for i := 0; i <= 513; i++ {
		sumEstimationErrorBF.Add(sumEstimationErrorBF, sumEstimationErrorIncrementBF)

		hardCalcBF, err := WIP2(sumEstimationBF, sumEstimationErrorBF, targetNumberNBF)
		if err != nil {
			fmt.Printf("Stopping due to %v", err)
			break
		}

		currentDiff.Sub(targetStartingValueHardCalcBF, hardCalcBF)

		if previousDiff.Cmp(currentDiff) == 0 {
			fmt.Printf("sumEstimationErrorBF: %.12f\n", sumEstimationErrorBF)
			break
		} else if hardCalcBF.Cmp(targetStartingValueHardCalcBF) < 0 {
			sumEstimationErrorBF.Sub(sumEstimationErrorBF, sumEstimationErrorIncrementBF)
			sumEstimationErrorIncrementBF.Quo(sumEstimationErrorIncrementBF, tenBF)
			fmt.Printf("sumEstimationErrorIncrementBF: %.12f\n", sumEstimationErrorIncrementBF)
		}
		previousDiff.Set(currentDiff)
		if i == 512 {
			panic("could not find base value")
		}
	}

	sumEstimationErrorIncrementBF.SetFloat64(0.1)

	fmt.Printf("\"run2\": %v\n", "run2")

	targetStartingValueHardCalc.Sub(targetStartingValueHardCalc, one)
	fmt.Printf("targetStartingValueSmallerPrime: %v\n", targetStartingValueHardCalc)
	targetStartingValueHardCalcBF.SetInt(targetStartingValueHardCalc)

	// find base value
	for i := 0; i <= 513; i++ {
		sumEstimationErrorBF.Add(sumEstimationErrorBF, sumEstimationErrorIncrementBF)

		hardCalcBF, err := WIP2(sumEstimationBF, sumEstimationErrorBF, targetNumberNBF)
		if err != nil {
			fmt.Printf("Stopping due to %v", err)
			break
		}

		currentDiff.Sub(targetStartingValueHardCalcBF, hardCalcBF)

		if previousDiff.Cmp(currentDiff) == 0 {
			fmt.Printf("sumEstimationErrorBF: %.12f\n", sumEstimationErrorBF)
			break
		} else if hardCalcBF.Cmp(targetStartingValueHardCalcBF) < 0 {
			sumEstimationErrorBF.Sub(sumEstimationErrorBF, sumEstimationErrorIncrementBF)
			sumEstimationErrorIncrementBF.Quo(sumEstimationErrorIncrementBF, tenBF)
			fmt.Printf("sumEstimationErrorIncrementBF: %.12f\n", sumEstimationErrorIncrementBF)
		}
		previousDiff.Set(currentDiff)
		if i == 512 {
			panic("could not find base value")
		}
	}
}

func TryToFindPrimes(finishedWaitGroup *sync.WaitGroup, targetNumberN *big.Int, sumEstimation *big.Int, iterationValue float64) {
	sumEstimationErrorStopValue := new(big.Int).SetInt64(math.MaxInt64)

	targetNumberNBF := new(big.Float).SetInt(targetNumberN)
	smallerPrimeBF := new(big.Float)
	biggerPrimeBF := new(big.Float)
	smallerPrime := new(big.Int)
	biggerPrime := new(big.Int)
	sumEstimationBF := new(big.Float).SetInt(sumEstimation)
	sumEstimationErrorBF := new(big.Float).SetInt64(0)
	sumEstimationError := new(big.Int).SetInt64(0)
	negOneBF := new(big.Float).SetInt64(-1)
	twoBF := new(big.Float).SetInt64(2)

	sumEstimationErrorIncrement := new(big.Float).SetFloat64(iterationValue)

	// DEBUG
	// hardcalc               period~0.122316, offset sumEstimation~0.061380
	// smallerPrimeBF         period~0.27873 , offset sumEstimation~0.069950
	// biggerPrimeBF          period~0.21797 , offset sumEstimation~0.05469

	// smallerPrimeBF y = 0.27873*x + 0.069950
	// biggerPrimeBF  y = 0.21797*x + 0.05469
	fmt.Printf("diff: %20v %20v %20v %20v %20v\n", "sumEstimationErrorBF", "hardCalc", "smallerPrimeBF", "biggerPrimeBF", "diff")
	debugAdjustsumEstimationErrorBF := new(big.Float).SetFloat64(0)
	sumEstimationErrorBF.Add(sumEstimationErrorBF, debugAdjustsumEstimationErrorBF)
	// sumEstimationErrorStopValue.SetInt64(100)

	for {
		// https://www.wolframalpha.com/input/?i=solve+%7Ba*b%3Dn%2C+a%2Bb%2Bc%3Dm%2C+a%3E0%2C+b%3E0%2C+n%3E0%2C+m%3E0%2C+b%3Ea%7D+for+%7Bc%2Ca%2Cb%7D+over+the+integers
		// c<m-2*sqrt(n)
		//
		// solve([p1*p2 = numberN, p1+p2+c1=sumEst], [p1, p2]);
		// p1(x) := -(sqrt(sumEst^2  - 2 * c1 * sumEst - 4 * numberN + c1^2 ) - sumEst + c1)/2;
		// p2(x) :=  (sqrt(sumEst^2  - 2 * c1 * sumEst - 4 * numberN + c1^2 ) + sumEst - c1)/2;

		// p1(x) := -(sqrt(67620573609986982832142^2  - 2 * x * 67620573609986982832142 - 4 * 1137501950159415429602603162643696967245621181 + x^2 ) - 67620573609986982832142 + x)/2;
		// p2(x) :=  (sqrt(67620573609986982832142^2  - 2 * x * 67620573609986982832142 - 4 * 1137501950159415429602603162643696967245621181 + x^2 ) + 67620573609986982832142 - x)/2;
		// ps(x) := p1(x) + p2(x);
		// pst(x) := 67551153748445296705730;
		// pm(x) := p1(x) * p2(x);
		// pt(x) := 1137501950159415429602603162643696967245621181;
		// plot2d([ps, pst], [x, 0, 79419861541686126412]);
		hardCalc, err := WIP2(sumEstimationBF, sumEstimationErrorBF, targetNumberNBF)
		if err != nil {
			fmt.Printf("Stopping due to %v\n", err)
			break
		}
		//    _____________________________________________
		// -╲╱sumEst² - 2 ⋅ c1 ⋅ sumEst - 4 ⋅ numberN + c1² - sumEst + c1
		//  ────────────────────────────────────────────────────────────
		//								2
		smallerPrimeBF.Sub(hardCalc, sumEstimationBF)
		smallerPrimeBF.Add(smallerPrimeBF, sumEstimationErrorBF)
		smallerPrimeBF.Quo(smallerPrimeBF, twoBF)
		smallerPrimeBF.Mul(smallerPrimeBF, negOneBF)
		//    _____________________________________________
		//  ╲╱sumEst² - 2 ⋅ c1 ⋅ sumEst - 4 ⋅ numberN + c1² + sumEst - c1
		//  ────────────────────────────────────────────────────────────
		//								2
		biggerPrimeBF.Add(hardCalc, sumEstimationBF)
		biggerPrimeBF.Sub(biggerPrimeBF, sumEstimationErrorBF)
		biggerPrimeBF.Quo(biggerPrimeBF, twoBF)

		// round calculated primes and check if they produce the desired target
		SetRoundInt(smallerPrime, smallerPrimeBF)
		SetRoundInt(biggerPrime, biggerPrimeBF)
		result := new(big.Int)
		result.Mul(smallerPrime, biggerPrime)

		diff := new(big.Int)
		diff.Sub(targetNumberN, result)
		diff.Abs(diff)
		fmt.Printf("diff: %20f %20f %20f %20f %20d\n", sumEstimationErrorBF, hardCalc, smallerPrimeBF, biggerPrimeBF, diff)

		if result.Cmp(targetNumberN) == 0 {
			fmt.Printf("\n##########################################\nfound it!\n")
			fmt.Printf("targetNumberN:        %v\n", targetNumberN)
			fmt.Printf("result:               %v\n", result)
			fmt.Printf("smallerPrime:         %v\n", smallerPrime)
			fmt.Printf("biggerPrime:          %v\n", biggerPrime)
			fmt.Printf("sumEstimationErrorBF: %v\n", sumEstimationError)
			fmt.Printf("##########################################\n")
			break
		}

		sumEstimationErrorBF.Add(sumEstimationErrorBF, sumEstimationErrorIncrement)
		sumEstimationErrorBF.Int(sumEstimationError)

		// stop condition
		if sumEstimationError.CmpAbs(sumEstimationErrorStopValue) >= 0 {
			fmt.Printf("Stopping due to reaching maximumum value\n")
			break
		}
	}
	finishedWaitGroup.Done()
}

// Sets destBigInt to the rounded value of sourceBigFloat
func SetRoundInt(destBigInt *big.Int, sourceBigFloat *big.Float) {
	bigFloatAdjusted := new(big.Float).SetFloat64(0.5)

	if sourceBigFloat.Sign() >= 0 {
		bigFloatAdjusted.Add(sourceBigFloat, bigFloatAdjusted)
	} else {
		bigFloatAdjusted.Sub(sourceBigFloat, bigFloatAdjusted)
	}
	bigFloatAdjusted.Int(destBigInt)
}

//   _____________________________________________
// ╲╱sumEst² - 2 * c1 * sumEst - 4 * numberN + c1²
func WIP2(sumEst *big.Float, c1 *big.Float, numberN *big.Float) (*big.Float, error) {
	// sumEst²
	sumEstimationSquared := new(big.Float)
	sumEstimationSquared.Mul(sumEst, sumEst)
	// 2 * c1 * sumEst
	twoC1SumEst := new(big.Float).SetInt64(2)
	twoC1SumEst.Mul(twoC1SumEst, c1)
	twoC1SumEst.Mul(twoC1SumEst, sumEst)
	// 4 * numberN
	fourNumberN := new(big.Float).SetInt64(4)
	fourNumberN.Mul(fourNumberN, numberN)
	// c1²
	c1Squared := new(big.Float)
	c1Squared.Mul(c1, c1)
	// sumEst² - 2 * c1 * sumEst - 4 * numberN + c1²
	sumToSquare := new(big.Float).SetFloat64(0)
	sumToSquare.Add(sumToSquare, sumEstimationSquared)
	sumToSquare.Sub(sumToSquare, twoC1SumEst)
	sumToSquare.Sub(sumToSquare, fourNumberN)
	sumToSquare.Add(sumToSquare, c1Squared)
	if sumToSquare.Sign() == -1 {
		return nil, errors.New("sumToSquare negative so calculation not possible\n")
	}
	//   _____________________________________________
	// ╲╱sumEst² - 2 * c1 * sumEst - 4 * numberN + c1²
	squaredRoot := new(big.Float)
	squaredRoot.Sqrt(sumToSquare)
	return squaredRoot, nil
}

func InverseTransform() {
	realNum, _ := new(big.Int).SetString("67551153748445296705730", 10)
	minValue, _ := new(big.Int).SetString("56706253667142884486418", 10)
	maxValue, _ := new(big.Int).SetString("75513164846306307011680", 10)

	minValueBF := new(big.Float).SetInt(minValue)
	maxValueBF := new(big.Float).SetInt(maxValue)

	//estimatedValue, _ := new(big.Float).SetString("0.5803355925313406")
	estimatedValue, _ := new(big.Float).SetString("0.5775671903588346")

	diff := new(big.Float)
	muldiff := new(big.Float)
	result := new(big.Float)

	diff.Sub(maxValueBF, minValueBF)
	muldiff.Mul(estimatedValue, diff)
	result.Add(muldiff, minValueBF)
	resultNum := new(big.Int)
	result.Int(resultNum)

	fmt.Printf("resultNum: %v\n", resultNum)
	fmt.Printf("resultNum: %v\n", realNum.Sub(realNum, resultNum))
}
