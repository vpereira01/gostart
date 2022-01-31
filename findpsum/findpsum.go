package findpsum

import (
	"errors"
	"fmt"
	"math"
	"math/big"
)

// Tries to find prime factors with a base guess of the factors sum
func TryToFindPrimeFactors(targetNumberN *big.Int, primeFactorsSumEstimatation *big.Int) TryToFindPrimeFactorsResult {
	var finalResult TryToFindPrimeFactorsResult
	results := make(chan TryToFindPrimeFactorsResult)
	stopFlag := false

	go tryToFindPrimeFactors(results, &stopFlag, targetNumberN, primeFactorsSumEstimatation, -1)
	go tryToFindPrimeFactors(results, &stopFlag, targetNumberN, primeFactorsSumEstimatation, 1)

	// Will wait for any result, if a factor is found then stopFlag is set and waits for remaing goroutines to stop
	for numberOfGoRoutinesRunning := 2; numberOfGoRoutinesRunning > 0; numberOfGoRoutinesRunning-- {
		result := <-results
		if result.foundFactors {
			stopFlag = true
			finalResult = result
		}
	}
	if stopFlag {
		return finalResult
	} else {
		return TryToFindPrimeFactorsResult{false, nil, nil}
	}
}

type TryToFindPrimeFactorsResult struct {
	foundFactors bool
	smallFactor  *big.Int
	bigFactor    *big.Int
}

