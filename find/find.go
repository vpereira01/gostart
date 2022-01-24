package find

import (
	"fmt"
	"math/big"
)

func FindPrimesSub(numberN *big.Int) *big.Int {
	numberNBigFloat := new(big.Float).SetInt(numberN)
	primesSubBigFloat := findPrimesSub(numberNBigFloat)
	primesSub := new(big.Int)
	primesSubBigFloat.Int(primesSub)
	return primesSub
}

// find 2 prime factors using the primes' subtraction as faster auxiliary
// base equations (maxima notation):
// solve([p1*p2 = numberN, p2-p1=subEst], [p1, p2]); (only positive solutions)
//   p1(subEst) :=  (sqrt(subEst^2  + 4 * numberN) - subEst)/2;
//   p2(subEst) :=  (sqrt(subEst^2  + 4 * numberN) + subEst)/2;
func findPrimesSub(numberN *big.Float) *big.Float {
	// in order to avoid big float precision errors lets assume a default precision double the number n
	// this is required because big floats initialize a precision and if a mul() requires a bigger precision
	//   to store its results the precision is not increase but the result is rounded to meet the precision
	bigFloatsDefaultPrecision := numberN.Prec() * 2

	twoBigFloat := new(big.Float).SetInt64(2)
	twoBigInt := new(big.Int).SetInt64(2)
	hundredBigInt := new(big.Int).SetInt64(100)
	// fourBigFloat := new(big.Float).SetInt64(4)
	// fourBigInt := new(big.Int).SetInt64(4)
	// we want to find perfect squares made by 4*numberN + subEst^2
	// so we first start by doing the startNum = math.ceil(math.sqrt(4*numberN))
	//   since there too big errors with calculating this we use math.floor instead
	//   also because the sqrt(subEst^2  + 4 * numberN) must be even, then startNum must be also even
	fourTimesNumberNBigFloat := new(big.Float).SetPrec(bigFloatsDefaultPrecision)
	fourTimesNumberNBigFloat.SetInt64(4)
	fourTimesNumberNBigFloat.Mul(fourTimesNumberNBigFloat, numberN)
	fourTimesNumberNBigInt := new(big.Int)
	fourTimesNumberNBigFloat.Int(fourTimesNumberNBigInt)

	startNumBigFloat := new(big.Float).SetPrec(bigFloatsDefaultPrecision)
	startNumBigFloat.Sqrt(fourTimesNumberNBigFloat)
	roundDownToPair(startNumBigFloat)
	startNumBigInt := new(big.Int)
	startNumBigFloat.Int(startNumBigInt)

	// then we check value1BigFloat = startNum**2 - 4*numberN so we get a candidate for value1BigFloat = subEst^2
	// if sqrt(value1BigFloat) % 1 != 0 ignore value, increment the startNum
	value1BigFloat := new(big.Float).SetPrec(bigFloatsDefaultPrecision)
	value1BigInt := new(big.Int)

	// DEBUG
	stopValueBigFloat := new(big.Float).SetPrec(bigFloatsDefaultPrecision)
	stopValueBigFloat.SetString("154357600259902")

	countTries := uint(0)
	countFailedSquareEstimation := uint(0)
	for {
		// value1BigFloat.Mul(startNumBigFloat, startNumBigFloat)       // startNum**2
		value1BigInt.Mul(startNumBigInt, startNumBigInt) // startNum**2
		// value1BigFloat.Sub(value1BigFloat, fourTimesNumberNBigFloat) // startNum**2 - 4*numberN
		value1BigInt.Sub(value1BigInt, fourTimesNumberNBigInt) // startNum**2 - 4*numberN
		// if the value is negative, skip the sqrt() and value check
		if value1BigInt.Sign() == -1 {
			fmt.Printf("-")
		} else {
			// check if possible perfect square
			// value1DivByFour := new(big.Float)
			// value1DivByFour.Quo(value1BigFloat, fourBigFloat)
			// value1ModByFour := new(big.Int)
			// value1ModByFour.Mod(value1BigInt, fourBigInt)
			// if value1DivByFour.IsInt() {
			// if value1ModByFour.Int64() == 0 {
			value1Mod10BigInt := new(big.Int)
			value1Mod10BigInt.Mod(value1BigInt, hundredBigInt)
			value1Mod10 := value1Mod10BigInt.Int64()
			// , , , , , , , , , , , , , , , , , , ,
			if value1Mod10 == 0 ||
				value1Mod10 == 1 ||
				value1Mod10 == 4 ||
				value1Mod10 == 9 ||
				value1Mod10 == 16 ||
				value1Mod10 == 21 ||
				value1Mod10 == 24 ||
				value1Mod10 == 25 ||
				value1Mod10 == 29 ||
				value1Mod10 == 36 ||
				value1Mod10 == 41 ||
				value1Mod10 == 44 ||
				value1Mod10 == 49 ||
				value1Mod10 == 56 ||
				value1Mod10 == 61 ||
				value1Mod10 == 64 ||
				value1Mod10 == 69 ||
				value1Mod10 == 76 ||
				value1Mod10 == 81 ||
				value1Mod10 == 84 ||
				value1Mod10 == 89 ||
				value1Mod10 == 96 {
				countFailedSquareEstimation++
				value1BigFloat.SetInt(value1BigInt)
				value1BigFloat.Sqrt(value1BigFloat) // sqrt(startNum**2 - 4*numberN)
				if value1BigFloat.IsInt() {
					fmt.Printf("\n##########################################\nfound it!\n")
					fmt.Printf("value1: %f countTries: %v countFailedSquareEstimation: %v\n", value1BigFloat, countTries, countFailedSquareEstimation)
					fmt.Printf("##########################################\n")
					return value1BigFloat
				}
			}
		}

		// DEBUG
		if value1BigFloat.Cmp(stopValueBigFloat) > 0 {
			panic("failed to find value")
		}

		countTries++
		startNumBigFloat.Add(startNumBigFloat, twoBigFloat)
		startNumBigInt.Add(startNumBigInt, twoBigInt)
		if countTries%500000 == 0 {
			fmt.Printf("*")
			// fmt.Printf("value1: %f\n", value1)
		}
	}
}

// Sets bigFloat to its ceiling value, rounding up
// warning: changes precision
func roundUp(bigFloat *big.Float) {
	// if already rounded, it's an int, nothing to do
	if bigFloat.IsInt() {
		return
	}

	oneBigInt := new(big.Int).SetInt64(1)
	tempBigInt := new(big.Int)

	bigFloat.Int(tempBigInt)
	tempBigInt.Add(tempBigInt, oneBigInt)
	bigFloat.SetInt(tempBigInt)
}

// Sets bigFloat to its floor value and then subtract one if odd
// warning: changes precision
func roundDownToPair(bigFloat *big.Float) {
	// if already rounded, it's an int, nothing to do
	if bigFloat.IsInt() {
		return
	}

	tempBigInt := new(big.Int)

	bigFloat.Int(tempBigInt)

	oneBigInt := new(big.Int).SetInt64(1)
	twoBigInt := new(big.Int).SetInt64(2)
	modBigInt := new(big.Int)
	modBigInt.Mod(tempBigInt, twoBigInt)
	if modBigInt.Int64() == 1 {
		tempBigInt.Sub(tempBigInt, oneBigInt)
	}
	bigFloat.SetInt(tempBigInt)
}
