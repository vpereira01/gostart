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
	fourBigFloat := new(big.Float).SetInt64(4)
	// we want to find perfect squares made by 4*numberN + subEst^2
	// so we first start by doing the startNum = math.ceil(math.sqrt(4*numberN))
	//   since there too big errors with calculating this we use math.floor instead
	//   also because the sqrt(subEst^2  + 4 * numberN) must be even, then startNum must be also even
	fourTimesNumberN := new(big.Float).SetPrec(bigFloatsDefaultPrecision)
	fourTimesNumberN.SetInt64(4)
	fourTimesNumberN.Mul(fourTimesNumberN, numberN)

	startNum := new(big.Float).SetPrec(bigFloatsDefaultPrecision)
	startNum.Sqrt(fourTimesNumberN)
	roundDownToPair(startNum)

	// then we check value1 = startNum**2 - 4*numberN so we get a candidate for value1 = subEst^2
	// if sqrt(value1) % 1 != 0 ignore value, increment the startNum
	value1 := new(big.Float).SetPrec(bigFloatsDefaultPrecision)

	// DEBUG
	stopValue := new(big.Float).SetPrec(bigFloatsDefaultPrecision)
	stopValue.SetString("154357600259902")

	countTries := uint(0)
	for {

		value1.Mul(startNum, startNum)       // startNum**2
		value1.Sub(value1, fourTimesNumberN) // startNum**2 - 4*numberN
		// if the value is negative, skip the sqrt() and value check
		if value1.Sign() == -1 {
			fmt.Printf("-")
		} else {
			value1DivByFour := new(big.Float)
			value1DivByFour.Quo(value1, fourBigFloat)
			if value1DivByFour.IsInt() {
				value1.Sqrt(value1) // sqrt(startNum**2 - 4*numberN)
				if value1.IsInt() {
					fmt.Printf("\n##########################################\nfound it!\n")
					fmt.Printf("value1: %f\n", value1)
					fmt.Printf("##########################################\n")
					return value1
				}
			}
		}

		// DEBUG
		if value1.Cmp(stopValue) > 0 {
			panic("failed to find value")
		}

		countTries++
		startNum.Add(startNum, twoBigFloat)
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