// Tries to find prime factors with a base guess of the factors sum
// Base formulas:
//    p1*p2 = numberN
//    p1+p2+c1=sumEst
//
// solve([p1*p2 = numberN, p1+p2+c1=sumEst], [p1, p2]);
// p1(x) := -(sqrt(sumEst^2  - 2 * c1 * sumEst - 4 * numberN + c1^2 ) - sumEst + c1)/2;
// p2(x) :=  (sqrt(sumEst^2  - 2 * c1 * sumEst - 4 * numberN + c1^2 ) + sumEst - c1)/2;
func tryToFindPrimeFactors(
	resultsChannel chan TryToFindPrimeFactorsResult,
	stopFlag *bool,
	targetNumberN *big.Int,
	sumEstimation *big.Int,
	iterationValue float64) {

	targetNumberNBF := new(big.Float).SetInt(targetNumberN)
	smallerPrimeBF := new(big.Float)
	biggerPrimeBF := new(big.Float)
	smallerPrime := new(big.Int)
	biggerPrime := new(big.Int)
	sumEstimationBF := new(big.Float).SetInt(sumEstimation)
	sumEstimationErrorBF := new(big.Float).SetInt64(0)
	negOneBF := new(big.Float).SetInt64(-1)
	twoBF := new(big.Float).SetInt64(2)

	sumEstimationErrorIncrement := new(big.Float).SetFloat64(iterationValue)

	// https://www.wolframalpha.com/input/?i=solve+%7Ba*b%3Dn%2C+a%2Bb%2Bc%3Dm%2C+a%3E0%2C+b%3E0%2C+n%3E0%2C+m%3E0%2C+b%3Ea%7D+for+%7Bc%2Ca%2Cb%7D+over+the+integers
	// c<m-2*sqrt(n)
	// c1 could have a maximum value
	sumEstimationErrorStopValueBF := new(big.Float)

	if iterationValue < 0 {
		sumEstimationErrorStopValueBF.SetInt64(math.MinInt64)
	} else {
		sumEstimationErrorStopValueBF.SetInt64(math.MaxInt64)
	}

	for {
		// p1(x) := -(sqrt(sumEst^2  - 2 * c1 * sumEst - 4 * numberN + c1^2 ) - sumEst + c1)/2;
		// p2(x) :=  (sqrt(sumEst^2  - 2 * c1 * sumEst - 4 * numberN + c1^2 ) + sumEst - c1)/2;
		//             _____________________________________________
		// calculate ╲╱sumEst² - 2 * c1 * sumEst - 4 * numberN + c1²
		baseValue, err := getPrimeFactorsBaseValue(sumEstimationBF, sumEstimationErrorBF, targetNumberNBF)
		if err == ErrSumToSquareIsNegative && iterationValue < 0 {
			sumEstimationErrorBF.Add(sumEstimationErrorBF, sumEstimationErrorIncrement)
			continue
		}
		if err != nil {
			fmt.Printf("Stopping due to %v\n", err)
			break
		}

		//    _____________________________________________
		// -╲╱sumEst² - 2 ⋅ c1 ⋅ sumEst - 4 ⋅ numberN + c1² - sumEst + c1
		//  ────────────────────────────────────────────────────────────
		//								2
		smallerPrimeBF.Sub(baseValue, sumEstimationBF)
		smallerPrimeBF.Add(smallerPrimeBF, sumEstimationErrorBF)
		smallerPrimeBF.Quo(smallerPrimeBF, twoBF)
		smallerPrimeBF.Mul(smallerPrimeBF, negOneBF)
		//    _____________________________________________
		//  ╲╱sumEst² - 2 ⋅ c1 ⋅ sumEst - 4 ⋅ numberN + c1² + sumEst - c1
		//  ────────────────────────────────────────────────────────────
		//								2
		biggerPrimeBF.Add(baseValue, sumEstimationBF)
		biggerPrimeBF.Sub(biggerPrimeBF, sumEstimationErrorBF)
		biggerPrimeBF.Quo(biggerPrimeBF, twoBF)

		// round calculated primes and check if they produce the desired target
		setRoundInt(smallerPrime, smallerPrimeBF)
		setRoundInt(biggerPrime, biggerPrimeBF)
		result := new(big.Int)
		result.Mul(smallerPrime, biggerPrime)

		if result.Cmp(targetNumberN) == 0 {
			fmt.Printf("\n##########################################\nfound it!\n")
			fmt.Printf("targetNumberN:        %v\n", targetNumberN)
			fmt.Printf("result:               %v\n", result)
			fmt.Printf("smallerPrime:         %v\n", smallerPrime)
			fmt.Printf("biggerPrime:          %v\n", biggerPrime)
			fmt.Printf("sumEstimationErrorBF: %f\n", sumEstimationErrorBF)
			fmt.Printf("##########################################\n")
			resultsChannel <- TryToFindPrimeFactorsResult{true, smallerPrime, biggerPrime}
			return
		}

		sumEstimationErrorBF.Add(sumEstimationErrorBF, sumEstimationErrorIncrement)

		// stop condition
		if iterationValue < 0 && sumEstimationErrorBF.Cmp(sumEstimationErrorStopValueBF) <= 0 {
			fmt.Printf("Stopping due to reaching maximumum value\n")
			break
		} else if iterationValue > 0 && sumEstimationErrorBF.Cmp(sumEstimationErrorStopValueBF) >= 0 {
			fmt.Printf("Stopping due to reaching maximumum value\n")
			break
		}

		// stop condition
		if *stopFlag {
			fmt.Printf("Stopping due to stopFlag\n")
			break
		}
	}
	resultsChannel <- TryToFindPrimeFactorsResult{false, nil, nil}
}

// Sets destBigInt to the rounded value of sourceBigFloat
func setRoundInt(destBigInt *big.Int, sourceBigFloat *big.Float) {
	bigFloatAdjusted := new(big.Float).SetFloat64(0.5)

	if sourceBigFloat.Sign() >= 0 {
		bigFloatAdjusted.Add(sourceBigFloat, bigFloatAdjusted)
	} else {
		bigFloatAdjusted.Sub(sourceBigFloat, bigFloatAdjusted)
	}
	bigFloatAdjusted.Int(destBigInt)
}

var ErrSumToSquareIsNegative = errors.New("sumToSquare negative so calculation not possible")

//   _____________________________________________
// ╲╱sumEst² - 2 * c1 * sumEst - 4 * numberN + c1²
func getPrimeFactorsBaseValue(sumEst *big.Float, c1 *big.Float, numberN *big.Float) (*big.Float, error) {
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
		return nil, ErrSumToSquareIsNegative
	}
	//   _____________________________________________
	// ╲╱sumEst² - 2 * c1 * sumEst - 4 * numberN + c1²
	squaredRoot := new(big.Float)
	squaredRoot.Sqrt(sumToSquare)
	return squaredRoot, nil
}
